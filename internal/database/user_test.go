package database

import (
	"big/internal/modals"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		user        *modals.User
		setupMock   func()
		expectedErr error
	}{
		{
			name: "Success - New user inserted",
			user: &modals.User{
				UUID:     "test-uuid",
				Name:     "testuser",
				Password: "password123",
			},
			setupMock: func() {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE name = \\$1\\)").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				mock.ExpectExec("INSERT INTO users\\(uuid, name, password\\)").
					WithArgs("test-uuid", "testuser", "password123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Failure - User already exists",
			user: &modals.User{
				UUID:     "test-uuid",
				Name:     "existinguser",
				Password: "password123",
			},
			setupMock: func() {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE name = \\$1\\)").
					WithArgs("existinguser").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedErr: ErrItemAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err := s.AddUser(tc.user)

			if err != tc.expectedErr {
				if err == nil {
					t.Errorf("Expected no error, got %v", err)
				} else {
					t.Errorf("Expected error %v got no err", tc.expectedErr)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
