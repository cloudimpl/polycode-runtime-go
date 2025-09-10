package context

import (
	"context"
	"fmt"
	"github.com/cloudimpl/byte-os/runtime"
	"github.com/cloudimpl/byte-os/sdk/errors"
)

type Memo struct {
	ctx           context.Context
	sessionId     string
	serviceClient runtime.ServiceClient
	getter        func() (any, error)
}

func (f Memo) Get() Response {
	req1 := runtime.ExecFuncRequest{
		Input: nil,
	}

	res1, err := f.serviceClient.ExecFunc(f.sessionId, req1)
	if err != nil {
		fmt.Printf("client: exec func error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   runtime.ErrTaskExecError.Wrap(err),
		}
	}

	if res1.IsCompleted {
		return Response{
			output:  res1.Output,
			isError: res1.IsError,
			error:   res1.Error,
		}
	}

	output, err := f.getter()
	var response Response
	if err != nil {
		response = Response{
			output:  nil,
			isError: true,
			error:   runtime.ErrTaskExecError.Wrap(err),
		}
	} else {
		response = Response{
			output:  output,
			isError: false,
			error:   errors.Error{},
		}
	}

	req2 := runtime.ExecFuncResult{
		Input:   nil,
		Output:  response.output,
		IsError: response.isError,
		Error:   response.error,
	}

	err = f.serviceClient.ExecFuncResult(f.sessionId, req2)
	if err != nil {
		fmt.Printf("client: exec func result error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   runtime.ErrTaskExecError.Wrap(err),
		}
	}

	return response
}
