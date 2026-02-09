package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shurcooL/githubv4"
)

// Client GitHub API 客户端
type Client struct {
	ghClient *githubv4.Client
}

// NewClient 创建 GitHub API 客户端
func NewClient() *Client {
	token := os.Getenv("GITHUB_TOKEN")
	var httpClient *http.Client

	if token != "" {
		httpClient = &http.Client{
			Transport: &authenticatedTransport{
				token:     token,
				transport: http.DefaultTransport,
			},
		}
	}

	return &Client{
		ghClient: githubv4.NewClient(httpClient),
	}
}

// authenticatedTransport 添加 Authorization header
type authenticatedTransport struct {
	token     string
	transport http.RoundTripper
}

func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.transport.RoundTrip(req)
}

// FetchIssue 获取指定 Issue 的完整数据
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error) {
	ctx := context.Background()

	// GraphQL 查询
	var q struct {
		Repository struct {
			Issue *struct {
				Title      string
				Body       *string
				Closed     bool
				CreatedAt  string
				URL        string
				Author     *struct {
					Login     string
					AvatarURL string
				}
				Reactions *struct {
					TotalCount int `graphql:"totalCount"`
				}
				Comments struct {
					Nodes []struct {
						Body      string
						CreatedAt string
						Author    *struct {
							Login     string
							AvatarURL string
						}
						Reactions *struct {
							TotalCount int `graphql:"totalCount"`
						}
					} `graphql:"nodes"`
				} `graphql:"comments(first: 100)"`
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repo),
		"number": githubv4.Int(number),
	}

	err := c.ghClient.Query(ctx, &q, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue: %w", err)
	}

	if q.Repository.Issue == nil {
		return nil, fmt.Errorf("resource not found: %s/%s/issues/%d", owner, repo, number)
	}

	issueData := q.Repository.Issue

	// 构建 Issue 对象
	issue := &Issue{
		Title:     issueData.Title,
		Body:      toString(issueData.Body),
		Author:    toLogin(issueData.Author),
		AuthorURL: toAvatarURL(issueData.Author),
		CreatedAt: toTime(issueData.CreatedAt),
		Status:    toStatus(issueData.Closed),
		URL:       issueData.URL,
	}

	// Comments
	for _, node := range issueData.Comments.Nodes {
		comment := Comment{
			Body:      node.Body,
			CreatedAt: toTime(node.CreatedAt),
			Author:    toLogin(node.Author),
			AuthorURL: toAvatarURL(node.Author),
		}

		issue.Comments = append(issue.Comments, comment)
	}

	return issue, nil
}

// FetchPullRequest 获取指定 Pull Request 的完整数据
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error) {
	ctx := context.Background()

	// GraphQL 查询
	var q struct {
		Repository struct {
			PullRequest *struct {
				Title      string
				Body       *string
				State      string
				Merged     bool
				CreatedAt  string
				URL        string
				Author     *struct {
					Login     string
					AvatarURL string
				}
				Reactions *struct {
					TotalCount int `graphql:"totalCount"`
				}
				Comments struct {
					Nodes []struct {
						Body      string
						CreatedAt string
						Author    *struct {
							Login     string
							AvatarURL string
						}
						Reactions *struct {
							TotalCount int `graphql:"totalCount"`
						}
					} `graphql:"nodes"`
				} `graphql:"comments(first: 100)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repo),
		"number": githubv4.Int(number),
	}

	err := c.ghClient.Query(ctx, &q, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull request: %w", err)
	}

	if q.Repository.PullRequest == nil {
		return nil, fmt.Errorf("resource not found: %s/%s/pull/%d", owner, repo, number)
	}

	prData := q.Repository.PullRequest

	// 构建 PullRequest 对象
	pr := &PullRequest{
		Title:     prData.Title,
		Body:      toString(prData.Body),
		Author:    toLogin(prData.Author),
		AuthorURL: toAvatarURL(prData.Author),
		CreatedAt: toTime(prData.CreatedAt),
		Status:    toPRStatus(prData.State, prData.Merged),
		URL:       prData.URL,
	}

	// Comments
	for _, node := range prData.Comments.Nodes {
		comment := Comment{
			Body:      node.Body,
			CreatedAt: toTime(node.CreatedAt),
			Author:    toLogin(node.Author),
			AuthorURL: toAvatarURL(node.Author),
		}

		pr.Comments = append(pr.Comments, comment)
	}

	return pr, nil
}

// FetchDiscussion 获取指定 Discussion 的完整数据
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error) {
	ctx := context.Background()

	// GraphQL 查询
	var q struct {
		Repository struct {
			Discussion *struct {
				Title      string
				Body       string
				Closed     bool
				CreatedAt  string
				URL        string
				Author     *struct {
					Login     string
					AvatarURL string
				}
				Reactions *struct {
					TotalCount int `graphql:"totalCount"`
				}
				Comments struct {
					Nodes []struct {
						Body      string
						CreatedAt string
						Author    *struct {
							Login     string
							AvatarURL string
						}
						IsAnswer bool `graphql:"isAnswer"`
						Reactions *struct {
							TotalCount int `graphql:"totalCount"`
						}
					} `graphql:"nodes"`
				} `graphql:"comments(first: 100)"`
			} `graphql:"discussion(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repo),
		"number": githubv4.Int(number),
	}

	err := c.ghClient.Query(ctx, &q, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discussion: %w", err)
	}

	if q.Repository.Discussion == nil {
		return nil, fmt.Errorf("resource not found: %s/%s/discussions/%d", owner, repo, number)
	}

	discussionData := q.Repository.Discussion

	// 构建 Discussion 对象
	discussion := &Discussion{
		Title:     discussionData.Title,
		Body:      discussionData.Body,
		Author:    toLogin(discussionData.Author),
		AuthorURL: toAvatarURL(discussionData.Author),
		CreatedAt: toTime(discussionData.CreatedAt),
		Status:    toStatus(discussionData.Closed),
		URL:       discussionData.URL,
	}

	// Comments
	for _, node := range discussionData.Comments.Nodes {
		comment := Comment{
			Body:      node.Body,
			CreatedAt: toTime(node.CreatedAt),
			Author:    toLogin(node.Author),
			AuthorURL: toAvatarURL(node.Author),
			IsAnswer:  node.IsAnswer,
		}

		discussion.Comments = append(discussion.Comments, comment)
	}

	return discussion, nil
}

// 辅助函数

func toString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

func toLogin(author *struct {
	Login     string
	AvatarURL string
}) string {
	if author == nil {
		return ""
	}
	return author.Login
}

func toAvatarURL(author *struct {
	Login     string
	AvatarURL string
}) string {
	if author == nil {
		return ""
	}
	return author.AvatarURL
}

func toStatus(closed bool) string {
	if closed {
		return "closed"
	}
	return "open"
}

func toPRStatus(state string, merged bool) string {
	if merged {
		return "merged"
	}
	if state == "CLOSED" {
		return "closed"
	}
	return "open"
}
