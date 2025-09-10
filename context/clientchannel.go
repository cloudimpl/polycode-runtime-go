package context

import "github.com/cloudimpl/byte-os/runtime"

type ClientChannel struct {
	name          string
	sessionId     string
	serviceClient runtime.ServiceClient
}

func (r ClientChannel) Emit(data any) error {
	req := runtime.RealtimeEventEmitRequest{
		Channel: r.name,
		Input:   data,
	}

	return r.serviceClient.EmitRealtimeEvent(r.sessionId, req)
}
