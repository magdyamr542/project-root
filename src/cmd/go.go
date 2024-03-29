package cmd

import (
	"fmt"
	"github.com/magdyamr542/project-root/src/fs"
	"github.com/magdyamr542/project-root/src/utils"
	"io"
	"os"
	"sort"
	"strings"
)

// Going to the root

type GoCmd struct{}

// Use the database of saved project roots to go to the root directory of the project
// Print the directory to the console such that the bash script can cd to it.
func (goCmd *GoCmd) Run(fs fs.FileSystemHandler) error {

	return Goto(fs, os.Stdout)
}

func Goto(fs fs.FileSystemHandler, writer io.Writer) error {
	// Get saved data
	storageFile, err := fs.GetStorageFile()
	if err != nil {
		return err
	}

	savedData := fs.GetContentOrEmptyString(storageFile)
	if len(savedData) == 0 {
		return fmt.Errorf("there are no registered paths")
	}

	savedEntries := utils.Filter(strings.Split(savedData, "\n"), func(entry string) bool {
		return len(entry) != 0
	})

	// get the cwd
	cwd, err := fs.Cwd()
	if err != nil {
		return err
	}

	pathMatches := utils.Filter(savedEntries, func(path string) bool {
		return strings.HasPrefix(cwd, path)
	})

	if len(pathMatches) == 0 {
		return fmt.Errorf("the current directory does not belong to a registered project. execute 'pr list' to see all paths")
	}

	// Take the longest match. should be only 1 path anyway as we are not handling nested roots.
	gotoPath := ""
	if len(pathMatches) == 1 {
		gotoPath = pathMatches[0]
	} else {
		sort.Slice(pathMatches, func(i, j int) bool {
			return len(pathMatches[i]) < len(pathMatches[j])
		})
		gotoPath = pathMatches[len(pathMatches)-1]
	}

	// save the cwd as the last path we were at before going to the root
	lastPathFile, err := fs.GetLastPathFile()
	if err != nil {
		return err
	}
	err = fs.WriteFile(lastPathFile, cwd, false)
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "%v", gotoPath)

	return nil
}
