package wpool

import (
	"context"
	"fmt"
)

type JobID string
type JobType string
type JobMetadata map[string]interface{}

type ExecutionFn[A, T any] func(ctx context.Context, args A) (T, error)

type JobDescriptor struct {
	ID       JobID
	JType    JobType
	Metadata map[string]interface{}
}

type Result[T any] struct {
	Value      T
	Err        error
	Descriptor JobDescriptor
}

type Job[A any, T any] struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn[A, T]
	Args       A
}

func (j Job[A, T]) Execute(ctx context.Context) Result[T] {
	fmt.Printf(">>> 🚀 Starting Job: '%v'\n", j.Descriptor.ID)

	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		fmt.Printf("!!! 🔥 Failed Job: '%v'\n", j.Descriptor.ID)
		return Result[T]{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}
	fmt.Printf("<<< ✨ Completed Job: '%v'\n", j.Descriptor.ID)
	return Result[T]{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}

type Fields[A, T any] struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn[A, T]
	Args       A
}
