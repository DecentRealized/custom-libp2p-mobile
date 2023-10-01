package file_handler

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"sync"
)

var openFileCache = &sync.Map{}
var downloadPath = "./"

func Reset() error {
	openFileCache = &sync.Map{}
	downloadPath = "./"
	return nil
}

// GetSHA256Sum gets sha25 sum
func GetSHA256Sum(file *os.File) (string, error) {
	// Seek to start
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	raw := h.Sum(nil)
	// Seek to start
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// GetFileSize returns file size
func GetFileSize(path string) (uint64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return uint64(stat.Size()), nil
}

// GetFile returns file
func GetFile(path string) (*os.File, error) {
	value, ok := openFileCache.Load(path)
	if ok && value != nil {
		f := value.(*os.File)
		return f, nil
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	openFileCache.Store(path, f)
	// Seek file to start
	if err != nil {
		return nil, err
	}
	return f, nil
}

// CloseFile closes file
func CloseFile(path string) error {
	value, ok := openFileCache.Load(path)
	if ok && value != nil {
		f := value.(*os.File)
		openFileCache.Delete(path)
		return f.Close()
	}
	openFileCache.Delete(path)
	return nil
}

// SetDownloadPath sets download
func SetDownloadPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return os.ErrInvalid
	}
	downloadPath = path
	return nil
}

// GetDownloadPath gets download path
func GetDownloadPath() string {
	return downloadPath
}
