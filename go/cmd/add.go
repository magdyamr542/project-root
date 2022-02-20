package cmd

import (
	"fmt"
	"project-root/fs"
	"project-root/utils"
	"strings"
)

func RegisterProject(path string) error {
	// Turn to abs path
	if fs.IsRelativePath(path) {
		absPath, err := fs.GetAbsPath(path)
		path = absPath
		if err != nil {
			return err
		}
	}

	// Check path exists
	exists, err := fs.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the path %s does not exist", path)
	}

	// Check storage dir exists
	storageDirPath, err := fs.GetStorageDir()
	if err != nil {
		return err
	}
	storageDirExists, err := fs.Exists(storageDirPath)
	if err != nil {
		return err
	}
	if !storageDirExists {
		// create the storage dir
		err := fs.MakeDir(storageDirPath)
		if err != nil {
			return err
		}
	}

	// Get saved data
	storageFile, err := fs.GetStorageFile()
	if err != nil {
		return err
	}
	savedData := fs.GetContentOrEmptyString(storageFile)

	// Don't handle nested case
	err = tryGetAlreadyRegisteredPath(savedData, path)
	if err != nil {
		return err
	}

	// Append the new path to the storage file
	err = fs.WriteFile(storageFile, path+"\n", true)
	if err != nil {
		return err
	}

	fmt.Printf("Success Added %s", path)

	return nil
}

func tryGetAlreadyRegisteredPath(savedData string, pathToRegister string) error {

	savedEntries := utils.Filter(strings.Split(savedData, "\n"), func(entry string) bool {
		return len(entry) != 0
	})

	for _, entry := range savedEntries {
		if strings.HasPrefix(pathToRegister, entry) {
			restOfPath := strings.Replace(pathToRegister, entry, "", 1)
			// Handle case of same path but two directories with same prefix name
			// /home/ok , /home/ok1 for example [this is an example of a positive case]
			if !strings.Contains(restOfPath, "/") && len(restOfPath) > 0 {
				continue
			}
			return fmt.Errorf("error the path %v is already a part of a registered project path %v. to see a list of all registered paths execute the list command ", pathToRegister, entry)
		}
	}

	return nil
}