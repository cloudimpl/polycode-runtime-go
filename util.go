package runtime

import (
	"errors"
	errors2 "github.com/cloudimpl/byte-os/sdk/errors"
	"github.com/gin-gonic/gin"
	"github.com/invopop/jsonschema"
	"log"
	"reflect"
)

func ErrorToServiceComplete(err errors2.Error, stacktraceStr string) ServiceCompleteEvent {
	var stacktrace errors2.Stacktrace
	if stacktraceStr != "" {
		stacktrace = errors2.Stacktrace{
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

func ErrorToApiComplete(err errors2.Error) ApiCompleteEvent {
	return ApiCompleteEvent{
		Response: ApiResponse{
			StatusCode:      500,
			Header:          make(map[string]string),
			Body:            err.ToJson(),
			IsBase64Encoded: false,
		},
	}
}

func ExtractServiceDescription(serviceMap map[string]Service) ([]ServiceDescription, error) {
	var services []ServiceDescription
	for srvName, srv := range serviceMap {
		serviceData := ServiceDescription{
			Name:  srvName,
			Tasks: make([]MethodDescription, 0),
		}

		res, err := srv.ExecuteService(nil, "@definition", nil)
		if err != nil {
			return nil, err
		}

		taskList := res.([]string)
		for _, taskName := range taskList {
			description, err := GetMethodDescription(srv, taskName)
			if err != nil {
				return nil, err
			}

			serviceData.Tasks = append(serviceData.Tasks, description)
		}

		services = append(services, serviceData)
	}

	return services, nil
}

func GetMethodDescription(service Service, method string) (MethodDescription, error) {
	description, err := service.GetDescription(method)
	if err != nil {
		return MethodDescription{}, err
	}

	isWorkflow := service.IsWorkflow(method)

	inputType, err := service.GetInputType(method)
	if err != nil {
		return MethodDescription{}, err
	}

	inputSchema, _, err := getSchema(inputType)
	if err != nil {
		log.Printf("Error getting method description: %s\n", err.Error())
		// skip schema extract errors
		//return MethodDescription{}, err
	}

	return MethodDescription{
		Name:        method,
		Description: description,
		IsWorkflow:  isWorkflow,
		Input:       inputSchema,
	}, nil
}

func getSchema(obj interface{}) (interface{}, any, error) {
	var schema interface{}
	for _, v := range jsonschema.Reflect(obj).Definitions {
		schema = v
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return nil, nil, errors.New("object must be a pointer")
	}

	pointsToValue := reflect.Indirect(reflect.ValueOf(obj))

	if pointsToValue.Kind() == reflect.Struct {
		return schema, obj, nil
	}

	if pointsToValue.Kind() == reflect.Slice {
		return nil, nil, errors.New("slice not supported as an input")
	}

	return schema, obj, nil
}

func LoadRoutes(httpHandler *gin.Engine) []RouteData {
	var routes = make([]RouteData, 0)
	if httpHandler != nil {
		for _, route := range httpHandler.Routes() {
			log.Printf("client: route found %s %s\n", route.Method, route.Path)

			routes = append(routes, RouteData{
				Method: route.Method,
				Path:   route.Path,
			})
		}
	}
	return routes
}
