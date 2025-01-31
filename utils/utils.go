package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/YasenMakioui/gostore/config"
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

func AddTrailingSlash(str string) string {

	len := len(str)

	if string(str[len-1]) == "/" {
		return str
	}

	return str + "/"
}

func GetLocalPath(contextPath string) string {
	baseDir := config.Config("BASEDIR")

	gostorePath := AddTrailingSlash(contextPath)
	gostorePath, _ = strings.CutPrefix(gostorePath, config.Config("GOSTOREPATH"))

	return path.Join(baseDir, gostorePath)
}

func IsFile(path string) (bool, error) {
	mode, err := os.Stat(path)

	if mode.IsDir() {
		return false, err
	}

	return true, err
}

// func ValidateFileMode(fileMode int) (bool) {
// 	// Following https://www.gnu.org/software/coreutils/manual/html_node/Mode-Structure.html
// 	// and https://www.gnu.org/software/coreutils/manual/html_node/Numeric-Modes.html

// 	// Use strings since its easier for the use of slices and loops

// 	validBits := []string{"4","2","1"}

// 	fileModeSlice := strings.Split(strconv.Itoa(fileMode),"")

// 	for _, v := range fileModeSlice {
// 		if slices.Contains(validBits, v) {
// 			// Incomplete.
// 		}
// 	}

// 	return true
// }
