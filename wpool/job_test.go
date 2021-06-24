package wpool

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	errDefault = errors.New("wrong argument type")
	descriptor = JobDescriptor{
		ID:    JobID("1"),
		JType: jobType("anyType"),
		Metadata: jobMetadata{
			"foo": "foo",
			"bar": "bar",
		},
	}
	execFn = func(ctx context.Context, args interface{}) (interface{}, error) {
		argVal, ok := args.(int)
		if !ok {
			return nil, errDefault
		}

		return argVal * 2, nil
	}
)

func Test_job_Execute(t *testing.T) {
	ctx := context.TODO()

	type fields struct {
		descriptor JobDescriptor
		execFn     ExecutionFn
		args       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Result
	}{
		{
			name: "job execution success",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       10,
			},
			want: Result{
				Value:      20,
				Descriptor: descriptor,
			},
		},
		{
			name: "job execution failure",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       "10",
			},
			want: Result{
				Err:        errDefault,
				Descriptor: descriptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Job{
				Descriptor: tt.fields.descriptor,
				ExecFn:     tt.fields.execFn,
				Args:       tt.fields.args,
			}

			got := j.execute(ctx)
			if tt.want.Err != nil {
				if !reflect.DeepEqual(got.Err, tt.want.Err) {
					t.Errorf("execute() = %v, wantError %v", got.Err, tt.want.Err)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
