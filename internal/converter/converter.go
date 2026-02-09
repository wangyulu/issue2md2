package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/wangyulu/issue2md2/internal/github"
)

// Reactions 反应统计（使用 github 包的类型）
type Reactions = github.Reactions

// Options Markdown 转换选项
type Options struct {
	EnableReactions bool // 是否启用 Reactions 显示
	EnableUserLinks bool // 是否将用户名渲染为链接
}

// DefaultOptions 返回默认转换选项
func DefaultOptions() *Options {
	return &Options{
		EnableReactions: false,
		EnableUserLinks: false,
	}
}

// renderReactions 渲染 Reactions 统计
func renderReactions(r *github.Reactions) string {
	if r == nil {
		return ""
	}

	var parts []string

	if r.ThumbsUp > 0 {
		parts = append(parts, fmt.Sprintf("👍 %d", r.ThumbsUp))
	}
	if r.ThumbsDown > 0 {
		parts = append(parts, fmt.Sprintf("👎 %d", r.ThumbsDown))
	}
	if r.Laugh > 0 {
		parts = append(parts, fmt.Sprintf("😄 %d", r.Laugh))
	}
	if r.Hooray > 0 {
		parts = append(parts, fmt.Sprintf("🎉 %d", r.Hooray))
	}
	if r.Confused > 0 {
		parts = append(parts, fmt.Sprintf("😕 %d", r.Confused))
	}
	if r.Heart > 0 {
		parts = append(parts, fmt.Sprintf("❤️ %d", r.Heart))
	}
	if r.Rocket > 0 {
		parts = append(parts, fmt.Sprintf("🚀 %d", r.Rocket))
	}
	if r.Eyes > 0 {
		parts = append(parts, fmt.Sprintf("👀 %d", r.Eyes))
	}

	return strings.Join(parts, " ")
}

// renderUser 渲染用户信息（根据选项决定是否添加链接）
func renderUser(username, userURL string, enableLinks bool) string {
	if enableLinks {
		return fmt.Sprintf("[@%s](%s)", username, userURL)
	}
	return fmt.Sprintf("@%s", username)
}

// ToMarkdown 将 Issue 转换为 Markdown 字符串
func ToMarkdown(issue *github.Issue, opts *Options) ([]byte, error) {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(generateFrontmatter(issue.Title, issue.URL, issue.Author, issue.AuthorURL, issue.CreatedAt, issue.Status, "issue"))
	sb.WriteString("\n")

	// 标题和正文
	sb.WriteString(fmt.Sprintf("# %s\n", issue.Title))
	sb.WriteString("\n")
	if issue.Body != "" {
		sb.WriteString(issue.Body)
		sb.WriteString("\n")
	}

	// Reactions
	if opts.EnableReactions && issue.Reactions != nil {
		sb.WriteString("## Reactions\n")
		sb.WriteString("\n")
		sb.WriteString(renderReactions(issue.Reactions))
		sb.WriteString("\n")
	}

	// Comments
	if len(issue.Comments) > 0 {
		sb.WriteString("---\n")
		sb.WriteString("\n")
		sb.WriteString("## Comments\n")
		sb.WriteString("\n")

		for _, comment := range issue.Comments {
			sb.WriteString(fmt.Sprintf("### %s commented at %s\n",
				renderUser(comment.Author, comment.AuthorURL, opts.EnableUserLinks),
				comment.CreatedAt.UTC().Format(time.RFC3339)))
			sb.WriteString("\n")
			if comment.Body != "" {
				sb.WriteString(comment.Body)
				sb.WriteString("\n")
			}
			if opts.EnableReactions && comment.Reactions != nil {
				sb.WriteString(renderReactions(comment.Reactions))
				sb.WriteString("\n")
			}
		}
	}

	return []byte(sb.String()), nil
}

// ToMarkdownPR 将 PullRequest 转换为 Markdown 字符串
func ToMarkdownPR(pr *github.PullRequest, opts *Options) ([]byte, error) {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(generateFrontmatter(pr.Title, pr.URL, pr.Author, pr.AuthorURL, pr.CreatedAt, pr.Status, "pull_request"))
	sb.WriteString("\n")

	// 标题和描述
	sb.WriteString(fmt.Sprintf("# %s\n", pr.Title))
	sb.WriteString("\n")
	if pr.Body != "" {
		sb.WriteString(pr.Body)
		sb.WriteString("\n")
	}

	// Reactions
	if opts.EnableReactions && pr.Reactions != nil {
		sb.WriteString("## Reactions\n")
		sb.WriteString("\n")
		sb.WriteString(renderReactions(pr.Reactions))
		sb.WriteString("\n")
	}

	// Comments
	if len(pr.Comments) > 0 {
		sb.WriteString("---\n")
		sb.WriteString("\n")
		sb.WriteString("## Comments\n")
		sb.WriteString("\n")

		for _, comment := range pr.Comments {
			sb.WriteString(fmt.Sprintf("### %s commented at %s\n",
				renderUser(comment.Author, comment.AuthorURL, opts.EnableUserLinks),
				comment.CreatedAt.UTC().Format(time.RFC3339)))
			sb.WriteString("\n")
			if comment.Body != "" {
				sb.WriteString(comment.Body)
				sb.WriteString("\n")
			}
			if opts.EnableReactions && comment.Reactions != nil {
				sb.WriteString(renderReactions(comment.Reactions))
				sb.WriteString("\n")
			}
		}
	}

	return []byte(sb.String()), nil
}

// ToMarkdownDiscussion 将 Discussion 转换为 Markdown 字符串
func ToMarkdownDiscussion(discussion *github.Discussion, opts *Options) ([]byte, error) {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(generateFrontmatter(discussion.Title, discussion.URL, discussion.Author, discussion.AuthorURL, discussion.CreatedAt, discussion.Status, "discussion"))
	sb.WriteString("\n")

	// 标题和正文
	sb.WriteString(fmt.Sprintf("# %s\n", discussion.Title))
	sb.WriteString("\n")
	if discussion.Body != "" {
		sb.WriteString(discussion.Body)
		sb.WriteString("\n")
	}

	// Reactions
	if opts.EnableReactions && discussion.Reactions != nil {
		sb.WriteString("## Reactions\n")
		sb.WriteString("\n")
		sb.WriteString(renderReactions(discussion.Reactions))
		sb.WriteString("\n")
	}

	// Comments
	if len(discussion.Comments) > 0 {
		sb.WriteString("---\n")
		sb.WriteString("\n")
		sb.WriteString("## Comments\n")
		sb.WriteString("\n")

		for _, comment := range discussion.Comments {
			sb.WriteString(fmt.Sprintf("### %s commented at %s\n",
				renderUser(comment.Author, comment.AuthorURL, opts.EnableUserLinks),
				comment.CreatedAt.UTC().Format(time.RFC3339)))
			sb.WriteString("\n")
			if comment.Body != "" {
				sb.WriteString(comment.Body)
				sb.WriteString("\n")
			}
			// Answer 标记
			if comment.IsAnswer {
				sb.WriteString("✅ **Answer**")
				sb.WriteString("\n")
			}
			if opts.EnableReactions && comment.Reactions != nil {
				sb.WriteString(renderReactions(comment.Reactions))
				sb.WriteString("\n")
			}
		}
	}

	return []byte(sb.String()), nil
}
