package server

import (
	"big/internal/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(s.corsMiddleware)

	api := api.NewApi()

	// Session/Authentication
	r.HandleFunc("/signup", api.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", api.LogoutHandler).Methods("POST")

	// Admin Session/Management
	r.HandleFunc("/admin/login", api.VerifyAdmin).Methods("POST")
	r.HandleFunc("/admin/recipe/delete/{id}", api.DeleteRecipeHandler).Methods("PATCH")
	r.HandleFunc("/admin/user/delete/{id}", api.DeleteUserHandler).Methods("PATCH")
	r.HandleFunc("/data/admin/allrecipes", api.GetAllRecipesHandler).Methods("GET")
	r.HandleFunc("/data/admin/allusers", api.GetAllUsersHandler).Methods("GET")

	// Recipe Management
	r.HandleFunc("/uploadrecipe", api.Auth(api.UploadRecipe)).Methods("POST")
	r.HandleFunc("/recipe/{id}", api.ServeRecipe).Methods("GET")
	r.HandleFunc("/recipe/{id}", api.EditRecipeHandler).Methods("PATCH")
	r.HandleFunc("/recipe/delete/{id}", api.DeleteRecipeHandler).Methods("PATCH")
	r.HandleFunc("/data/recipe/mostviewed", api.MostViewedRecipesHandler).Methods("GET")
	r.HandleFunc("/data/recipe/search", api.SearchRecipeHandler).Methods("GET")
	r.HandleFunc("/data/recipe/name/{id}", api.RecipeInfoHandler).Methods("GET")
	r.HandleFunc("/data/recipe/content/{id}", api.RecipeMdContent).Methods("GET")

	// User Management
	r.HandleFunc("/user/{id}", api.DeleteUserHandler).Methods("DELETE")

	// Data Handlers
	r.HandleFunc("/data/login", api.LoginInfoHandler).Methods("GET")
	r.HandleFunc("/data/login/recipes", api.LoginRecipeInfoHandler).Methods("GET")
	r.HandleFunc("/data/admin/login", api.AdminLoginInfoHandler).Methods("GET")
	r.HandleFunc("/data/admin/dashboard", api.AdminDashboardDataHandler).Methods("GET")

	// Debugging
	r.HandleFunc("/health", s.healthHandler)

	return r
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from frontend origin (replace with actual frontend URL)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Set a specific origin

		// Allowed HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// Allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

		// If you want to allow credentials (cookies, auth headers), set this to true
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
