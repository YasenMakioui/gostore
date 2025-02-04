package handlers

import (
	//"encoding/json"
	//stdErrors "errors"
	//"log"
	//"path"
	//"strconv"
	//"strings"

	//"fmt"

	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/YasenMakioui/gostore/config"
	. "github.com/YasenMakioui/gostore/internal/service" // adding a dot to not use service.something
	"github.com/YasenMakioui/gostore/pkg/errors"
	"github.com/YasenMakioui/gostore/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func GetObject(requestContext *fiber.Ctx) error {

	key := utils.GetLocalPath(requestContext.Path())

	if err := utils.CheckPath(key); err != nil {
		return requestContext.Status(fiber.StatusNotFound).JSON(
			errors.FormatError("Not found"),
		)
	}

	mode, err := utils.GetFileMode(key)

	if err != nil {
		return requestContext.Status(fiber.StatusInternalServerError).JSON(
			errors.FormatError("Could not retrieve information from file/dir"),
		)
	}

	isFile, err := utils.IsFile(key)

	filesystemObject, err := NewFilesystemObject(key, mode, isFile)

	if isFile {
		if err != nil {
			return requestContext.Status(fiber.StatusInternalServerError).JSON(
				errors.FormatError("Could not map to file object"),
			)
		}

		contents, err := filesystemObject.Read()

		if err != nil {
			return requestContext.Status(fiber.StatusInternalServerError).JSON(
				errors.FormatError("Could not read the file"),
			)
		}

		return requestContext.Send(contents)
	}

	entries, err := filesystemObject.List()

	// Create returnable object.
	var responseData []FilesystemObjectPayload

	for _, value := range entries {

		// Extract only the app root. We don't want the absolute path given to the user.
		fmt.Println(value)
		name := strings.Replace(value.GetName(), config.Config("BASEDIR"), "", -1)

		responseObject := &FilesystemObjectPayload{
			Name: name,
			Mode: strconv.FormatInt(int64(value.GetMode()), 8),
			File: value.GetFile(),
		}

		responseData = append(
			responseData, *responseObject,
		)
	}

	if err != nil {
		return requestContext.Status(fiber.StatusInternalServerError).JSON(
			errors.FormatError(err.Error()),
		)
	}

	return requestContext.JSON(
		responseData,
	)

}

// Maps to POST method
func CreateObject(requestContext *fiber.Ctx) error {
	// Create the payload object
	payload := new(FilesystemObjectPayload)

	if err := requestContext.BodyParser(payload); err != nil {
		return requestContext.Status(fiber.StatusBadRequest).JSON(
			errors.FormatError("Invalid payload"),
		)
	}

	// Create local path. We get something like /dir/file and we need to merge it with basedir
	// /
	targetPath := strings.Replace(requestContext.Path(), config.Config("GOSTOREPATH"), "", -1)

	name := filepath.Join(config.Config("BASEDIR"), targetPath)

	payload.Name = name

	// File does not exist. We can create it.

	// We need to create the filemode
	// We get the mode as a string and we convert it to octal

	mode, err := strconv.ParseInt(payload.Mode, 8, 32)

	if err != nil {
		return requestContext.Status(fiber.StatusBadRequest).JSON(
			errors.FormatError("Invalid mode"),
		)
	}

	if _, err := NewFilesystemObject(payload.Name, os.FileMode(mode), payload.File); err != nil {
		return requestContext.Status(fiber.StatusBadRequest).JSON(
			errors.FormatError(fmt.Sprintf("Internal server error: %v", err)),
		)
	}

	return requestContext.JSON(
		payload,
	)

	// NOTES - filemode not aplyinh correctly
	// This post can be made anywhere.
	// Modify to use the path of the request to create a file
}

// A payload specifying the type is needed. Default will be file but if file: false in payload then create a dir
// On success, we return the object
// The payload can have all data related to the object
// The only data that is a must is the name, the path will be the one targeted by the request

// object, err := createObject(c.Body())

// if err != nil {
// 	c.Status(fiber.StatusInternalServerError).JSON(
// 		errors.FormatError(err.Error()),
// 	)
// }

// // We check the path from the request and using the file name we add the path attr to the object

// localPath := utils.GetLocalPath(c.Path())

// targetPath := path.Join(localPath, object.Name) // gettter

// object.Path = targetPath // setter

// if exists, _ := utils.CheckPath(targetPath); exists {
// 	return c.Status(fiber.StatusConflict).JSON(
// 		errors.FormatError("File or directory already exists"),
// 	)
// }

// // Check if its a file or a directory in the host since we can't create a file inside a file

// isFile, err := utils.IsFile(localPath)

// if err != nil {
// 	log.Printf("Error checking file mode: %v", err)

// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 		"error": "Error checking file mode",
// 	})
// }

// if isFile {
// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 		"error": "Cannot create an object inside a file",
// 	})
// }

// // checks were done, persist the object

// object, err = object.Persist()

// if err != nil {
// 	log.Printf("Error creating directory: %v", err)

// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 		"error": "Error creating file",
// 	})
// }

// return c.JSON(object)
//}

// Maps to DELETE method
// func DeleteOjbect(c *fiber.Ctx) error {

// 	object := new(Object)

// 	pathSlice := strings.Split(c.Path(), "/")

// 	object.Name = pathSlice[len(pathSlice)-1]
// 	object.Path = utils.GetLocalPath(c.Path())

// 	if _, err := object.Delete(); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Could not delete de file",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Deleted successfully",
// 		"name":    object.Name,
// 	})
// }

// Maps to PUT method
// func ModifyObject(c *fiber.Ctx) error {
// 	// move object to another dir
// 	// change object name
// 	// modify object permissions
// 	// Payload
// 	/*
// 				{
// 					target: "",
// 					attributes: {
// 						file: ""
// 						name: "",
// 		                permission: "",
// 					    path: "",
// 					}
// 				}
// 	*/
// 	// 1. read the target to gather information about the object
// 	// 2. use the target informationto create the object
// 	// 3. change the attributes
// 	// 4. persist
// 	// We update the attributes and then we persist it

// 	//object := new(Object)

// 	// bind the object data

// 	//modifiedObject := NewObject()

// 	return c.SendString("modify")
// }

/*
	Private functions used by the handlers
*/

// func createObject(body []byte) (*Object, error) {
// 	// Converts the given body to a new object of type Object
// 	var data map[string]string

// 	if err := json.Unmarshal(body, &data); err != nil { // passing the reference we don't need a return value
// 		return nil, stdErrors.New("Invalid JSON")
// 	}

// 	if err := validateObjectPayload(data); err != nil {
// 		return nil, stdErrors.New(`
// 		Invalid Format. Follow the format:
// 		{
// 			"file": false <---- optional
// 			"name": "filename" <---- required
// 			"path": "filepath" <---- required
// 			"mode": 0755 <---- optional
// 		 }
// 		`)
// 	}

// 	file, _ := strconv.ParseBool(data["file"])
// 	name := data["name"]
// 	objectPath := data["path"]
// 	mode, _ := strconv.Atoi(data["mode"])

// 	object, err := NewObject(
// 		file,
// 		name,
// 		objectPath,
// 		mode,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return object, err
// }
