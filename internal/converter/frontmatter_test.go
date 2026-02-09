package converter

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		url      string
		author   string
		authorURL string
		createdAt time.Time
		status   string
		typ      string
		expected string
	}{
		{
			name:     "Issue Frontmatter",
			title:    "Test Issue",
			url:      "https://github.com/owner/repo/issues/123",
			author:   "octocat",
			authorURL: "https://github.com/octocat",
			createdAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			status:   "open",
			typ:      "issue",
			expected: `---
title: "Test Issue"
url: "https://github.com/owner/repo/issues/123"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "issue"
---
`,
		},
		{
			name:     "PR Frontmatter",
			title:    "Test PR",
			url:      "https://github.com/owner/repo/pull/42",
			author:   "octocat",
			authorURL: "https://github.com/octocat",
			createdAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			status:   "open",
			typ:      "pull_request",
			expected: `---
title: "Test PR"
url: "https://github.com/owner/repo/pull/42"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "pull_request"
---
`,
		},
		{
			name:     "Discussion Frontmatter",
			title:    "Test Discussion",
			url:      "https://github.com/owner/repo/discussions/7",
			author:   "octocat",
			authorURL: "https://github.com/octocat",
			createdAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			status:   "open",
			typ:      "discussion",
			expected: `---
title: "Test Discussion"
url: "https://github.com/owner/repo/discussions/7"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "discussion"
---
`,
		},
		{
			name:     "Title with special characters",
			title:    `Test "Issue" with 'quotes' and \backslashes\`,
			url:      "https://github.com/owner/repo/issues/123",
			author:   "octocat",
			authorURL: "https://github.com/octocat",
			createdAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			status:   "open",
			typ:      "issue",
			expected: `---
title: 'Test "Issue" with ''quotes'' and \backslashes\'
url: "https://github.com/owner/repo/issues/123"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "issue"
---
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFrontmatter(tt.title, tt.url, tt.author, tt.authorURL, tt.createdAt, tt.status, tt.typ)

			if result != tt.expected {
				// Split into lines for better error reporting
				resultLines := strings.Split(result, "\n")
				expectedLines := strings.Split(tt.expected, "\n")

				t.Errorf("generateFrontmatter() output mismatch:\nGot:\n%s\n\nWant:\n%s", result, tt.expected)

				// Find the first differing line
				maxLines := len(resultLines)
				if len(expectedLines) > maxLines {
					maxLines = len(expectedLines)
				}
				for i := 0; i < maxLines; i++ {
					var got, want string
					if i < len(resultLines) {
						got = resultLines[i]
					}
					if i < len(expectedLines) {
						want = expectedLines[i]
					}
					if got != want {
						t.Errorf("Line %d:\n  Got:  %q\n  Want: %q", i+1, got, want)
						break
					}
				}
			}
		})
	}
}
