package cmd_test

import (
	"errors"
	"project-root/src/cmd"
	"project-root/src/fs"
	"strconv"
	"testing"
)

type fsMockClearCmd struct {
	fs.FileSystemHandler
	storageFile string
	err         error
	clearArgs   []string
}

func (fs *fsMockClearCmd) GetStorageFile() (string, error) {
	return fs.storageFile, fs.err
}

func (fs *fsMockClearCmd) WriteFile(path string, content string, shouldAppend bool) error {
	fs.clearArgs = append(fs.clearArgs, path, content, strconv.FormatBool(shouldAppend))
	return fs.err
}

func TestClearCmd(t *testing.T) {

	t.Run("clears the saved paths on clear cmd", func(t *testing.T) {
		fsMock := fsMockClearCmd{
			err:         nil,
			storageFile: "/some/path",
			clearArgs:   []string{},
		}
		cmd.Clear(&fsMock)

		filePath := fsMock.clearArgs[0]
		contentToWrite := fsMock.clearArgs[1]
		shouldAppend := fsMock.clearArgs[2]

		if filePath != "/some/path" {
			t.Errorf("got wrong path %v", filePath)
		}

		if len(contentToWrite) != 0 {
			t.Errorf("got wrong content to write %v", contentToWrite)
		}

		if shouldAppend != "false" {
			t.Errorf("got wrong shouldAppend %v", shouldAppend)
		}

	})

	t.Run("returns error on error from the clear func", func(t *testing.T) {
		fsMock := fsMockClearCmd{
			err:         errors.New("some error"),
			storageFile: "/some/path",
			clearArgs:   []string{},
		}
		err := cmd.Clear(&fsMock)

		if err.Error() != "couldn't clear the database of saved project roots. Does file exist ?" {
			t.Errorf(err.Error())
		}

	})

}
