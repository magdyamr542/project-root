package cmd_test

import (
	"errors"
	"project-root/src/cmd"
	"project-root/src/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fsMockClearCmd struct {
	fs.FileSystemHandler
	mock.Mock
}

func (fs *fsMockClearCmd) GetStorageFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func (fs *fsMockClearCmd) WriteFile(path string, content string, shouldAppend bool) error {
	args := fs.Called(path, content, shouldAppend)
	return args.Error(0)
}

func TestClearCmd(t *testing.T) {

	t.Run("clears the saved paths on clear cmd", func(t *testing.T) {
		fsMock := new(fsMockClearCmd)
		fsMock.On("GetStorageFile").Return("/file/path", nil)
		fsMock.On("WriteFile", "/file/path", "", false).Return(nil)
		cmd.Clear(fsMock)

	})

	t.Run("returns error on error from the clear func", func(t *testing.T) {
		fsMock := new(fsMockClearCmd)
		fsMock.On("GetStorageFile").Return("/file/path", nil)
		fsMock.On("WriteFile", "/file/path", "", false).Return(errors.New("some error from clear func"))
		err := cmd.Clear(fsMock)
		assert.EqualError(t, err, "couldn't clear the database of saved project roots. Does file exist ?", "wrong err msg")

	})

}
