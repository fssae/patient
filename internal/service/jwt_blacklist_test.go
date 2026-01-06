package service

import (
	"testing"
)

func TestJWTBlacklistService_MakeKey(t *testing.T) {
	service := &JWTBlacklistService{}

	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "Simple token",
			token:    "token123",
			expected: "jwt:blacklist:token123",
		},
		{
			name:     "JWT token",
			token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0",
			expected: "jwt:blacklist:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0",
		},
		{
			name:     "Empty token",
			token:    "",
			expected: "jwt:blacklist:",
		},
		{
			name:     "Token with special characters",
			token:    "token-with-dashes_and_underscores.123",
			expected: "jwt:blacklist:token-with-dashes_and_underscores.123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.makeKey(tt.token)
			if result != tt.expected {
				t.Errorf("makeKey(%s) = %s, want %s", tt.token, result, tt.expected)
			}
		})
	}
}
