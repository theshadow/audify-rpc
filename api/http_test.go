package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus/hooks/test"
)

var ErrorExpected = errors.New("test failure, this is expected")

// Test that Retrying() will make the correct attempts and return the expected errors.
func TestRetryMaxRetry(t *testing.T) {
	retries := 0
	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		retries++
		return nil, ErrorExpected
	}

	_, err := Retrying(3, doer)(context.Background(), &http.Client{}, &http.Request{})

	if retries != 3 {
		t.Logf("expected 3 attempts, instead had %d", retries)
		t.Fail()
	}

	e, ok := err.(ErrorMaxRetryAttempts)
	if !ok {
		t.Logf("expected error of %T instead received %T", ErrorMaxRetryAttempts{}, err)
		t.Fail()
	}

	for i, e := range e.Errors {
		if e != ErrorExpected {
			t.Logf("unexpected error, expected %T instead received %T at offset %d", ErrorMaxRetryAttempts{}, e, i)
			t.Fail()
		}
	}
}

// Test that Retrying() returns on the first available success.
func TestRetryReturnOnSuccess(t *testing.T) {
	expected := &http.Response{}
	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		return expected, nil
	}

	actual, err := Retrying(1, doer)(context.Background(), &http.Client{}, &http.Request{})

	if err != nil {
		t.Logf("unexpected error received.")
		t.Fail()
	}

	if actual != expected {
		t.Logf("unexpected response instance returned")
		t.Fail()
	}
}

// Test that BackingOff() waits the specified number of nanoseconds.
func TestBackOffWaitOnError(t *testing.T) {
	delay := uint(1000)

	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		return nil, ErrorExpected
	}

	start := time.Now()
	_, err := BackingOff(delay, doer)(context.Background(), &http.Client{}, &http.Request{})
	elapsed := time.Since(start)

	if err != ErrorExpected {
		t.Logf("unexpected error type %T", err)
		t.Fail()
	}

	if elapsed.Nanoseconds() <= int64(delay) {
		t.Logf("Expected to wait %d nanoseconds, instead waited: %d", delay, elapsed.Nanoseconds())
	}
}

func TestBackOffReturnsOnSuccess(t *testing.T) {
	expected := &http.Response{}
	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		return expected, nil
	}

	actual, err := BackingOff(1000, doer)(context.Background(), &http.Client{}, &http.Request{})

	if err != nil {
		t.Logf("unexpected error: %s", err)
		t.Fail()
	}

	if actual != expected {
		t.Logf("unexpected response instance")
		t.Fail()
	}
}

func TestLogging(t *testing.T) {
	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		return nil, ErrorExpected
	}

	l, _ := test.NewNullLogger()
	u, _ := url.Parse("https://www.example.com")

	actual, err := Logging(l, doer)(context.Background(), &http.Client{}, &http.Request{URL:u})

	if actual != nil {
		t.Logf("response was expected to be nil, instead received %#v", err)
		t.Fail()
	}

	if err != ErrorExpected {
		t.Logf("error was expected to be ErrorExpected, instead received %#v", err)
		t.Fail()
	}
}

type MockCacher struct {
	GetFn func(key string) (interface{}, bool, error)
	SetFn func(key string, x interface{}, d time.Duration) error
}

func (c MockCacher) Get(key string) (interface{}, bool, error) {
	return c.GetFn(key)
}
func (c MockCacher) Set(key string, x interface{}, d time.Duration) error {
	return c.SetFn(key, x, d)
}

func TestCachingHit(t *testing.T) {
	cached := &http.Response{Status: "200 OK"}

	called := new(bool)
	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		*called = true
		return nil, nil
	}

	c := MockCacher{
		GetFn: func(key string) (interface{}, bool, error) {
			return *cached, true, nil
		},
	}

	l, _ := test.NewNullLogger()
	u, _ := url.Parse("https://www.example.com")

	actual, err := Caching(c, l, doer)(context.Background(), &http.Client{}, &http.Request{URL:u})

	if actual.Status != "200 OK" {
		t.Logf("expected Status to be '200 OK', instead received '%s'", actual.Status)
		t.Fail()
	}

	if err != nil {
		t.Logf("unexpected error '%s'", err)
		t.Fail()
	}

	if *called {
		t.Logf("unexpected call of doer()")
		t.Fail()
	}
}

type Dummy struct {}

func TestCachingMissCastingFailure(t *testing.T) {
	expected := &http.Response{Status:"200 OK"}
	cached := &Dummy{}
	u, _ := url.Parse("https://www.example.com")

	called := new(bool)
	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		*called = true
		return expected, nil
	}

	c := MockCacher{
		GetFn: func(key string) (interface{}, bool, error) {
			return *cached, true, nil
		},
		SetFn: func(key string, x interface{}, d time.Duration) error {
			if key != u.String() {
				t.Logf("expected key of '%s' instead received '%s'", u.String(), key)
				t.Fail()
			}

			r, _ := x.(http.Response)
			if r.Status != expected.Status {
				t.Logf("expected test object, received %#v instead", x)
				t.Fail()
			}
			return nil
		},
	}

	l, _ := test.NewNullLogger()

	actual, err := Caching(c, l, doer)(context.Background(), &http.Client{}, &http.Request{URL:u})

	if actual.Status != "200 OK" {
		t.Logf("expected Status to be '200 OK', instead received '%s'", actual.Status)
		t.Fail()
	}

	if err != nil {
		t.Logf("unexpected error '%s'", err)
		t.Fail()
	}

	if !*called {
		t.Logf("exepected call of doer() never occurred")
		t.Fail()
	}
}