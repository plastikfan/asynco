package wpool

import (
	"context"
)

type JobID string
type JobType string
type JobMetadata map[string]interface{}

type ExecutionFn func(ctx context.Context, args interface{}) (interface{}, error)

type JobDescriptor struct {
	ID       JobID
	JType    JobType
	Metadata map[string]interface{}
}

type Result struct {
	Value      interface{}
	Err        error
	Descriptor JobDescriptor
}

type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       interface{}
}

func (j Job) Execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}

type Fields struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       interface{}
}
