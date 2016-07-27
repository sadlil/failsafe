package failsafe

import "time"

// Configurations of a failsafe instance or failsafe chain
type Config struct {
	// Number of maximum retries that will be done in case
	// of any failure, before running the error handlers.
	// value is set to 1 by default.
	MaxRetry int

	// Time to wait before retrying.
	RetryDelayDuration time.Duration

	// Retry only if the returned error of an Run() method
	// matches any elements of this list
	RetryOnlyOnErr []error

	// Skip Running the failsafe methods it the returned
	// err is in this list.
	SkipErrIfIn []error
	// Turn off the panic recovery, panics in case of the
	// Run() or error handler panics.
	// value is set to false by default
	CrashRecoveryOff bool

	// Panic Handlers, those will runs in case of any panic
	PanicHandlers []PanicHandler

	// Time to wait before retrying any handlers
	RetryDelay time.Duration
}

type setterInterface interface {
	Put(string, interface{}) *value
}

type getterInterface interface {
	Get(string) (*abstractValue, bool)
	Len() int
}

type valueInterface interface {
	setterInterface
	getterInterface
	Clear() *value
	Remove(string) *value
}

type RunHandler func(*FailsafeContext) error
type OnErrHandler func(*FailsafeContext) error
type OnFailHandler func(*FailsafeContext) error

type failsafeElement interface {
	Run(*FailsafeContext) error
	OnErr(*FailsafeContext) error
}

type failsafeChainElement interface {
	failsafeElement
	OnFail(*FailsafeContext) error
}

type Result struct {
	returned getterInterface
}
