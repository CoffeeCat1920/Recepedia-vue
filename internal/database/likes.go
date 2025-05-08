package database

import (
	"big/internal/modals"
	"database/sql"
)

func (s *service) AddLike(like *modals.Like) error {
	q := `
  INSERT INTO likes(id, userid, recipeid)
  VALUES($1, $2, $3)
  `
	_, err := s.db.Exec(q, like.UUID, like.UserId, like.RecipeId)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteLikeFromUserRecipeId(userid string, recipeid string) error {
	q := "delete from likes where userid = $1 and recipeid = $2"

	res, err := s.db.Exec(q, userid, recipeid)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return ErrItemNotFound
	}

	return nil
}

func (s *service) IsLiked(userid string, recipeid string) error {
	q := "SELECT 1 FROM likes WHERE userid = $1 AND recipeid = $2 LIMIT 1"

	var exists int
	err := s.db.QueryRow(q, userid, recipeid).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrItemNotFound
		}
		return err
	}

	return nil
}
