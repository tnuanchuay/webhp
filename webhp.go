package webhp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DEFAULT_DURATION_CHANNEL_BUFFER_SIZE = 1000000
)

type LoadGenerator struct {
	Url             string
	actualRate      float64
	delay           time.Duration
	durationChannel chan time.Duration
	callCount       int64
}

func NewLoadGenerator(url string, request_per_sec float64) LoadGenerator {
	period := 1 / request_per_sec * 1e9
	return LoadGenerator{
		Url:             url,
		delay:           time.Duration(period),
		durationChannel: make(chan time.Duration, DEFAULT_DURATION_CHANNEL_BUFFER_SIZE),
	}
}

func (lg LoadGenerator) Test() {

	go lg.durationCalculator()

	for {
		start := time.Now()
		<-time.Tick(lg.delay)
		go lg.httpCaller()
		lg.durationChannel <- time.Now().Sub(start)
	}
}

func (lg *LoadGenerator) durationCalculator() {
	//var sum int64
	for {
		d := <- lg.durationChannel
		fmt.Println(d)
	}
}

func (lg LoadGenerator) httpCaller() {
	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	done := make(chan bool)

	go func(d chan bool) {
		defer func() { d <- true }()
		r, err := http.Get(lg.Url)
		if err == nil {
			ioutil.ReadAll(r.Body)
		}
	}(done)

	select {
	case <-ctx.Done():
	case <-done:
	}
}