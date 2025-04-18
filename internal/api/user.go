package api

import (
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *api) SignUpHandler(w http.ResponseWriter, r *http.Request) {

	var userReq UserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)

	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		fmt.Printf("Can't decode user cause, %s \n", err.Error())
		return
	}

	user := modals.NewUser(userReq.Name, userReq.Password)

	db := api.db

	err = db.AddUser(user)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		fmt.Printf("Can't add user to the db case, %s \n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	// deleting all the users
	db := api.db
	err := db.DeleteRecipeByUser(userId)
	if err != nil {
		fmt.Printf("Can't find user of id %s \n", userId)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// delete sessions
	err = db.DeleteSessionByUser(userId)
	if err != nil {
		fmt.Printf("Can't delete session of  id %s \n", userId)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// deleting the user
	err = db.DeleteUserByUUid(userId)
	if err != nil {
		fmt.Printf("Can't find user of id %s \n", userId)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	db := api.db

	users, err := db.GetAllUsers()

	if err != nil {
		http.Error(w, "Can't get any recipes", http.StatusNotFound)
		fmt.Printf("Can't get all the recipes cause, %s", err.Error())
		return
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Can't Marshall json", http.StatusInternalServerError)
		fmt.Printf("Can't get all the recipes cause, %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
