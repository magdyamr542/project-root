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

type fsMockLsCmd struct {
	fs.FileSystemHandler
	mock.Mock
}

func (fs *fsMockLsCmd) GetContentOrEmptyString(path string) string {
	args := fs.Called(path)
	return args.String(0)
}

func (fs *fsMockLsCmd) GetStorageFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func (fs *fsMockLsCmd) Cwd() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func TestLsCmd(t *testing.T) {

	t.Run("lists all saved paths", func(t *testing.T) {
		fsMock := new(fsMockLsCmd)
		buffer := bytes.Buffer{}
		fsMock.On("GetContentOrEmptyString", "").Return("path1\npath2\npath3\n")
		fsMock.On("GetStorageFile").Return("", nil)
		fsMock.On("Cwd").Return("", nil)
		cmd.ListProjects(fsMock, &buffer)
		want := `path1
path2
path3
`
		got := buffer.String()

		assert.Equal(t, want, got)

	})

	t.Run("returns empty string on empty storage", func(t *testing.T) {
		fsMock := new(fsMockLsCmd)
		buffer := bytes.Buffer{}
		fsMock.On("GetContentOrEmptyString", "").Return("")
		fsMock.On("GetStorageFile").Return("", nil)
		cmd.ListProjects(fsMock, &buffer)
		want := ""
		got := buffer.String()
		assert.Equal(t, want, got)

	})

	t.Run("returns error if storage file not found", func(t *testing.T) {
		fsMock := new(fsMockLsCmd)
		fsMock.On("GetStorageFile").Return("", errors.New("storage file not found"))
		buffer := bytes.Buffer{}
		err := cmd.ListProjects(fsMock, &buffer)

		assert.Len(t, buffer.String(), 0)
		assert.EqualErrorf(t, err, "storage file not found", "wrong error")

	})

	t.Run("appends current prefix to path if we are inside it", func(t *testing.T) {
		fsMock := new(fsMockLsCmd)
		fsMock.On("GetStorageFile").Return("", nil)
		fsMock.On("Cwd").Return("path2/some/file.go", nil)
		fsMock.On("GetContentOrEmptyString", "").Return(`path1
path2
path3
path4
`)
		buffer := bytes.Buffer{}

		cmd.ListProjects(fsMock, &buffer)
		want := `path1
path2 (current)
path3
path4
`

		got := buffer.String()
		assert.Equal(t, want, got)

	})

	t.Run("appends current prefix to path if we are inside it second case", func(t *testing.T) {

		fsMock := new(fsMockLsCmd)
		fsMock.On("GetStorageFile").Return("", nil)
		fsMock.On("Cwd").Return("path2", nil)
		fsMock.On("GetContentOrEmptyString", "").Return(`path1
path2
path3
path4
`)
		buffer := bytes.Buffer{}

		cmd.ListProjects(fsMock, &buffer)
		want := `path1
path2 (current)
path3
path4
`
		got := buffer.String()
		assert.Equal(t, want, got)
	})

}
