package cmd_test

import (
	"bytes"
	"errors"
	"project-root/src/cmd"
	"project-root/src/fs"
	"strconv"
	"testing"
)

type fsMockAddCmd struct {
	fs.FileSystemHandler
	storageFile       string
	getStorageFileErr error

	makeDirError   error
	writeFileArgs  []string
	isRelativePath bool

	exists      bool
	existsError error

	storageDir       string
	getStorageDirErr error

	content        string
	writeFileError error
	absPath        string
	absPathError   error
}

func (fs *fsMockAddCmd) GetStorageFile() (string, error) {
	return fs.storageFile, fs.getStorageFileErr
}

func (fs *fsMockAddCmd) GetStorageDir() (string, error) {
	return fs.storageDir, fs.getStorageDirErr
}
func (fs *fsMockAddCmd) MakeDir(path string) error {
	return fs.makeDirError
}
func (fs *fsMockAddCmd) Exists(path string) (bool, error) {
	return fs.exists, fs.existsError
}

func (fs *fsMockAddCmd) GetContentOrEmptyString(path string) string {
	return fs.content
}
func (fs *fsMockAddCmd) IsRelativePath(path string) bool {
	return fs.isRelativePath
}
func (fs *fsMockAddCmd) WriteFile(path string, content string, shouldAppend bool) error {
	fs.writeFileArgs = append(fs.writeFileArgs, path, content, strconv.FormatBool(shouldAppend))
	return fs.writeFileError
}

func (fs *fsMockAddCmd) GetAbsPath(path string) (string, error) {
	return fs.absPath, fs.absPathError
}

func TestAddCmd(t *testing.T) {

	t.Run("returns error if there is an error making the relative path to abs path", func(t *testing.T) {
		fsMock := fsMockAddCmd{
			isRelativePath: true,
			absPathError:   errors.New("could not make relative path to abs path"),
		}

		pathToAdd := "/path2"
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, &fsMock, &buffer)

		if err.Error() != "could not make relative path to abs path" {
			t.Errorf("got wrong error %v", err.Error())
		}

	})

	t.Run("returns error if exists func returns an error", func(t *testing.T) {
		fsMock := fsMockAddCmd{
			isRelativePath: false,
			existsError:    errors.New("no such path in the fs"),
		}

		pathToAdd := "/path2"
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, &fsMock, &buffer)

		if err.Error() != "no such path in the fs" {
			t.Errorf("got wrong error %v", err.Error())
		}

	})

	t.Run("returns error if the path is not in the fs", func(t *testing.T) {
		fsMock := fsMockAddCmd{
			isRelativePath: false,
			exists:         false,
		}

		pathToAdd := "/path2"
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, &fsMock, &buffer)

		if err.Error() != "the path /path2 does not exist" {
			t.Errorf("got wrong error %v", err.Error())
		}

	})

	t.Run("adds the path if there is no nesting", func(t *testing.T) {
		fsMock := fsMockAddCmd{
			content:     "/path1",
			exists:      true,
			storageFile: "/path/to/storage/file",
		}

		pathToAdd := "/path2"
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, &fsMock, &buffer)

		if err != nil {
			t.Errorf(err.Error())
		}

		if buffer.String() != "added /path2" {
			t.Errorf("wrong buffer content %v", buffer.String())
		}

		storageFile := fsMock.writeFileArgs[0]
		addedPath := fsMock.writeFileArgs[1]
		shouldAppend := fsMock.writeFileArgs[2]

		if storageFile != "/path/to/storage/file" {
			t.Errorf("wrong storage file %v", storageFile)
		}

		if addedPath != "/path2\n" {
			t.Errorf("wrong added path %v", addedPath)
		}

		if shouldAppend != "true" {
			t.Errorf("wrong should append %v", shouldAppend)
		}

	})

	t.Run("errors when there is nesting when adding the path", func(t *testing.T) {
		fsMock := fsMockAddCmd{
			content:     "/path1",
			exists:      true,
			storageFile: "/path/to/storage/file",
		}

		pathToAdd := "/path1/nested/dir"
		buffer := bytes.Buffer{}
		err := cmd.RegisterProject(pathToAdd, &fsMock, &buffer)

		if err.Error() != "the path /path1/nested/dir is already a part of a registered project path /path1. to see a list of all registered paths execute the list command " {
			t.Errorf(err.Error())
		}

		if len(buffer.String()) != 0 {
			t.Errorf("wrong buffer content %v", buffer.String())
		}

	})

}
