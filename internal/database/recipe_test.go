package database

import (
	"big/internal/modals"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

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
