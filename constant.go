package budget

import (
	"time"

	"github.com/xh3b4sd/tracer"
)

type ConstantConfig struct {
	// Budget is the amount of attempts that can be used up when consuming the
	// error budget. The configured operation is being executed until it
	// succeeds or the error budget is used up. A budget of 3 means the
	// configured operation will be executed up to 3 times.
	Budget int
	// Duration is the time to wait after any given retry. Given a Duration of 5
	// seconds and a Budget of 3 the execution would happen as follows.
	//
	//     * first execution fails
	//     * wait 5 seconds
	//     * second execution fails
	//     * wait 5 seconds
	//     * third execution fails
	//     * return error
	//
	Duration time.Duration
}

type Constant struct {
	budget   int
	duration time.Duration
}

func NewConstant(config ConstantConfig) (*Constant, error) {
	if config.Budget < 1 {
		return nil, tracer.Maskf(invalidConfigError, "%T.Budget must be greater than zero", config.Budget)
	}
	if config.Duration < 0 {
		return nil, tracer.Maskf(invalidConfigError, "%T.Duration must not be negative", config.Duration)
	}

	c := &Constant{
		budget:   config.Budget,
		duration: config.Duration,
	}

	return c, nil
}

func (c *Constant) Execute(o func() error) error {
	b := c.budget

	for {
		err := o()
		if err == nil {
			return nil
		}

		if _, ok := err.(Stop); ok {
			return nil
		}

		b--
		if b <= 0 {
			return tracer.Mask(err)
		}

		<-time.After(c.duration)
	}
}
