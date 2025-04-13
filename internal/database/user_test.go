package database

import (
	"big/internal/modals"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestService_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
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
			name: "Success - User added successfully",
			user: &modals.User{
				UUID:     "test-uuid",
				Name:     "testuser",
				Password: "password123",
			},
			setupMock: func() {
				// Mock the doesExists check - user doesn't exist
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE name = \\$1\\)").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				// Mock the successful insert operation
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
				// Mock the doesExists check - user exists
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE name = \\$1\\)").
					WithArgs("existinguser").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			},
			expectedErr: ErrUserAlreadyExists,
		},
		{
			name: "Failure - Database query error during existence check",
			user: &modals.User{
				UUID:     "test-uuid",
				Name:     "testuser",
				Password: "password123",
			},
			setupMock: func() {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE name = \\$1\\)").
					WithArgs("testuser").
					WillReturnError(fmt.Errorf("connection refused"))

				mock.ExpectExec("INSERT INTO users\\(uuid, name, password\\)").
					WithArgs("test-uuid", "testuser", "password123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Failure - Database insert constraint violation",
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
					WillReturnError(fmt.Errorf("unique constraint violation"))
			},
			expectedErr: fmt.Errorf("unique constraint violation"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := s.AddUser(tc.user)

			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tc.expectedErr != nil && err == nil {
				t.Errorf("Expected error %v, got no error", tc.expectedErr)
			} else if tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error() {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetUserByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		userName    string
		setupMock   func()
		expectedErr bool
	}{

		{
			name:     "Success",
			userName: "test",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"uuid", "name", "password"}).
					AddRow("test-uuid-123", "test", "hashed_password_here")

				mock.ExpectQuery("SELECT \\* FROM users WHERE name = \\$1").
					WithArgs("test").
					WillReturnRows(rows)
			},
			expectedErr: false,
		},
		{
			name:     "Failure",
			userName: "test",
			setupMock: func() {
				mock.ExpectQuery("SELECT \\* FROM users WHERE name = \\$1").
					WithArgs("test").
					WillReturnRows(sqlmock.NewRows([]string{"uuid", "name", "password"}))
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			user, err := s.GetUserByName(tc.userName)

			if (err != nil) != tc.expectedErr {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				} else {
					t.Errorf("Expected error %v got noting", err)
				}
			}

			if !tc.expectedErr {
				if user.Name != tc.userName {
					t.Errorf("Expected in GetUserByName(), %s but got %s", tc.userName, user.Name)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

}

func TestGetUserByUUid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		userUUid    string
		setupMock   func()
		expectedErr bool
	}{
		{
			name:     "Success",
			userUUid: "test-uuid-123",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"uuid", "name", "password"}).
					AddRow("test-uuid-123", "test", "hashed_password_here")

				mock.ExpectQuery("SELECT \\* FROM users WHERE uuid = \\$1").
					WithArgs("test-uuid-123").
					WillReturnRows(rows)
			},
			expectedErr: false,
		},
		{
			name:     "Failure",
			userUUid: "test-uuid-123",
			setupMock: func() {
				mock.ExpectQuery("SELECT \\* FROM users WHERE uuid = \\$1").
					WithArgs("test-uuid-123").
					WillReturnRows(sqlmock.NewRows([]string{"uuid", "name", "password"}))
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			user, err := s.GetUserByUUid(tc.userUUid)

			if (err != nil) != tc.expectedErr {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				} else {
					t.Errorf("Expected error %v got noting", err)
				}
			}

			if !tc.expectedErr {
				if user.UUID != tc.userUUid {
					t.Errorf("Expected in GetUserByUUid(), %s but got %s", tc.userUUid, user.UUID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

}

func TestDeleteUserByUUid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	s := &service{db: db}

	testCases := []struct {
		name        string
		userUUid    string
		setupMock   func()
		expectedErr bool
	}{
		{
			name:     "Success",
			userUUid: "test-uuid-123",
			setupMock: func() {

				mock.ExpectExec("DELETE FROM users WHERE uuid = \\$1").
					WithArgs("test-uuid-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: false,
		},
		{
			name:     "Failure",
			userUUid: "test-uuid-123",
			setupMock: func() {
				mock.ExpectExec("DELETE FROM users WHERE uuid = \\$1").
					WithArgs("test-uuid-123").
					WillReturnError(errors.New("db Error"))
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := s.DeleteUserByUUid(tc.userUUid)

			if (err != nil) != tc.expectedErr {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				} else {
					t.Errorf("Expected error %v got noting", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

}

func TestNumberOfUsers(t *testing.T) {

}
