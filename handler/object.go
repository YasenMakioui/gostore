package handler

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/YasenMakioui/gostore/utils"
	"github.com/gofiber/fiber/v2"
)

type Object struct {
	File bool   `json:"file"`
	Name string `json:"name"`
}

func GetObject(c *fiber.Ctx) error {
	// Create the file list

	var objectList []*Object

	c.Set("Access-Control-Allow-Origin", "*")

	localPath := utils.GetLocalPath(c.Path())

	// Check if file exists

	if _, err := utils.CheckPath(localPath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}
	// Check if its a file or a directory

	isFile, err := utils.IsFile(localPath)

	if err != nil {
		// Error returned, something went wrong checking the file mode
		log.Printf("Error checking file mode: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not check the file mode",
		})
	}

	// If its a file, read the file

	if isFile {

		contents, err := utils.ReadObject(localPath)

		if err != nil {
			log.Printf("Error reading file: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not read the specified file",
			})
		}

		return c.JSON(contents)
	}

	// If its a dir, return the files inside the dir

	entries, err := os.ReadDir(localPath)

	if err != nil {
		log.Printf("Error getting file info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not read the specified directory",
		})
	}

	for _, entry := range entries {

		//o := new(Object) // https://www.freecodecamp.org/news/new-vs-make-functions-in-go/

		info, err := entry.Info()

		if err != nil {
			// Could not get the path info
			log.Printf("Error getting file info: %v", err)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing directory entry",
			})
		}

		object := &Object{
			Name: info.Name(),
			File: !info.IsDir(),
		}

		objectList = append(objectList, object)

	}

	return c.JSON(objectList)
}

func CreateObject(c *fiber.Ctx) error {
	// A payload specifying the type is needed. Default will be file but if file: false in payload then create a dir
	object := new(Object)

	if err := c.BodyParser(object); err != nil { // Bind the request body to the Object struct
		log.Printf("Error parsing request body %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request paylad",
		})
	}

	// Create the file on the specified path

	localPath := utils.GetLocalPath(c.Path())

	targetPath := path.Join(localPath, object.Name)

	if exists, _ := utils.CheckPath(targetPath); exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "File or directory already exists",
		})
	}

	// Check if its a file or a directory in the local storage

	isFile, err := utils.IsFile(localPath)

	if err != nil {
		// Error returned, something went wrong checking the file mode
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

	// If its a dir, then create the specified dir

	if !object.File {

		if err := os.MkdirAll(targetPath, 0755); err != nil {
			log.Printf("Error creating directory: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating file",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Directory created successfully",
			"name":    object.Name,
		})
	}

	// If its a file, then create the specified file

	if err := os.WriteFile(targetPath, []byte(""), 0644); err != nil {
		log.Printf("Error creating file: %v", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "File created successfully",
		"name":    object.Name,
	})
}

func DeleteOjbect(c *fiber.Ctx) error {

	pathSlice := strings.Split(c.Path(), "/")

	objectName := pathSlice[len(pathSlice)-1]

	localPath := utils.GetLocalPath(c.Path())

	if _, err := utils.CheckPath(localPath); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "The specified path does not exist",
		})
	}

	if err := os.RemoveAll(localPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete de file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
		"name":    objectName,
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
			move: "",
		}

	*/

	//object := new(Object)

	return c.SendString("modify")
}
