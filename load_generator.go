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
	DEFAULT_DURATION_CHANNEL_BUFFER_SIZE = 100
)

type LoadGenerator struct {
	Method string
	Url    *url.URL
	Data   DataContainer

	actualRate            float64
	delay                 time.Duration
	responseResultChannel chan *responseResult
	callCount             int64
	done                  chan struct{}
}

func NewLoadGenerator(method, rawurl string, request_per_sec float64) LoadGenerator {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return LoadGenerator{
		Url:    u,
		Method: method,
		Data:   newDataContainer(),

		delay:                 time.Duration(1 / request_per_sec * 1e9),
		responseResultChannel: make(chan *responseResult, DEFAULT_DURATION_CHANNEL_BUFFER_SIZE),
		done:                  make(chan struct{}, 1),
	}
}

func (lg LoadGenerator) Execute() {
	go lg.responseResultCalculator()
	for {
		<-time.Tick(lg.delay)
		select {
		case <- lg.done:
			fmt.Println("Average response time", lg.Data.AverageResponseTime())
			lg.Data.PrintHttpStatus()
			return
		default:
			go lg.httpCallHandler()
		}
	}
}

func (lg LoadGenerator) Done() {
	lg.done <- struct{}{}
}

func (lg *LoadGenerator) responseResultCalculator() {
	for {
		result := <-lg.responseResultChannel
		lg.Data.Add(*result)
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

func (lg LoadGenerator) invokeHttpCall(ctx context.Context) <-chan *responseResult {
	done := make(chan *responseResult, 1)

	rr := newErrorResponseResult(errors.New("http request timeout"))
	defer func() { done <- &rr }()

	req := http.Request{
		Method: lg.Method,
		URL:    lg.Url,
	}
	client := http.Client{}

	start := time.Now()
	res, err := client.Do(&req)
	if err != nil {
		newErrorResponseResult(err)
		return done
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		newErrorResponseResult(err)
		return done
	}

	rr = newResponseResult(time.Now().Sub(start), int64(len(b)), res.StatusCode)

	return done
}
