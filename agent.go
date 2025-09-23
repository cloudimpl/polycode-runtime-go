package runtime

import (
	"context"
	"fmt"
	"github.com/cloudimpl/polycode-sdk-go"
)

type AgentBuilder struct {
	ctx           context.Context
	sessionId     string
	envId         string
	agent         string
	serviceClient ServiceClient
}

func (r AgentBuilder) WithEnvId(envId string) polycode.AgentBuilder {
	r.envId = envId
	return r
}

func (r AgentBuilder) Get() polycode.Agent {
	return Agent{
		ctx:           r.ctx,
		sessionId:     r.sessionId,
		envId:         r.envId,
		agent:         r.agent,
		serviceClient: r.serviceClient,
	}
}

type Agent struct {
	ctx           context.Context
	sessionId     string
	envId         string
	agent         string
	serviceClient ServiceClient
}

func (r Agent) Call(options polycode.TaskOptions, input polycode.AgentInput) polycode.Response {
	req := ExecServiceRequest{
		EnvId:   r.envId,
		Service: "agent-service",
		Method:  "CallAgent",
		Options: options.WithSequenceKey(r.agent + ":" + input.SessionKey),
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
