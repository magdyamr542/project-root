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

type fsMockGoCmd struct {
	fs.FileSystemHandler
	mock.Mock
}

func (fs *fsMockGoCmd) GetContentOrEmptyString(path string) string {
	args := fs.Called(path)
	return args.String(0)
}

func (fs *fsMockGoCmd) GetStorageFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func (fs *fsMockGoCmd) GetLastPathFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func (fs *fsMockGoCmd) WriteFile(path string, content string, shouldAppend bool) error {
	args := fs.Called(path, content, shouldAppend)
	return args.Error(0)
}

func (fs *fsMockGoCmd) Cwd() (string, error) {

	args := fs.Called()
	return args.String(0), args.Error(1)
}

func TestGoCmd(t *testing.T) {

	t.Run("throws an error if the go cmd is used inside a dir that was not saved before", func(t *testing.T) {
		fsMock := new(fsMockGoCmd)
		fsMock.On("GetStorageFile").Return("file", nil)
		fsMock.On("GetContentOrEmptyString", "file").Return("path1\npath2\n")
		fsMock.On("Cwd").Return("/path3", nil)
		buffer := bytes.Buffer{}
		err := cmd.Goto(fsMock, &buffer)
		assert.Len(t, buffer.String(), 0)
		assert.EqualError(t, err, "the current directory does not belong to a registered project. execute 'pr list' to see all paths", "Wrongs error msg")

	})

	t.Run("throws an error if the go cmd is used when there are no saved entries", func(t *testing.T) {
		fsMock := new(fsMockGoCmd)
		fsMock.On("GetStorageFile").Return("file", nil)
		fsMock.On("GetContentOrEmptyString", "file").Return("")
		fsMock.On("Cwd").Return("/path3", nil)
		buffer := bytes.Buffer{}
		err := cmd.Goto(fsMock, &buffer)
		assert.Len(t, buffer.String(), 0)
		assert.EqualError(t, err, "there are no registered paths", "Wrongs error msg")

	})

	t.Run("writes the path to cd into to the writer if it exists", func(t *testing.T) {
		fsMock := new(fsMockGoCmd)
		fsMock.On("GetStorageFile").Return("file", nil)
		fsMock.On("GetContentOrEmptyString", "file").Return("/home/1\n/home/2\n")
		fsMock.On("Cwd").Return("/home/1/some/nested/dir", nil)
		fsMock.On("GetLastPathFile").Return("lastPathFile", nil)
		fsMock.On("WriteFile", "lastPathFile", "/home/1/some/nested/dir", false).Return(nil)
		buffer := bytes.Buffer{}

		err := cmd.Goto(fsMock, &buffer)
		if len(buffer.String()) == 0 {
			assert.Fail(t, "wrong len for buffer")
		}

		assert.ErrorIs(t, err, nil)
		assert.Equal(t, buffer.String(), "/home/1")

	})

	t.Run("takes the longest match when having multiple nested paths", func(t *testing.T) {
		fsMock := new(fsMockGoCmd)
		fsMock.On("GetStorageFile").Return("file", nil)
		fsMock.On("GetContentOrEmptyString", "file").Return(`/home/1/2
/home/1/2/3
/home/1/2/3/4`)
		fsMock.On("Cwd").Return("/home/1/2/3/5", nil)
		fsMock.On("GetLastPathFile").Return("lastPathFile", nil)
		fsMock.On("WriteFile", "lastPathFile", "/home/1/2/3/5", false).Return(nil)
		buffer := bytes.Buffer{}

		err := cmd.Goto(fsMock, &buffer)
		if len(buffer.String()) == 0 {
			assert.Fail(t, "wrong buffer len")
		}

		assert.ErrorIs(t, err, nil)
		assert.Equal(t, buffer.String(), "/home/1/2/3")

	})

	t.Run("writes the cwd to the lastPathFile", func(t *testing.T) {
		fsMock := new(fsMockGoCmd)
		fsMock.On("GetStorageFile").Return("file", nil)
		fsMock.On("GetContentOrEmptyString", "file").Return(`/home/1/2
/home/1/2/3
/home/1/2/3/4`)
		fsMock.On("Cwd").Return("/home/1/2/3/5", nil)
		fsMock.On("GetLastPathFile").Return("lastPathFile.txt", nil)
		fsMock.On("WriteFile", "lastPathFile.txt", "/home/1/2/3/5", false).Return(nil)

		buffer := bytes.Buffer{}
		err := cmd.Goto(fsMock, &buffer)
		assert.ErrorIs(t, err, nil)

	})

	t.Run("returns error when failing to write the cwd to the last path file", func(t *testing.T) {
		fsMock := new(fsMockGoCmd)
		fsMock.On("GetStorageFile").Return("file", nil)
		fsMock.On("GetContentOrEmptyString", "file").Return(`/home/1/2
/home/1/2/3
/home/1/2/3/4`)
		fsMock.On("Cwd").Return("/home/1/2/3/5", nil)
		fsMock.On("GetLastPathFile").Return("lastPathFile.txt", nil)
		fsMock.On("WriteFile", "lastPathFile.txt", "/home/1/2/3/5", false).Return(errors.New("could not write to lastPathFile.txt"))

		buffer := bytes.Buffer{}
		err := cmd.Goto(fsMock, &buffer)
		assert.EqualError(t, err, "could not write to lastPathFile.txt", "Wrongs error msg")

	})

}
