package webhp

import (
	"time"
)

type responseResult struct {
	duration   time.Duration
	httpStatus int
	size       int64
	err        error
}

func newResponseResult(duration time.Duration, size int64, httpStatus int) responseResult {
	return responseResult{
		duration:   duration,
		httpStatus: httpStatus,
		size:       size,
	}
}

func newErrorResponseResult(err error) responseResult {
	return responseResult{
		err: err,
	}
}
