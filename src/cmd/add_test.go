package cmd_test

import (
	"bytes"
	"errors"
	"project-root/src/cmd"
	"project-root/src/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fsMockAddCmd struct {
	fs.FileSystemHandler
	mock.Mock
}

func (fs *fsMockAddCmd) GetStorageFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func (fs *fsMockAddCmd) GetStorageDir() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}
func (fs *fsMockAddCmd) MakeDir(path string) error {
	args := fs.Called(path)
	return args.Error(0)
}
func (fs *fsMockAddCmd) Exists(path string) (bool, error) {
	args := fs.Called(path)
	return args.Bool(0), args.Error(1)
}

func (fs *fsMockAddCmd) GetContentOrEmptyString(path string) string {
	args := fs.Called(path)
	return args.String(0)
}
func (fs *fsMockAddCmd) IsRelativePath(path string) bool {
	args := fs.Called(path)
	return args.Bool(0)
}
func (fs *fsMockAddCmd) WriteFile(path string, content string, shouldAppend bool) error {
	args := fs.Called(path, content, shouldAppend)
	return args.Error(0)
}

func (fs *fsMockAddCmd) GetAbsPath(path string) (string, error) {
	args := fs.Called(path)
	return args.String(0), args.Error(1)
}

func TestAddCmd(t *testing.T) {

	t.Run("returns error if there is an error making the relative path to abs path", func(t *testing.T) {
		pathToAdd := "/path2"
		fsMock := new(fsMockAddCmd)
		fsMock.On("IsRelativePath", pathToAdd).Return(true)
		fsMock.On("GetAbsPath", pathToAdd).Return("", errors.New("could not make relative path to abs path"))
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, fsMock, &buffer)
		assert.EqualError(t, err, "could not make relative path to abs path")

	})

	t.Run("returns error if exists func returns an error", func(t *testing.T) {
		pathToAdd := "/path2"
		fsMock := new(fsMockAddCmd)
		fsMock.On("IsRelativePath", pathToAdd).Return(false)
		fsMock.On("GetAbsPath", pathToAdd).Return(true, errors.New("could not make relative path to abs path"))
		fsMock.On("Exists", pathToAdd).Return(false, errors.New("no such file in the fs"))

		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, fsMock, &buffer)

		assert.EqualError(t, err, "no such file in the fs")

	})

	t.Run("returns error if the path is not in the fs", func(t *testing.T) {
		pathToAdd := "/path2"
		fsMock := new(fsMockAddCmd)
		fsMock.On("IsRelativePath", pathToAdd).Return(false)
		fsMock.On("GetAbsPath", pathToAdd).Return(true, nil)
		fsMock.On("Exists", pathToAdd).Return(false, nil)
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, fsMock, &buffer)
		assert.EqualError(t, err, "the path /path2 does not exist")

	})

	t.Run("errors when there is nesting when adding the path", func(t *testing.T) {
		pathToAdd := "/path1/nested/dir"
		storageDir := "/path/to/storage/dir"
		fsMock := new(fsMockAddCmd)
		fsMock.On("IsRelativePath", pathToAdd).Return(false)
		fsMock.On("GetStorageFile").Return("/path/to/storage/file", nil)
		fsMock.On("GetContentOrEmptyString", "/path/to/storage/file").Return("/path1")
		fsMock.On("Exists", pathToAdd).Return(true, nil)
		fsMock.On("Exists", storageDir).Return(true, nil)
		fsMock.On("GetStorageDir").Return(storageDir, nil)

		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, fsMock, &buffer)

		if err.Error() != "the path /path1/nested/dir is already a part of a registered project path /path1. to see a list of all registered paths execute the list command " {
			t.Errorf(err.Error())
		}

		assert.EqualError(t, err, "the path /path1/nested/dir is already a part of a registered project path /path1. to see a list of all registered paths execute the list command ", "wrong error msg")
		assert.Len(t, buffer.String(), 0)

	})

	t.Run("adds the path if there is no nesting", func(t *testing.T) {

		pathToAdd := "/path2"
		storageDir := "/path/to/storage/"
		storageFile := "/path/to/storage/file"
		fsMock := new(fsMockAddCmd)
		fsMock.On("IsRelativePath", pathToAdd).Return(false)
		fsMock.On("GetStorageFile").Return(storageFile, nil)
		fsMock.On("GetContentOrEmptyString", storageFile).Return("/path1")
		fsMock.On("Exists", pathToAdd).Return(true, nil)
		fsMock.On("Exists", storageDir).Return(true, nil)
		fsMock.On("GetStorageDir").Return(storageDir, nil)
		fsMock.On("WriteFile", storageFile, "/path2\n", true).Return(nil)

		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, fsMock, &buffer)

		assert.ErrorIs(t, err, nil, "error should have been nil")
		assert.Equal(t, buffer.String(), "added /path2", "wrong buffer content")
	})

}
