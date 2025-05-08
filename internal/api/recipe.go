package api

import (
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/spf13/afero"
)

func (api *api) UploadRecipe(w http.ResponseWriter, r *http.Request) {

	var recipeInfo RecipeInfo
	err := json.NewDecoder(r.Body).Decode(&recipeInfo)

	if err != nil {
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
	db := api.db

	session, err := db.GetSession(c.Value)
	if err != nil {
		fmt.Printf("Can't find Session %s, cause %s\n", c.Value, err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	recipe := modals.NewRecipe(recipeInfo.Name, session.OwnerId)
	err = db.AddRecipe(recipe)
	if err != nil {
		fmt.Printf("Can't Add the recipes cause, %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create directory for recipe
	directoryPath := "upload/recipes/" + recipe.UUID
	if api.fileExists(directoryPath) {
		fmt.Println("Recipe Directory Already Exists")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = api.fs.Mkdir(directoryPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Creating an html file
	err = api.mdFileGenreator(recipeInfo.Content, recipe.UUID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Creating an html file
	err = api.htmlFileGenerator(recipe.Name, recipeInfo.Content, recipe.UUID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) EditRecipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeUUID := vars["id"]

	var recipeInfo RecipeInfo

	err := json.NewDecoder(r.Body).Decode(&recipeInfo)
	if err != nil {
		fmt.Printf("Can't parse Json cause %s \n", err)
		http.Error(w, "Bad JSON provided", http.StatusBadRequest)
		return
	}

	err = api.deleteRecipeFiles(recipeUUID)
	if err != nil {
		fmt.Printf("Can't delete filese cause, %s \n", err)
		http.Error(w, "Can't delete recipe files", http.StatusInternalServerError)
		return
	}

	err = api.mdFileGenreator(recipeInfo.Content, recipeUUID)
	if err != nil {
		fmt.Printf("Can't generate md cause %s \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = api.htmlFileGenerator(recipeInfo.Name, recipeInfo.Content, recipeUUID)
	if err != nil {
		fmt.Printf("Can't generate html cause, %s \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	db := api.db

	err = db.EditRecipeName(recipeUUID, recipeInfo.Name)
	fmt.Printf("Provided recipe name, %s", recipeInfo.Name)
	if err != nil {
		fmt.Printf("Can't change recipe name cause, %s \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeUUID := vars["id"]

	directoryPath := "upload/recipes/" + recipeUUID

	if !(api.authSameUser(r)) {
		fmt.Println("Doesn't have the permission to edit recipe")
		http.Error(w, "Doesn't have the permission to edit recipe", http.StatusInternalServerError)
		return
	}

	db := api.db

	err := db.DeleteRecipe(recipeUUID)
	if err != nil {
		fmt.Printf("Can't find Recipe To Delete cause, %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !api.fileExists(directoryPath) {
		fmt.Printf("Can't find recipe directory  to delete in path, %s\n", directoryPath)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = os.RemoveAll(directoryPath)
	if err != nil {
		fmt.Printf("Can't find Delete recipe cause, %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) SearchRecipeHandler(w http.ResponseWriter, r *http.Request) {

	searchterm := r.URL.Query().Get("searchTerm")

	fmt.Printf("\n%s\n", searchterm)

	db := api.db

	recipes, err := db.SearchRecipe(searchterm)
	if err != nil {
		http.Error(w, "No recipes found", http.StatusInternalServerError)
		fmt.Printf("Can't get the searched recipes, %s", err.Error())
		return
	}

	jsonData, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, "Can't Marshall json", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *api) ServeRecipe(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	db := api.db

	recipe, err := db.GetRecipe(vars["id"])

	if err != nil {
		http.Error(w, "Can't Render the Recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", err.Error())
		return
	}

	directoryPath := "upload/recipes/" + recipe.UUID + "/recipe.html"

	if !api.fileExists(directoryPath) {
		http.Error(w, "Can't find this recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", recipe.UUID)
		return
	}

	// Wrap Afero FS as io/fs.FS
	fsWrapper := afero.NewIOFS(api.fs)

	tmpl, err := template.ParseFS(fsWrapper, directoryPath)

	err = tmpl.Execute(w, tmpl)
	if err != nil {
		http.Error(w, "Can't find this recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", recipe.UUID)
		return
	}

	err = db.IncreaseRecipeViews(recipe)
	if err != nil {
		http.Error(w, "Can't increase recipe's views", http.StatusInternalServerError)
		fmt.Printf("Can't increase recipe's views cause, %s", err.Error())
		return
	}

}

func (api *api) RecipeInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db := api.db

	recipe, err := db.GetRecipe(vars["id"])

	if err != nil {
		http.Error(w, "Can't Render the Recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", err.Error())
		return
	}

	filePath := "upload/recipes/" + recipe.UUID + "/recipe.md"

	if !api.fileExists(filePath) {
		http.Error(w, "Can't find this recipe", http.StatusInternalServerError)
		fmt.Printf("\nCan't get the most viewed recipes cause, %s\n", recipe.UUID)
		return
	}

	jsonData, err := json.Marshal(recipe.Name)
	if err != nil {
		http.Error(w, `{error: "Failed to fetch recipes"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *api) RecipeMdContent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	db := api.db

	recipe, err := db.GetRecipe(vars["id"])

	if err != nil {
		http.Error(w, "Can't Render the Recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", err.Error())
		return
	}

	fmt.Printf("Serving Recipe %s", recipe.Name)

	directoryPath := "upload/recipes/" + recipe.UUID + "/recipe.md"

	if !api.fileExists(directoryPath) {
		http.Error(w, "Can't find this recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", recipe.UUID)
		return
	}

	content, err := os.ReadFile(directoryPath)
	mdContent := string(content)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mdContent))
}

func (api *api) GetRecipeByUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	var recipes []modals.Recipe

	db := api.db

	recipes, err := db.GetRecipesByUser(name)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch recipes"}`, http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *api) MostViewedRecipesHandler(w http.ResponseWriter, r *http.Request) {
	recipes, err := api.db.MostViewedRecipes()
	if err != nil {
		http.Error(w, "Can't get any recipes", http.StatusNotFound)
	}

	jsonData, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, "Can't Marshall json", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *api) GetAllRecipesHandler(w http.ResponseWriter, r *http.Request) {
	recipes, err := api.db.GetAllRecipes()
	if err != nil {
		http.Error(w, "Can't get any recipes", http.StatusNotFound)
		fmt.Printf("Can't get all the recipes cause, %s", err.Error())
		return
	}

	jsonData, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, "Can't Marshall json", http.StatusInternalServerError)
		fmt.Printf("Can't get all the recipes cause, %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
