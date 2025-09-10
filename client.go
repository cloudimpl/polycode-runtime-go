package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudimpl/byte-os/sdk"
	"github.com/cloudimpl/byte-os/sdk/errors"
	"log"
	"net/http"
	"time"
)

const (
	Insert DbAction = "insert"
	Update DbAction = "update"
	Upsert DbAction = "upsert"
	Delete DbAction = "delete"
)

type TaskStatus int8

type DbAction string

type StartAppRequest struct {
	AppName    string               `json:"appName"`
	AppPort    uint                 `json:"appPort"`
	Services   []ServiceDescription `json:"services"`
	ApiHandler string               `json:"apiHandler"`
	Routes     []RouteData          `json:"routes"`
}

type ExecServiceRequest struct {
	EnvId         string            `json:"envId"`
	Service       string            `json:"service"`
	TenantId      string            `json:"tenantId"`
	PartitionKey  string            `json:"partitionKey"`
	Method        string            `json:"method"`
	Options       sdk.TaskOptions   `json:"options"`
	FireAndForget bool              `json:"fireAndForget"`
	Headers       map[string]string `json:"headers"`
	Input         any               `json:"input"`
}

type ExecAppRequest struct {
	EnvId         string          `json:"envId"`
	AppName       string          `json:"service"`
	Method        string          `json:"method"`
	Options       sdk.TaskOptions `json:"options"`
	FireAndForget bool            `json:"fireAndForget"`
	Input         any             `json:"input"`
}

type ExecServiceExtendedRequest struct {
	EnvId              string             `json:"envId"`
	ExecServiceRequest ExecServiceRequest `json:"execServiceRequest"`
}

type ExecAppExtendedRequest struct {
	EnvId          string         `json:"envId"`
	ExecAppRequest ExecAppRequest `json:"execAppRequest"`
}

type ExecServiceResponse struct {
	IsAsync bool         `json:"isAsync"`
	Output  any          `json:"output"`
	IsError bool         `json:"isError"`
	Error   errors.Error `json:"error"`
}

type ExecAppResponse struct {
	IsAsync bool         `json:"isAsync"`
	Output  any          `json:"output"`
	IsError bool         `json:"isError"`
	Error   errors.Error `json:"error"`
}

type ExecApiRequest struct {
	EnvId         string          `json:"envId"`
	Controller    string          `json:"controller"`
	Path          string          `json:"path"`
	Options       sdk.TaskOptions `json:"options"`
	FireAndForget bool            `json:"fireAndForget"`
	Request       sdk.ApiRequest  `json:"request"`
}

type ExecApiExtendedRequest struct {
	EnvId          string         `json:"envId"`
	ExecApiRequest ExecApiRequest `json:"execApiRequest"`
}

type ExecApiResponse struct {
	IsAsync  bool            `json:"isAsync"`
	Response sdk.ApiResponse `json:"response"`
	IsError  bool            `json:"isError"`
	Error    errors.Error    `json:"error"`
}

type ExecFuncRequest struct {
	Input any `json:"input"`
}

type ExecFuncResult struct {
	Input   any          `json:"input"`
	Output  any          `json:"output"`
	IsError bool         `json:"isError"`
	Error   errors.Error `json:"error"`
}

type ExecFuncResponse struct {
	IsAsync     bool         `json:"isAsync"`
	IsCompleted bool         `json:"isCompleted"`
	Output      any          `json:"output"`
	IsError     bool         `json:"isError"`
	Error       errors.Error `json:"error"`
}

// PutRequest represents the JSON structure for put operations
type PutRequest struct {
	Action     DbAction `json:"action"`
	IsGlobal   bool     `json:"isGlobal"`
	Collection string   `json:"collection"`
	Key        string   `json:"key"`
	Item       any      `json:"item"`
	TTL        int64    `json:"TTL"`
}

type UnsafePutRequest struct {
	TenantId     string     `json:"tenantId"`
	PartitionKey string     `json:"partitionKey"`
	PutRequest   PutRequest `json:"putRequest"`
}

// QueryRequest represents the JSON structure for query operations
type QueryRequest struct {
	IsGlobal   bool          `json:"isGlobal"`
	Collection string        `json:"collection"`
	Key        string        `json:"key"`
	Filter     string        `json:"filter"`
	Args       []interface{} `json:"args"`
	Limit      int           `json:"limit"`
}

type UnsafeQueryRequest struct {
	TenantId     string       `json:"tenantId"`
	PartitionKey string       `json:"partitionKey"`
	QueryRequest QueryRequest `json:"queryRequest"`
}

// GetFileRequest represents the JSON structure for get file operations
type GetFileRequest struct {
	Key string `json:"key"`
}

type GetUploadLinkRequest struct {
	Key      string `json:"key"`
	TempFile bool   `json:"tempFile"`
}

// GetFileResponse represents the JSON structure for get file response
type GetFileResponse struct {
	Content string `json:"content"`
}

type GetLinkResponse struct {
	Link string `json:"link"`
}

// PutFileRequest represents the JSON structure for put file operations
type PutFileRequest struct {
	Key      string `json:"key"`
	TempFile bool   `json:"tempFile"`
	Content  string `json:"content"`
	FilePath string `json:"filePath"`
}

type DeleteFileRequest struct {
	Key string `json:"key"`
}

type RenameFileRequest struct {
	OldKey   string `json:"oldKey"`
	NewKey   string `json:"newKey"`
	TempFile bool   `json:"tempFile"`
}

type CreateFolderRequest struct {
	Folder string `json:"folder"`
}

// ListFilePageRequest carries the usual params plus the ContinuationToken from the previous page.
type ListFilePageRequest struct {
	Prefix            string  `json:"prefix"`            // optional sub‐folder under partitionKey
	MaxKeys           int32   `json:"maxKeys"`           // how many items to return
	ContinuationToken *string `json:"continuationToken"` // nil for first page; set to NextToken from prior response
}

// ListFileResponse is one S3 object entry
type ListFileResponse struct {
	Key          string    `json:"key"` // relative to the provided Prefix
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

// ListFilePageResponse returns one page of results plus the token for the next page.
type ListFilePageResponse struct {
	Files                 []ListFileResponse `json:"files"`
	NextContinuationToken *string            `json:"nextContinuationToken"`
	IsTruncated           bool               `json:"isTruncated"`
}

type SignalEmitRequest struct {
	TaskId     string       `json:"taskId"`
	SignalName string       `json:"signalName"`
	Output     any          `json:"output"`
	IsError    bool         `json:"isError"`
	Error      errors.Error `json:"error"`
}

type RealtimeEventEmitRequest struct {
	Channel string `json:"channel"`
	Input   any    `json:"input"`
}

type SignalWaitRequest struct {
	SignalName string `json:"signalName"`
}

type SignalWaitResponse struct {
	IsAsync bool         `json:"isAsync"`
	Output  any          `json:"output"`
	IsError bool         `json:"isError"`
	Error   errors.Error `json:"error"`
}

type GetMetaDataRequest struct {
	Group string `json:"group"`
	Type  string `json:"type"`
	Key   string `json:"key"`
}

type IncrementCounterRequest struct {
	Group string `json:"group"`
	Name  string `json:"name"`
	Count uint64 `json:"count"`
	Limit uint64 `json:"limit"`
	TTL   int64  `json:"TTL"`
}

type AcquireLockRequest struct {
	Key string `json:"key"`
	TTL int64  `json:"TTL"`
}

type ReleaseLockRequest struct {
	Key string `json:"key"`
}

type IncrementCounterResponse struct {
	Value       uint64 `json:"value"`
	Incremented bool   `json:"incremented"`
}

// NewServiceClient creates a new ServiceClient with a reusable HTTP client
func NewServiceClient(baseURL string) ServiceClient {
	return &ServiceClientImpl{
		httpClient: &http.Client{
			Timeout: time.Second * 30, // Set a reasonable timeout for HTTP requests
		},
		baseURL: baseURL,
	}
}

type ServiceClient interface {
	StartApp(req StartAppRequest) error

	ExecService(sessionId string, req ExecServiceRequest) (ExecServiceResponse, error)
	ExecApp(sessionId string, req ExecAppRequest) (ExecAppResponse, error)
	ExecApi(sessionId string, req ExecApiRequest) (ExecApiResponse, error)
	ExecFunc(sessionId string, req ExecFuncRequest) (ExecFuncResponse, error)
	ExecFuncResult(sessionId string, req ExecFuncResult) error

	GetItem(sessionId string, req QueryRequest) (map[string]interface{}, error)
	UnsafeGetItem(sessionId string, req UnsafeQueryRequest) (map[string]interface{}, error)
	QueryItems(sessionId string, req QueryRequest) ([]map[string]interface{}, error)
	UnsafeQueryItems(sessionId string, req UnsafeQueryRequest) ([]map[string]interface{}, error)

	PutItem(sessionId string, req PutRequest) error
	UnsafePutItem(sessionId string, req UnsafePutRequest) error

	GetFile(sessionId string, req GetFileRequest) (GetFileResponse, error)
	GetFileDownloadLink(sessionId string, req GetFileRequest) (GetLinkResponse, error)
	PutFile(sessionId string, req PutFileRequest) error
	GetFileUploadLink(sessionId string, req GetUploadLinkRequest) (GetLinkResponse, error)
	DeleteFile(sessionId string, req DeleteFileRequest) error
	RenameFile(sessionId string, req RenameFileRequest) error
	ListFile(sessionId string, req ListFilePageRequest) (ListFilePageResponse, error)
	CreateFolder(sessionId string, req CreateFolderRequest) error

	EmitSignal(sessionId string, req SignalEmitRequest) error
	WaitForSignal(sessionId string, req SignalWaitRequest) (SignalWaitResponse, error)
	EmitRealtimeEvent(sessionId string, req RealtimeEventEmitRequest) error

	AcquireLock(sessionId string, req AcquireLockRequest) error
	ReleaseLock(sessionId string, req ReleaseLockRequest) error

	IncrementCounter(sessionId string, req IncrementCounterRequest) (IncrementCounterResponse, error)
	GetMeta(sessionId string, req GetMetaDataRequest) (map[string]interface{}, error)

	Acknowledge(sessionId string) error
}

// ServiceClientImpl is a reusable client for calling the service API
type ServiceClientImpl struct {
	httpClient *http.Client
	baseURL    string
}

// StartApp starts the app
func (sc *ServiceClientImpl) StartApp(req StartAppRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, "", "v1/system/app/start", req)
}

// ExecService executes a service with the given request
func (sc *ServiceClientImpl) ExecService(sessionId string, req ExecServiceRequest) (ExecServiceResponse, error) {
	var res ExecServiceResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/service/exec", req, &res)
	if err != nil {
		return ExecServiceResponse{}, err
	}

	if res.IsAsync {
		panic(ErrTaskStopped)
	}

	return res, nil
}

func (sc *ServiceClientImpl) ExecApp(sessionId string, req ExecAppRequest) (ExecAppResponse, error) {
	var res ExecAppResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/app/exec", req, &res)
	if err != nil {
		return ExecAppResponse{}, err
	}

	if res.IsAsync {
		panic(ErrTaskStopped)
	}

	return res, nil
}

func (sc *ServiceClientImpl) ExecApi(sessionId string, req ExecApiRequest) (ExecApiResponse, error) {
	var res ExecApiResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/api/exec", req, &res)
	if err != nil {
		return ExecApiResponse{}, err
	}

	if res.IsAsync {
		panic(ErrTaskStopped)
	}

	return res, nil
}

func (sc *ServiceClientImpl) ExecFunc(sessionId string, req ExecFuncRequest) (ExecFuncResponse, error) {
	var res ExecFuncResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/func/exec", req, &res)
	if err != nil {
		return ExecFuncResponse{}, err
	}

	if res.IsAsync {
		panic(ErrTaskStopped)
	}

	return res, nil
}

func (sc *ServiceClientImpl) ExecFuncResult(sessionId string, req ExecFuncResult) error {
	var res ExecFuncResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/func/exec/result", req, &res)
	if err != nil {
		return err
	}

	if res.IsAsync {
		panic(ErrTaskStopped)
	}

	return nil
}

// GetItem gets an item from the database
func (sc *ServiceClientImpl) GetItem(sessionId string, req QueryRequest) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/get", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) UnsafeGetItem(sessionId string, req UnsafeQueryRequest) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/unsafe-get", req, &res)
	return res, err
}

// QueryItems queries items from the database
func (sc *ServiceClientImpl) QueryItems(sessionId string, req QueryRequest) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/query", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) UnsafeQueryItems(sessionId string, req UnsafeQueryRequest) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/unsafe-query", req, &res)
	return res, err
}

// PutItem puts an item into the database
func (sc *ServiceClientImpl) PutItem(sessionId string, req PutRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/put", req)
}

func (sc *ServiceClientImpl) UnsafePutItem(sessionId string, req UnsafePutRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/unsafe-put", req)
}

// GetFile gets a file from the file store
func (sc *ServiceClientImpl) GetFile(sessionId string, req GetFileRequest) (GetFileResponse, error) {
	var res GetFileResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/get", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) GetFileDownloadLink(sessionId string, req GetFileRequest) (GetLinkResponse, error) {
	var res GetLinkResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/get-download-link", req, &res)
	return res, err
}

// PutFile puts a file into the file store
func (sc *ServiceClientImpl) PutFile(sessionId string, req PutFileRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/put", req)
}

func (sc *ServiceClientImpl) GetFileUploadLink(sessionId string, req GetUploadLinkRequest) (GetLinkResponse, error) {
	var res GetLinkResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/get-upload-link", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) DeleteFile(sessionId string, req DeleteFileRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/delete", req)
}

func (sc *ServiceClientImpl) RenameFile(sessionId string, req RenameFileRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/rename", req)
}

func (sc *ServiceClientImpl) ListFile(sessionId string, req ListFilePageRequest) (ListFilePageResponse, error) {
	var res ListFilePageResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/list", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) CreateFolder(sessionId string, req CreateFolderRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/create-folder", req)
}

func (sc *ServiceClientImpl) EmitSignal(sessionId string, req SignalEmitRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/signal/emit", req)
}

func (sc *ServiceClientImpl) WaitForSignal(sessionId string, req SignalWaitRequest) (SignalWaitResponse, error) {
	res := SignalWaitResponse{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/signal/await", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) EmitRealtimeEvent(sessionId string, req RealtimeEventEmitRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/realtime/event/emit", req)
}

func (sc *ServiceClientImpl) AcquireLock(sessionId string, req AcquireLockRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/lock/acquire", req)
}

func (sc *ServiceClientImpl) ReleaseLock(sessionId string, req ReleaseLockRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/lock/release", req)
}

func (sc *ServiceClientImpl) IncrementCounter(sessionId string, req IncrementCounterRequest) (IncrementCounterResponse, error) {
	var res IncrementCounterResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/elevated/context/counter/increment", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) GetMeta(sessionId string, req GetMetaDataRequest) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/elevated/context/meta/get", req, &res)
	return res, err
}

func (sc *ServiceClientImpl) Acknowledge(sessionId string) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/acknowledge", nil)
}

func executeApiWithoutResponse(httpClient *http.Client, baseUrl string, sessionId string, path string, req any) error {
	log.Printf("client: exec api without response from %s with session id %s", path, sessionId)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", baseUrl, path), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-polycode-task-session-id", sessionId)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error, status: %v", resp.Status)
	}

	return nil
}

func executeApiWithResponse[T any](httpClient *http.Client, baseUrl string, sessionId string, path string, req any, res *T) error {
	log.Printf("client: exec api with response from %s with session id %s\n", path, sessionId)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", baseUrl, path), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-polycode-task-session-id", sessionId)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if res == nil {
		return ErrSidecarClientFailed.With("response is null")
	}

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(res)
		if err != nil {
			return err
		}
		return nil
	} else {
		errorEvent := ErrorEvent{}
		err = json.NewDecoder(resp.Body).Decode(&errorEvent)
		if err != nil {
			return err
		}
		return errorEvent.Error
	}
}
