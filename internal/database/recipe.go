package database

import "big/internal/modals"

func (s *service) AddRecipe(recipe *modals.Recipe) error {
	q := `
  INSERT INTO recipes(uuid, name, ownerid, views)
  VALUES($1, $2, $3, -1)
  `
	_, err := s.db.Exec(q, recipe.UUID, recipe.Name, recipe.OwnerId, recipe.Views)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetRecipe(UUID string) (*modals.Recipe, error) {
	var recipe *modals.Recipe

	q := `
  SELECT * FROM recipes
  WHERE uuid = $1;
  `

	row := s.db.QueryRow(q, UUID)
	err := row.Scan(&recipe.UUID, &recipe.Name, &recipe.OwnerId, &recipe.Views)

	if err != nil {
		return nil, err
	}

	return recipe, nil
}
