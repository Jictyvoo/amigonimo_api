package runners

import "testing"

type Runner interface {
	Run(ctx RunnerContext) error
}

type MultiRunner struct {
	Runners []Runner
}

func (mr MultiRunner) Run(t testing.TB) error {
	ctx := NewRunnerContext(t)
	for _, runner := range mr.Runners {
		if err := runner.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
