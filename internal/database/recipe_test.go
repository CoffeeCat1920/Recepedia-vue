package database

import (
	"big/internal/modals"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRecipeModule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to initialize sqlmock:", err)
	}
	defer db.Close()

	s := &service{db: db}

	// Step 0: Define initial recipe
	recipe := &modals.Recipe{
		UUID:    "recipe-uuid-123",
		Name:    "Chocolate Cake",
		OwnerId: "user-123",
	}

	// Step 1: Expect AddRecipe
	mock.ExpectExec(`INSERT INTO recipes\(uuid, name, ownerid, views\) VALUES\(\$1, \$2, \$3, -1\)`).
		WithArgs(recipe.UUID, recipe.Name, recipe.OwnerId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute Add
	if err := s.AddRecipe(recipe); err != nil {
		t.Fatalf("AddRecipe failed: %v", err)
	}

	// Step 2: Expect EditRecipe (change name)
	updatedName := "Vanilla Cake"
	mock.ExpectExec(`UPDATE recipes SET name = \$1`).
		WithArgs(updatedName, recipe.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute Edit
	if err := s.EditRecipeName(recipe.UUID, updatedName); err != nil {
		t.Fatalf("EditRecipe failed: %v", err)
	}

	// Step 3: Expect GetRecipe (after edit)
	mock.ExpectQuery(`SELECT \* FROM recipes WHERE uuid = \$1;`).
		WithArgs(recipe.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"}).
			AddRow(recipe.UUID, updatedName, recipe.OwnerId, -1))

	// Execute Get
	updatedRecipe, err := s.GetRecipe(recipe.UUID)
	if err != nil {
		t.Fatalf("GetRecipe failed: %v", err)
	}

	if updatedRecipe.Name != updatedName {
		t.Errorf("Expected recipe name to be %s, got %s", updatedName, updatedRecipe.Name)
	}

	// Step 4: Expect DeleteRecipe
	mock.ExpectExec(`DELETE FROM recipes WHERE uuid = \$1`).
		WithArgs(recipe.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute Delete
	if err := s.DeleteRecipe(recipe.UUID); err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	emptyRows := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"})
	mock.ExpectQuery("SELECT \\* FROM recipes WHERE uuid = \\$1;").
		WithArgs(recipe.UUID).
		WillReturnRows(emptyRows)

	_, err = s.GetRecipe(recipe.UUID)
	if !errors.Is(err, ErrItemNotFound) {
		t.Errorf("Expected ErrItemNotFound, got: %v", err)
	}

	// Final: Ensure all expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAddRecipe(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Can't Initalize database")
	}

	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		recipe      *modals.Recipe
		setupMock   func(recipe *modals.Recipe)
		expectedErr error
	}{
		{
			name: "Success - Recipe Added Successfully",
			recipe: &modals.Recipe{
				UUID:    "test-uuid-123",
				Name:    "test-recipe",
				OwnerId: "owner-123",
			},
			setupMock: func(recipe *modals.Recipe) {
				mock.ExpectExec(`INSERT INTO recipes\(uuid, name, ownerid, views\) VALUES\(\$1, \$2, \$3, -1\)`).
					WithArgs(recipe.UUID, recipe.Name, recipe.OwnerId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		tc.setupMock(tc.recipe)

		err := s.AddRecipe(tc.recipe)

		if err != tc.expectedErr {
			if err != nil && tc.expectedErr == nil {
				t.Errorf("In %s, Expected no error but got %v", tc.name, err)
			} else if err != nil && tc.expectedErr != nil {
				t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
			} else if err == nil && tc.expectedErr != nil {
				t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
			}
		}

	}

}

func TestGetRecipe(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Can't initialize sqlmock")
	}

	s := &service{db: db}

	testCases := []struct {
		name        string
		uuid        string
		setupMock   func(recipeName string)
		expectedErr error
	}{
		{
			name: "Success - Got a Recipe",
			uuid: "test-uuid-123",
			setupMock: func(recipeName string) {
				row := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"}).
					AddRow(`test-uuid-123`, `test`, `test-uuid-123`, 4)

				mock.ExpectQuery("SELECT \\* FROM recipes WHERE uuid = \\$1;").
					WithArgs(recipeName).
					WillReturnRows(row)
			},
			expectedErr: nil,
		},
		{
			name: "Failure - No Recipe Found",
			uuid: "non-existent-uuid",
			setupMock: func(_ string) {
				emptyRows := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"})
				mock.ExpectQuery("SELECT \\* FROM recipes WHERE uuid = \\$1;").
					WithArgs("non-existent-uuid").
					WillReturnRows(emptyRows)
			},
			expectedErr: ErrItemNotFound,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			tc.setupMock(tc.uuid)

			recipe, err := s.GetRecipe(tc.uuid)

			if err != tc.expectedErr {
				if err != nil && tc.expectedErr == nil {
					t.Errorf("In %s, Expected no error but got %v", tc.name, err)
				} else if err != nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
				} else if err == nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
				}
			}

			if recipe != nil && recipe.UUID != tc.uuid {
				t.Error(ErrItemMismatch.Error())
			}

		})

	}

}

func TestDeleteRecipe(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Can't Initalize database")
	}

	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		uuid        string
		setupMock   func(uuid string)
		expectedErr error
	}{
		{
			name: "Success - Recipe Deleted Successfully",
			uuid: "test-uuid-123",
			setupMock: func(uuid string) {
				mock.ExpectExec("DELETE FROM recipes WHERE uuid = \\$1").
					WithArgs(uuid).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Failure - Recipe Not Found",
			uuid: "test-uuid-123",
			setupMock: func(uuid string) {
				mock.ExpectExec("DELETE FROM recipes WHERE uuid = \\$1").
					WithArgs(uuid).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: ErrItemNotFound,
		},
	}

	for _, tc := range testCases {
		tc.setupMock(tc.uuid)

		err := s.DeleteRecipe(tc.uuid)

		if err != tc.expectedErr {
			if err != nil && tc.expectedErr == nil {
				t.Errorf("In %s, Expected no error but got %v", tc.name, err)
			} else if err != nil && tc.expectedErr != nil {
				t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
			} else if err == nil && tc.expectedErr != nil {
				t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
			}
		}

	}

}

func TestDeleteRecipeByUser(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Can't Initalize database")
	}

	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		ownerid     string
		setupMock   func(ownerid string)
		expectedErr error
	}{
		{
			name:    "Success - Recipe Deleted Successfully",
			ownerid: "test-user",
			setupMock: func(ownerid string) {
				mock.ExpectExec("DELETE FROM recipes WHERE ownerid = \\$1").
					WithArgs(ownerid).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name:    "Success - No Recipe Found",
			ownerid: "test-user",
			setupMock: func(ownerid string) {
				mock.ExpectExec("DELETE FROM recipes WHERE ownerid = \\$1").
					WithArgs(ownerid).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			tc.setupMock(tc.ownerid)

			err := s.DeleteRecipeByUser(tc.ownerid)

			if err != tc.expectedErr {
				if err != nil && tc.expectedErr == nil {
					t.Errorf("In %s, Expected no error but got %v", tc.name, err)
				} else if err != nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
				} else if err == nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
				}
			}

		})

	}

}

func MostViewedRecipes(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Can't initialize sqlmock")
	}

	s := &service{db: db}

	testCases := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "Success - Got Most Viewed recipes",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"})
				rows.AddRow("test-uuid-123", "test-recipe-1", "test-uuid-123", 0)
				rows.AddRow("test-uuid-123", "test-recipe-2", "test-uuid-123", 3)
				rows.AddRow("test-uuid-123", "test-recipe-3", "test-uuid-123", 5)

				mock.ExpectQuery("SELECT * FROM recipes ORDER BY views LIMIT 10").
					WillReturnRows(rows)
			},
			expectedErr: nil,
		},
		{
			name: "Failure - Got no recipes",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"})
				mock.ExpectQuery("SELECT * FROM recipes ORDER BY views LIMIT 10").
					WillReturnRows(rows)
			},
			expectedErr: ErrItemNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			_, err := s.MostViewedRecipes()

			if err != tc.expectedErr {
				if err != nil && tc.expectedErr == nil {
					t.Errorf("In %s, Expected no error but got %v", tc.name, err)
				} else if err != nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
				} else if err == nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
				}
			}

		})
	}

}

func TestSearchRecipe(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Can't initialize sqlmock")
	}

	s := &service{db: db}

	testCases := []struct {
		name        string
		searchTerm  string
		setupMock   func(searchTerm string)
		expectedErr error
	}{
		{
			name:       "Success - Got Searched Recipe",
			searchTerm: "test",
			setupMock: func(searchTerm string) {
				rows := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"}).
					AddRow(`test-uuid-123`, `test1`, `test-uuid-123`, 0).
					AddRow(`test-uuid-123`, `test2`, `test-uuid-123`, 1).
					AddRow(`test-uuid-123`, `test3`, `test-uuid-123`, 2)

				mock.ExpectQuery("SELECT \\* FROM recipes WHERE name ILIKE \\$1").
					WithArgs("%" + searchTerm + "%").
					WillReturnRows(rows)
			},
			expectedErr: nil,
		},
		{
			name:       "Success - No Recipe Found",
			searchTerm: "test",
			setupMock: func(searchTerm string) {
				rows := sqlmock.NewRows([]string{"uuid", "name", "ownerid", "views"})

				mock.ExpectQuery("SELECT \\* FROM recipes WHERE name ILIKE \\$1").
					WithArgs("%" + searchTerm + "%").
					WillReturnRows(rows)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(tc.searchTerm)

			_, err := s.SearchRecipe(tc.searchTerm)

			if err != tc.expectedErr {
				if err != nil && tc.expectedErr == nil {
					t.Errorf("In %s, Expected no error but got %v", tc.name, err)
				} else if err != nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
				} else if err == nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
				}
			}

		})
	}

}

func TestIncreaseRecipeViews(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Can't initialize sqlmock")
	}

	s := &service{db: db}

	testCases := []struct {
		name        string
		recipeName  string
		ownerId     string
		setupMock   func(recipeUUid string)
		expectedErr error
	}{
		{
			name:       "Success - Increased Recipe view by one",
			recipeName: "test-recipe",
			ownerId:    "test-uuid-123",
			setupMock: func(recipeUUid string) {
				mock.ExpectExec("UPDATE recipes SET views = views \\+ 1 WHERE uuid = \\$1").
					WithArgs(recipeUUid).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name:       "Failure - No recipe increased Recipe view by one",
			recipeName: "test-recipe",
			ownerId:    "test-uuid-123",
			setupMock: func(recipeUUid string) {
				mock.ExpectExec("UPDATE recipes SET views = views \\+ 1 WHERE uuid = \\$1").
					WithArgs(recipeUUid).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: ErrItemNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			recipe := modals.NewRecipe(tc.name, tc.ownerId)

			tc.setupMock(recipe.UUID)

			err := s.IncreaseRecipeViews(recipe)

			if err != tc.expectedErr {
				if err != nil && tc.expectedErr == nil {
					t.Errorf("In %s, Expected no error but got %v", tc.name, err)
				} else if err != nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
				} else if err == nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
				}
			}

		})
	}
}

func TestEditRecipeName(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Can't initialize mock database")
	}

	s := &service{db: db}

	testCases := []struct {
		name          string
		recipeUUid    string
		recipeNewName string
		setupMock     func(recipeNewName string, recipeUUid string)
		expectedErr   error
	}{
		{
			name:          "Success - Edited Recipe Name",
			recipeNewName: "test",
			recipeUUid:    "test-uuid-123",
			setupMock: func(recipeNewName string, recipeUUid string) {
				mock.ExpectExec("UPDATE recipes SET name = \\$1 WHERE uuid = \\$2").
					WithArgs(recipeNewName, recipeUUid).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name:          "Failure - Null Value",
			recipeNewName: "",
			recipeUUid:    "test-uuid-123",
			setupMock: func(recipeNewName string, recipeUUid string) {
				mock.ExpectExec("UPDATE recipes SET name = \\$1 WHERE uuid = \\$2").
					WithArgs(recipeNewName, recipeUUid).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: ErrItemNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupMock(tc.recipeNewName, tc.recipeUUid)

			err := s.EditRecipeName(tc.recipeUUid, tc.recipeNewName)

			if err != tc.expectedErr {
				if err != nil && tc.expectedErr == nil {
					t.Errorf("In %s, Expected no error but got %v", tc.name, err)
				} else if err != nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got %v", tc.name, tc.expectedErr, err)
				} else if err == nil && tc.expectedErr != nil {
					t.Errorf("In %s, Expected error %v but got nil", tc.name, tc.expectedErr)
				}
			}

		})
	}

}
