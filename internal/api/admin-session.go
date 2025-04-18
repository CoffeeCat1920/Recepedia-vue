package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"
)

// Helper function to get the admin session from the request
func (api *api) authAdmin(r *http.Request) (*modals.AdminSession, error) {
	cookie, err := r.Cookie("admin-session-token")
	if err != nil {
		fmt.Printf("\nCan't find cookie\n")
		return nil, err
	}

	sessionid := cookie.Value
	session, err := database.New().GetAdminSession(sessionid)

	if err != nil {
		fmt.Printf("\nCan't find session %s in db case, %s\n", sessionid, err.Error())
		return nil, err
	}

	fmt.Printf("\nFound session %s in db\n", sessionid)
	return session, nil
}

// Function to create middleware for admin authentication
func (api *api) AdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := api.authAdmin(r)
		if err != nil {
			http.Redirect(w, r, "/view/admin-login", 302)
			fmt.Printf("Can't log the admin in cause, %s", err.Error())
			return
		}

		next(w, r)
	}
}

// Function to verify admin credentials
func (api *api) VerifyAdmin(w http.ResponseWriter, r *http.Request) {
	var adminReq AdminRequest
	err := json.NewDecoder(r.Body).Decode(&adminReq)

	ad := modals.NewAdmin()
	if !(ad.CheckPassword(adminReq.Password)) {
		http.Error(w, "Wrong Password", http.StatusUnauthorized)
		return
	}

	db := api.db

	session := api.createAdminCookie(w)
	err = db.AddAdminSession(session)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSession created and cookie set. Session ID: %s\n", session.SessionId)
	http.Redirect(w, r, "/view/admin-dashboard", 302)
}

// Function to handle admin login info requests
func (api *api) AdminLoginInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := api.authAdmin(r)

	loginInfo := &AdminLoginInfo{
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

	loginInfo = &AdminLoginInfo{
		UUID:     session.SessionId,
		LoggedIn: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginInfo)
}
