package runtime

import (
	"context"
	"fmt"
	"github.com/cloudimpl/byte-os/sdk"
)

type App struct {
	ctx           context.Context
	sessionId     string
	envId         string
	appName       string
	serviceClient ServiceClient
}

func (r App) RequestReply(options sdk.TaskOptions, method string, input any) Response {
	req := ExecAppRequest{
		EnvId:   r.envId,
		AppName: r.appName,
		Method:  method,
		Options: options,
		Input:   input,
	}

	output, err := r.serviceClient.ExecApp(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	fmt.Printf("client: exec task output: %v\n", output)
	return Response{
		output:  output.Output,
		isError: output.IsError,
		error:   output.Error,
	}
}

func (r App) Send(options sdk.TaskOptions, method string, input any) error {
	req := ExecAppRequest{
		EnvId:         r.envId,
		AppName:       r.appName,
		Method:        method,
		Options:       options,
		FireAndForget: true,
		Input:         input,
	}

	output, err := r.serviceClient.ExecApp(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return ErrTaskExecError.Wrap(err)
	}

	fmt.Printf("client: exec task output: %v\n", output)
	if output.IsError {
		return output.Error
	} else {
		return nil
	}
}
