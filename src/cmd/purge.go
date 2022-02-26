package cmd

import (
	"fmt"
	"project-root/fs"
	"project-root/utils"
	"strings"
)

// Purge
type PurgeCmd struct{}

// read the storage file and for each path that does not exist
// in the fs anymore rm it.
func (purgeCmd *PurgeCmd) Run() error {
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

	entriesToDelete := []string{}

	// get entries to delete
	for _, entry := range savedEntries {
		doesExist, err := fs.Exists(entry)
		if err != nil {
			return err
		}
		if !doesExist {
			entriesToDelete = append(entriesToDelete, entry)
		}
	}

	if len(entriesToDelete) == 0 {
		return nil
	}

	// get filtered entries that should not be deleted
	filteredEntries := utils.Filter(savedEntries, func(entry string) bool {
		return !utils.Contains(entriesToDelete, entry)
	})

	// write the storage file again
	err = fs.WriteFile(storageFile, strings.Join(filteredEntries, "\n"), false)
	if err != nil {
		return err
	}

	// show the entries that have been deleted
	pathStrings := ""
	if len(entriesToDelete) == 1 {
		pathStrings = "path"
	} else {
		pathStrings = "paths"
	}

	fmt.Println("deleted", len(entriesToDelete), pathStrings)
	for _, deletedPath := range entriesToDelete {
		fmt.Printf("%v", deletedPath)
	}

	return nil

}
