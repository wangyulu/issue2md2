package config

import (
	"os"
	"testing"
)

func TestGetGitHubToken(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		expected string
	}{
		{
			name:     "Token 未设置",
			setup:    func() { os.Unsetenv("GITHUB_TOKEN") },
			expected: "",
		},
		{
			name:     "Token 已设置",
			setup:    func() { os.Setenv("GITHUB_TOKEN", "ghp_xxx") },
			expected: "ghp_xxx",
		},
		{
			name:     "Token 为空字符串",
			setup:    func() { os.Setenv("GITHUB_TOKEN", "") },
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setup()
			defer os.Unsetenv("GITHUB_TOKEN")

			// Execute
			result := GetGitHubToken()

			// Verify
			if result != tt.expected {
				t.Errorf("GetGitHubToken() = %q, want %q", result, tt.expected)
			}
		})
	}
}
