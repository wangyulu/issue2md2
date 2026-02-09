package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeIssue       ResourceType = "issue"
	ResourceTypePullRequest ResourceType = "pull_request"
	ResourceTypeDiscussion  ResourceType = "discussion"
)

// Resource 解析后的 GitHub 资源信息
type Resource struct {
	Type     ResourceType
	Owner    string
	Repo     string
	Number   int
	Original string // 原始 URL
}

// GitHub URL 正则表达式
var (
	issueURLPattern       = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/issues/(\d+)`)
	pullRequestURLPattern = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/pull/(\d+)`)
	discussionURLPattern  = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/discussions/(\d+)`)
)

// ParseURL 解析 GitHub URL 并返回 Resource 信息
// 如果 URL 格式无效或不是支持的类型，返回错误
//
// 支持的 URL 格式:
//   - https://github.com/owner/repo/issues/{number}
//   - https://github.com/owner/repo/pull/{number}
//   - https://github.com/owner/repo/discussions/{number}
func ParseURL(url string) (*Resource, error) {
	// 尝试匹配 Issue URL
	if matches := issueURLPattern.FindStringSubmatch(url); matches != nil {
		return parseMatches(url, matches, ResourceTypeIssue)
	}

	// 尝试匹配 Pull Request URL
	if matches := pullRequestURLPattern.FindStringSubmatch(url); matches != nil {
		return parseMatches(url, matches, ResourceTypePullRequest)
	}

	// 尝试匹配 Discussion URL
	if matches := discussionURLPattern.FindStringSubmatch(url); matches != nil {
		return parseMatches(url, matches, ResourceTypeDiscussion)
	}

	// 检查是否为 GitHub URL 但不是支持的资源类型
	if strings.HasPrefix(url, "https://github.com/") {
		return nil, fmt.Errorf("unsupported resource type: %s", url)
	}

	return nil, fmt.Errorf("invalid GitHub URL: %s", url)
}

// parseMatches 解析正则匹配结果
func parseMatches(url string, matches []string, typ ResourceType) (*Resource, error) {
	owner := matches[1]
	repo := matches[2]
	numberStr := matches[3]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", numberStr)
	}

	return &Resource{
		Type:     typ,
		Owner:    owner,
		Repo:     repo,
		Number:   number,
		Original: url,
	}, nil
}
