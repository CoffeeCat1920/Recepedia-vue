package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Called")
	vars := mux.Vars(r)
	userId := vars["id"]

	// deleting all the users
	err := database.New().DeleteRecipeByUser(userId)
	if err != nil {
		fmt.Printf("Can't find user of id %s \n", userId)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// delete sessions
	err = database.New().DeleteSessionByUser(userId)
	if err != nil {
		fmt.Printf("Can't delete session of  id %s \n", userId)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// deleting the user
	err = database.New().DeleteUserByUUid(userId)
	if err != nil {
		fmt.Printf("Can't find user of id %s \n", userId)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
