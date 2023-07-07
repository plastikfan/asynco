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
const JobCount = 8

var doubler = func(ctx context.Context, args Arguments) (OutputType, error) {
	output := OutputType(args * 2)

	return output, nil
}

func makeLoad(size int, workload []wpool.Job[Arguments, OutputType]) []wpool.Job[Arguments, OutputType] {
	offset := 0

	for i := 0; i < size; i++ {
		num := i + offset
		descriptor := wpool.JobDescriptor{
			ID:    wpool.JobID(strconv.Itoa(num)),
			JType: wpool.JobType("anyType"),
			Metadata: wpool.JobMetadata{
				"foo": "foo",
				"bar": "bar",
			},
		}
		j := wpool.Job[Arguments, OutputType]{
			Descriptor: descriptor,
			ExecFn:     doubler,
			Args:       Arguments(10 * i),
		}
		workload = append(workload, j)
	}

	return workload
}

func main() {
	pool := wpool.New[Arguments, OutputType](WorkerCount)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*1)
	defer cancel()

	workload := makeLoad(JobCount, []wpool.Job[Arguments, OutputType]{})

	fmt.Printf("💠 ===[ running pool ->  jobs:'%v', workers:'%v' ]=== 💠\n", len(workload), WorkerCount)

	go pool.GenerateFrom(workload)
	go pool.Run(ctx)

	resultCount := 0

	for {
		select {
		case r, ok := <-pool.Results():
			if !ok {
				fmt.Println("---> 🍒 CONTINUE")
				continue
			}

			i, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				fmt.Printf("💢 unexpected error: %v", err)
				return
			}

			resultCount++
			fmt.Printf("---> 🧙‍♂️ RESULT(descriptor: %v): '%v'\n", i, r.Value)
		case <-pool.Done:
			fmt.Printf("🎯 We're done now!!! (result count: '%v')\n", resultCount)
			return

			// don't need to use a default clause here; the loop will spin:
			// default:
		}
	}
}
