package runtime

import (
	"context"
	"fmt"
	"github.com/cloudimpl/byte-os/sdk"
	"github.com/cloudimpl/byte-os/sdk/errors"
)

type Memo struct {
	ctx           context.Context
	sessionId     string
	serviceClient ServiceClient
	getter        func() (any, error)
}

func (f Memo) Get() sdk.Response {
	req1 := ExecFuncRequest{
		Input: nil,
	}

	res1, err := f.serviceClient.ExecFunc(f.sessionId, req1)
	if err != nil {
		fmt.Printf("client: exec func error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
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
			error:   ErrTaskExecError.Wrap(err),
		}
	} else {
		response = Response{
			output:  output,
			isError: false,
			error:   errors.Error{},
		}
	}

	req2 := ExecFuncResult{
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
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	return response
}
