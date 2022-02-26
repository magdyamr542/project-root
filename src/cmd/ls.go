package cmd

import (
	"fmt"
	"project-root/fs"
	"project-root/utils"
	"strings"
)

// Listing all paths
type LsCmd struct{}

func (lsCmd *LsCmd) Run() error {
	return listProjects()
}

func logEntry(index int, entry string) {
	fmt.Println(index, entry)
}

func listProjects() error {
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
	// If the current dir is part of a saved path then mark it with [current]
	for index, entry := range savedEntries {
		if strings.HasPrefix(cwd, entry) {
			// Determine if this is the current dir
			restOfPath := strings.Replace(cwd, entry, "", 1)
			if len(restOfPath) == 0 || strings.HasPrefix(restOfPath, "/") {
				logEntry(index, fmt.Sprintf("%v [current]", entry))
			} else {
				logEntry(index, entry)
			}

		} else {
			logEntry(index, entry)
		}
	}

	return nil

}
