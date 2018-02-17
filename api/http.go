package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// A Doer defines the interface compatible with ctxhttp.Do
type Doer func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error)

// Logging will log any error returned from do() and then return.
func Logging(l *log.Logger, do Doer) Doer {
	return func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
		}

		l.Infof("Req: %s", req.URL.String())
		resp, err := do(ctx, client, req)
		if err != nil {
			l.Warn(err)
		}

		if resp != nil {
			l.Infof("Resp: %s", resp.Status)
		}
		return resp, err
	}
}

// Caching will attempt to pull the result from the cache, if it's a cache miss it will make the actual request.
// A failure during a cache read will result in cache miss behavior.
func Caching(c Cacher, l *log.Logger, do Doer) Doer {
	return func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
		}
		// try to grab the cached result if we encounter an error
		// gracefully fail.
		data, found, err := c.Get(req.URL.String())
		if err != nil {
			l.Warnf("unable to access cache, %s", err)
		}

		// cache found
		if found == true {
			// can we cast the data?
			if re, ok := data.(http.Response); data != nil && ok {
				// cache hit
				return &re, nil
			} else {
				l.Warnf("unable to cast cached data to response: %#v", data)
			}
		}

		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
		}
		// cache miss, perform the request
		r, err := do(ctx, client, req)
		if err != nil {
			return nil, err
		}

		// attempt to cache for next time.
		if err = c.Set(req.URL.String(), *r, time.Minute * 15); err != nil {
			return nil, err
		}

		return r, nil
	}
}

// BackingOff() will call do() and if it returns an error it will wait nsecs nanoseconds before returning.
func BackingOff(nsecs uint, do Doer) Doer {
	return func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
		}
		resp, err := do(ctx, client, req)
		if err != nil {
			time.Sleep(time.Millisecond * time.Duration(nsecs))
		}
		return resp, err
	}
}

// Retrying() will call do() and if an error is returned it will make an additional number of attempts equal to attempts.
func Retrying(attempts uint, do Doer) Doer {
	return func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		errors := ErrorMaxRetryAttempts{}
		for i := uint(0); i < attempts; i++ {
			select {
			case <-ctx.Done():
				return nil, context.Canceled
			default:
			}
			resp, err := do(ctx, client, req)
			if err == nil {
				return resp, nil
			}
			errors.Append(err)
		}
		return nil, errors
	}
}

// ErrorMaxRetryAttempts is returned when Retrying() exhausts all of its attempts.
type ErrorMaxRetryAttempts struct {
	Attempts uint
	Errors []error
	Message string
}

func (e ErrorMaxRetryAttempts) Error() string {
	return fmt.Sprintf("failed to make request made %d attempts", e.Attempts)
}

func (e ErrorMaxRetryAttempts) Append(err error) {
	e.Errors = append(e.Errors, err)
}
