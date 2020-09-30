package budget

// Interface implements any error budget mechanism in order to retry arbitrary
// operations. Note that certain implementations may not be thread safe. Note
// that all implementations have to guarantee consecutive executions without
// causing conflicts in using up the configured error budget. Any instance of
// any budget implementation must therefore be reusable. That means that Execute
// must be callable multiple times while respecting the configured budget during
// execution.
type Interface interface {
	Execute(o func() error) error
}
