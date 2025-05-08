package api

import (
	"big/internal/modals"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *api) LikeRecipeHandler(w http.ResponseWriter, r *http.Request) {
	info := &struct {
		Userid   string `json:"userid"`
		Recipeid string `json:"recipeid"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Form Error", http.StatusBadRequest)
		fmt.Print(err.Error())
		return
	}

	like := modals.NewLike(info.Userid, info.Recipeid)

	err = api.db.AddLike(like)
	if err != nil {
		http.Error(w, "Can't Add Like to db", http.StatusInternalServerError)
		fmt.Print(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) IsLikedHandler(w http.ResponseWriter, r *http.Request) {
	info := &struct {
		Userid   string `json:"userid"`
		Recipeid string `json:"recipeid"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Form Error", http.StatusBadRequest)
		fmt.Print(err.Error())
		return
	}

	err = api.db.IsLiked(info.Userid, info.Recipeid)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		fmt.Print(err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"liked": true})
}

func (api *api) DeleteLikeHandler(w http.ResponseWriter, r *http.Request) {
	info := &struct {
		Userid   string `json:"userid"`
		Recipeid string `json:"recipeid"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Form Error", http.StatusBadRequest)
		fmt.Print(err.Error())
		return
	}

	err = api.db.DeleteLikeFromUserRecipeId(info.Userid, info.Recipeid)
	if err != nil {
		http.Error(w, "Can't Add Like to db", http.StatusInternalServerError)
		fmt.Print(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
