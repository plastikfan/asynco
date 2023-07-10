package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/plastikfan/asynco/workers/wpool"
)

type (
	Inputs     int
	OutputType int
	IntResult  wpool.Result[OutputType]
)

const (
	DefWorkerCount = 3
	DefJobCount    = 8
	MaxWorkerCount = 3
	MaxJobCount    = 100
)

var doubler = func(ctx context.Context, args Inputs) (OutputType, error) {
	output := OutputType(args * 2)

	return output, nil
}

func sum(load ...int) int {
	s := 0
	for _, b := range load {
		s += b
	}
	return s
}

func makeLoad(size int, workload []wpool.Job[Inputs, OutputType]) []wpool.Job[Inputs, OutputType] {
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
		j := wpool.Job[Inputs, OutputType]{
			Descriptor: descriptor,
			ExecFn:     doubler,
			Args:       Inputs(10 * i),
		}
		workload = append(workload, j)
	}

	return workload
}

func processArgs(args []string) (int, []int) {
	DefBatch := []int{DefJobCount}
	wc := DefWorkerCount
	batches := []int{}
	length := len(args)

	if length <= 1 {
		return wc, DefBatch
	}

	if parse, err := strconv.Atoi(args[1]); err == nil {
		if parse >= 1 && parse <= runtime.NumCPU() {
			wc = parse
		}
	}

	if len(args) >= 3 {
		args = args[2:]
		total := 0
		for _, a := range args {
			if parse, err := strconv.Atoi(a); err == nil {
				total += parse
				if total <= MaxJobCount {
					batches = append(batches, parse)
				}
			}
		}
	} else {
		batches = DefBatch
	}
	return wc, batches
}

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*1)
	defer cancel()

	args := os.Args
	workerCount, batches := processArgs(args)
	pool := wpool.New[Inputs, OutputType](workerCount)

	fmt.Printf("💧 ARGUMENTS(%v): '%v', workerCount: '%v' \n", len(args), args, workerCount)

	go func() {
		for _, b := range batches {
			fmt.Printf("💠 ===[ running pool ->  jobs-batch:'%v', workers:'%v' ]=== 💠\n", b, workerCount)

			workload := makeLoad(b, []wpool.Job[Inputs, OutputType]{})
			pool.Dispatch(workload)
		}
	}()

	fmt.Printf("🍤 TOTAL WORKLOAD: '%v'\n", sum(batches...))

	go pool.Run(ctx)

	resultCount := 0

	for running := true; running; {
		select {
		case r, ok := <-pool.Results():
			if !ok {
				fmt.Println("---> 🍒 CONTINUE")
				continue
			}

			resultCount++
			fmt.Printf("---> 🧙‍♂️ RESULT(descriptor: %v): '%v'\n", r.Descriptor.ID, r.Value)
		case <-pool.Done:
			running = false

			// don't need to use a default clause here; the loop will spin:
			// default:
		}
	}

	pool.Finish()
	fmt.Printf("🎯 We're done now!!! (result count: '%v')\n", resultCount)
}
