package api

import (
	"errors"
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/spf13/afero"
)

func (api *api) fileExists(path string) bool {
	AppFs := api.fs
	_, err := AppFs.Stat(path)
	return (err == nil)
}

func (api *api) mdFileGenreator(content string, uuid string) error {

	AppFs := api.fs

	directoryPath := "upload/recipes/" + uuid
	if !api.fileExists(directoryPath) {
		fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath)
	}

	mdPath := directoryPath + "/recipe.md"
	mdFile, err := AppFs.Create(mdPath)
	if err != nil {
		fmt.Printf("Error creating markdown file: %s\n", err.Error())
		return err
	}
	defer mdFile.Close()

	err = afero.WriteFile(api.fs, mdPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing to markdown file: %s\n", err.Error())
		return err
	}

	return nil
}

func (api *api) htmlFileGenerator(recipeName string, content string, uuid string) error {

	AppFs := api.fs

	directoryPath := "upload/recipes/" + uuid
	if !api.fileExists(directoryPath) {
		fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath)
	}

	htmlFile := directoryPath + "/recipe.html"

	if api.fileExists(htmlFile) {
		return errors.New("Html file already Exists")
	}

	htmlContent := markdown.ToHTML([]byte(content), nil, nil)

	templateContent := fmt.Sprintf(
		`<title> %s </title>
  `+string(htmlContent), recipeName)

	err := afero.WriteFile(AppFs, htmlFile, []byte(templateContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (api *api) deleteRecipeFiles(uuid string) error {
	AppFs := api.fs

	directoryPath := "upload/recipes/" + uuid
	if !api.fileExists(directoryPath) {
		fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath)
	}

	mdPath := directoryPath + "/recipe.md"
	if !api.fileExists(mdPath) {
		fmt.Printf("Recipe mdFile Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe mdFile Doesn't Exists, %s", directoryPath)
	}
	err := AppFs.Remove(mdPath)
	if err != nil {
		return err
	}

	htmlFile := directoryPath + "/recipe.html"
	if !api.fileExists(htmlFile) {
		fmt.Printf("Recipe HTML Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe HTML Doesn't Exists, %s", directoryPath)
	}
	err = AppFs.Remove(htmlFile)
	if err != nil {
		return err
	}

	return nil
}
