package budget

import (
	"github.com/xh3b4sd/tracer"
)

// Single executes the given operation exactly once, regardless if it fails or
// succeeds. It is usually only used for testing.
type Single struct{}

func NewSingle() *Single {
	return &Single{}
}

func (s *Single) Execute(o func() error) error {
	err := o()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
