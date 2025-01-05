package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetFilePath(gostorePath, uriPath string, prefix string) string {
	// Removes the prefix of the given path and adds it to the gostore path
	// path /api/file/images/image.png with prefix /api/file/ will return images/image.png

	targetPath := strings.TrimPrefix(uriPath, prefix) // /api/gostore/ /api/gostore/images/image.png -> images/image.png

	path := filepath.Join(gostorePath, targetPath)

	fmt.Println(path)

	return path
}

func CheckPath(path string) (bool, error) {

	_, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return true, nil
}
