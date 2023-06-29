package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/plastikfan/asynco/workers/wpool"
)

func main() {
	fmt.Println("---> SUCCESS ✔️")

	ctx := context.TODO()
	errDefault := errors.New("wrong argument type")
	descriptor := wpool.JobDescriptor{
		ID:    wpool.JobID("1"),
		JType: wpool.JobType("anyType"),
		Metadata: wpool.JobMetadata{
			"foo": "foo",
			"bar": "bar",
		},
	}

	execFn := func(ctx context.Context, args interface{}) (interface{}, error) {
		argVal, ok := args.(int)
		if !ok {
			return nil, errDefault
		}

		return argVal * 2, nil
	}

	o := struct {
		name   string
		fields wpool.Fields
		want   wpool.Result
	}{
		name: "job execution success",
		fields: wpool.Fields{
			Descriptor: descriptor,
			ExecFn:     execFn,
			Args:       10,
		},
		want: wpool.Result{
			Value:      20,
			Descriptor: descriptor,
		},
	}

	j := wpool.Job{
		ExecFn: func(ctx context.Context, args interface{}) (interface{}, error) {
			argVal, ok := args.(int)
			if !ok {
				return nil, errDefault
			}

			return argVal * 2, nil
		},
		Args: 10,
	}

	got := j.Execute(ctx)

	if o.want.Err != nil {
		if !reflect.DeepEqual(got.Err, o.want.Err) {
			fmt.Printf("execute() = %v, wantError %v\n", got.Err, o.want.Err)
		}
		return
	}

	if !reflect.DeepEqual(got, o.want) {
		fmt.Printf("execute() = %v, want %v\n", got, o.want)
	}
}
