package modals_test

import (
	"big/internal/modals"
	"testing"
)

func TestAdminCheckPassword(t *testing.T) {

	tests := []struct {
		test_password string
		expected      bool
	}{
		{"123", true},
		{"fsasd", false},
	}

	for _, tt := range tests {
		user := modals.NewAdmin()
		isCorrect := user.CheckPassword(tt.test_password)

		if isCorrect != tt.expected {
			t.Errorf("(a *admin)CheckPassword() returned = %v, expected = %v", isCorrect, tt.expected)
		}
	}
}
