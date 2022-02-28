package fs

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

var storageDir string = ".proot"
var storageFile string = "storage.txt"

func ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Exists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func WriteFile(filePath string, content string, append bool) error {
	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(filePath, flags, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return err
	}

	return nil
}

func MakeDir(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func GetAbsPath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

func Cwd() (string, error) {
	return os.Getwd()
}

func IsRelativePath(filePath string) bool {
	return !filepath.IsAbs(filePath)
}

func GetStorageDir() (string, error) {
	homePath, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homePath, storageDir), nil
}

func GetStorageFile() (string, error) {
	homePath, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homePath, storageDir, storageFile), nil
}

// - reads the content of the passed path
func GetContentOrEmptyString(path string) string {
	text, err := ReadFile(path)
	if err != nil {
		return ""
	}
	return text
}
