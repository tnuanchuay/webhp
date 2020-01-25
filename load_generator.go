package webhp

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
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

	actualRate            float64
	delay                 time.Duration
	responseResultChannel chan *ResponseResult
	callCount             int64
}

func NewLoadGenerator(method, rawurl string, request_per_sec float64) LoadGenerator {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return LoadGenerator{
		Url:                   u,
		Method:                method,
		delay:                 time.Duration(1 / request_per_sec * 1e9),
		responseResultChannel: make(chan *ResponseResult, DEFAULT_DURATION_CHANNEL_BUFFER_SIZE),
	}
}

func (lg LoadGenerator) Execute() {
	go lg.responseResultCalculator()
	for {
		<- time.Tick(lg.delay)
		go lg.httpCallHandler()
	}
}

func (lg *LoadGenerator) responseResultCalculator() {
	for {
		result := <-lg.responseResultChannel
		if result == nil {
			log.Println("timeout")
			continue
		}
		if result.err != nil {
			log.Println(result.err)
			continue
		}

		fmt.Println(result.duration)
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

func (lg LoadGenerator) invokeHttpCall(ctx context.Context) <-chan *ResponseResult {
	done := make(chan *ResponseResult, 1)
	var responseResult *ResponseResult = nil
	defer func() { done <- responseResult }()

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

	r := newResponseResult(time.Now().Sub(start), int64(len(b)), res.StatusCode)
	responseResult = &r

	return done
}
