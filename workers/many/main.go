package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/plastikfan/asynco/workers/wpool"
)

type Arguments int
type OutputType int
type IntResult wpool.Result[OutputType]

const WorkerCount = 3
const JobCount = 4

type MyContext struct {
	CompleteBy time.Time
	Complete   <-chan struct{}
	Error      error
}

func (c MyContext) Deadline() (time.Time, bool) {
	return c.CompleteBy, c.CompleteBy.IsZero()
}

func (c MyContext) Done() <-chan struct{} {
	return c.Complete
}

func (c MyContext) Err() error {
	return c.Error
}

func (c MyContext) Value(key any) any {
	return 42
}

func main() {
	doubler := func(ctx context.Context, args Arguments) (OutputType, error) {
		output := OutputType(args * 2)

		return output, nil
	}

	deadline := time.Now().Add(time.Second * 2)

	ctx := MyContext{
		CompleteBy: deadline,
	}

	pool := wpool.New[Arguments, OutputType](WorkerCount)

	workload := []wpool.Job[Arguments, OutputType]{}

	// This deadlocks when the job count > worker count:
	//
	// 💠 ===[ running pool ->  jobs:'4', workers:'3' ]=== 💠
	// fatal error: all goroutines are asleep - deadlock!
	//
	for i := 0; i < JobCount; i++ {
		descriptor := wpool.JobDescriptor{
			ID:    wpool.JobID(strconv.Itoa(i)),
			JType: wpool.JobType("anyType"),
			Metadata: wpool.JobMetadata{
				"foo": "foo",
				"bar": "bar",
			},
		}
		j := wpool.Job[Arguments, OutputType]{
			Descriptor: descriptor,
			ExecFn:     doubler,
			Args:       10,
		}
		workload = append(workload, j)
	}
	fmt.Printf("💠 ===[ running pool ->  jobs:'%v', workers:'%v' ]=== 💠\n", len(workload), WorkerCount)

	pool.GenerateFrom(workload)

	pool.Run(ctx)
}
