package runtime

import (
	"context"
	"fmt"
	"github.com/cloudimpl/byte-os/sdk"
)

type AgentBuilder struct {
	ctx           context.Context
	sessionId     string
	envId         string
	agent         string
	serviceClient ServiceClient
	tenantId      string
}

func (r AgentBuilder) WithTenantId(tenantId string) sdk.AgentBuilder {
	r.tenantId = tenantId
	return r
}

func (r AgentBuilder) Get() sdk.Agent {
	return Agent{
		ctx:           r.ctx,
		sessionId:     r.sessionId,
		envId:         r.envId,
		agent:         r.agent,
		serviceClient: r.serviceClient,
		tenantId:      r.tenantId,
	}
}

type Agent struct {
	ctx           context.Context
	sessionId     string
	envId         string
	agent         string
	serviceClient ServiceClient
	tenantId      string
}

func (r Agent) Call(options sdk.TaskOptions, input sdk.AgentInput) sdk.Response {
	req := ExecServiceRequest{
		EnvId:        r.envId,
		Service:      "agent-service",
		TenantId:     r.tenantId,
		PartitionKey: r.agent + ":" + input.SessionKey,
		Method:       "CallAgent",
		Options:      options,
		Headers: map[string]string{
			AgentNameHeader: r.agent,
		},
		Input: input,
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
