package fs

import (
	"errors"
	"os"
	"path/filepath"
)

func ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Exists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
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

func IsRelativePath(filePath string) bool {
	return !filepath.IsAbs(filePath)
}
