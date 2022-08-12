package discovery

import "time"

// Option is builder option.
type Option func(o *option)

type option struct {
	timeout          time.Duration
	insecure         bool
	debugLogDisabled bool
}

// WithTimeout with timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.timeout = timeout
	}
}

// WithInsecure with isSecure option.
func WithInsecure(insecure bool) Option {
	return func(o *option) {
		o.insecure = insecure
	}
}

// DisableDebugLog disables update instances log.
func DisableDebugLog() Option {
	return func(o *option) {
		o.debugLogDisabled = true
	}
}
