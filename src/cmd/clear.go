package cmd

import (
	"fmt"
	"github.com/magdyamr542/project-root/src/fs"
)

// Clearing the db

type ClearCmd struct{}

func (clrCmd *ClearCmd) Run(fs fs.FileSystemHandler) error {
	return Clear(fs)
}

func Clear(fs fs.FileSystemHandler) error {

	if err := clear(fs); err != nil {
		return fmt.Errorf("couldn't clear the database of saved project roots. Does file exist ?")
	}
	return nil
}

func clear(fs fs.FileSystemHandler) error {
	// Get saved data
	storageFile, err := fs.GetStorageFile()
	if err != nil {
		return err
	}
	return fs.WriteFile(storageFile, "", false)
}
