package runtime

import (
	"context"
	"fmt"
	"github.com/cloudimpl/polycode-sdk-go"
)

type ServiceBuilder struct {
	ctx           context.Context
	sessionId     string
	envId         string
	service       string
	serviceClient ServiceClient
}

func (r ServiceBuilder) WithEnvId(envId string) polycode.ServiceBuilder {
	r.envId = envId
	return r
}

func (r ServiceBuilder) Get() polycode.Service {
	return Service{
		ctx:           r.ctx,
		sessionId:     r.sessionId,
		envId:         r.envId,
		service:       r.service,
		serviceClient: r.serviceClient,
	}
}

type Service struct {
	ctx           context.Context
	sessionId     string
	envId         string
	service       string
	serviceClient ServiceClient
}

func (r Service) RequestReply(options polycode.TaskOptions, method string, input any) polycode.Response {
	req := ExecServiceRequest{
		EnvId:   r.envId,
		Service: r.service,
		Method:  method,
		Options: options,
		Input:   input,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
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

func (r Service) Send(options polycode.TaskOptions, method string, input any) error {
	req := ExecServiceRequest{
		EnvId:         r.envId,
		Service:       r.service,
		Method:        method,
		Options:       options,
		FireAndForget: true,
		Input:         input,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
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
