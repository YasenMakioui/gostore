package service

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/YasenMakioui/gostore/utils"
)

// Logging should be handled by the handler, do not add logs on object methods.

// Object maps directly to a file or directory in the host system
type Object struct { // make it private! object. The only way to get the object should be with the constructor
	File       bool   `json:"file"`
	Name       string `json:"name"`
	Path       string `json:"path"` // Maybe create a struct for Path since we can't export it https://stackoverflow.com/questions/50319404/has-json-tag-but-not-exported
	Permission int    `json:"permission"`
} // also the properties of the object can be private

// Constructor
func NewObject(file bool, name string, path string, permission int) (*Object, error) {
	// validate
	object := new(Object)

	object.File = file
	object.Name = name

	// check if the path is actually a valid path

	object.Path = path
	// TODO: Validate path. Carefull with using something like checking if exists.

	// Check if its four digits length
	if permission < 1000 || permission > 9999 {
		// provisional. TODO: Add a better validation for UNIX bytes
		permission = 0644 // Since its incorrect we add a default value
	}
	object.Permission = permission

	return object, nil
}

// Maybe we can create a builder pattern adding a constructor
// that hides the logic that checks if the file is or not a dir and so on.

func (o *Object) Delete() (string, error) {
	// method used to delete the object thus deleting the file or dir in the filesystem
	if _, err := utils.CheckPath(o.Path); err != nil {
		fmt.Println(err)
		return o.Name, err
	}

	if err := os.RemoveAll(o.Path); err != nil {
		fmt.Println(err)
		return o.Name, err
	}

	return o.Name, nil
}

func (o *Object) List() ([]*Object, error) {
	// Performs an ls to retrieve the information. If its a file, returns the empty list and an error
	var objectList []*Object

	entries, err := os.ReadDir(o.Path)

	if err != nil {
		return objectList, err
	}

	for _, entry := range entries {

		info, err := entry.Info()

		if err != nil {
			return objectList, err
		}

		object := &Object{
			Name:       info.Name(),
			File:       !info.IsDir(),
			Path:       o.Path,
			Permission: int(info.Mode().Perm()),
		}

		objectList = append(objectList, object)

	}

	return objectList, nil
}

func (o *Object) Read() (map[string]string, error) {
	// reads the file and returns a map with the key being the result and the value being its contents
	contentMap := make(map[string]string)

	content, err := os.ReadFile(o.Path)

	if err != nil {
		return contentMap, err
	}

	contentMap["res"] = string(content)

	return contentMap, err
}

func (o *Object) Persist() (*Object, error) {
	// Given the object, we persist it in the filesystem

	if o.File {
		if err := os.MkdirAll(o.Path, fs.FileMode(o.Permission)); err != nil {
			return o, err
		}
	}

	// If its a file, then create the specified file

	if err := os.WriteFile(o.Path, []byte(""), fs.FileMode(o.Permission)); err != nil {
		return o, err
	}

	return o, nil
}
