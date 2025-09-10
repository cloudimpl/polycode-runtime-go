package context

import (
	"context"
)

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
