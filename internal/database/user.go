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

func (s *service) GetUserByName(name string) (*modals.User, error) {
	var user modals.User

	query := "SELECT * FROM users WHERE name = $1;"
	row := s.db.QueryRow(query, name)

	err := row.Scan(&user.UUID, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) GetUserByUUid(uuid string) (*modals.User, error) {
	var user modals.User

	query := "SELECT * FROM users WHERE uuid = $1;"
	row := s.db.QueryRow(query, uuid)

	err := row.Scan(&user.UUID, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) DeleteUserByUUid(uuid string) error {
	q := `
		DELETE FROM users 
		WHERE uuid = $1;
	`
	_, err := s.db.Exec(q, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetUserUUid(name string) (string, error) {
	var user modals.User

	query := "SELECT * FROM users WHERE name = $1"
	err := s.db.QueryRow(query, name).Scan(&user.UUID, &user.Name, &user.Password)
	if err != nil {
		return "", err
	}

	return user.UUID, nil
}

func (s *service) NumberOfUsers() int {
	var numberOfUsers int

	q := `
		SELECT COUNT(*) FROM users; 
	`
	s.db.QueryRow(q).Scan(&numberOfUsers)

	return numberOfUsers
}

func (s *service) GetAllUsers() ([]modals.User, error) {
	var users []modals.User

	q := `SELECT * FROM users;`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user modals.User
		err := rows.Scan(&user.UUID, &user.Name, &user.Password) // Adjust fields as per your Recipe struct
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of recipes
	return users, nil

}
