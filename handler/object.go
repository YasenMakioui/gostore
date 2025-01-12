package handler

import (
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
	// A payload specifying the type is needed. Default will be file but if file: false in payload then create a dir
	o := new(Object)

	if err := c.BodyParser(o); err != nil { // Bind the request body to the Object struct
		return c.SendStatus(500)
	}

	// Create the file on the specified path

	baseDir := config.Config("BASEDIR")

	contextPath := c.Path() // This contains the url without /api/v1/gostore/store

	gostorePath := utils.AddTrailingSlash(contextPath)

	gostorePath, _ = strings.CutPrefix(gostorePath, "/api/v1/gostore/store")

	localPath := path.Join(baseDir, gostorePath)

	exists, _ := utils.CheckPath(path.Join(localPath, o.Name))

	if exists {
		return c.SendString("Object already exists")
	}

	// Check if its a file or a directory

	mode, err := os.Stat(localPath)

	if err != nil {
		return c.SendStatus(500)
	}

	if !mode.IsDir() {
		return c.SendString("Can't create an object into a file type")
	}

	// If its a dir, then create the specified file

	if !o.File {
		err := os.MkdirAll(path.Join(localPath, o.Name), 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
			return c.SendStatus(500)
		}

		return c.SendString("Created dir successfully")
	}

	err = os.WriteFile(path.Join(localPath, o.Name), []byte(""), 0666)

	if err != nil {
		log.Fatal(err)
		c.SendStatus(500)
	}

	return c.JSON(o)
}

func DeleteOjbect(c *fiber.Ctx) error {

	baseDir := config.Config("BASEDIR")

	contextPath := c.Path() // This contains the url without /api/v1/gostore/store

	gostorePath := utils.AddTrailingSlash(contextPath)

	gostorePath, _ = strings.CutPrefix(gostorePath, "/api/v1/gostore/store")

	localPath := path.Join(baseDir, gostorePath)

	_, err := utils.CheckPath(path.Join(localPath))

	if err != nil {
		// File does not exist, send a 404
		return c.SendStatus(404)
	}

	err = os.RemoveAll(localPath)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendString("Object deleted")
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
