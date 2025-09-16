package runtime

import "github.com/cloudimpl/polycode-sdk-go/errors"

var ErrInternal = errors.DefineError("polycode.client.runtime", 1, "internal error")
var ErrPanic = errors.DefineError("polycode.client.runtime", 2, "task in progress")
var ErrSidecarClientFailed = errors.DefineError("polycode.client.runtime", 3, "sidecar client failed, reason: [%s]")
var ErrServiceExecError = errors.DefineError("polycode.client", 4, "service exec error")
var ErrApiExecError = errors.DefineError("polycode.client", 5, "api exec error")
var ErrBadRequest = errors.DefineError("polycode.client", 6, "bad request")
var ErrTaskExecError = errors.DefineError("polycode.client", 7, "task execution error")
var ErrTaskStopped = &ErrPanic
