package context

import (
	"context"
	"github.com/cloudimpl/byte-os/runtime"
	"github.com/cloudimpl/byte-os/sdk"
	"time"
)

type Context struct {
	ctx        context.Context
	sessionId  string
	client     runtime.ServiceClient
	meta       sdk.HandlerContextMeta
	callerMeta sdk.CallerContextMeta
	authCtx    sdk.AuthContext
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
		ctx: c,
	}
}

func (c Context) UnsafeDb() sdk.DataStoreBuilder {
	//TODO implement me
	panic("implement me")
}

func (c Context) FileStore() sdk.Folder {
	//TODO implement me
	panic("implement me")
}

func (c Context) TempFileStore() sdk.Folder {
	//TODO implement me
	panic("implement me")
}

type BaseContext interface {
	context.Context
	Meta() HandlerContextMeta
	AuthContext() AuthContext
	Logger() Logger
	UnsafeDb() DataStoreBuilder
	FileStore() Folder
	TempFileStore() Folder
}

type ServiceContext interface {
	BaseContext
	Db() DataStore
}

type WorkflowContext interface {
	BaseContext
	Service(service string) *ServiceBuilder
	Agent(agent string) *AgentBuilder
	ServiceEx(envId string, service string) *ServiceBuilder
	AgentEx(envId string, agent string) *AgentBuilder
	App(appName string) Service
	AppEx(envId string, appName string) Service
	Controller(controller string) Controller
	ControllerEx(envId string, controller string) Controller
	Memo(getter func() (any, error)) Response
	Signal(signalName string) Signal
	ClientChannel(channelName string) ClientChannel
	Lock(key string) Lock
}

type ApiContext interface {
	WorkflowContext
}

type RawContext interface {
	BaseContext
	GetMeta(group string, typeName string, key string) (map[string]interface{}, error)
}
