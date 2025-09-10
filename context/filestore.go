package context

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cloudimpl/byte-os/runtime"
)

type FileStore struct {
	client    runtime.ServiceClient
	sessionId string
}

func (d FileStore) NewFolder(name string) (Folder, error) {
	req := runtime.CreateFolderRequest{
		Folder: name,
	}

	err := d.client.CreateFolder(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to create folder: %s\n", err.Error())
		return Folder{}, err
	}

	return d.Folder(name), nil
}

func (d FileStore) Get(path string) (bool, []byte, error) {
	req := runtime.GetFileRequest{
		Key: path,
	}

	res, err := d.client.GetFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file: %s\n", err.Error())
		return false, nil, err
	}

	if res.Content == "" {
		return false, nil, nil
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		fmt.Printf("failed to decode base64: %s\n", err.Error())
		return true, nil, err
	}

	return true, data, nil
}

func (d FileStore) GetDownloadLink(path string) (string, error) {
	req := runtime.GetFileRequest{
		Key: path,
	}

	res, err := d.client.GetFileDownloadLink(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (d FileStore) Save(path string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := runtime.PutFileRequest{
		Key:      path,
		TempFile: false,
		Content:  base64Data,
	}

	err := d.client.PutFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) SaveTemp(path string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := runtime.PutFileRequest{
		Key:      path,
		TempFile: true,
		Content:  base64Data,
	}

	err := d.client.PutFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) Upload(path string, filePath string) error {
	req := runtime.PutFileRequest{
		Key:      path,
		TempFile: false,
		FilePath: filePath,
	}

	err := d.client.PutFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) UploadTemp(path string, filePath string) error {
	req := runtime.PutFileRequest{
		Key:      path,
		TempFile: true,
		FilePath: filePath,
	}

	err := d.client.PutFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) GetUploadLink(path string) (string, error) {
	req := runtime.GetUploadLinkRequest{
		Key:      path,
		TempFile: false,
	}

	res, err := d.client.GetFileUploadLink(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (d FileStore) GetTempUploadLink(path string) (string, error) {
	req := runtime.GetUploadLinkRequest{
		Key:      path,
		TempFile: true,
	}

	res, err := d.client.GetFileUploadLink(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (d FileStore) Delete(path string) error {
	req := runtime.DeleteFileRequest{
		Key: path,
	}

	err := d.client.DeleteFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to delete file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) Move(oldPath string, newPath string) error {
	req := runtime.RenameFileRequest{
		OldKey: oldPath,
		NewKey: newPath,
	}

	err := d.client.RenameFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to rename file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) Folder(name string) Folder {
	return Folder{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
	}
}

type Folder struct {
	client    runtime.ServiceClient
	sessionId string
	name      string
}

func (f Folder) Load(name string) (bool, []byte, error) {
	req := runtime.GetFileRequest{
		Key: f.name + "/" + name,
	}

	res, err := f.client.GetFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file: %s\n", err.Error())
		return false, nil, err
	}

	if res.Content == "" {
		return false, nil, nil
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		fmt.Printf("failed to decode base64: %s\n", err.Error())
		return true, nil, err
	}

	return true, data, nil
}

func (f Folder) Save(name string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := runtime.PutFileRequest{
		Key:      f.name + "/" + name,
		TempFile: false,
		Content:  base64Data,
	}

	err := f.client.PutFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (f Folder) SaveTemp(name string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := runtime.PutFileRequest{
		Key:      f.name + "/" + name,
		TempFile: true,
		Content:  base64Data,
	}

	err := f.client.PutFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (f Folder) Upload(name string, filePath string) error {
	req := runtime.PutFileRequest{
		Key:      f.name + "/" + name,
		TempFile: false,
		FilePath: filePath,
	}

	err := f.client.PutFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (f Folder) UploadTemp(name string, filePath string) error {
	req := runtime.PutFileRequest{
		Key:      f.name + "/" + name,
		TempFile: true,
		FilePath: filePath,
	}

	err := f.client.PutFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func newFileStore(client runtime.ServiceClient, sessionId string) FileStore {
	return FileStore{
		client:    client,
		sessionId: sessionId,
	}
}
