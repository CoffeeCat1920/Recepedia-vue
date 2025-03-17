package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

  var userReq UserRequest 
  err := json.NewDecoder(r.Body).Decode(&userReq)

  if err != nil {
    http.Error(w, "Invalid JSON input", http.StatusBadRequest)
    fmt.Printf("Can't decode user cause, %s \n", err.Error())
    return
  }

  user := modals.NewUser(userReq.Name, userReq.Password)

  err = database.New().AddUser(user)
  if err != nil {
    http.Error(w, "Database Error", http.StatusInternalServerError)
    fmt.Printf("Can't add user to the db case, %s \n", err.Error())
    return
  }

  w.WriteHeader(http.StatusOK)

}
