package api

import (
	"big/internal/modals"
	"encoding/json"
	"net/http"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
  response := modals.Response{Message: "Hello World"}
  w.Header().Set("content-type", "application/json")
  json.NewEncoder(w).Encode(response)
}
