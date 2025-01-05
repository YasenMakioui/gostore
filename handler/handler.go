package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/YasenMakioui/gostore/config"
	"github.com/YasenMakioui/gostore/utils"
	"github.com/gofiber/fiber/v2"
)

func ListObject(c *fiber.Ctx) error {
	// Return a map or error
	// The map will be in this format
	// {
	//    name: true/false
	//}
	// key is the name and value is true if its a dir
	fileType := "f"
	var targetPath string
	var dir strings.Builder

	len, _ := dir.WriteString(config.Config("BASEDIR"))

	// Add a trailing slash if not on path

	if dir.String()[len-1:] != "/" {
		dir.WriteString("/")
	}

	//var files map[string]bool
	path := c.Path()

	files := make(map[string]string) //https://go.dev/blog/maps use make since var files map[string]bool creates a reference type and the value is nil

	targetPath = utils.GetFilePath(dir.String(), path, "/api/gostore/")

	_, err := utils.CheckPath(targetPath)

	if err != nil {
		return c.SendStatus(404)
	}

	//check if its a file

	mode, err := os.Stat(targetPath)

	if err != nil {
		return c.SendStatus(500)
	}

	if !mode.IsDir() {

		contents, err := readObject(targetPath)

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(contents)
	}

	entries, err := os.ReadDir(targetPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		info, err := e.Info()

		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {
			fileType = "d"
		}

		files[info.Name()] = fileType

		fileType = "f"
	}

	fmt.Println(targetPath)

	return c.JSON(files)
}

func CreateObject(c *fiber.Ctx) error {
	return c.SendString("Creating file")
}

func DeleteOjbect(c *fiber.Ctx) error {
	return c.SendString("Deleting object")
}

func GetObjectInfo() {}

func readObject(path string) (map[string]string, error) {

	contentMap := make(map[string]string)

	content, err := os.ReadFile(path)

	if err != nil {
		return contentMap, err
	}

	contentMap["res"] = string(content)

	return contentMap, err
}
