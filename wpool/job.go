package wpool

import "context"

type jobID string
type jobType string
type jobMetadata map[string]interface{}

type executionFn func(ctx context.Context, args interface{}) (interface{}, error)

type runnable interface {
	Execute() result
}

type job struct {
	descriptor jobDescriptor
	execFn     executionFn
	ctx        context.Context
	args       interface{}
}

type jobDescriptor struct {
	id       jobID
	jType    jobType
	metadata map[string]interface{}
}

type result struct {
	value         interface{}
	err           error
	jobDescriptor jobDescriptor
}

func (j job) Execute() result {
	value, err := j.execFn(j.ctx, j.args)
	if err != nil {
		return result{
			err:           err,
			jobDescriptor: j.descriptor,
		}
	}

	return result{
		value:         value,
		jobDescriptor: j.descriptor,
	}
}
