package webhp

import (
	"fmt"
	"strconv"
	"time"
)

type DataContainer struct {
	Data []responseResult
}

func newDataContainer() DataContainer {
	return DataContainer{
		Data: make([]responseResult, 0, 0),
	}
}

func (c *DataContainer) Add(result responseResult) {
	c.Data = append(c.Data, result)
}

func (c *DataContainer) Count() int {
	return len(c.Data)
}

func (c *DataContainer) AverageResponseTime() time.Duration {
	var sum int64 = 0
	var count int64 = 0
	for _, d := range c.Data {
		if d.err != nil {
			continue
		}

		sum = sum + d.duration.Nanoseconds()
		count = count + 1
	}

	if count == 0 {
		return 0
	}else{
		return time.Duration(float64(sum) / float64(count))
	}
}

func (c *DataContainer) PrintHttpStatus() {
	count2xx := 0
	count3xx := 0
	count4xx := 0
	count5xx := 0
	httpError := 0
	successful := 0

	for _, d := range c.Data {
		if d.err != nil {
			httpError = httpError + 1
			continue
		}

		successful = successful + 1

		switch strconv.Itoa(d.httpStatus)[0] {
		case '2':
			count2xx = count2xx + 1
		case '3':
			count3xx = count3xx + 1
		case '4':
			count4xx = count4xx + 1
		case '5':
			count5xx = count5xx + 1
		default:
		}
	}
	fmt.Println("total requests", c.Count())
	fmt.Println("2xx responses", count2xx, "responses", "percentage", fmt.Sprintf("%.3f%%", float64(count2xx)/float64(c.Count())*100.0))
	fmt.Println("3xx responses", count3xx, "responses", "percentage", fmt.Sprintf("%.3f%%", float64(count3xx)/float64(c.Count())*100.0))
	fmt.Println("4xx responses", count4xx, "responses", "percentage", fmt.Sprintf("%.3f%%", float64(count4xx)/float64(c.Count())*100.0))
	fmt.Println("5xx responses", count5xx, "responses", "percentage", fmt.Sprintf("%.3f%%", float64(count5xx)/float64(c.Count())*100.0))
	fmt.Println("successful requests", successful, "requests", "percentage", fmt.Sprintf("%.3f%%", float64(successful)/float64(c.Count())*100.0))
	fmt.Println("http error requests", httpError, "requests", "percentage", fmt.Sprintf("%.3f%%", float64(httpError)/float64(c.Count())*100.0))
}
