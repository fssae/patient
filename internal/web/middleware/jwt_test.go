package middleware

import (
	"testing"
)

func TestMatchPathPattern(t *testing.T) {
	builder := &LoginJWTMiddlewareBuilder{}

	tests := []struct {
		name        string
		requestPath string
		pattern     string
		expected    bool
	}{
		{
			name:        "Exact match",
			requestPath: "/login",
			pattern:     "/login",
			expected:    true,
		},
		{
			name:        "No match",
			requestPath: "/api/users",
			pattern:     "/login",
			expected:    false,
		},
		{
			name:        "Wildcard match - swagger",
			requestPath: "/swagger/index.html",
			pattern:     "/swagger/*",
			expected:    true,
		},
		{
			name:        "Wildcard match - nested path",
			requestPath: "/swagger/v2/doc.json",
			pattern:     "/swagger/*",
			expected:    true,
		},
		{
			name:        "Wildcard match - api prefix",
			requestPath: "/api/user/login",
			pattern:     "/api/*",
			expected:    true,
		},
		{
			name:        "Wildcard match - debug",
			requestPath: "/debug/config",
			pattern:     "/debug/*",
			expected:    true,
		},
		{
			name:        "Wildcard match - exact prefix",
			requestPath: "/swagger",
			pattern:     "/swagger/*",
			expected:    true,
		},
		{
			name:        "No wildcard match",
			requestPath: "/swagge",
			pattern:     "/swagger/*",
			expected:    false,
		},
		{
			name:        "Register exact match",
			requestPath: "/register",
			pattern:     "/register",
			expected:    true,
		},
		{
			name:        "Metrics exact match",
			requestPath: "/metrics",
			pattern:     "/metrics",
			expected:    true,
		},
		{
			name:        "Health exact match",
			requestPath: "/health",
			pattern:     "/health",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.matchPathPattern(tt.requestPath, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchPathPattern(%s, %s) = %v, want %v",
					tt.requestPath, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	builder := &LoginJWTMiddlewareBuilder{}

	tests := []struct {
		name        string
		tokenHeader string
		expected    string
	}{
		{
			name:        "Bearer token with uppercase",
			tokenHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:        "Bearer token with lowercase",
			tokenHeader: "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:        "Token without Bearer prefix",
			tokenHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:        "Invalid format with multiple spaces",
			tokenHeader: "Bearer token extra",
			expected:    "",
		},
		{
			name:        "Empty token",
			tokenHeader: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.extractToken(tt.tokenHeader)
			if result != tt.expected {
				t.Errorf("extractToken(%s) = %s, want %s",
					tt.tokenHeader, result, tt.expected)
			}
		})
	}
}

func TestIsIgnorePath(t *testing.T) {
	builder := &LoginJWTMiddlewareBuilder{
		paths: []string{"/custom/path"},
		ignorePaths: []string{
			"/login",
			"/register",
			"/swagger/*",
			"/metrics",
		},
	}

	tests := []struct {
		name        string
		requestPath string
		expected    bool
	}{
		{
			name:        "Login path should be ignored",
			requestPath: "/login",
			expected:    true,
		},
		{
			name:        "Register path should be ignored",
			requestPath: "/register",
			expected:    true,
		},
		{
			name:        "Swagger path should be ignored",
			requestPath: "/swagger/index.html",
			expected:    true,
		},
		{
			name:        "Metrics path should be ignored",
			requestPath: "/metrics",
			expected:    true,
		},
		{
			name:        "Custom path should be ignored",
			requestPath: "/custom/path",
			expected:    true,
		},
		{
			name:        "Protected path should not be ignored",
			requestPath: "/api/users",
			expected:    false,
		},
		{
			name:        "Another protected path should not be ignored",
			requestPath: "/api/customer/list",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.isIgnorePath(tt.requestPath)
			if result != tt.expected {
				t.Errorf("isIgnorePath(%s) = %v, want %v",
					tt.requestPath, result, tt.expected)
			}
		})
	}
}
