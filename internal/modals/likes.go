package modals

import "github.com/google/uuid"

type Like struct {
	UUID     uuid.UUID `json:"uuid"`
	UserId   string    `json:"userid"`
	RecipeId string    `json:"recipeid"`
}

func NewLike(userid string, recipe string) *Like {
	return &Like{
		UUID:     uuid.New(),
		UserId:   userid,
		RecipeId: recipe,
	}
}
