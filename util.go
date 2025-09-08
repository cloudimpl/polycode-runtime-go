package _go

import "github.com/cloudimpl/byte-os/sdk/errors"

func ErrorToServiceComplete(err errors.Error, stacktraceStr string) ServiceCompleteEvent {
	var stacktrace errors.Stacktrace
	if stacktraceStr != "" {
		stacktrace = errors.Stacktrace{
			Stacktrace:   stacktraceStr,
			IsAvailable:  true,
			IsCompressed: false,
		}
		_ = stacktrace.Compress()
	}

	return ServiceCompleteEvent{
		Output:     nil,
		IsError:    true,
		Error:      err,
		Stacktrace: stacktrace,
	}
}

func ErrorToApiComplete(err errors.Error) ApiCompleteEvent {
	return ApiCompleteEvent{
		Response: ApiResponse{
			StatusCode:      500,
			Header:          make(map[string]string),
			Body:            err.ToJson(),
			IsBase64Encoded: false,
		},
	}
}
