package cmd

import (
	"fmt"
	"io"
	"os"

	"project-root/src/fs"
	"project-root/src/utils"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// ToCmd uses the provided path prefix to go to a possible saved project
type ToCmd struct {
	Path string `arg:"" name:"path prefix" help:"Path prefix to go to"`
}

// Run implements the ToCmd
func (toCmd *ToCmd) Run(fs fs.FileSystemHandler) error {
	return To(fs, toCmd.Path, os.Stdout)
}

// To is the impl of the To Cmd
func To(fs fs.FileSystemHandler, pathPrefix string, writer io.Writer) error {
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
		pathLower := strings.ToLower(entry)
		prefixLower := strings.ToLower(pathPrefix)
		return len(entry) != 0 && fuzzy.Match(prefixLower, pathLower)
	})

	if len(savedEntries) == 0 {
		return fmt.Errorf("the prefix '%v' does not appear in any of the saved paths", pathPrefix)
	}

	if len(savedEntries) > 1 {
		return fmt.Errorf("the prefix '%s' appears in more than one path. be more specific.\n%v", pathPrefix, strings.Join(savedEntries, "\n"))
	}

	fmt.Fprintf(writer, "%v", savedEntries[0])

	return nil
}
