package webhp

import "time"

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

func (c *DataContainer) AverageResponseTime() time.Duration{
	var sum int64 = 0
	var count int64 = 0
	for _, d := range c.Data{
		if d.err != nil {
			continue
		}

		sum = sum + d.duration.Nanoseconds()
		count = count + 1
	}

	return time.Duration(float64(sum) / float64(count))
}