# API Sketch - Core Package Interfaces

**Purpose:** 本文档定义 `internal/converter` 和 `internal/github` 包对外暴露的主要接口，作为后续开发的参考。

**Version:** 1.0
**Created:** 2026-02-08

---

## 1. internal/parser

URL 解析与类型识别包。

### 1.1 数据结构

```go
// Resource 表示解析后的 GitHub 资源
type Resource struct {
    Type     ResourceType // issue, pull_request, discussion
    Owner    string
    Repo     string
    Number   int
    Original string // 原始 URL
}

// ResourceType 资源类型
type ResourceType string

const (
    ResourceTypeIssue       ResourceType = "issue"
    ResourceTypePullRequest ResourceType = "pull_request"
    ResourceTypeDiscussion  ResourceType = "discussion"
)
```

### 1.2 函数接口

```go
// ParseURL 解析 GitHub URL 并返回 Resource 信息
// 如果 URL 格式无效或不是支持的类型，返回错误
func ParseURL(url string) (*Resource, error)
```

---

## 2. internal/github

GitHub API 交互包。负责从 GitHub 获取 Issue/PR/Discussion 数据。

### 2.1 配置

```go
// Config GitHub API 配置
type Config struct {
    Token string // 从环境变量 GITHUB_TOKEN 读取，可为空
}

// NewConfig 创建配置，自动从环境变量读取 token
func NewConfig() *Config
```

### 2.2 数据结构

```go
// Issue 表示 GitHub Issue 数据
type Issue struct {
    Title      string
    Body       string
    Author     string
    AuthorURL  string
    CreatedAt  time.Time
    Status     string // open, closed
    URL        string
    Reactions  *Reactions // 可能为 nil
    Comments   []Comment
}

// PullRequest 表示 GitHub Pull Request 数据
type PullRequest struct {
    Title      string
    Body       string
    Author     string
    AuthorURL  string
    CreatedAt  time.Time
    Status     string // open, closed, merged
    URL        string
    Reactions  *Reactions // 可能为 nil
    Comments   []Comment  // 包含普通评论和 Review Comments
}

// Discussion 表示 GitHub Discussion 数据
type Discussion struct {
    Title      string
    Body       string
    Author     string
    AuthorURL  string
    CreatedAt  time.Time
    Status     string // open, closed
    URL        string
    Reactions  *Reactions // 可能为 nil
    Comments   []Comment
}

// Comment 表示评论数据
type Comment struct {
    Author     string
    AuthorURL  string
    Body       string
    CreatedAt  time.Time
    Reactions  *Reactions // 可能为 nil
    IsAnswer   bool       // Discussion 特有
}

// Reactions 表示反应统计
type Reactions struct {
    ThumbsUp   int
    ThumbsDown int
    Laugh      int
    Hooray     int
    Confused   int
    Heart      int
    Rocket     int
    Eyes       int
}
```

### 2.3 函数接口

```go
// Client GitHub API 客户端
type Client struct {
    config *Config
    // 内部包含 github.Client 实例
}

// NewClient 创建 GitHub API 客户端
func NewClient(cfg *Config) *Client

// FetchIssue 获取指定 Issue 的完整数据
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)

// FetchPullRequest 获取指定 Pull Request 的完整数据
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)

// FetchDiscussion 获取指定 Discussion 的完整数据
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
```

---

## 3. internal/converter

数据转换为 Markdown 包。负责将获取的数据转换为标准 Markdown 格式。

### 3.1 配置选项

```go
// Options Markdown 转换选项
type Options struct {
    EnableReactions   bool  // 是否启用 Reactions 显示
    EnableUserLinks   bool  // 是否将用户名渲染为链接
}

// DefaultOptions 返回默认选项
func DefaultOptions() *Options
```

### 3.2 函数接口

```go
// ToMarkdown 将 Issue 转换为 Markdown 字符串
func ToMarkdown(issue *Issue, opts *Options) ([]byte, error)

// ToMarkdownPR 将 PullRequest 转换为 Markdown 字符串
func ToMarkdownPR(pr *PullRequest, opts *Options) ([]byte, error)

// ToMarkdownDiscussion 将 Discussion 转换为 Markdown 字符串
func ToMarkdownDiscussion(discussion *Discussion, opts *Options) ([]byte, error)
```

### 3.3 内部辅助函数（不对外暴露）

```go
// generateFrontmatter 生成 YAML Frontmatter
func generateFrontmatter(title, url, author, authorURL string, createdAt time.Time, status, typ string) string

// renderReactions 渲染 Reactions 统计
func renderReactions(r *Reactions) string

// renderUser 渲染用户信息（根据选项决定是否添加链接）
func renderUser(username, userURL string, enableLinks bool) string
```

---

## 4. internal/cli

命令行接口包。负责参数解析和错误输出。

### 4.1 数据结构

```go
// Flags 命令行标志
type Flags struct {
    EnableReactions bool
    EnableUserLinks bool
}

// Args 命令行参数
type Args struct {
    URL        string // 必需
    OutputFile string // 可选，为空表示输出到 stdout
}
```

### 4.2 函数接口

```go
// ParseArgs 解析命令行参数
func ParseArgs(os.Args []string) (*Flags, *Args, error)

// PrintHelp 打印使用帮助
func PrintHelp(w io.Writer)

// PrintError 打印错误信息到 stderr
func PrintError(w io.Writer, err error)
```

---

## 5. internal/config

配置管理包。

### 5.1 函数接口

```go
// GetGitHubToken 从环境变量读取 GitHub Token
func GetGitHubToken() string

// SetGitHubToken 设置 GitHub Token（用于测试）
func SetGitHubToken(token string)
```

---

## 6. 数据流转示意

```
CLI Input (URL)
    ↓
internal/parser.ParseURL()
    ↓
Resource {Type, Owner, Repo, Number}
    ↓
internal/github.Client.Fetch*()
    ↓
Issue / PullRequest / Discussion
    ↓
internal/converter.ToMarkdown*()
    ↓
[]byte (Markdown)
    ↓
Output (File / Stdout)
```

---

## 7. 设计原则说明

### 7.1 简单性
- 避免过度抽象，每个包职责单一
- 不使用不必要的接口，直接使用具体类型
- 错误处理显式，使用 `%w` 包装

### 7.2 解耦
- `parser` 不依赖任何其他包
- `github` 只依赖 GitHub API 客户端库
- `converter` 只依赖数据结构，不依赖 `github` 包
- `cli` 是胶水层，协调其他包

### 7.3 测试友好
- 所有函数都是纯函数（除了 Client，但可以通过接口模拟）
- 数据结构简单，易于构造测试数据
- 依赖通过参数注入，便于 mock

---

## 8. 待定事项

以下细节需要在实现时进一步确定：

1. **Reactions 渲染格式**: 具体的 emoji + 数量展示方式
2. **Discussion Answer 标记**: 使用引用块还是 emoji (✅)
3. **时间格式**: Frontmatter 中时间使用什么格式（ISO 8601）
4. **错误类型**: 是否定义自定义 error 类型
5. **GitHub API 客户端**: 是否使用 `google/go-github` 库，还是自己实现

---

## 9. 示例：典型调用流程

```go
// 1. 解析 URL
res, err := parser.ParseURL("https://github.com/owner/repo/issues/123")
if err != nil {
    return err
}

// 2. 创建 GitHub 客户端
cfg := config.NewConfig()
client := github.NewClient(cfg)

// 3. 获取数据
var data interface{}
switch res.Type {
case parser.ResourceTypeIssue:
    data, err = client.FetchIssue(res.Owner, res.Repo, res.Number)
case parser.ResourceTypePullRequest:
    data, err = client.FetchPullRequest(res.Owner, res.Repo, res.Number)
case parser.ResourceTypeDiscussion:
    data, err = client.FetchDiscussion(res.Owner, res.Repo, res.Number)
}
if err != nil {
    return err
}

// 4. 转换为 Markdown
opts := &converter.Options{
    EnableReactions: true,
    EnableUserLinks: true,
}
var markdown []byte
switch v := data.(type) {
case *github.Issue:
    markdown, err = converter.ToMarkdown(v, opts)
case *github.PullRequest:
    markdown, err = converter.ToMarkdownPR(v, opts)
case *github.Discussion:
    markdown, err = converter.ToMarkdownDiscussion(v, opts)
}

// 5. 输出
fmt.Println(string(markdown))
```
