package wpool

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	errDefault = errors.New("wrong argument type")
	descriptor = jobDescriptor{
		id:    jobID("1"),
		jType: jobType("anyType"),
		metadata: jobMetadata{
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
	type fields struct {
		descriptor jobDescriptor
		execFn     executionFn
		ctx        context.Context
		args       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   result
	}{
		{
			name: "job execution success",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				ctx:        context.TODO(),
				args:       10,
			},
			want: result{
				value:         20,
				err:           nil,
				jobDescriptor: descriptor,
			},
		},
		{
			name: "job execution failure",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				ctx:        context.TODO(),
				args:       "10",
			},
			want: result{
				value:         nil,
				err:           errDefault,
				jobDescriptor: descriptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := job{
				descriptor: tt.fields.descriptor,
				execFn:     tt.fields.execFn,
				ctx:        tt.fields.ctx,
				args:       tt.fields.args,
			}
			if got := j.Execute(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
