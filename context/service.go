package context

import (
	"context"
	"fmt"
	"github.com/cloudimpl/byte-os/runtime"
	"github.com/cloudimpl/byte-os/sdk"
)

type ServiceBuilder struct {
	ctx           context.Context
	sessionId     string
	envId         string
	service       string
	serviceClient runtime.ServiceClient
	tenantId      string
	partitionKey  string
}

func (r *ServiceBuilder) WithTenantId(tenantId string) *ServiceBuilder {
	r.tenantId = tenantId
	return r
}

func (r *ServiceBuilder) WithPartitionKey(partitionKey string) *ServiceBuilder {
	r.partitionKey = partitionKey
	return r
}

func (r *ServiceBuilder) Get() Service {
	return Service{
		ctx:           r.ctx,
		sessionId:     r.sessionId,
		envId:         r.envId,
		service:       r.service,
		serviceClient: r.serviceClient,
		tenantId:      r.tenantId,
		partitionKey:  r.partitionKey,
	}
}

type Service struct {
	ctx           context.Context
	sessionId     string
	envId         string
	service       string
	serviceClient runtime.ServiceClient
	tenantId      string
	partitionKey  string
}

func (r Service) RequestReply(options sdk.TaskOptions, method string, input any) Response {
	req := runtime.ExecServiceRequest{
		EnvId:        r.envId,
		Service:      r.service,
		TenantId:     r.tenantId,
		PartitionKey: r.partitionKey,
		Method:       method,
		Options:      options,
		Input:        input,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   runtime.ErrTaskExecError.Wrap(err),
		}
	}

	fmt.Printf("client: exec task output: %v\n", output)
	return Response{
		output:  output.Output,
		isError: output.IsError,
		error:   output.Error,
	}
}

func (r Service) Send(options sdk.TaskOptions, method string, input any) error {
	req := runtime.ExecServiceRequest{
		EnvId:         r.envId,
		Service:       r.service,
		TenantId:      r.tenantId,
		PartitionKey:  r.partitionKey,
		Method:        method,
		Options:       options,
		FireAndForget: true,
		Input:         input,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return runtime.ErrTaskExecError.Wrap(err)
	}

	fmt.Printf("client: exec task output: %v\n", output)
	if output.IsError {
		return output.Error
	} else {
		return nil
	}
}
