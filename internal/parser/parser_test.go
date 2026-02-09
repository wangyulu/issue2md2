package parser

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedType   ResourceType
		expectedOwner  string
		expectedRepo   string
		expectedNumber int
		expectError    bool
	}{
		{
			name:           "Issue URL",
			url:            "https://github.com/owner/repo/issues/123",
			expectedType:   ResourceTypeIssue,
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedNumber: 123,
			expectError:    false,
		},
		{
			name:           "Issue URL with comment",
			url:            "https://github.com/owner/repo/issues/123#issuecomment-456",
			expectedType:   ResourceTypeIssue,
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedNumber: 123,
			expectError:    false,
		},
		{
			name:           "Pull Request URL",
			url:            "https://github.com/owner/repo/pull/42",
			expectedType:   ResourceTypePullRequest,
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedNumber: 42,
			expectError:    false,
		},
		{
			name:           "Pull Request URL with discussion",
			url:            "https://github.com/owner/repo/pull/42#discussion_r123",
			expectedType:   ResourceTypePullRequest,
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedNumber: 42,
			expectError:    false,
		},
		{
			name:           "Discussion URL",
			url:            "https://github.com/owner/repo/discussions/7",
			expectedType:   ResourceTypeDiscussion,
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedNumber: 7,
			expectError:    false,
		},
		{
			name:           "Discussion URL with comment",
			url:            "https://github.com/owner/repo/discussions/7#discussioncomment-890",
			expectedType:   ResourceTypeDiscussion,
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedNumber: 7,
			expectError:    false,
		},
		{
			name:        "Invalid URL (not github)",
			url:         "https://example.com/invalid",
			expectError: true,
		},
		{
			name:        "Invalid URL (tree page)",
			url:         "https://github.com/owner/repo/tree/main",
			expectError: true,
		},
		{
			name:        "Invalid URL (missing number)",
			url:         "https://github.com/owner/repo/issues",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseURL(tt.url)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseURL(%q) expected error, got nil", tt.url)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseURL(%q) unexpected error: %v", tt.url, err)
				return
			}

			if result.Type != tt.expectedType {
				t.Errorf("ParseURL(%q).Type = %q, want %q", tt.url, result.Type, tt.expectedType)
			}
			if result.Owner != tt.expectedOwner {
				t.Errorf("ParseURL(%q).Owner = %q, want %q", tt.url, result.Owner, tt.expectedOwner)
			}
			if result.Repo != tt.expectedRepo {
				t.Errorf("ParseURL(%q).Repo = %q, want %q", tt.url, result.Repo, tt.expectedRepo)
			}
			if result.Number != tt.expectedNumber {
				t.Errorf("ParseURL(%q).Number = %d, want %d", tt.url, result.Number, tt.expectedNumber)
			}
			if result.Original != tt.url {
				t.Errorf("ParseURL(%q).Original = %q, want %q", tt.url, result.Original, tt.url)
			}
		})
	}
}
