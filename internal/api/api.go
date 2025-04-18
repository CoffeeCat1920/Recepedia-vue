package api

import (
	"big/internal/database"
	"big/internal/modals"
	"database/sql"
	"fmt"
	"net/http"
)

type UserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AdminRequest struct {
	Password string `json:"password"`
}

type LoginInfo struct {
	LoggedIn bool   `json:"loggedIn"`
	UUID     string `json:"uuid"`
	User     string `json:"name"`
}

type AdminLoginInfo struct {
	LoggedIn bool   `json:"loggedIn"`
	UUID     string `json:"uuid"`
}

type RecipeInfo struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Api interface {
	Auth(next http.HandlerFunc) http.HandlerFunc

	SignUpHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)

	VerifyAdmin(w http.ResponseWriter, r *http.Request)
	GetAllRecipesHandler(w http.ResponseWriter, r *http.Request)
	GetAllUsersHandler(w http.ResponseWriter, r *http.Request)

	UploadRecipe(w http.ResponseWriter, r *http.Request)
	ServeRecipe(w http.ResponseWriter, r *http.Request)
	EditRecipeHandler(w http.ResponseWriter, r *http.Request)
	DeleteRecipeHandler(w http.ResponseWriter, r *http.Request)
	MostViewedRecipesHandler(w http.ResponseWriter, r *http.Request)
	SearchRecipeHandler(w http.ResponseWriter, r *http.Request)
	RecipeInfoHandler(w http.ResponseWriter, r *http.Request)
	RecipeMdContent(w http.ResponseWriter, r *http.Request)

	DeleteUserHandler(w http.ResponseWriter, r *http.Request)

	LoginInfoHandler(w http.ResponseWriter, r *http.Request)
	LoginRecipeInfoHandler(w http.ResponseWriter, r *http.Request)
	AdminLoginInfoHandler(w http.ResponseWriter, r *http.Request)
	AdminDashboardDataHandler(w http.ResponseWriter, r *http.Request)
}

type api struct {
	db database.Service
}

func NewApi() Api {
	return &api{db: database.New()}
}

func NewTest(db *sql.DB) Api {
	return &api{db: database.NewTest(db)}
}

func (api *api) createCookie(w http.ResponseWriter, ownerId string) *modals.Session {
	session := modals.NewSession(ownerId)
	exp, err := session.GetExpTime()

	if err != nil {
		fmt.Printf("Can't make a new cookie\n")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    session.SessionId,
		Expires:  exp,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return session
}
