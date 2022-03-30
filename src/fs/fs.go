package fs

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

var storageDir string = ".proot"
var storageFile string = "storage.txt"
var lastPathFile string = "lastPath.txt"

type FileSystemHandler interface {
	ReadFile(path string) (string, error)
	Exists(path string) (bool, error)
	WriteFile(filepath, content string, append bool) error
	MakeDir(dirPath string) error
	GetHomeDir() (string, error)
	GetAbsPath(path string) (string, error)
	Cwd() (string, error)
	IsRelativePath(path string) bool
	GetStorageDir() (string, error)
	GetStorageFile() (string, error)
	GetLastPathFile() (string, error)
	GetContentOrEmptyString(path string) string
}

type FileSystem struct{}

func (fs *FileSystem) ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (fs *FileSystem) Exists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fs *FileSystem) WriteFile(filePath string, content string, append bool) error {
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

func (fs *FileSystem) MakeDir(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileSystem) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (fs *FileSystem) GetAbsPath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

func (fs *FileSystem) Cwd() (string, error) {
	return os.Getwd()
}

func (fs *FileSystem) IsRelativePath(filePath string) bool {
	return !filepath.IsAbs(filePath)
}

func (fs *FileSystem) GetStorageDir() (string, error) {
	homePath, err := fs.GetHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homePath, storageDir), nil
}

func (fs *FileSystem) GetStorageFile() (string, error) {
	homePath, err := fs.GetHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homePath, storageDir, storageFile), nil
}

func (fs *FileSystem) GetLastPathFile() (string, error) {
	homePath, err := fs.GetHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homePath, storageDir, lastPathFile), nil
}

// - reads the content of the passed path
func (fs *FileSystem) GetContentOrEmptyString(path string) string {
	text, err := fs.ReadFile(path)
	if err != nil {
		return ""
	}
	return text
}

var DefaultFileSystem FileSystemHandler = &FileSystem{}
