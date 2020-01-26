package webhp

import (
	"context"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultDurationChannelBufferSize = 100
	DefaultMaximumConcurrency        = 10000
)

type LoadGenerator struct {
	Method   string
	Url      *url.URL
	Data     DataContainer
	Duration time.Duration

	startTestingTime      time.Time
	stopTestingTime       time.Time
	actualRate            float64
	delay                 time.Duration
	responseResultChannel chan *responseResult
	callCount             int64
	concurrentMeasurementChannel chan struct{}
	maximumConcurrency int
}

func NewLoadGenerator(method, rawurl string, request_per_sec float64, duration time.Duration) LoadGenerator {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return LoadGenerator{
		Url:      u,
		Method:   method,
		Data:     newDataContainer(),
		Duration: duration,

		delay:                 time.Duration(1 / request_per_sec * 1e9),
		responseResultChannel: make(chan *responseResult, DefaultDurationChannelBufferSize),
		concurrentMeasurementChannel: make(chan struct{}, DefaultMaximumConcurrency),
		maximumConcurrency: 0,
	}
}

func (lg LoadGenerator) PrintInfo() {
	duration := time.Now().Sub(lg.startTestingTime)
	fmt.Println("Average response time", lg.Data.AverageResponseTime())
	fmt.Println("Maximum concurrency", lg.maximumConcurrency)
	lg.Data.PrintHttpStatus()
	fmt.Println("Testing Duration", duration)
	fmt.Println("Throughput", lg.Data.Count()/int(duration.Seconds()), "request/sec")
}

func (lg *LoadGenerator) Execute() {
	go lg.startBackgroundProcess()
	go lg.captureMaximumConcurrency()

	lg.startTestingTime = time.Now()
	done := time.After(lg.Duration * time.Second)

	for {
		select {
		case <-done:
			lg.stopTestingTime = time.Now()
			lg.PrintInfo()

			for len(lg.responseResultChannel) != 0 {
			}
			return
		default:
			go lg.httpCallHandler()
			<-time.Tick(lg.delay)
		}
	}
}

func (lg *LoadGenerator) startBackgroundProcess() {
	done := time.After(lg.Duration * time.Second)

	for {
		select {
		case <-done:
			return
		default:
			result := <-lg.responseResultChannel
			lg.Data.Add(*result)
		}

	}
}

func (lg *LoadGenerator) captureMaximumConcurrency() {
	done := time.After(lg.Duration * time.Second)

	for {
		select {
		case <-done:
			return
		default:
			l := len(lg.concurrentMeasurementChannel)
			if lg.maximumConcurrency < l {
				lg.maximumConcurrency = l
			}
		}
	}
}

func (lg LoadGenerator) httpCallHandler() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	select {
	case result := <-lg.invokeHttpCall(ctx):
		lg.responseResultChannel <- result
	case <-ctx.Done():
	}
}

func (lg *LoadGenerator) invokeHttpCall(ctx context.Context) <-chan *responseResult {
	done := make(chan *responseResult, 1)

	rr := newErrorResponseResult(errors.New("http request timeout"))
	defer func() { done <- &rr }()

	req := http.Request{
		Method: lg.Method,
		URL:    lg.Url,
	}

	lg.concurrentMeasurementChannel <- struct{}{}

	client := http.Client{}
	defer client.CloseIdleConnections()

	start := time.Now()
	res, err := client.Do(&req)

	defer res.Body.Close()
	if err != nil {
		newErrorResponseResult(err)
		return done
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		newErrorResponseResult(err)
		return done
	}

	<-lg.concurrentMeasurementChannel

	rr = newResponseResult(time.Now().Sub(start), int64(len(b)), res.StatusCode)

	return done
}
