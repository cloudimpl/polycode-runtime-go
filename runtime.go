package runtime

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudimpl/polycode-sdk-go"
	errors2 "github.com/cloudimpl/polycode-sdk-go/errors"
	"github.com/gin-gonic/gin"
	"log"
	"runtime/debug"
)

var CurrentRuntime Runtime

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
	IsError    bool               `json:"isError"`
	Output     any                `json:"output"`
	Error      errors2.Error      `json:"error"`
	Stacktrace errors2.Stacktrace `json:"stacktrace"`
	Logs       []LogMsg           `json:"logs"`
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

type ClientService interface {
	GetName() string
	GetDescription(method string) (string, error)
	GetInputType(method string) (any, error)
	GetOutputType(method string) (any, error)
	IsWorkflow(method string) bool
	ExecuteService(ctx polycode.ServiceContext, method string, input any) (any, error)
	ExecuteWorkflow(ctx polycode.WorkflowContext, method string, input any) (any, error)
}

type Runtime interface {
	RegisterService(service Service) error
	RegisterApi(httpHandler *gin.Engine) error
	RegisterValidator(validator polycode.Validator) error
	GetValidator() polycode.Validator
	RunService(ctx context.Context, event ServiceStartEvent) (evt ServiceCompleteEvent)
	RunApi(ctx context.Context, event ApiStartEvent) (evt ApiCompleteEvent)
	Start() error
}

type ClientRuntime struct {
	env         ClientEnv
	serviceMap  map[string]ClientService
	httpHandler *gin.Engine
	client      ServiceClient
	validator   polycode.Validator
}

func (c ClientRuntime) getService(serviceName string) (ClientService, error) {
	service := c.serviceMap[serviceName]
	if service == nil {
		return nil, fmt.Errorf("client: service %s not registered", serviceName)
	}
	return service, nil
}

func (c ClientRuntime) getApi() (*gin.Engine, error) {
	if c.httpHandler == nil {
		return nil, errors.New("client: api not registered")
	}
	return c.httpHandler, nil
}

func (c ClientRuntime) RegisterService(service ClientService) error {
	log.Println("client: register service ", service.GetName())

	if c.serviceMap[service.GetName()] != nil {
		return fmt.Errorf("client: service %s already registered", service.GetName())
	} else {
		c.serviceMap[service.GetName()] = service
		return nil
	}
}

func (c ClientRuntime) RegisterApi(httpHandler *gin.Engine) error {
	if c.httpHandler != nil {
		return errors.New("client: api already registered")
	}

	c.httpHandler = httpHandler
	return nil
}

func (c ClientRuntime) RegisterValidator(validator polycode.Validator) error {
	if c.httpHandler != nil {
		return errors.New("client: validator already registered")
	}

	c.validator = validator
	return nil
}

func (c ClientRuntime) GetValidator() polycode.Validator {
	return c.validator
}

func (c ClientRuntime) Start() error {
	services, err := ExtractServiceDescription(c.serviceMap)
	if err != nil {
		return fmt.Errorf("client: failed to extract service description: %w", err)
	}

	req := StartAppRequest{
		AppName:  c.env.AppName,
		AppPort:  c.env.AppPort,
		Services: services,
		Routes:   LoadRoutes(c.httpHandler),
	}

	for {
		err = c.client.StartApp(req)
		if err == nil {
			break
		}
	}

	return nil
}

func (c ClientRuntime) RunService(ctx context.Context, event ServiceStartEvent) (evt ServiceCompleteEvent) {
	fmt.Printf("service started %s.%s", event.Service, event.Method)

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			recovered, ok := r.(error)

			if ok {
				if errors.Is(recovered, ErrTaskStopped) {
					fmt.Printf("service stopped %s.%s", event.Service, event.Method)
					evt = ValueToServiceComplete(nil)
				} else {
					stackTrace := string(debug.Stack())
					fmt.Printf("stack trace %s\n", stackTrace)
					fmt.Printf("recoverted %v", r)
					err2 := ErrInternal.Wrap(recovered)
					evt = ErrorToServiceComplete(err2, stackTrace)
				}
			} else {
				stackTrace := string(debug.Stack())
				fmt.Printf("stack trace %s\n", stackTrace)

				errorStr := fmt.Sprintf("recoverted %v", r)
				fmt.Printf("recoverted %v", r)
				err2 := ErrInternal.Wrap(errors.New(errorStr))
				evt = ErrorToServiceComplete(err2, stackTrace)
			}
		}
	}()

	service, err := c.getService(event.Service)
	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		fmt.Printf("failed to get service %s\n", err.Error())
		return ErrorToServiceComplete(err2, "")
	}

	inputObj, err := service.GetInputType(event.Method)
	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		fmt.Printf("failed to get input type %s\n", err.Error())
		return ErrorToServiceComplete(err2, "")
	}

	err = ConvertType(event.Input, inputObj)
	if err != nil {
		err2 := ErrBadRequest.Wrap(err)
		fmt.Printf("failed to convert input %s\n", err.Error())
		return ErrorToServiceComplete(err2, "")
	}

	err = c.validator.Validate(inputObj)
	if err != nil {
		err2 := ErrBadRequest.Wrap(err)
		fmt.Printf("failed to validate input %s\n", err.Error())
		return ErrorToServiceComplete(err2, "")
	}

	ctxImpl := &Context{
		ctx:       ctx,
		sessionId: event.SessionId,
		client:    c.client,
		meta:      event.Meta,
		authCtx:   event.AuthContext,
	}

	var ret any
	if service.IsWorkflow(event.Method) {
		fmt.Printf("service %s exec workflow %s with session id %s", event.Service, event.Method, event.SessionId)
		ret, err = service.ExecuteWorkflow(ctxImpl, event.Method, inputObj)
	} else {
		fmt.Printf("service %s exec handler %s with session id %s", event.Service, event.Method, event.SessionId)
		ret, err = service.ExecuteService(ctxImpl, event.Method, inputObj)
	}

	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		fmt.Printf("failed to execute service %s\n", err.Error())
		return ErrorToServiceComplete(err2, "")
	}

	fmt.Printf("service %s exec success %s\n", event.Service, event.Method)
	serviceCompleteEvent := ValueToServiceComplete(ret)
	return serviceCompleteEvent
}

func (c ClientRuntime) RunApi(ctx context.Context, event ApiStartEvent) (evt ApiCompleteEvent) {
	fmt.Printf("api started %s %s", event.Request.Method, event.Request.Path)

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			recovered, ok := r.(error)

			if ok {
				if errors.Is(recovered, ErrTaskStopped) {
					fmt.Printf("api stopped %s %s", event.Request.Method, event.Request.Path)
					evt = ApiCompleteEvent{
						Response: polycode.ApiResponse{
							StatusCode:      202,
							Header:          make(map[string]string),
							Body:            "",
							IsBase64Encoded: false,
						},
					}
				} else {
					stackTrace := string(debug.Stack())
					fmt.Printf("stack trace %s\n", stackTrace)
					fmt.Printf("recovered %v", r)
					err2 := ErrInternal.Wrap(recovered)
					evt = ApiCompleteEvent{
						Response: polycode.ApiResponse{
							StatusCode:      500,
							Header:          make(map[string]string),
							Body:            err2.ToJson(),
							IsBase64Encoded: false,
						},
					}
				}
			} else {
				stackTrace := string(debug.Stack())
				fmt.Printf("stack trace %s\n", stackTrace)

				errorStr := fmt.Sprintf("recoverted %v", r)
				fmt.Printf("recoverted %v", r)
				err2 := ErrInternal.Wrap(errors.New(errorStr))
				evt = ApiCompleteEvent{
					Response: polycode.ApiResponse{
						StatusCode:      500,
						Header:          make(map[string]string),
						Body:            err2.ToJson(),
						IsBase64Encoded: false,
					},
				}
			}
		}
	}()

	if c.httpHandler == nil {
		err2 := ErrApiExecError.Wrap(errors.New("http handler not set"))
		fmt.Printf("api stopped %s %s, reason: %s", event.Request.Method, event.Request.Path, err2.Error())
		return ErrorToApiComplete(err2)
	}

	ctxImpl := &Context{
		ctx:       ctx,
		sessionId: event.SessionId,
		client:    c.client,
		meta:      event.Meta,
		authCtx:   event.AuthContext,
	}

	newCtx := context.WithValue(ctx, "polycode.context", ctxImpl)
	httpReq, err := ConvertToHttpRequest(newCtx, event.Request)
	if err != nil {
		err2 := ErrApiExecError.Wrap(err)
		fmt.Printf("failed to convert http request %s\n", err.Error())
		return ErrorToApiComplete(err2)
	}

	res := ManualInvokeHandler(c.httpHandler, httpReq)
	fmt.Printf("api completed %s %s\n", event.Request.Method, event.Request.Path)
	return ApiCompleteEvent{
		Response: res,
	}
}

func RegisterService(service Service) error {
	return CurrentRuntime.RegisterService(service)
}

func RegisterApi(httpHandler *gin.Engine) error {
	return CurrentRuntime.RegisterApi(httpHandler)
}

func RegisterValidator(validator polycode.Validator) error {
	return CurrentRuntime.RegisterValidator(validator)
}

func GetValidator() polycode.Validator {
	return CurrentRuntime.GetValidator()
}

func Start() error {
	return CurrentRuntime.Start()
}
