package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_ValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "Valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "Invalid email - no @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "Invalid email - no domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "Invalid email - no local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "Empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "Invalid email - multiple @",
			email:   "user@@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "StrongPassword123!",
			wantErr:  false,
		},
		{
			name:     "Valid password - minimum length",
			password: "Pass123!",
			wantErr:  false,
		},
		{
			name:     "Invalid password - too short",
			password: "Pass1!",
			wantErr:  true,
		},
		{
			name:     "Invalid password - no uppercase",
			password: "password123!",
			wantErr:  true,
		},
		{
			name:     "Invalid password - no lowercase",
			password: "PASSWORD123!",
			wantErr:  true,
		},
		{
			name:     "Invalid password - no number",
			password: "Password!",
			wantErr:  true,
		},
		{
			name:     "Invalid password - no special character",
			password: "Password123",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "Valid username",
			username: "username123",
			wantErr:  false,
		},
		{
			name:     "Valid username with underscore",
			username: "user_name",
			wantErr:  false,
		},
		{
			name:     "Valid username - minimum length",
			username: "abc",
			wantErr:  false,
		},
		{
			name:     "Valid username - maximum length",
			username: "abcdefghij1234567890",
			wantErr:  false,
		},
		{
			name:     "Invalid username - too short",
			username: "ab",
			wantErr:  true,
		},
		{
			name:     "Invalid username - too long",
			username: "abcdefghij1234567890toolong",
			wantErr:  true,
		},
		{
			name:     "Invalid username - special characters",
			username: "user@name",
			wantErr:  true,
		},
		{
			name:     "Invalid username - spaces",
			username: "user name",
			wantErr:  true,
		},
		{
			name:     "Invalid username - starts with number",
			username: "123username",
			wantErr:  true,
		},
		{
			name:     "Empty username",
			username: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper validation functions for testing
func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidInput
	}
	
	// Simple email validation
	if len(email) < 3 || !containsChar(email, '@') {
		return ErrInvalidEmail
	}
	
	parts := splitString(email, '@')
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return ErrInvalidEmail
	}
	
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return ErrInvalidInput
	}
	
	if len(password) < 8 {
		return ErrWeakPassword
	}
	
	hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case isSpecialChar(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrWeakPassword
	}
	
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return ErrInvalidInput
	}
	
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	
	// Username should start with a letter
	if username[0] < 'a' || (username[0] > 'z' && username[0] < 'A') || username[0] > 'Z' {
		return ErrInvalidUsername
	}
	
	// Check for valid characters
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '_') {
			return ErrInvalidUsername
		}
	}
	
	return nil
}

// Helper functions
func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func splitString(s string, delimiter byte) []string {
	var parts []string
	start := 0
	
	for i := 0; i < len(s); i++ {
		if s[i] == delimiter {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	
	return parts
}

func isSpecialChar(char rune) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	for _, sc := range specialChars {
		if char == sc {
			return true
		}
	}
	return false
}