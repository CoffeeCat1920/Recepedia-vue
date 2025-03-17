package api

import (
	"big/internal/database"
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"
)

func auth(r *http.Request) (*modals.Session, error) {
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {

  var userReq UserRequest
  err := json.NewDecoder(r.Body).Decode(&userReq)

  if err != nil {
    http.Error(w, "Invalid JSON input", http.StatusBadRequest)
    fmt.Printf("Can't login cause, %s \n", err.Error())
  }

  db := database.New()
  user, err := db.GetUserByName(userReq.Name)
  if err != nil {
    fmt.Printf("\nCan't find user\n")
    http.Redirect(w, r, "/view/login", 302)
    return
  }
  if !user.CheckPassword(userReq.Password) {
    http.Error(w, "Incorrect Password", http.StatusInternalServerError)
    fmt.Printf("Password Incorrect\n")
    return
  } 

  session := createCookie(w, user.UUID)
  err = db.AddSession(session)
  if err != nil {
    http.Error(w, "Database Error", http.StatusInternalServerError)
    fmt.Printf("Can't add the session to db cause, %s\n", err.Error())
    return
  }
  
  w.WriteHeader(http.StatusOK)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
  
  session, err := auth(r)
  if err != nil {
    fmt.Printf("\nCan't Autherize session with id %s in db", session.SessionId)
    http.Error(w, "Can't Logout", http.StatusInternalServerError)
    return 
  }

  err = database.New().DeleteSession(session.SessionId)
  if err != nil {
    fmt.Printf("\nCan't delete Session %s in db cause,\n%s", session.SessionId, err.Error())
    http.Error(w, "Can't Logout", http.StatusInternalServerError)
    return 
  }

}


func LoginInfoHandler(w http.ResponseWriter, r *http.Request) {
  session, err := auth(r) 
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    fmt.Printf("Can't get session in db cause, %s\n", err.Error())
    return
  }
 
  user, err := database.New().GetUserByUUid(session.OwnerId)
  if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    fmt.Printf("Can't get user in db cause, %s\n", err.Error())
    return
  }

  loginInfo := &LoginInfo{
    User: user.Name,  
    LoggedIn: true, 
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(loginInfo)
} 
