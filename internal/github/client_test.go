package github

import (
	"testing"
)

// TestFetchIssue 测试获取 Issue 数据
func TestFetchIssue(t *testing.T) {
	// 这是一个集成测试，需要真实的 GitHub API 访问
	// 如果没有 GITHUB_TOKEN，可能会遇到限流问题

	client := NewClient()

	// 使用公开的 Issue 进行测试
	issue, err := client.FetchIssue("octocat", "Hello-World", 348)

	if err != nil {
		t.Fatalf("FetchIssue() failed: %v", err)
	}

	if issue == nil {
		t.Fatal("FetchIssue() returned nil issue")
	}

	// 验证基本字段
	if issue.Title == "" {
		t.Error("Issue.Title is empty")
	}
	if issue.Author == "" {
		t.Error("Issue.Author is empty")
	}
	if issue.AuthorURL == "" {
		t.Error("Issue.AuthorURL is empty")
	}
	if issue.URL == "" {
		t.Error("Issue.URL is empty")
	}
	if issue.Status == "" {
		t.Error("Issue.Status is empty")
	}
}

// TestFetchPullRequest 测试获取 Pull Request 数据
func TestFetchPullRequest(t *testing.T) {
	t.Skip("Skipping - requires a known PR ID")

	client := NewClient()

	// 使用公开的 PR 进行测试（使用一个确实存在的 PR）
	pr, err := client.FetchPullRequest("golang", "go", 62140)

	if err != nil {
		t.Fatalf("FetchPullRequest() failed: %v", err)
	}

	if pr == nil {
		t.Fatal("FetchPullRequest() returned nil PR")
	}

	// 验证基本字段
	if pr.Title == "" {
		t.Error("PullRequest.Title is empty")
	}
	if pr.Author == "" {
		t.Error("PullRequest.Author is empty")
	}
	if pr.AuthorURL == "" {
		t.Error("PullRequest.AuthorURL is empty")
	}
	if pr.URL == "" {
		t.Error("PullRequest.URL is empty")
	}
	if pr.Status == "" {
		t.Error("PullRequest.Status is empty")
	}
}

// TestFetchDiscussion 测试获取 Discussion 数据
func TestFetchDiscussion(t *testing.T) {
	client := NewClient()

	// 使用公开的 Discussion 进行测试
	// 注意：需要真实的 Discussion ID
	discussion, err := client.FetchDiscussion("community", "community", 12345)

	if err != nil {
		t.Fatalf("FetchDiscussion() failed: %v", err)
	}

	if discussion == nil {
		t.Fatal("FetchDiscussion() returned nil discussion")
	}

	// 验证基本字段
	if discussion.Title == "" {
		t.Error("Discussion.Title is empty")
	}
	if discussion.Author == "" {
		t.Error("Discussion.Author is empty")
	}
	if discussion.AuthorURL == "" {
		t.Error("Discussion.AuthorURL is empty")
	}
	if discussion.URL == "" {
		t.Error("Discussion.URL is empty")
	}
	if discussion.Status == "" {
		t.Error("Discussion.Status is empty")
	}
}

// TestFetchNonExistentIssue 测试获取不存在的 Issue
func TestFetchNonExistentIssue(t *testing.T) {
	client := NewClient()

	_, err := client.FetchIssue("octocat", "Hello-World", 999999)

	if err == nil {
		t.Error("Expected error for non-existent issue, got nil")
	}
}

// TestFetchNonExistentPR 测试获取不存在的 PR
func TestFetchNonExistentPR(t *testing.T) {
	client := NewClient()

	_, err := client.FetchPullRequest("octocat", "Hello-World", 999999)

	if err == nil {
		t.Error("Expected error for non-existent PR, got nil")
	}
}

// TestFetchNonExistentDiscussion 测试获取不存在的 Discussion
func TestFetchNonExistentDiscussion(t *testing.T) {
	client := NewClient()

	_, err := client.FetchDiscussion("community", "community", 999999)

	if err == nil {
		t.Error("Expected error for non-existent discussion, got nil")
	}
}
