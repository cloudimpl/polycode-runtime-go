package runtime

import "github.com/cloudimpl/byte-os/sdk/errors"

var ErrInternal = errors.DefineError("polycode.client.runtime", 1, "internal error")
var ErrPanic = errors.DefineError("polycode.client.runtime", 2, "task in progress")
var ErrSidecarClientFailed = errors.DefineError("polycode.client.runtime", 3, "sidecar client failed, reason: [%s]")
var ErrTaskStopped = &ErrPanic
