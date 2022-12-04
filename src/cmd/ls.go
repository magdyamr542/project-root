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

	sort.Slice(savedEntries, func(i, j int) bool {
		return len(savedEntries[i]) < len(savedEntries[j])
	})
	// If the current dir is part of a saved path then mark it with [current]
	for index := range savedEntries {
		entry := savedEntries[index]
		baseName := filepath.Base(entry)
		var toLog string
		if strings.HasPrefix(cwd, entry) {
			// Determine if this is the current dir
			restOfPath := strings.Replace(cwd, entry, "", 1)
			if len(restOfPath) == 0 || strings.HasPrefix(restOfPath, "/") {
				toLog = fmt.Sprintf("%v [current]", baseName)
			} else {
				toLog = baseName
			}

		} else {
			toLog = baseName
		}

		logEntry(toLog, writer)
	}

	return nil

}
