package budget

// Interface implements any error budget mechanism in order to retry arbitrary
// operations. Note that certain implementations may not be thread safe.
type Interface interface {
	// Operation is the function being executed until it does not return an
	// error or the configured budget got used up.
	Execute(o func() error) error
}
