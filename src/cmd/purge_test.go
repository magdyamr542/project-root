package cmd_test

import (
	"bytes"
	"testing"

	"github.com/magdyamr542/project-root/src/cmd"
	"github.com/magdyamr542/project-root/src/fs"

	"github.com/stretchr/testify/mock"
)

type fsMockPurgeCmd struct {
	fs.FileSystemHandler
	mock.Mock
}

func (fs *fsMockPurgeCmd) GetContentOrEmptyString(path string) string {
	args := fs.Called(path)
	return args.String(0)
}

func (fs *fsMockPurgeCmd) GetStorageFile() (string, error) {
	args := fs.Called()
	return args.String(0), args.Error(1)
}

func (fs *fsMockPurgeCmd) Exists(path string) (bool, error) {
	args := fs.Called(path)
	return args.Bool(0), args.Error(1)
}

func (fs *fsMockPurgeCmd) WriteFile(path string, content string, append bool) error {
	args := fs.Called(path, content, append)
	return args.Error(0)
}

func TestPurgeCmd(t *testing.T) {

	t.Run("does not do anything if all saved paths still exist in the fs", func(t *testing.T) {
		fsMock := new(fsMockPurgeCmd)
		fsMock.On("GetContentOrEmptyString", "").Return(`path1
path2
path3
path4
`)
		fsMock.On("Exists", "path1").Return(true, nil)
		fsMock.On("Exists", "path2").Return(true, nil)
		fsMock.On("Exists", "path3").Return(true, nil)
		fsMock.On("Exists", "path4").Return(true, nil)
		fsMock.On("GetStorageFile").Return("", nil)
		buffer := bytes.Buffer{}

		cmd.Purge(fsMock, &buffer)
		want := ""
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})

	t.Run("deletes all paths which were saved before but does not exist in the fs anymore", func(t *testing.T) {

		fsMock := new(fsMockPurgeCmd)

		fsMock.On("GetStorageFile").Return("", nil)
		fsMock.On("GetContentOrEmptyString", "").Return(`path1
path2
path3
path4
`)
		fsMock.On("Exists", "path1").Return(false, nil)
		fsMock.On("Exists", "path2").Return(true, nil)
		fsMock.On("Exists", "path3").Return(false, nil)
		fsMock.On("Exists", "path4").Return(true, nil)

		fsMock.On("WriteFile", "", "path2\npath4", false).Return(nil)
		buffer := bytes.Buffer{}

		cmd.Purge(fsMock, &buffer)
		want := `deleted 2 paths
path1
path3
`
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})

	t.Run("deletes one path which was saved before but does not exist in the fs anymore", func(t *testing.T) {

		fsMock := new(fsMockPurgeCmd)
		fsMock.On("GetStorageFile").Return("", nil)
		fsMock.On("GetContentOrEmptyString", "").Return(`path1
path2
path3
path4
`)

		fsMock.On("Exists", "path1").Return(true, nil)
		fsMock.On("Exists", "path2").Return(true, nil)
		fsMock.On("Exists", "path3").Return(false, nil)
		fsMock.On("Exists", "path4").Return(true, nil)

		fsMock.On("WriteFile", "", "path1\npath2\npath4", false).Return(nil)
		buffer := bytes.Buffer{}

		cmd.Purge(fsMock, &buffer)
		want := `deleted 1 path
path3
`
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})
}
