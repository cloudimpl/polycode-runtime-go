package runtime

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"runtime/debug"
)

type ClientRuntime interface {
	RegisterService(service Service) error
	RegisterApi(engine *gin.Engine) error
	Start() error
	RunService(ctx context.Context, event ServiceStartEvent) (evt ServiceCompleteEvent, err error)
	RunApi(ctx context.Context, event ApiStartEvent) (evt ApiCompleteEvent, err error)
}

type ClientRuntimeImpl struct {
	env         ClientEnv
	serviceMap  map[string]Service
	httpHandler *gin.Engine
	client      ServiceClient
}

func (c ClientRuntimeImpl) getService(serviceName string) (Service, error) {
	service := c.serviceMap[serviceName]
	if service == nil {
		return nil, fmt.Errorf("client: service %s not registered", serviceName)
	}
	return service, nil
}

func (c ClientRuntimeImpl) getApi() (*gin.Engine, error) {
	if c.httpHandler == nil {
		return nil, errors.New("client: api not registered")
	}
	return c.httpHandler, nil
}

func (c ClientRuntimeImpl) RegisterService(service Service) error {
	log.Println("client: register service ", service.GetName())

	if c.serviceMap[service.GetName()] != nil {
		return fmt.Errorf("client: service %s already registered", service.GetName())
	} else {
		c.serviceMap[service.GetName()] = service
		return nil
	}
}

func (c ClientRuntimeImpl) RegisterApi(engine *gin.Engine) error {
	if c.httpHandler != nil {
		return errors.New("client: api already registered")
	}

	c.httpHandler = engine
	return nil
}

func (c ClientRuntimeImpl) Start() error {
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

func (c ClientRuntimeImpl) RunService(ctx context.Context, event ServiceStartEvent) (evt ServiceCompleteEvent, err error) {
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
				err2 := ErrInternal.Wrap(fmt.Errorf(errorStr))
				evt = ErrorToServiceComplete(err2, stackTrace)
			}
		}
	}()

	service, err := c.getService(event.Service)
	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2, "")
	}

	inputObj, err := service.GetInputType(event.Method)
	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2, "")
	}

	err = ConvertType(event.Input, inputObj)
	if err != nil {
		err2 := ErrBadRequest.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2, "")
	}

	err = currentValidator.Validate(inputObj)
	if err != nil {
		err2 := ErrBadRequest.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2, "")
	}

	ctxImpl := &ContextImpl{
		ctx:           ctx,
		sessionId:     event.SessionId,
		dataStore:     newDatabase(serviceClient, event.SessionId),
		fileStore:     newFileStore(serviceClient, event.SessionId),
		config:        AppConfig{},
		serviceClient: serviceClient,
		logger:        taskLogger,
		meta:          event.Meta,
		authCtx:       event.AuthContext,
	}

	var ret any
	if service.IsWorkflow(event.Method) {
		taskLogger.Info().Msg(fmt.Sprintf("service %s exec workflow %s with session id %s", event.Service,
			event.Method, event.SessionId))
		ret, err = service.ExecuteWorkflow(ctxImpl, event.Method, inputObj)
	} else {
		taskLogger.Info().Msg(fmt.Sprintf("service %s exec handler %s with session id %s", event.Service,
			event.Method, event.SessionId))
		ret, err = service.ExecuteService(ctxImpl, event.Method, inputObj)
	}

	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2, "")
	}

	taskLogger.Info().Msg("service completed")
	serviceCompleteEvent := ValueToServiceComplete(ret)
	serviceCompleteEvent.Meta = meta

	return serviceCompleteEvent
}

func (c ClientRuntimeImpl) RunApi(ctx context.Context, event ApiStartEvent) (evt ApiCompleteEvent, err error) {
	taskLogger.Info().Msg(fmt.Sprintf("api started %s %s", event.Request.Method, event.Request.Path))

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			recovered, ok := r.(error)

			if ok {
				if errors.Is(recovered, ErrTaskStopped) {
					taskLogger.Info().Msg("api stopped")
					evt = ApiCompleteEvent{
						Response: ApiResponse{
							StatusCode:      202,
							Header:          make(map[string]string),
							Body:            "",
							IsBase64Encoded: false,
						},
					}
				} else {
					stackTrace := string(debug.Stack())
					fmt.Printf("stack trace %s\n", stackTrace)

					taskLogger.Error().Msg(recovered.Error())
					err2 := ErrInternal.Wrap(recovered)
					evt = ApiCompleteEvent{
						Response: ApiResponse{
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
				taskLogger.Error().Msg(errorStr)
				err2 := ErrInternal.Wrap(fmt.Errorf(errorStr))
				evt = ApiCompleteEvent{
					Response: ApiResponse{
						StatusCode:      500,
						Header:          make(map[string]string),
						Body:            err2.ToJson(),
						IsBase64Encoded: false,
					},
				}
			}
		}
	}()

	if httpHandler == nil {
		err2 := ErrApiExecError.Wrap(errors.New("http handler not set"))
		taskLogger.Error().Msg(err2.Error())
		return ErrorToApiComplete(err2)
	}

	ctxImpl := &ContextImpl{
		ctx:           ctx,
		sessionId:     event.SessionId,
		dataStore:     newDatabase(serviceClient, event.SessionId),
		fileStore:     newFileStore(serviceClient, event.SessionId),
		config:        AppConfig{},
		serviceClient: serviceClient,
		logger:        taskLogger,
		meta:          event.Meta,
		authCtx:       event.AuthContext,
	}

	newCtx := context.WithValue(ctx, "polycode.context", ctxImpl)
	httpReq, err := ConvertToHttpRequest(newCtx, event.Request)
	if err != nil {
		err2 := ErrApiExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToApiComplete(err2)
	}

	res := ManualInvokeHandler(httpHandler, httpReq)
	taskLogger.Info().Msg("api completed")
	return ApiCompleteEvent{
		Response: res,
	}
}
