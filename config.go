package failsafe

import (
	"fmt"
	"runtime"
	"time"

	"github.com/golang/glog"
)

// Abstract PanicHandlers, handles any panics happens inside
// the failsafe call. Receives the interface provided to the
// panic as parameter.
type PanicHandler func(interface{})

func NewConfig() *Config {
	return &Config{
		MaxRetry:         1,
		CrashRecoveryOff: false,
		RetryOnlyOnErr:   make([]error, 0),
		PanicHandlers: []PanicHandler{
			logPanic,
		},
	}
}

func (c *Config) Retries(t int) *Config {
	if t <= 0 {
		t = 0
	}
	c.MaxRetry = t
	return c
}

func (c *Config) RetryDelay(t time.Duration) *Config {
	c.RetryDelayDuration = t
	return c
}

func (c *Config) RetryOnError(err ...error) *Config {
	c.RetryOnlyOnErr = append(c.RetryOnlyOnErr, err...)
	return c
}

func (c *Config) SkipErrOn(err ...error) *Config {
	c.SkipErrIfIn = append(c.SkipErrIfIn, err...)
	return c
}

func (c *Config) PanicRecoveryOff() *Config {
	c.CrashRecoveryOff = true
	return c
}

func (c *Config) PanicRecoveryOn() *Config {
	c.CrashRecoveryOff = false
	return c
}

func (c *Config) PanicHandler(i ...PanicHandler) *Config {
	c.PanicHandlers = append(c.PanicHandlers, i...)
	return c
}

func logPanic(r interface{}) {
	callers := ""
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = callers + fmt.Sprintf("%v:%v\n", file, line)
	}
	glog.Errorf("Recovered from panic: (%v)\n%v", r, callers)
}
