package failsafe

import (
	"errors"
	"sync"
	"time"
)

type failsafe struct {
	config   *Config
	contexts *FailsafeContext

	elm failsafeElement
}

func New() *failsafe {
	return &failsafe{
		config:   NewConfig(),
		contexts: newFailsafeContext(),
	}
}

func NewWithConfig(c *Config) *failsafe {
	return &failsafe{
		config:   c,
		contexts: newFailsafeContext(),
	}
}

func (f *failsafe) AddHandler(r RunHandler, e OnErrHandler) *failsafe {
	return f.AddTask(&failsafeHandler{
		runHandler:   r,
		onErrHandler: e,
	})
}

func (f *failsafe) AddTask(t failsafeElement) *failsafe {
	f.elm = t
	return f
}

func (f *failsafe) WithParam(key string, value interface{}) *failsafe {
	f.contexts.Params.Put(key, value)
	return f
}

func (f *failsafe) WithParams(params ...map[string]interface{}) *failsafe {
	for _, ps := range params {
		for key, val := range ps {
			f.contexts.Params.Put(key, val)
		}
	}
	return f
}

func (f *failsafe) Run() getterInterface {
	var wg sync.WaitGroup

	wg.Add(1)
	var runErr error
	var runPanic bool
	var recovery interface{}
	go func() {
		defer wg.Done()
		if !f.config.CrashRecoveryOff {
			defer func() {
				if r := recover(); r != nil {
					runErr = errors.New("run() paniced")
					runPanic = true
					recovery = r
				}
			}()
		}

		for i := 1; i <= f.config.MaxRetry; i++ {
			err := f.elm.Run(f.contexts)
			if err == nil {
				break
			}
			if err != nil {
				runErr = err
			}

			if f.config.RetryDelayDuration >= time.Nanosecond {
				time.Sleep(f.config.RetryDelayDuration)
			}
		}
	}()
	wg.Wait()
	if runPanic {
		for _, handler := range f.config.PanicHandlers {
			handler(recovery)
		}
	}

	var errErr error
	var errPanic bool
	if runErr != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !f.config.CrashRecoveryOff {
				defer func() {
					if r := recover(); r != nil {
						errErr = errors.New("onErr() paniced")
						errPanic = true
						recovery = r
					}
				}()
			}

			for i := 1; i <= f.config.MaxRetry; i++ {
				err := f.elm.OnErr(f.contexts)
				if err == nil {
					break
				}
				if err != nil {
					errErr = err
				}

				if f.config.RetryDelayDuration >= time.Nanosecond {
					time.Sleep(f.config.RetryDelayDuration)
				}
			}

		}()
	}
	if errPanic {
		for _, handler := range f.config.PanicHandlers {
			handler(recovery)
		}
	}
	wg.Wait()
	return f.contexts.Return.(getterInterface)
}

type failsafeHandler struct {
	runHandler   RunHandler
	onErrHandler OnErrHandler
}

func (f *failsafeHandler) Run(c *FailsafeContext) error {
	return f.runHandler(c)
}

func (f *failsafeHandler) OnErr(c *FailsafeContext) error {
	return f.onErrHandler(c)
}
