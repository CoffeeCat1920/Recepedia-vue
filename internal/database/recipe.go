package database

import "big/internal/modals"

func (s *service) AddRecipe(recipe *modals.Recipe) error {
	q := `
  INSERT INTO recipes(uuid, name, ownerid, views)
  VALUES($1, $2, $3, -1)
  `
	_, err := s.db.Exec(q, recipe.UUID, recipe.Name, recipe.OwnerId)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetRecipe(UUID string) (*modals.Recipe, error) {
	var recipe modals.Recipe

	q := `
  SELECT * FROM recipes
  WHERE uuid = $1;
  `

	row := s.db.QueryRow(q, UUID)
	err := row.Scan(&recipe.UUID, &recipe.Name, &recipe.OwnerId, &recipe.Views)

	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (s *service) DeleteRecipe(uuid string) error {
  q := "DELETE FROM recipes WHERE uuid = $1" 

  _, err := s.db.Query(q, uuid) 

  if err != nil {
    return  err
  }

  return nil
}

func (s *service) MostViewedRecipes() ([]modals.Recipe, error) {
  var recipes []modals.Recipe

  rows, err := s.db.Query("SELECT * FROM recipes ORDER BY views LIMIT 10;")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  for rows.Next() {
    var recipe modals.Recipe
    err := rows.Scan(&recipe.UUID, &recipe.Name, &recipe.OwnerId, &recipe.Views)
    if err != nil {
      return nil, err
    }
    recipes = append(recipes, recipe) 
  }

  return recipes, nil
} 

func (s *service) GetRecipeByName(name string) (*modals.Recipe, error) {
  var recipe *modals.Recipe

	q := `
  SELECT * FROM recipes
  WHERE name = $1;
  `
	row := s.db.QueryRow(q, name)
	err := row.Scan(&recipe.UUID, &recipe.Name, &recipe.OwnerId, &recipe.Views)

	if err != nil {
		return nil, err
	}

	return recipe, nil
} 

func (s *service) SearchRecipe(name string) ([]modals.Recipe, error) {
  var recipes []modals.Recipe

  searchTerm := "%" + name + "%" 
  query := "SELECT * FROM recipes WHERE name ILIKE $1"

  rows, err := s.db.Query(query, searchTerm)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  if err := rows.Err(); err != nil {
    return nil, err
  }

  for rows.Next() {
    var recipe modals.Recipe
    err := rows.Scan(&recipe.UUID, &recipe.Name, &recipe.OwnerId, &recipe.Views)
    if err != nil {
      return nil, err
    }
    recipes = append(recipes, recipe)
  }

  return recipes, nil
}

func (s *service) IncreaseRecipeViews(recipe *modals.Recipe) (error) {
  q := `
    UPDATE recipes 
    SET views = views + 1
    WHERE uuid = $1
  `

  _, err := s.db.Query(q, recipe.UUID)

  if err != nil {
    return err
  }

  return nil
} 

func (s *service) EditRecipeName(uuid string, name string) (error) {
	q := `
		UPDATE recipes
		SET name = $1
		WHERE uuid = $2
	`	

	_, err := s.db.Exec(q, name, uuid)

  if err != nil {
    return err
  }

  return nil
}   

func (s *service) GetRecipesByUser(name string) ([]modals.Recipe, error) {
  var recipes []modals.Recipe

  user, err := s.GetUserByName(name)
  if err != nil {
    return nil, err
  }

  q := "SELECT * FROM recipes WHERE ownerid = $1" 
  rows, err := s.db.Query(q, user.UUID)

  if err != nil {
    return nil, err
  }

  for rows.Next() {
    var recipe modals.Recipe 
    err := rows.Scan(&recipe.UUID, &recipe.Name, &recipe.OwnerId, &recipe.Views)

    if err != nil {
      return nil, err
    }

    recipes = append(recipes, recipe)
  }

  return recipes, nil
}
