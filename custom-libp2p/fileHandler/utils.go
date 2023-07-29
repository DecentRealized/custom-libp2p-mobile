package fileHandler

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
)

func GetSHA256Sum(file *os.File) (string, error) {
	h := sha256.New()
	_, err := io.Copy(h, file)
	if err != nil {
		return "", err
	}
	raw := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
