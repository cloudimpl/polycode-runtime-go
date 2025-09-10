package context

import (
	"context"
	"github.com/cloudimpl/byte-os/runtime"
	"github.com/cloudimpl/byte-os/sdk"
)

type Controller struct {
	ctx           context.Context
	sessionId     string
	envId         string
	controller    string
	serviceClient runtime.ServiceClient
}

func (r Controller) RequestReply(options sdk.TaskOptions, path string, apiReq sdk.ApiRequest) (sdk.ApiResponse, error) {
	req := runtime.ExecApiRequest{
		EnvId:      r.envId,
		Controller: r.controller,
		Path:       path,
		Options:    options,
		Request:    apiReq,
	}

	output, err := r.serviceClient.ExecApi(r.sessionId, req)
	if err != nil {
		return sdk.ApiResponse{}, err
	}

	if output.IsError {
		return sdk.ApiResponse{}, output.Error
	}

	return output.Response, nil
}

func (r Controller) Send(options sdk.TaskOptions, path string, apiReq sdk.ApiRequest) error {
	req := runtime.ExecApiRequest{
		EnvId:         r.envId,
		Controller:    r.controller,
		Path:          path,
		Options:       options,
		FireAndForget: true,
		Request:       apiReq,
	}

	output, err := r.serviceClient.ExecApi(r.sessionId, req)
	if err != nil {
		return err
	}

	if output.IsError {
		return output.Error
	}

	return nil
}
