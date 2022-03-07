package cmd_test

import (
	"bytes"
	"project-root/src/cmd"
	"project-root/src/fs"
	"testing"
)

type fsMockPurgeCmd struct {
	fs.FileSystemHandler
	storageFile string
	content     string
	pathsMap    map[string]bool // maps to false if the path does not exist
	err         error
}

func (fs *fsMockPurgeCmd) GetContentOrEmptyString(path string) string {
	return fs.content
}

func (fs *fsMockPurgeCmd) GetStorageFile() (string, error) {
	return fs.storageFile, fs.err
}

func (fs *fsMockPurgeCmd) Exists(path string) (bool, error) {
	exists, ok := fs.pathsMap[path]
	return exists && ok, fs.err
}

func (fs *fsMockPurgeCmd) WriteFile(path string, content string, append bool) error {
	return fs.err
}

func TestPurgeCmd(t *testing.T) {

	t.Run("does not do anything if all saved paths still exist in the fs", func(t *testing.T) {
		fsMock := fsMockPurgeCmd{
			content: `path1
path2
path3
path4
`,
			err: nil,
			pathsMap: map[string]bool{
				"path1": true,
				"path2": true,
				"path3": true,
				"path4": true,
			},
		}
		buffer := bytes.Buffer{}

		cmd.Purge(&fsMock, &buffer)
		want := ""
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})

	t.Run("deletes all paths which were saved before but does not exist in the fs anymore", func(t *testing.T) {
		fsMock := fsMockPurgeCmd{
			content: `path1
path2
path3
path4
`,
			err: nil,
			pathsMap: map[string]bool{
				"path1": false,
				"path2": true,
				"path3": false,
				"path4": true,
			},
		}
		buffer := bytes.Buffer{}

		cmd.Purge(&fsMock, &buffer)
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
		fsMock := fsMockPurgeCmd{
			content: `path1
path2
path3
path4
`,
			err: nil,
			pathsMap: map[string]bool{
				"path1": true,
				"path2": true,
				"path3": false,
				"path4": true,
			},
		}
		buffer := bytes.Buffer{}

		cmd.Purge(&fsMock, &buffer)
		want := `deleted 1 path
path3
`
		got := buffer.String()
		if want != got {
			t.Fatalf("expected %v but got %v", want, got)
		}

	})
}
