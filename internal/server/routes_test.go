package server

import (
	"big/internal/api"
	"big/internal/modals"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
)

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestAddRecipeHandler(t *testing.T) {
	// Making the mock info
	sessionInfo := struct {
		SessionId string `json:"sessionid"`
		OwnerId   string `json:"ownerid"`
		Exp       string `json:"exp"`
	}{
		SessionId: "testid-123",
		OwnerId:   "ownerid-123",
		Exp:       "test",
	}

	recipeInfo := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		Name:    "Test",
		Content: "Test",
	}
	recipe := modals.NewRecipe(recipeInfo.Name, sessionInfo.OwnerId)

	// Mocking the DB result
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Can't Initalize database")
	}

	// Getting Session
	rows := sqlmock.NewRows([]string{"sessionid", "ownerid", "exp"})
	rows.AddRow(sessionInfo.SessionId, sessionInfo.OwnerId, sessionInfo.Exp)

	mock.ExpectQuery(`SELECT \* FROM sessions WHERE sessionid = \$1;`).
		WithArgs(sessionInfo.SessionId).
		WillReturnRows(rows)

	// Inserting Recipe
	mock.ExpectExec(`INSERT INTO recipes\(uuid, name, ownerid, views\) VALUES\(\$1, \$2, \$3, -1\)`).
		WithArgs(sqlmock.AnyArg(), recipe.Name, recipe.OwnerId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	defer db.Close()

	// Initializing the api and mock server
	memFs := afero.NewMemMapFs()
	api := api.NewTestWith(db, memFs)
	server := httptest.NewServer(http.HandlerFunc(api.UploadRecipe))
	defer server.Close()

	body, err := json.Marshal(recipeInfo)
	if err != nil {
		t.Fatalf("error making response body. Err: %v", err)
	}

	// Defining the request
	req, err := http.NewRequest(http.MethodPost, server.URL+"/uploadrecipe", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("error making New Request. Err: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session-token",
		Value: sessionInfo.SessionId,
	})

	// Executing it on mock client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("error reading response. Err: %v", err)
	}

	// Asserting Failures
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("error making New Request. Err: %v", err)
	}
}

func TestGetRecipeHandler(t *testing.T) {
	// creating mock data
	recipeUUID := "test-uuid-123"
	recipeHTMLPath := fmt.Sprintf("upload/recipes/%s/recipe.html", recipeUUID)

	// mocking the file system
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(recipeHTMLPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	htmlContent := "<html><body>Mock Recipe Page</body></html>"
	err = afero.WriteFile(fs, recipeHTMLPath, []byte(htmlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// creating mock database
	db, mock, err := sqlmock.New()
	row := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"}).
		AddRow(`test-uuid-123`, `test`, `test-uuid-123`, 4)

	mock.ExpectQuery("SELECT \\* FROM recipes WHERE uuid = \\$1;").
		WithArgs(recipeUUID).
		WillReturnRows(row)

	mock.ExpectExec("UPDATE recipes SET views = views \\+ 1 WHERE uuid = \\$1").
		WithArgs(recipeUUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	defer db.Close()

	// Initializing the api and server
	api := api.NewTestWith(db, fs)
	r := mux.NewRouter()
	r.HandleFunc("/recipe/{id}", api.ServeRecipe)
	server := httptest.NewServer(r)
	defer server.Close()

	// Creating the request
	req, err := http.NewRequest(http.MethodGet, server.URL+"/recipe/"+recipeUUID, nil)

	// Executing it on mock client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("error reading response. Err: %v", err)
	}

	// Asserting Failures
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("error making New Request. Err: %v", err)
	}

}

func TestEdiRecipeHandler(t *testing.T) {
	// Making the mock data
	RecipeInfo := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		Name:    "new-name",
		Content: "test-content",
	}
	recipeUUID := "test-uuid-123"

	recipeHTMLPath := fmt.Sprintf("upload/recipes/%s/recipe.html", recipeUUID)
	recipeMdPath := fmt.Sprintf("upload/recipes/%s/recipe.md", recipeUUID)

	// mocking the file system
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(recipeHTMLPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	htmlContent := "<html><body>Mock Recipe Page</body></html>"
	err = afero.WriteFile(fs, recipeHTMLPath, []byte(htmlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	mdContent := "* Mock Recipe Page"
	err = afero.WriteFile(fs, recipeMdPath, []byte(mdContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Mocking the DB
	db, mock, err := sqlmock.New()

	mock.ExpectExec("UPDATE recipes SET name = \\$1 WHERE uuid = \\$2").
		WithArgs(RecipeInfo.Name, recipeUUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Making the mock api and server
	api := api.NewTestWith(db, fs)
	r := mux.NewRouter()
	r.HandleFunc("/recipe/{id}", api.EditRecipeHandler)
	server := httptest.NewServer(r)
	defer server.Close()

	// Creating the request
	body, err := json.Marshal(RecipeInfo)
	if err != nil {
		t.Fatalf("error making response body. Err: %v", err)
	}
	req, err := http.NewRequest(http.MethodPatch, server.URL+"/recipe/"+recipeUUID, bytes.NewBuffer(body))

	// Executing the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("error reading response. Err: %v", err)
	}

	// Asserting Failures
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("error making New Request. Err: %v", err)
	}
}

func TestDeleteRecipeHandler(t *testing.T) {
	// Making the mock data
	userUUID := "test-uuid-123"
	recipeUUID := "test-recipe-123"

	recipeHTMLPath := fmt.Sprintf("upload/recipes/%s/recipe.html", recipeUUID)
	recipeMdPath := fmt.Sprintf("upload/recipes/%s/recipe.md", recipeUUID)

	sessionInfo := struct {
		sessionId string
		ownerid   string
	}{
		sessionId: "test-session-123",
		ownerid:   userUUID,
	}

	recipeInfo := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		Name:    "new-name",
		Content: "test-content",
	}
	recipe := modals.NewRecipe(recipeInfo.Name, userUUID)

	// Making the mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Can't Init database")
	}

	row := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"}).
		AddRow(recipeUUID, `test`, userUUID, -1)

	mock.ExpectQuery("SELECT \\* FROM recipes WHERE uuid = \\$1;").
		WithArgs(recipeUUID).
		WillReturnRows(row)

	rows := sqlmock.NewRows([]string{"sessionid", "ownerid", "exp"})
	rows.AddRow(sessionInfo.sessionId, userUUID, "2025-06-04 22:52:37")

	mock.ExpectQuery("SELECT \\* FROM sessions WHERE ownerid = \\$1;").
		WithArgs(recipe.UUID).
		WillReturnRows(row)

	// Mocking the filesystem
	fs := afero.NewMemMapFs()

	err = fs.Mkdir(recipeHTMLPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	htmlContent := "<html><body>Mock Recipe Page</body></html>"
	err = afero.WriteFile(fs, recipeHTMLPath, []byte(htmlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	mdContent := "* Mock Recipe Page"
	err = afero.WriteFile(fs, recipeMdPath, []byte(mdContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Making the mock api and server
	api := api.NewTestWith(db, fs)
	r := mux.NewRouter()
	r.HandleFunc("/recipe/{id}", api.DeleteRecipeHandler)
	server := httptest.NewServer(r)
	defer server.Close()

	// Making the mock request
	req, err := http.NewRequest(http.MethodPatch, server.URL+"/recipe/"+recipeUUID, nil)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session-token",
		Value: sessionInfo.sessionId,
	})

	// Executing the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("error reading response. Err: %v", err)
	}

	// Asserting Failures
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("error making New Request. Err: %v", err)
	}

}
