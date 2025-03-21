package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func mdFileGenreator(content string, uuid string) (error) {
  
	directoryPath := "upload/recipes/" + uuid 
	if !fileExists(directoryPath) {
		fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
    return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath) 
	}

  mdPath := directoryPath + "/recipe.md"
	mdFile, err := os.Create(mdPath)
	if err != nil {
		fmt.Printf("Error creating markdown file: %s\n", err.Error())
		return err
	}
	defer mdFile.Close()

	_, err = mdFile.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to markdown file: %s\n", err.Error())
		return err
	}

  return nil
} 

func htmlFileGenerator(recipeName string, content string, uuid string) (error) {
	directoryPath := "upload/recipes/" + uuid 
	if !fileExists(directoryPath) {
    fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
    return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath) 
	}

  htmlFile := directoryPath + "/recipe.html"

  htmlContent := markdown.ToHTML([]byte(content), nil, nil)
  
  templateContent := fmt.Sprintf(
  `<template>  
   <title> %s </title>
  ` + string(htmlContent) + `
   </template>`, recipeName)

  if !fileExists(htmlFile) {
    err := os.WriteFile(htmlFile, []byte(templateContent), fs.ModePerm)
    if err != nil {
      panic(err)
    }
  }
  
  return nil
}

func UploadRecipe(w http.ResponseWriter, r *http.Request) {
  
  var recipeInfo RecipeInfo 
  err := json.NewDecoder(r.Body).Decode(&recipeInfo) 

  if (err != nil) {
    fmt.Printf("Can't get recipe Info cause, %s", err)
    http.Error(w, "Invalid Request", http.StatusBadRequest)
  }

	// Get session token
	c, err := r.Cookie("session-token")
	if err != nil {
		fmt.Println("Can't find Cookie")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get user session
	session, err := database.New().GetSession(c.Value)
	if err != nil {
		fmt.Printf("Can't find Session %s, cause %s\n", c.Value, err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

  recipe := modals.NewRecipe(recipeInfo.Name, session.OwnerId)
  err = database.New().AddRecipe(recipe)
	if err != nil {
		fmt.Printf("Can't Add the recipes cause, %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create directory for recipe
	directoryPath := "upload/recipes/" + recipe.UUID
	if fileExists(directoryPath) {
		fmt.Println("Recipe Directory Already Exists")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = os.Mkdir(directoryPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

  // Creating an html file
  err = mdFileGenreator(recipeInfo.Content, recipe.UUID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

  // Creating an html file
  err = htmlFileGenerator(recipe.Name, recipeInfo.Content, recipe.UUID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

  w.WriteHeader(http.StatusOK)

}
