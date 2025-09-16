package runtime

import (
	"context"
	"github.com/cloudimpl/polycode-sdk-go"
	"time"
)

type Context struct {
	ctx       context.Context
	sessionId string
	client    ServiceClient
	meta      sdk.HandlerContextMeta
	authCtx   sdk.AuthContext
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c Context) Err() error {
	return c.ctx.Err()
}

func (c Context) Value(key any) any {
	return c.ctx.Value(key)
}

func (c Context) Meta() sdk.HandlerContextMeta {
	return c.meta
}

func (c Context) AuthContext() sdk.AuthContext {
	return c.authCtx
}

func (c Context) Logger() sdk.Logger {
	return JsonLogger{
		section: "task",
	}
}

func (c Context) UnsafeDb() sdk.DataStoreBuilder {
	return UnsafeDataStoreBuilder{
		sessionId: c.sessionId,
		client:    c.client,
	}
}

func (c Context) FileStore() sdk.Folder {
	return Folder{
		sessionId: c.sessionId,
		client:    c.client,
		parent:    nil,
		name:      "files",
	}
}

func (c Context) TempFileStore() sdk.Folder {
	return Folder{
		sessionId: c.sessionId,
		client:    c.client,
		parent:    nil,
		name:      "temp-files",
	}
}

func (c Context) Db() sdk.DataStore {
	return DataStore{
		sessionId: c.sessionId,
		client:    c.client,
	}
}

func (c Context) Service(service string) sdk.ServiceBuilder {
	return ServiceBuilder{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		service:       service,
		serviceClient: c.client,
	}
}

func (c Context) Agent(agent string) sdk.AgentBuilder {
	return AgentBuilder{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		agent:         agent,
		serviceClient: c.client,
	}
}

func (c Context) ServiceEx(envId string, service string) sdk.ServiceBuilder {
	return ServiceBuilder{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		service:       service,
		serviceClient: c.client,
		envId:         envId,
	}
}

func (c Context) AgentEx(envId string, agent string) sdk.AgentBuilder {
	return AgentBuilder{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		agent:         agent,
		serviceClient: c.client,
		envId:         envId,
	}
}

func (c Context) App(appName string) sdk.Service {
	// new method added with v2
	panic("implement me")
}

func (c Context) AppEx(envId string, appName string) sdk.Service {
	// new method added with v2
	panic("implement me")
}

func (c Context) Controller(controller string) sdk.Controller {
	return Controller{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		controller:    controller,
		serviceClient: c.client,
	}
}

func (c Context) ControllerEx(envId string, controller string) sdk.Controller {
	return Controller{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		controller:    controller,
		serviceClient: c.client,
		envId:         envId,
	}
}

func (c Context) Memo(getter func() (any, error)) sdk.Response {
	return Memo{
		ctx:           c.ctx,
		sessionId:     c.sessionId,
		serviceClient: c.client,
		getter:        getter,
	}.Get()
}

func (c Context) Signal(signalName string) sdk.Signal {
	// not implemented, old architecture is not good
	panic("implement me")
}

func (c Context) ClientChannel(channelName string) sdk.ClientChannel {
	return ClientChannel{
		name:          channelName,
		sessionId:     c.sessionId,
		serviceClient: c.client,
	}
}

func (c Context) Lock(key string) sdk.Lock {
	return Lock{
		client:    c.client,
		sessionId: c.sessionId,
		key:       key,
	}
}
