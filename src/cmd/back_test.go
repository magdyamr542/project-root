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

type fsMockBackCmd struct {
	fs.FileSystemHandler
	mock.Mock
}

func (fs *fsMockBackCmd) GetContentOrEmptyString(path string) string {
	args := fs.Called(path)
	return args.String(0)
}

func (fs *fsMockBackCmd) GetLastPathFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func TestBackCmd(t *testing.T) {

	t.Run("throws an error if the retreiving the last path file errors", func(t *testing.T) {
		fsMock := new(fsMockBackCmd)
		fsMock.On("GetLastPathFile").Return("", errors.New("error getting the last path file"))
		buffer := bytes.Buffer{}
		err := cmd.Back(fsMock, &buffer)
		assert.Len(t, buffer.String(), 0)
		assert.EqualError(t, err, "error getting the last path file", "Wrongs error msg")

	})

	t.Run("throws an error if there is no last path the user went from to a project root", func(t *testing.T) {
		fsMock := new(fsMockBackCmd)
		fsMock.On("GetLastPathFile").Return("lastPathFile", nil)
		fsMock.On("GetContentOrEmptyString", "lastPathFile").Return("")
		buffer := bytes.Buffer{}
		err := cmd.Back(fsMock, &buffer)
		assert.Len(t, buffer.String(), 0)
		assert.EqualError(t, err, "there is no last path you went from recently to a project root", "Wrongs error msg")
	})

	t.Run("writes the last path from which the user went to a project root to the writer", func(t *testing.T) {
		fsMock := new(fsMockBackCmd)
		fsMock.On("GetLastPathFile").Return("lastPathFile", nil)
		fsMock.On("GetContentOrEmptyString", "lastPathFile").Return("/some/last/path/")
		buffer := bytes.Buffer{}
		err := cmd.Back(fsMock, &buffer)
		if len(buffer.String()) == 0 {
			assert.Fail(t, "wrong len for buffer")
		}

		assert.ErrorIs(t, err, nil)
		assert.Equal(t, buffer.String(), "/some/last/path/")

	})

}
