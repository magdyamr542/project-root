package cmd_test

import (
	"bytes"
	"project-root/src/cmd"
	"project-root/src/fs"
	"testing"
)

type fsMockGoCmd struct {
	fs.FileSystemHandler
	storageFile string
	content     string
	cwd         string
	err         error
}

func (fs *fsMockGoCmd) GetContentOrEmptyString(path string) string {
	return fs.content
}

func (fs *fsMockGoCmd) GetStorageFile() (string, error) {
	return fs.storageFile, fs.err
}

func (fs *fsMockGoCmd) Cwd() (string, error) {
	return fs.cwd, nil
}

func TestGoCmd(t *testing.T) {

	t.Run("throws an error if the go cmd is used inside a dir that was not saved before", func(t *testing.T) {
		fsMock := fsMockGoCmd{
			content: `/path1
/path2
`,
			err: nil,
			cwd: "/path3",
		}
		buffer := bytes.Buffer{}
		err := cmd.Goto(&fsMock, &buffer)
		if len(buffer.String()) != 0 {
			t.Fail()
		}

		if err.Error() != "the current directory does not belong to a registered project. execute 'pr list' to see all paths" {
			t.Errorf(err.Error())
		}

	})

	t.Run("throws an error if the go cmd is used when there are no saved entries", func(t *testing.T) {
		fsMock := fsMockGoCmd{
			content: "",
			err:     nil,
			cwd:     "/path3",
		}
		buffer := bytes.Buffer{}
		err := cmd.Goto(&fsMock, &buffer)
		if len(buffer.String()) != 0 {
			t.Fail()
		}

		if err.Error() != "there are no registered paths" {
			t.Errorf(err.Error())
		}

	})

	t.Run("writes the path to cd into to the writer if it exists", func(t *testing.T) {
		fsMock := fsMockGoCmd{
			content: `/home/1
/home/2`,
			err: nil,
			cwd: "/home/1/some/nested/dir",
		}
		buffer := bytes.Buffer{}

		err := cmd.Goto(&fsMock, &buffer)
		if len(buffer.String()) == 0 {
			t.Fail()
		}

		if err != nil {
			t.Errorf(err.Error())
		}

		if buffer.String() != "/home/1" {
			t.Errorf("got path %v ", buffer.String())
		}

	})

	t.Run("takes the longest match when having multiple nested paths", func(t *testing.T) {
		fsMock := fsMockGoCmd{
			content: `/home/1/2
/home/1/2/3
/home/1/2/3/4/`,
			err: nil,
			cwd: "/home/1/2/3/5",
		}
		buffer := bytes.Buffer{}

		err := cmd.Goto(&fsMock, &buffer)
		if len(buffer.String()) == 0 {
			t.Fail()
		}

		if err != nil {
			t.Errorf(err.Error())
		}

		if buffer.String() != "/home/1/2/3" {
			t.Errorf("got path %v ", buffer.String())
		}

	})

}
