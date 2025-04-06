package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func mdFileGenreator(content string, uuid string) error {

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

func htmlFileGenerator(recipeName string, content string, uuid string) error {
	directoryPath := "upload/recipes/" + uuid
	if !fileExists(directoryPath) {
		fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath)
	}

	htmlFile := directoryPath + "/recipe.html"

	htmlContent := markdown.ToHTML([]byte(content), nil, nil)

	templateContent := fmt.Sprintf(
		`<title> %s </title>
  `+string(htmlContent), recipeName)

	if !fileExists(htmlFile) {
		err := os.WriteFile(htmlFile, []byte(templateContent), fs.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func deleteRecipeFiles(uuid string) error {
	directoryPath := "upload/recipes/" + uuid
	if !fileExists(directoryPath) {
		fmt.Printf("Recipe Directory Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe Directory Doesn't Exists, %s", directoryPath)
	}

	mdPath := directoryPath + "/recipe.md"
	if !fileExists(mdPath) {
		fmt.Printf("Recipe mdFile Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe mdFile Doesn't Exists, %s", directoryPath)
	}
	err := os.Remove(mdPath)
	if err != nil {
		return err
	}

	htmlFile := directoryPath + "/recipe.html"
	if !fileExists(htmlFile) {
		fmt.Printf("Recipe mdFile Doesn't Exists, %s \n", directoryPath)
		return fmt.Errorf("Recipe mdFile Doesn't Exists, %s", directoryPath)
	}
	err = os.Remove(htmlFile)
	if err != nil {
		return err
	}

	return nil
}

func UploadRecipe(w http.ResponseWriter, r *http.Request) {

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

func EditRecipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeUUID := vars["id"]

	fmt.Printf("%s", recipeUUID)
	var recipeInfo RecipeInfo

	err := json.NewDecoder(r.Body).Decode(&recipeInfo)
	if err != nil {
		fmt.Printf("Can't parse Json cause %s \n", err)
		http.Error(w, "Bad JSON provided", http.StatusBadRequest)
		return
	}

	err = deleteRecipeFiles(recipeUUID)
	if err != nil {
		fmt.Printf("Can't delete filese cause, %s \n", err)
		http.Error(w, "Can't delete recipe files", http.StatusInternalServerError)
		return
	}

	err = mdFileGenreator(recipeInfo.Content, recipeUUID)
	if err != nil {
		fmt.Printf("Can't generate md cause %s \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = htmlFileGenerator(recipeInfo.Name, recipeInfo.Content, recipeUUID)
	if err != nil {
		fmt.Printf("Can't generate html cause, %s \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = database.New().EditRecipeName(recipeUUID, recipeInfo.Name)
	fmt.Printf("Provided recipe name, %s", recipeInfo.Name)
	if err != nil {
		fmt.Printf("Can't change recipe name cause, %s \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeUUID := vars["id"]

	directoryPath := "upload/recipes/" + recipeUUID

	if !(authSameUser(r)) {
		fmt.Println("Doesn't have the permission to edit recipe")
		http.Error(w, "Doesn't have the permission to edit recipe", http.StatusInternalServerError)
		return
	}

	err := database.New().DeleteRecipe(recipeUUID)
	if err != nil {
		fmt.Printf("Can't find Recipe To Delete cause, %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !fileExists(directoryPath) {
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

func SearchRecipeHandler(w http.ResponseWriter, r *http.Request) {

	searchterm := r.URL.Query().Get("searchTerm")

	fmt.Printf("\n%s\n", searchterm)

	recipes, err := database.New().SearchRecipe(searchterm)
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

func ServeRecipe(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	db := database.New()

	recipe, err := db.GetRecipe(vars["id"])

	if err != nil {
		http.Error(w, "Can't Render the Recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", err.Error())
		return
	}

	fmt.Printf("Serving Recipe %s", recipe.Name)

	directoryPath := "upload/recipes/" + recipe.UUID + "/recipe.html"

	if !fileExists(directoryPath) {
		http.Error(w, "Can't find this recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", recipe.UUID)
		return
	}

	tmpl, err := template.ParseFiles(directoryPath)

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

func RecipeInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db := database.New()

	recipe, err := db.GetRecipe(vars["id"])

	if err != nil {
		http.Error(w, "Can't Render the Recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", err.Error())
		return
	}

	fmt.Printf("Editing Recipe %s", recipe.Name)

	filePath := "upload/recipes/" + recipe.UUID + "/recipe.md"

	if !fileExists(filePath) {
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

func RecipeMdContent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	db := database.New()

	recipe, err := db.GetRecipe(vars["id"])

	if err != nil {
		http.Error(w, "Can't Render the Recipe", http.StatusInternalServerError)
		fmt.Printf("Can't get the most viewed recipes cause, %s", err.Error())
		return
	}

	fmt.Printf("Serving Recipe %s", recipe.Name)

	directoryPath := "upload/recipes/" + recipe.UUID + "/recipe.md"

	if !fileExists(directoryPath) {
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

func GetRecipeByUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	var recipes []modals.Recipe

	db := database.New()

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

func MostViewedRecipesHandler(w http.ResponseWriter, r *http.Request) {
	recipes, err := database.New().MostViewedRecipes()
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

func GetAllRecipesHandler(w http.ResponseWriter, r *http.Request) {
	recipes, err := database.New().GetAllRecipes()
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
