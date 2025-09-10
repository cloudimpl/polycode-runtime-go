package runtime

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cloudimpl/byte-os/sdk"
)

type Folder struct {
	client    ServiceClient
	sessionId string
	parent    sdk.Folder
	name      string
}

func (f Folder) Parent() sdk.Folder {
	return f.parent
}

func (f Folder) Name() string {
	return f.name
}

func (f Folder) Path() string {
	return f.parent.Path() + "/" + f.name
}

func (f Folder) Folder(name string) sdk.Folder {
	return Folder{
		client:    f.client,
		sessionId: f.sessionId,
		parent:    f,
		name:      name,
	}
}

func (f Folder) CreateNewFolder(name string) (sdk.Folder, error) {
	req := CreateFolderRequest{
		Folder: name,
	}

	err := f.client.CreateFolder(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to create folder: %s\n", err.Error())
		return Folder{}, err
	}

	return f.Folder(name), nil
}

func (f Folder) File(name string) sdk.File {
	return File{
		client:    f.client,
		sessionId: f.sessionId,
		parent:    f,
		name:      name,
	}
}

type File struct {
	client    ServiceClient
	sessionId string
	parent    sdk.Folder
	name      string
}

func (f File) Parent() sdk.Folder {
	return f.parent
}

func (f File) Name() string {
	return f.name
}

func (f File) Path() string {
	return f.parent.Path() + "/" + f.name
}

func (f File) Get() (bool, []byte, error) {
	req := GetFileRequest{
		Key: f.Path(),
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

func (f File) Download(filePath string) error {
	// new method added with v2
	panic("implement me")
}

func (f File) GetDownloadLink() (string, error) {
	req := GetFileRequest{
		Key: f.Path(),
	}

	res, err := f.client.GetFileDownloadLink(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (f File) Save(data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := PutFileRequest{
		Key:      f.Path(),
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

func (f File) Upload(filePath string) error {
	// new method added with v2
	panic("implement me")
}

func (f File) GetUploadLink() (string, error) {
	req := GetUploadLinkRequest{
		Key:      f.Path(),
		TempFile: false,
	}

	res, err := f.client.GetFileUploadLink(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (f File) Delete() error {
	req := DeleteFileRequest{
		Key: f.Path(),
	}

	err := f.client.DeleteFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to delete file: %s\n", err.Error())
		return err
	}

	return nil
}

func (f File) Rename(newName string) error {
	req := RenameFileRequest{
		OldKey: f.Path(),
		NewKey: f.parent.Path() + "/" + newName,
	}

	err := f.client.RenameFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to rename file: %s\n", err.Error())
		return err
	}

	f.name = newName
	return nil
}

func (f File) MoveTo(dest sdk.Folder) error {
	req := RenameFileRequest{
		OldKey: f.Path(),
		NewKey: dest.Path() + "/" + f.name,
	}

	err := f.client.RenameFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to rename file: %s\n", err.Error())
		return err
	}

	f.parent = dest
	return nil
}

func (f File) CopyTo(dest sdk.Folder) error {
	// new method added with v2
	panic("implement me")
}
