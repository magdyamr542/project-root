package cmd_test

import (
	"bytes"
	"errors"
	"project-root/src/cmd"
	"project-root/src/fs"
	"testing"
)

type fsMockLsCmd struct {
	fs.FileSystemHandler
	cwd         string
	storageFile string
	content     string
	err         error
}

func (fs *fsMockLsCmd) GetContentOrEmptyString(path string) string {
	return fs.content
}

func (fs *fsMockLsCmd) GetStorageFile() (string, error) {
	return fs.storageFile, fs.err
}

func (fs *fsMockLsCmd) Cwd() (string, error) {
	return fs.cwd, fs.err
}

func TestLsCmd(t *testing.T) {

	t.Run("lists all saved paths", func(t *testing.T) {

		fsMock := fsMockLsCmd{
			content: "path1\npath2\npath3\n",
		}
		buffer := bytes.Buffer{}
		cmd.ListProjects(&fsMock, &buffer)
		want := `0 path1
1 path2
2 path3
`
		got := buffer.String()

		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}
	})

	t.Run("returns empty string on empty storage", func(t *testing.T) {
		fsMock := fsMockLsCmd{
			content: "",
		}
		buffer := bytes.Buffer{}
		cmd.ListProjects(&fsMock, &buffer)
		want := ""
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})

	t.Run("returns error if storage file not found", func(t *testing.T) {
		fsMock := fsMockLsCmd{
			content: "",
			err:     errors.New("storage file not found"),
		}
		buffer := bytes.Buffer{}
		err := cmd.ListProjects(&fsMock, &buffer)

		if buffer.Len() != 0 {
			t.Fatalf("writer should be empty. file not found but it's not. it's value is %v", buffer.String())
		}

		if err.Error() != "storage file not found" {
			t.Fail()
		}

	})

	t.Run("appends current prefix to path if we are inside it", func(t *testing.T) {
		fsMock := fsMockLsCmd{
			content: `path1
path2
path3
path4
`,
			err: nil,
			cwd: "path2/some/file.go",
		}
		buffer := bytes.Buffer{}

		cmd.ListProjects(&fsMock, &buffer)
		want := `0 path1
1 path2 [current]
2 path3
3 path4
`
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})

	t.Run("appends current prefix to path if we are inside it second case", func(t *testing.T) {
		fsMock := fsMockLsCmd{
			content: `path1
path2
path3
path4
`,
			err: nil,
			cwd: "path2",
		}
		buffer := bytes.Buffer{}

		cmd.ListProjects(&fsMock, &buffer)
		want := `0 path1
1 path2 [current]
2 path3
3 path4
`
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}
	})

	t.Run("appends current prefix to path if we are inside it third case", func(t *testing.T) {

		fsMock := fsMockLsCmd{
			content: `path1
path2
path2/dir1
path3
path4
`,
			err: nil,
			cwd: "path2/dir1/dir2/dir3/",
		}
		buffer := bytes.Buffer{}

		cmd.ListProjects(&fsMock, &buffer)
		want := `0 path1
1 path2 [current]
2 path2/dir1 [current]
3 path3
4 path4
`
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})

}
