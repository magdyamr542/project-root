package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"project-root/src/fs"
)

type BackCmd struct{}

// Run go back to the path you went to the project root form
func (backCmd *BackCmd) Run(fs fs.FileSystemHandler) error {
	return Back(fs, os.Stdout)
}

// Back prints the last path from which the user went to a saved project root
func Back(fs fs.FileSystemHandler, writer io.Writer) error {

	lastPathFile, err := fs.GetLastPathFile()
	if err != nil {
		return err
	}

	lastPath := fs.GetContentOrEmptyString(lastPathFile)
	if len(lastPath) == 0 {
		return errors.New("there is no last path you went from recently to a project root")
	}

	fmt.Fprintf(writer, "%v", lastPath)

	return nil

}
