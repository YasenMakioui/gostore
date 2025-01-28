package handler

import (
	"log"
	"path"
	"strings"

	. "github.com/YasenMakioui/gostore/service" // adding a dot to not use service.something
	"github.com/YasenMakioui/gostore/utils"
	"github.com/gofiber/fiber/v2"
)

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

		contents, err := object.Read()

		if err != nil {
			log.Printf("Error reading file: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not read the specified file",
			})
		}

		return c.JSON(contents)
	}

	// If its a dir, return the files inside the dir

	objectList, err = object.List()

	if err != nil {
		log.Printf("Error reading dir: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not read the specified directory",
		})
	}

	return c.JSON(objectList)
}

// Maps to POST method
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

	object, err = object.Persist()

	if err != nil {
		log.Printf("Error creating directory: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating file",
		})
	}

	return c.JSON(object)
}

// Maps to DELETE method
func DeleteOjbect(c *fiber.Ctx) error {

	object := new(Object)

	pathSlice := strings.Split(c.Path(), "/")

	object.Name = pathSlice[len(pathSlice)-1]
	object.Path = utils.GetLocalPath(c.Path())

	if _, err := object.Delete(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete de file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
		"name":    object.Name,
	})
}

// Maps to PUT method
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
