package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (api *api) auth(r *http.Request) (*modals.Session, error) {
	cookie, err := r.Cookie("session-token")
	if err != nil {
		return nil, err
	}

	sessionId := cookie.Value
	session, err := database.New().GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (api *api) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var userReq UserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)

	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		fmt.Printf("Can't login cause, %s \n", err.Error())
	}

	fmt.Printf("\nSession Cookie, name: %s\n", userReq.Name)

	db := api.db
	user, err := db.GetUserByName(userReq.Name)

	if err != nil {
		http.Error(w, "Incorrect Password", http.StatusInternalServerError)
		fmt.Printf("\nCan't find user cause, %s\n", err)
		return
	}
	if !user.CheckPassword(userReq.Password) {
		http.Error(w, "Incorrect Password", http.StatusInternalServerError)
		fmt.Printf("Password Incorrect\n")
		return
	}

	session := api.createCookie(w, user.UUID)
	err = db.AddSession(session)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		fmt.Printf("Can't add the session to db cause, %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := api.auth(r)
	if err != nil {
		fmt.Printf("\nCan't Autherize session with id %s in db", session.SessionId)
		http.Error(w, "Can't Logout", http.StatusInternalServerError)
		return
	}

	err = api.db.DeleteSession(session.SessionId)
	if err != nil {
		fmt.Printf("\nCan't delete Session %s in db cause,\n%s", session.SessionId, err.Error())
		http.Error(w, "Can't Logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire immediately
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func (api *api) LoginInfoHandler(w http.ResponseWriter, r *http.Request) {

	session, err := api.auth(r)

	loginInfo := &LoginInfo{
		User:     "",
		UUID:     "",
		LoggedIn: false,
	}

	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(loginInfo)

		fmt.Printf("Can't get session in db cause, %s\n", err.Error())
		return
	}

	user, err := api.db.GetUserByUUid(session.OwnerId)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(loginInfo)

		fmt.Printf("Can't get user in db cause, %s\n", err.Error())
		return
	}

	loginInfo = &LoginInfo{
		User:     user.Name,
		UUID:     user.UUID,
		LoggedIn: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginInfo)
}

func (api *api) LoginRecipeInfoHandler(w http.ResponseWriter, r *http.Request) {

	session, err := api.auth(r)

	var recipes []modals.Recipe

	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(recipes)

		fmt.Printf("Can't get session in db cause, %s\n", err.Error())
		return
	}

	db := api.db

	user, err := db.GetUserByUUid(session.OwnerId)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(recipes)

		fmt.Printf("Can't get user in db cause, %s\n", err.Error())
		return
	}

	recipes, err = db.GetRecipesByUser(user.Name)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(recipes)

		fmt.Printf("Can't get recipes in db cause, %s\n", err.Error())
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

func (api *api) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := api.auth(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (api *api) authSameUser(r *http.Request) bool {
	// Getting the recipe
	vars := mux.Vars(r)
	recipeUUID := vars["id"]
	recipe, err := database.New().GetRecipe(recipeUUID)
	if err != nil {
		fmt.Println("Can't find Recipe To Edit")
		return false
	}

	// Get session token
	c, err := r.Cookie("session-token")
	if err != nil {
		fmt.Println("Can't find Cookie")
		return false
	}

	// Get user session
	db := api.db
	session, err := db.GetSession(c.Value)
	if err != nil {
		fmt.Printf("Can't find Session %s\n", c.Value)
		return false
	}

	// Verify the permission
	if !(recipe.OwnerId == session.OwnerId) {
		fmt.Printf("You don't have ther permission to edit the recipe\n")
		return false
	}

	return true
}

func (api *api) createAdminCookie(w http.ResponseWriter) *modals.AdminSession {
	session := modals.NewAdminSession()
	exp, err := session.GetExpTime()
	if err != nil {
		panic(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin-session-token",
		Value:    session.SessionId,
		Expires:  exp,
		Path:     "/", // Add this
		Domain:   "",
		HttpOnly: true,                 // Add this
		SameSite: http.SameSiteLaxMode, // Add this
	})
	return session
}
