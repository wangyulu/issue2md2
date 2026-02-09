package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/wangyulu/issue2md2/internal/github"
)

func TestRenderReactions(t *testing.T) {
	tests := []struct {
		name     string
		r        *Reactions
		expected string
	}{
		{
			name: "Multiple reactions",
			r: &Reactions{
				ThumbsUp:   5,
				ThumbsDown: 2,
				Heart:      3,
			},
			expected: "👍 5 👎 2 ❤️ 3",
		},
		{
			name: "Single reaction",
			r: &Reactions{
				ThumbsUp: 1,
			},
			expected: "👍 1",
		},
		{
			name:     "All zero",
			r:        &Reactions{},
			expected: "",
		},
		{
			name:     "Nil reactions",
			r:        nil,
			expected: "",
		},
		{
			name: "All reaction types",
			r: &Reactions{
				ThumbsUp:   1,
				ThumbsDown: 1,
				Laugh:      1,
				Hooray:     1,
				Confused:   1,
				Heart:      1,
				Rocket:     1,
				Eyes:       1,
			},
			expected: "👍 1 👎 1 😄 1 🎉 1 😕 1 ❤️ 1 🚀 1 👀 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderReactions(tt.r)

			if result != tt.expected {
				t.Errorf("renderReactions() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRenderUser(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		userURL     string
		enableLinks bool
		expected    string
	}{
		{
			name:        "With links enabled",
			username:    "octocat",
			userURL:     "https://github.com/octocat",
			enableLinks: true,
			expected:    "[@octocat](https://github.com/octocat)",
		},
		{
			name:        "With links disabled",
			username:    "octocat",
			userURL:     "https://github.com/octocat",
			enableLinks: false,
			expected:    "@octocat",
		},
		{
			name:        "Empty username",
			username:    "",
			userURL:     "https://github.com/ghost",
			enableLinks: true,
			expected:    "[@](https://github.com/ghost)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderUser(tt.username, tt.userURL, tt.enableLinks)

			if result != tt.expected {
				t.Errorf("renderUser() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.EnableReactions != false {
		t.Errorf("DefaultOptions().EnableReactions = %v, want false", opts.EnableReactions)
	}
	if opts.EnableUserLinks != false {
		t.Errorf("DefaultOptions().EnableUserLinks = %v, want false", opts.EnableUserLinks)
	}
}

func TestToMarkdown(t *testing.T) {
	tests := []struct {
		name           string
		issue          *github.Issue
		enableReactions bool
		enableUserLinks bool
		contains       []string // 检查输出是否包含某些关键字
	}{
		{
			name: "简单 Issue（无评论）",
			issue: &github.Issue{
				Title:     "Test Issue",
				Body:      "Issue body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/octocat/Hello-World/issues/123",
				Comments: []github.Comment{},
			},
			enableReactions: false,
			enableUserLinks: false,
			contains:       []string{"---", "Test Issue", "Issue body"},
		},
		{
			name: "带 Reactions 的 Issue",
			issue: &github.Issue{
				Title:     "Test Issue",
				Body:      "Issue body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/octocat/Hello-World/issues/123",
				Reactions: &github.Reactions{
					ThumbsUp: 5,
					Heart:      3,
				},
			},
			enableReactions: true,
			enableUserLinks: false,
			contains:       []string{"## Reactions", "👍 5", "❤️ 3"},
		},
		{
			name: "带评论的 Issue",
			issue: &github.Issue{
				Title:     "Test Issue",
				Body:      "Issue body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/octocat/Hello-World/issues/123",
				Comments: []github.Comment{
					{
						Body:      "First comment",
						Author:    "user1",
						AuthorURL: "https://github.com/user1",
						CreatedAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
					},
					{
						Body:      "Second comment",
						Author:    "user2",
						AuthorURL: "https://github.com/user2",
						CreatedAt: time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC),
					},
				},
			},
			enableReactions: false,
			enableUserLinks: false,
			contains: []string{"### @user1", "### @user2"},
		},
		{
			name: "带 Reactions 的评论",
			issue: &github.Issue{
				Title:     "Test Issue",
				Body:      "Issue body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/octocat/Hello-World/issues/123",
				Comments: []github.Comment{
					{
						Body:      "Comment with reactions",
						Author:    "user1",
						AuthorURL: "https://github.com/user1",
						CreatedAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
						Reactions: &github.Reactions{ThumbsUp: 2},
					},
				},
			},
			enableReactions: true,
			enableUserLinks: false,
			contains: []string{"👍 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{
				EnableReactions: tt.enableReactions,
				EnableUserLinks: tt.enableUserLinks,
			}

			result, err := ToMarkdown(tt.issue, opts)

			if err != nil {
				t.Fatalf("ToMarkdown() failed: %v", err)
			}

			if result == nil {
				t.Fatal("ToMarkdown() returned nil")
			}

			resultStr := string(result)

			// 检查是否包含预期的内容
			for _, expected := range tt.contains {
				if !strings.Contains(resultStr, expected) {
					t.Errorf("ToMarkdown() output missing expected content: %q", expected)
				}
			}

			// 检查 Frontmatter
			if !strings.HasPrefix(resultStr, "---") {
				t.Error("ToMarkdown() output should start with frontmatter")
			}
		})
	}
}

func TestToMarkdownPR(t *testing.T) {
	tests := []struct {
		name     string
		pr        *github.PullRequest
		contains  []string
	}{
		{
			name: "简单 PR（无评论）",
			pr: &github.PullRequest{
				Title:     "Test PR",
				Body:      "PR description",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/octocat/Hello-World/pull/42",
				Comments: []github.Comment{},
			},
			contains: []string{"type: \"pull_request\"", "Test PR", "PR description"},
		},
		{
			name: "带 Review Comments 的 PR",
			pr: &github.PullRequest{
				Title:     "Test PR",
				Body:      "PR description",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/octocat/Hello-World/pull/42",
				Comments: []github.Comment{
					{
						Body:      "Review comment 1",
						Author:    "reviewer1",
						AuthorURL: "https://github.com/reviewer1",
						CreatedAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
					},
					{
						Body:      "Regular comment",
						Author:    "commenter1",
						AuthorURL: "https://github.com/commenter1",
						CreatedAt: time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC),
					},
				},
			},
			contains: []string{"### @reviewer1", "### @commenter1"},
		},
		{
			name: "Merged 状态的 PR",
			pr: &github.PullRequest{
				Title:     "Test PR",
				Body:      "PR description",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "merged",
				URL:       "https://github.com/octocat/Hello-World/pull/42",
			},
			contains: []string{"status: \"merged\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions()

			result, err := ToMarkdownPR(tt.pr, opts)

			if err != nil {
				t.Fatalf("ToMarkdownPR() failed: %v", err)
			}

			if result == nil {
				t.Fatal("ToMarkdownPR() returned nil")
			}

			resultStr := string(result)

			// 检查是否包含预期的内容
			for _, expected := range tt.contains {
				if !strings.Contains(resultStr, expected) {
					t.Errorf("ToMarkdownPR() output missing expected content: %q", expected)
				}
			}
		})
	}
}

func TestToMarkdownDiscussion(t *testing.T) {
	tests := []struct {
		name     string
		discussion *github.Discussion
		contains  []string
	}{
		{
			name: "简单 Discussion（无评论）",
			discussion: &github.Discussion{
				Title:     "Test Discussion",
				Body:      "Discussion body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/community/community/discussions/12345",
				Comments: []github.Comment{},
			},
			contains: []string{"type: \"discussion\"", "Test Discussion", "Discussion body"},
		},
		{
			name: "带 Answer 的 Discussion",
			discussion: &github.Discussion{
				Title:     "Test Discussion",
				Body:      "Discussion body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/community/community/discussions/12345",
				Comments: []github.Comment{
					{
						Body:      "This is the answer",
						Author:    "expert",
						AuthorURL: "https://github.com/expert",
						CreatedAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
						IsAnswer:  true,
					},
				},
			},
			contains: []string{"✅ **Answer**"},
		},
		{
			name: "带 Reactions 的 Answer",
			discussion: &github.Discussion{
				Title:     "Test Discussion",
				Body:      "Discussion body",
				Author:    "octocat",
				AuthorURL: "https://github.com/octocat",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Status:    "open",
				URL:       "https://github.com/community/community/discussions/12345",
				Comments: []github.Comment{
					{
						Body:      "This is the answer",
						Author:    "expert",
						AuthorURL: "https://github.com/expert",
						CreatedAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
						IsAnswer:  true,
					},
				},
			},
			contains: []string{"✅ **Answer**"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions()

			result, err := ToMarkdownDiscussion(tt.discussion, opts)

			if err != nil {
				t.Fatalf("ToMarkdownDiscussion() failed: %v", err)
			}

			if result == nil {
				t.Fatal("ToMarkdownDiscussion() returned nil")
			}

			resultStr := string(result)

			// 检查是否包含预期的内容
			for _, expected := range tt.contains {
				if !strings.Contains(resultStr, expected) {
					t.Errorf("ToMarkdownDiscussion() output missing expected content: %q", expected)
				}
			}
		})
	}
}
