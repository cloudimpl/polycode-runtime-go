package runtime

import (
	"github.com/cloudimpl/polycode-sdk-go"
	"github.com/cloudimpl/polycode-sdk-go/errors"
)

const AgentNameHeader = "X-Sidecar-Agent-Name"

type MethodStartEvent struct {
	SessionId   string                      `json:"sessionId"`
	Method      string                      `json:"method"`
	Meta        polycode.HandlerContextMeta `json:"meta"`
	AuthContext polycode.AuthContext        `json:"authContext"`
	Input       any                         `json:"input"`
}

type ServiceStartEvent struct {
	SessionId   string                      `json:"sessionId"`
	Service     string                      `json:"service"`
	Method      string                      `json:"method"`
	Meta        polycode.HandlerContextMeta `json:"meta"`
	AuthContext polycode.AuthContext        `json:"authContext"`
	Input       any                         `json:"input"`
}

type ServiceCompleteEvent struct {
	IsError    bool              `json:"isError"`
	Output     any               `json:"output"`
	Error      errors.Error      `json:"error"`
	Stacktrace errors.Stacktrace `json:"stacktrace"`
	Logs       []LogMsg          `json:"logs"`
}

type ApiStartEvent struct {
	SessionId   string                      `json:"sessionId"`
	Meta        polycode.HandlerContextMeta `json:"meta"`
	AuthContext polycode.AuthContext        `json:"authContext"`
	Request     polycode.ApiRequest         `json:"request"`
}

type ApiCompleteEvent struct {
	Path     string               `json:"path"`
	Response polycode.ApiResponse `json:"response"`
	Logs     []LogMsg             `json:"logs"`
}

type ErrorEvent struct {
	Error errors.Error `json:"error"`
}

type RouteData struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type ServiceDescription struct {
	Name  string              `json:"name"`
	Tasks []MethodDescription `json:"tasks"`
}

type MethodDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IsWorkflow  bool        `json:"isWorkflow"`
	Input       interface{} `json:"input"`
}

type ClientEnv struct {
	AppName string `json:"appName"`
	AppPort uint   `json:"appPort"`
}
