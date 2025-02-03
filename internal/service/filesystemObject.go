package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/YasenMakioui/gostore/pkg/utils"
)

// Since we want to bind the payload to a struct but we have the fsobject to remain private
// We need this struct.
type FilesystemObjectPayload struct {
	Name string `json:"name"`
	Mode string `json:"mode"` // We put it as a string to use the octal format
	File bool   `json:"file"`
}

// Object maps directly to a file or directory in the host system
type FilesystemObject struct { // make this object private
	name string // This is the key. The name represents the full path (file included)
	mode fs.FileMode
	file bool
}

func NewFilesystemObject(name string, mode fs.FileMode, file bool) (*FilesystemObject, error) {

	// TODO control the fact that maybe we do not want to create the object since those objects already exist.
	filesystemObject := &FilesystemObject{
		name: name,
		mode: mode,
		file: file,
	}

	if err := utils.CheckPath(name); err != nil {
		if file {
			if err := createFile(filesystemObject); err != nil {
				return filesystemObject, err // Already formatted error from createFile
			}
		} else {
			if err := createDirectory(filesystemObject); err != nil {
				return filesystemObject, err
			}
		}
	}

	return filesystemObject, nil
}

func (o *FilesystemObject) GetName() string {
	return o.name
}

func (o *FilesystemObject) GetMode() fs.FileMode {
	return o.mode
}

func (o *FilesystemObject) GetFile() bool {
	return o.file
}

func (o *FilesystemObject) SetName(name string) error {

	if err := os.Rename(o.name, name); err != nil {
		return fmt.Errorf("Could not perform mv operation: %v", err)
	}

	o.name = name

	// Refresh mode
	info, err := os.Stat(name)
	if err == nil {
		o.mode = info.Mode()
	}

	return nil
}

func (o *FilesystemObject) SetMode(mode fs.FileMode) error {

	if err := os.Chmod(o.name, mode); err != nil {
		return fmt.Errorf("Could not change file mode %v", err)
	}

	o.mode = mode

	return nil
}

func (o *FilesystemObject) Delete() (string, error) {
	// Deletes the file or dir and returns the name

	if err := os.RemoveAll(o.name); err != nil {
		return o.name, fmt.Errorf("Could not perform rm operation: %v", err)
	}

	return o.name, nil
}

func (file *FilesystemObject) Read() ([]byte, error) {

	content, err := os.ReadFile(file.name)

	if err != nil {
		return content, fmt.Errorf("Could not perform read operation on file %v: %v", file.name, err)
	}

	return content, nil
}

func (dir *FilesystemObject) List() ([]FilesystemObject, error) {

	var directoryContents []FilesystemObject
	var fsObject FilesystemObject

	entries, err := os.ReadDir(dir.name)

	if err != nil {
		return directoryContents, fmt.Errorf("Could not perform list opera tion in directory: %v", err)
	}

	for _, entry := range entries {

		info, err := entry.Info()

		if err != nil {
			return directoryContents, fmt.Errorf("Something happened while checking entry information: %v", err)
		}

		name := filepath.Join(dir.name, info.Name())
		mode := info.Mode().Perm() // only the 3 bits
		isFile := !info.IsDir()

		fsObject = FilesystemObject{
			name: name,
			mode: mode,
			file: isFile,
		}

		directoryContents = append(directoryContents, fsObject)
	}

	return directoryContents, nil
}

func createFile(file *FilesystemObject) error { // private since it will be executed by the constructor once

	if _, err := os.Stat(file.name); err != nil {
		return fmt.Errorf("File %v already exists", file.name)
	}

	f, err := os.Create(file.name)

	if err != nil {
		return fmt.Errorf("Could not create file %v: %v", file.name, err)
	}
	defer f.Close()
	return nil
}

func createDirectory(dir *FilesystemObject) error {
	if err := os.Mkdir(dir.name, dir.mode); err != nil {
		return fmt.Errorf("Could not create dir %v: %v", dir.name, err)
	}

	return nil
}
