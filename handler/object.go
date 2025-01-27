package handler

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/YasenMakioui/gostore/utils"
	"github.com/gofiber/fiber/v2"
)

// Logging should be handled by the handler, do not add logs on object methods.

// Object maps directly to a file or directory in the host system
type Object struct {
	File       bool   `json:"file"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Permission int    `json:"permission"`
}

// Maybe we can create a builder pattern adding a constructor
// that hides the logic that checks if the file is or not a dir and so on.

func (o *Object) delete() (string, error) {
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

func (o *Object) list() ([]*Object, error) {
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

func (o *Object) read() (map[string]string, error) {
	// reads the file and returns a map with the key being the result and the value being its contents
	contentMap := make(map[string]string)

	content, err := os.ReadFile(o.Path)

	if err != nil {
		return contentMap, err
	}

	contentMap["res"] = string(content)

	return contentMap, err
}

func (o *Object) persist() (*Object, error) {
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

func GetObject(c *fiber.Ctx) error {
	// Define the slice that will contain the objecs
	var objectList []*Object

	// Construct the object
	pathSlice := strings.Split(c.Path(), "/")

	objectName := pathSlice[len(pathSlice)-1]

	localPath := utils.GetLocalPath(c.Path())

	isFile, err := utils.IsFile(localPath)

	if err != nil {
		// Error returned, something went wrong checking the file mode
		log.Printf("Error checking file mode: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not check the file mode",
		})
	}

	object := &Object{
		Name: objectName,
		Path: localPath,
		File: isFile,
		// No need to define permission here since we don't need it
	}

	// temp, for debug, pending delete
	c.Set("Access-Control-Allow-Origin", "*")

	// Check if file exists

	if _, err := utils.CheckPath(localPath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	if isFile {

		contents, err := object.read()

		if err != nil {
			log.Printf("Error reading file: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not read the specified file",
			})
		}

		return c.JSON(contents)
	}

	// If its a dir, return the files inside the dir

	objectList, err = object.list()

	if err != nil {
		log.Printf("Error reading dir: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not read the specified directory",
		})
	}

	return c.JSON(objectList)
}

func CreateObject(c *fiber.Ctx) error {
	// A payload specifying the type is needed. Default will be file but if file: false in payload then create a dir
	// On success, we return the object
	// The payload can have all data related to the object
	// The only data that is a must is the name, the path will be the one targeted by the request

	// we create the blank object just to bind it to the request
	object := new(Object)

	if err := c.BodyParser(object); err != nil { // Bind the request body to the Object struct
		log.Printf("Error parsing request body %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request paylad",
		})
	}

	// We check the path from the request and using the file name we add the path attr to the object

	localPath := utils.GetLocalPath(c.Path())

	targetPath := path.Join(localPath, object.Name)

	object.Path = targetPath

	if exists, _ := utils.CheckPath(targetPath); exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "File or directory already exists",
		})
	}

	// Check if its a file or a directory in the host since we can't create a file inside a file

	isFile, err := utils.IsFile(localPath)

	if err != nil {
		log.Printf("Error checking file mode: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking file mode",
		})
	}

	if isFile {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot create an object inside a file",
		})
	}

	// checks were done, persist the object

	object, err = object.persist()

	if err != nil {
		log.Printf("Error creating directory: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating file",
		})
	}

	return c.JSON(object)
}

func DeleteOjbect(c *fiber.Ctx) error {

	object := new(Object)

	pathSlice := strings.Split(c.Path(), "/")

	object.Name = pathSlice[len(pathSlice)-1]
	object.Path = utils.GetLocalPath(c.Path())

	if _, err := object.delete(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete de file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
		"name":    object.Name,
	})
}

func ModifyObject(c *fiber.Ctx) error {
	// move object to another dir
	// change object name
	// modify object permissions
	// Payload
	/*
		{
			object: "",
			name: "",
			permission: "",
			path: "",
		}

	*/

	//object := new(Object)

	// bind the object data

	//modifiedObject := NewObject()

	return c.SendString("modify")
}
