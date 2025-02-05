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

func GetFilesystemObject(requestContext *fiber.Ctx) error {

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
func CreateFilesystemObject(requestContext *fiber.Ctx) error {
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

	// field name contains the full system path and should be the root of the app
	return requestContext.JSON(
		payload,
	)

	// NOTES - directory permissions are not set correctly.
}

func DeleteFilesystemOjbect(requestContext *fiber.Ctx) error {
	// We need the localpath
	targetPath := strings.Replace(requestContext.Path(), config.Config("GOSTOREPATH"), "", -1)

	name := filepath.Join(config.Config("BASEDIR"), targetPath)

	fsObject := NewFilesystemObjectMapper(name)

	if _, err := fsObject.Delete(); err != nil {
		return requestContext.Status(fiber.StatusInternalServerError).JSON(
			errors.FormatError(err.Error()),
		)
	}

	// This is returning the full path of the system. Should return the root of the program.
	return requestContext.JSON(fiber.Map{
		"name": name,
	})
}

func ModifyFilesystemObject(requestContext *fiber.Ctx) error {

	// PUT to /dir/testdir/file and payload with modifications
	// Ignore other data in the payload that is not mode or name

	payload := struct {
		Name string `json:"name"`
		Mode string `json:"mode"`
	}{}

	if err := requestContext.BodyParser(&payload); err != nil {
		fmt.Println(err)
		return requestContext.Status(fiber.StatusBadRequest).JSON(
			errors.FormatError("Invalid payload. Add fields name and mode"),
		)
	}

	octalBits, err := strconv.ParseInt(payload.Mode, 8, 32)

	if err != nil {
		requestContext.JSON(
			errors.FormatError("Invalid file mode"),
		)
	}

	// Careful here. Control if the user is using the correct path in the app context and not the system context.
	//mode := os.FileMode(octalBits)
	//name := payload.Name
	// convert payload.mode into a valid fs.FileMode
	fmt.Println(os.FileMode(octalBits))

	return requestContext.JSON(payload)
}
