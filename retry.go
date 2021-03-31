package retry

import (
	"errors"
	"math"
	"net/http"
	"time"
)

// DurationFunc describes how long to wait between retries. It receives the
// retry count, starting from 1. Retry count is passed as a time.Duration value
// to avoid type conversion.
type DurationFunc func(time.Duration) time.Duration

// ExpDuration returns an exponential DurationFunc, starting from given base.
func ExpDuration(base time.Duration) DurationFunc {
	return func(i time.Duration) time.Duration {
		res := float64(base) * math.Exp(float64(i-1))
		return time.Duration(res)
	}
}

// Retryer is an interface for handling retries.
//
// The Next method returns false when given error is nil, or when maximum
// number of attempts is reached.
// It always returns true on the first call, making it usable in a for loop.
// HttpNext behaves just like Next, except it will retry on server errors (status >= 500)
type Retryer interface {
	Next(error) bool
	HttpNext(*http.Response, error) bool
}

// New returns a Retryer which allows up to n attempts, waiting between
// subsequent attempts according to given DurationFunc.
func New(n int, df DurationFunc) Retryer {
	if n < 1 {
		n = 1
	}
	return &retryer{
		maxRetires:   time.Duration(n),
		durationFunc: df,
	}
}

// Exp returns a new Retryer with exponential DurationFunc, starting from given
// base.
func Exp(n int, base time.Duration) Retryer {
	return New(n, ExpDuration(base))
}

type RetryerFactory func() Retryer

// Factory returns a reusable function producing Retryers with given
// configuration.
func Factory(n int, df DurationFunc) RetryerFactory {
	return func() Retryer {
		return New(n, df)
	}
}

type retryer struct {
	i            time.Duration
	maxRetires   time.Duration
	durationFunc DurationFunc
}

func (r *retryer) Next(err error) bool {
	// Always run the first iteration
	if r.i == 0 {
		r.i++
		return true
	}

	// Exit on success, or if max retires reached
	if err == nil || r.i >= r.maxRetires {
		return false
	}

	// Sleep on error
	d := r.durationFunc(r.i)
	time.Sleep(d)
	r.i++
	return true
}

func (r *retryer) HttpNext(res *http.Response, err error) bool {
	if err == nil && res != nil && res.StatusCode >= 500 {
		// Create a non-nil error, so that Next will cause a retry
		err = errors.New("server error")
	}
	return r.Next(err)
}
