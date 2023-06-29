package main

import (
	"context"
	"fmt"

	"github.com/plastikfan/asynco/workers/wpool"
)

type Arguments int
type OutputType int
type IntResult wpool.Result[OutputType]

func main() {
	ctx := context.TODO()
	descriptor := wpool.JobDescriptor{
		ID:    wpool.JobID("1"),
		JType: wpool.JobType("anyType"),
		Metadata: wpool.JobMetadata{
			"foo": "foo",
			"bar": "bar",
		},
	}

	doubler := func(ctx context.Context, args Arguments) (OutputType, error) {
		output := OutputType(args * 2)

		return output, nil
	}

	o := struct {
		name string
		want IntResult
	}{
		name: "✔️ job execution success (doubler)",
		want: IntResult{
			Value:      20,
			Descriptor: descriptor,
		},
	}

	j := wpool.Job[Arguments, OutputType]{
		Descriptor: descriptor,
		ExecFn:     doubler,
		Args:       10,
	}

	fmt.Printf("===[ running job '%v' ]===\n", o.name)

	got := j.Execute(ctx)

	if o.want.Err != nil {
		if got.Err != o.want.Err {
			fmt.Printf("execute() = %v, wantError %v\n", got.Err, o.want.Err)
			return
		}
	}

	if o.want.Value == got.Value {
		fmt.Printf("🎯 SUCCESS(result): %v\n", got.Value)
	} else {
		fmt.Printf("💥 FAILED(result): %v\n", got.Value)
	}
}
