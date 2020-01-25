package webhp

import (
	"time"
)

type ResponseResult struct {
	duration   time.Duration
	httpStatus int
	size       int64
	err        error
}

func newResponseResult(duration time.Duration, size int64, httpStatus int) ResponseResult {
	return ResponseResult{
		duration:   duration,
		httpStatus: httpStatus,
		size:       size,
	}
}

func newErrorResponseResult(err error) ResponseResult {
	return ResponseResult{
		err: err,
	}
}
