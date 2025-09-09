package runtime

import "github.com/cloudimpl/byte-os/sdk"

type Service interface {
	GetName() string
	GetDescription(method string) (string, error)
	GetInputType(method string) (any, error)
	GetOutputType(method string) (any, error)
	IsWorkflow(method string) bool
	ExecuteService(ctx sdk.ServiceContext, method string, input any) (any, error)
	ExecuteWorkflow(ctx sdk.WorkflowContext, method string, input any) (any, error)
}
