package context

import "github.com/cloudimpl/byte-os/sdk/errors"

type Signal interface {
	Await() Response
	EmitValue(taskId string, data any) error
	EmitError(taskId string, err errors.Error) error
}
