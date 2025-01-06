package handler

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/YasenMakioui/gostore/config"
	"github.com/YasenMakioui/gostore/utils"
	"github.com/gofiber/fiber/v2"
)

type Object struct {
	File bool   `json:"file"`
	Name string `json:"name"`
}

func GetObject(c *fiber.Ctx) error {

	baseDir := config.Config("BASEDIR")

	contextPath := c.Path() // This contains the url without /api/v1/gostore/store

	gostorePath := utils.AddTrailingSlash(contextPath)

	gostorePath, _ = strings.CutPrefix(gostorePath, "/api/v1/gostore/store")

	localPath := path.Join(baseDir, gostorePath)

	// Create the file list

	var objects []*Object

	// Check if file exists

	_, err := utils.CheckPath(localPath)

	if err != nil {
		// File does not exist, send a 404
		return c.SendStatus(404)
	}

	// Check if its a file or a directory

	mode, err := os.Stat(localPath)

	if err != nil {
		return c.SendStatus(500)
	}

	// If its a file, read the file

	if !mode.IsDir() {

		contents, err := readObject(localPath)

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(contents)
		//return c.JSON(contents)
	}

	// If its a dir, return the files inside the dir

	entries, err := os.ReadDir(localPath)

	if err != nil {
		return c.SendStatus(500)
	}

	for _, e := range entries {

		o := new(Object) // https://www.freecodecamp.org/news/new-vs-make-functions-in-go/

		info, err := e.Info()

		if err != nil {
			log.Fatal(err)
			c.SendStatus(500)
		}

		o.Name = info.Name()
		o.File = true

		if info.IsDir() {
			o.File = false
		}

		objects = append(objects, o)
	}

	return c.JSON(objects) //c.JSON(files)
}

func CreateObject(c *fiber.Ctx) error {

	o := new(Object)

	if err := c.BodyParser(o); err != nil {
		return c.SendStatus(500)
	}

	return c.SendString("Creating file")
}

func DeleteOjbect(c *fiber.Ctx) error {

	o := new(Object)

	if err := c.BodyParser(o); err != nil {
		return c.SendStatus(500)
	}

	fmt.Println(o.File)

	return c.SendString("Deleting objectt")
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
