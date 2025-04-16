package database

import (
	"big/internal/modals"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Can't initialize mock database")
	}

	s := &service{db: db}

	testCases := []struct {
		name      string
		sessionId string
		ownerUUid string
		setupMock func(exp string, sessionId string)
		expectErr bool
	}{
		{
			name:      "Success - Recipe added successfully",
			sessionId: "test-session",
			ownerUUid: "test-uuid-123",
			setupMock: func(exp string, sessionId string) {
				mock.ExpectExec("INSERT INTO sessions\\(sessionid, ownerid, exp\\) VALUES\\(\\$1, \\$2, \\$3\\)").
					WithArgs(sessionId, "test-uuid-123", exp).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectErr: false,
		},
		{
			name:      "Failure - Internal Database Error",
			sessionId: "test-session",
			ownerUUid: "test-uuid-123",
			setupMock: func(exp string, sessionId string) {
				mock.ExpectExec("INSERT INTO sessions\\(sessionid, ownerid, exp\\) VALUES\\(\\$1, \\$2, \\$3\\)").
					WithArgs(sessionId, "test-uuid-123", exp).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {

		session := modals.NewSession("test-uuid-123")

		tc.setupMock(session.Exp, session.SessionId)

		err = s.AddSession(session)

		if tc.expectErr && err == nil {
			t.Errorf("Expected error but got nil")
		}
		if !tc.expectErr && err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}

	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteSessionByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Can't initialize mock database")
	}

	s := &service{db: db}

	testCases := []struct {
		name      string
		setupMock func(ownerid string)
		expectErr bool
	}{
		{
			name: "Succes - Session deleted successfully",
			setupMock: func(ownerid string) {
				mock.ExpectExec("DELETE FROM sessions WHERE ownerid = \\$1").
					WithArgs(ownerid).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name: "Error - Session deleted successfully",
			setupMock: func(ownerid string) {
				mock.ExpectExec("DELETE FROM sessions WHERE ownerid = \\$1").
					WithArgs(ownerid).
					WillReturnError(errors.New("db Error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {

		ownerid := "test-uuid-123"

		tc.setupMock(ownerid)

		err := s.DeleteSessionByUser(ownerid)

		if tc.expectErr && err == nil {
			t.Errorf("Expected error but got nil")
		}
		if !tc.expectErr && err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}

	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

}
