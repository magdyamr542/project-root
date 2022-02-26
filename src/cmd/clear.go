package cmd

import (
	"fmt"
	"project-root/src/fs"
)

// Clearing the db

type ClearCmd struct{}

func (clrCmd *ClearCmd) Run() error {
	if err := clear(); err != nil {
		return fmt.Errorf("error couldn't clear the database of saved project roots. Does file exist ?")
	}
	return nil
}

func clear() error {
	// Get saved data
	storageFile, err := fs.GetStorageFile()
	if err != nil {
		return err
	}
	return fs.WriteFile(storageFile, "", false)
}
