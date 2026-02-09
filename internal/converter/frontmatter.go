package converter

import (
	"fmt"
	"strings"
	"time"
)

// generateFrontmatter 生成 YAML Frontmatter
func generateFrontmatter(title, url, author, authorURL string, createdAt time.Time, status, typ string) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", quoteYAML(title)))
	sb.WriteString(fmt.Sprintf("url: %s\n", quoteYAML(url)))
	sb.WriteString(fmt.Sprintf("author: %s\n", quoteYAML(author)))
	sb.WriteString(fmt.Sprintf("author_url: %s\n", quoteYAML(authorURL)))
	sb.WriteString(fmt.Sprintf("created_at: %q\n", createdAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("status: %s\n", quoteYAML(status)))
	sb.WriteString(fmt.Sprintf("type: %s\n", quoteYAML(typ)))
	sb.WriteString("---\n")

	return sb.String()
}

// quoteYAML 为 YAML 字符串添加引号
func quoteYAML(s string) string {
	// 如果字符串包含单引号，使用单引号包裹并转义单引号
	if strings.Contains(s, "'") {
		return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
	}
	// 否则使用双引号包裹（默认）
	return fmt.Sprintf("%q", s)
}
