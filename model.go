package _go

import (
	"github.com/cloudimpl/byte-os/sdk"
	"github.com/cloudimpl/byte-os/sdk/errors"
)

const AgentNameHeader = "X-Sidecar-Agent-Name"

type MethodStartEvent struct {
	SessionId   string                 `json:"sessionId"`
	Method      string                 `json:"method"`
	Meta        sdk.HandlerContextMeta `json:"meta"`
	AuthContext sdk.AuthContext        `json:"authContext"`
	Input       any                    `json:"input"`
}

type ServiceStartEvent struct {
	SessionId   string                 `json:"sessionId"`
	Service     string                 `json:"service"`
	Method      string                 `json:"method"`
	Meta        sdk.HandlerContextMeta `json:"meta"`
	AuthContext sdk.AuthContext        `json:"authContext"`
	Input       any                    `json:"input"`
}

type ServiceCompleteEvent struct {
	IsError    bool              `json:"isError"`
	Output     any               `json:"output"`
	Error      errors.Error      `json:"error"`
	Stacktrace errors.Stacktrace `json:"stacktrace"`
	Logs       []LogMsg          `json:"logs"`
}

type ApiStartEvent struct {
	SessionId   string                 `json:"sessionId"`
	Meta        sdk.HandlerContextMeta `json:"meta"`
	AuthContext sdk.AuthContext        `json:"authContext"`
	Request     ApiRequest             `json:"request"`
}

type ApiCompleteEvent struct {
	Path     string      `json:"path"`
	Response ApiResponse `json:"response"`
	Logs     []LogMsg    `json:"logs"`
}

type ErrorEvent struct {
	Error errors.Error `json:"error"`
}

type ApiRequest struct {
	Id              string            `json:"id"`
	Host            string            `json:"host"`
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	Query           map[string]string `json:"query"`
	Header          map[string]string `json:"header"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type ApiResponse struct {
	StatusCode      int               `json:"statusCode"`
	Header          map[string]string `json:"header"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
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
