package context

import (
	"github.com/cloudimpl/byte-os/runtime"
	"github.com/cloudimpl/byte-os/sdk/errors"
)

type Response struct {
	output  any
	isError bool
	error   errors.Error
}

func (r Response) IsError() bool {
	return r.isError
}

func (r Response) HasResult() bool {
	return r.output != nil
}

func (r Response) Get(ret any) error {
	if r.isError {
		return r.error
	}

	return runtime.ConvertType(r.output, ret)
}

func (r Response) GetAny() (any, error) {
	if r.isError {
		return nil, r.error
	} else {
		return r.output, nil
	}
}
