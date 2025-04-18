package modals_test

import (
	"big/internal/modals"
	"testing"
)

func TestUserCheckPassword(t *testing.T) {
	tests := []struct {
		name          string
		real_password string
		test_password string
		expected      bool
	}{
		{"Ali", "123", "456", false},
		{"Awais", "123", "123", true},
		{"Abdul Rehamn", "fsasd", "fsasd", true},
		{"Ayan", "AsdF", "asdf", false},
	}

	for _, tt := range tests {
		user := modals.NewUser(tt.name, tt.real_password)
		isCorrect := user.CheckPassword(tt.test_password)

		if isCorrect != tt.expected {
			t.Errorf("(u *user)CheckPassword() returned = %v, expected = %v", isCorrect, tt.expected)
		}
	}
}
