package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"project-root/src/fs"
	"project-root/src/utils"
	"sort"
	"strings"
)

// Listing all paths
type LsCmd struct {
}

func (lsCmd *LsCmd) Run(fs fs.FileSystemHandler) error {
	return ListProjects(fs, os.Stdout)
}

func logEntry(entry string, writer io.Writer) {
	fmt.Fprintln(writer, entry)
}

func ListProjects(fs fs.FileSystemHandler, writer io.Writer) error {
	// Get saved data
	storageFile, err := fs.GetStorageFile()
	if err != nil {
		return err
	}
	savedData := fs.GetContentOrEmptyString(storageFile)

	if len(savedData) == 0 {
		return nil
	}

	savedEntries := utils.Filter(strings.Split(savedData, "\n"), func(entry string) bool {
		return len(entry) != 0
	})

	cwd, err := fs.Cwd()
	if err != nil {
		return err
	}

	savedEntries = utils.Map(savedEntries, func(entry string) string {
		return filepath.Base(entry)
	})
	sort.Strings(savedEntries)

	for index := range savedEntries {
		toLog := savedEntries[index]
		if strings.Contains(cwd, toLog) {
			toLog += " (current)"
		}

		logEntry(toLog, writer)
	}

	return nil

}
