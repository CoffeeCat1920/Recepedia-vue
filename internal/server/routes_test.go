package server

import (
	"big/internal/api"
	"big/internal/modals"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	db, mock, err := sqlmock.New()

	recipe := modals.NewRecipe("test")

	mock.ExpectExec(`INSERT INTO recipes\(uuid, name, ownerid, views\) VALUES\(\$1, \$2, \$3, -1\)`).
		WithArgs(recipe.UUID, recipe.Name, recipe.OwnerId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err != nil {
		t.Fatal("Can't Initalize database")
	}

	defer db.Close()

	api := api.NewTest(db)
	server := httptest.NewServer(http.HandlerFunc(api.UploadRecipe))
	defer server.Close()

	recipeInfo := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		Name:    "Test",
		Content: "Test",
	}

	body, err := json.Marshal(recipeInfo)
	if err != nil {
		t.Fatalf("error making response body. Err: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/uploadrecipe", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("error making New Request. Err: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session-token",
		Value: "gAX-QMAoZCKf72-Jq8oJdXaXXo92HwPyZobvK8NHnT4=",
	})

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("error reading response. Err: %v", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("error making New Request. Err: %v", err)
	}
}

func TestGetRecipeHandler(t *testing.T) {

}

func TestEdiRecipeHandler(t *testing.T) {

}

func TestDeleteRecipeHandler(t *testing.T) {

}
