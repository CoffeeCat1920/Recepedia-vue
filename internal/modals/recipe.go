package modals

import "github.com/google/uuid"

type Recipe struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	OwnerId string `json:"ownerId"`
	Views   int    `json:"views"`
}

func NewRecipe(name, ownerId string) *Recipe {
	return &Recipe{
		UUID:    uuid.NewString(),
		Name:    name,
		OwnerId: ownerId,
		Views:   0,
	}
}
