package api

import (
	"encoding/json"
	"net/http"
)

func (api *api) AdminDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	db := api.db
	numberOfRecipes := db.NumberOfRecipes()
	numberOfUsers := db.NumberOfUsers()

	data := struct {
		NumberOfRecipes int `json:"numberOfRecipes"`
		NumberOfUsers   int `json:"numberOfUsers"`
	}{
		NumberOfRecipes: numberOfRecipes,
		NumberOfUsers:   numberOfUsers,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
