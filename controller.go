package runtime

import (
	"context"
	"github.com/cloudimpl/polycode-sdk-go"
)

type Controller struct {
	ctx           context.Context
	sessionId     string
	envId         string
	controller    string
	serviceClient ServiceClient
}

func (r Controller) Call(options sdk.TaskOptions, path string, apiReq sdk.ApiRequest) (sdk.ApiResponse, error) {
	req := ExecApiRequest{
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
