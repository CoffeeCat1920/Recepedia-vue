package api

import (
	"big/internal/modals"
	"fmt"
	"net/http"
)

type UserRequest struct {
  Name string `json:"name"`
  Password string `json:"password"`
}

type LoginInfo struct {
  LoggedIn bool `json:"loggedIn"`
  User string `json:"name"`
}

type RecipeInfo struct {
  Name string `json:"name"`
  Content string `json:"content"` 
}

func createCookie(w http.ResponseWriter, ownerId string) (*modals.Session) {
  session := modals.NewSession(ownerId)
  exp, err := session.GetExpTime()

  if err != nil {
    fmt.Printf("Can't make a new cookie\n")
  }

  http.SetCookie(w, &http.Cookie{
    Name: "session-token",
    Value: session.SessionId,
    Expires: exp,
    Path:     "/",         
    Domain:   "",
    HttpOnly: true,         
    SameSite: http.SameSiteLaxMode,   
  })
  return session
}
