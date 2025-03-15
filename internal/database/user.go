package database

import (
	"big/internal/modals"
)

func (s *service) AddUser(user *modals.User) error {
	if s.doesExists(user.Name, "name", "users") {
		return ErrUserAlreadyExists
	}

	q := `
  INSERT INTO users(uuid, name, password)
  VALUES($1, $2, $3)
  `
	_, err := s.db.Exec(q, user.UUID, user.Name, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetUser(UUID string) (*modals.User, error) {
	var user *modals.User

	q := `
  SELECT * FROM users 
  WHERE uuid = $1
  `

	row := s.db.QueryRow(q, UUID)
	err := row.Scan(&user.UUID, &user.Name, &user.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}
