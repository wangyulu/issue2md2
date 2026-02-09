# Technical Implementation Plan - issue2md

**Document ID:** PLAN-001
**Status:** Draft
**Created:** 2026-02-08
**Version:** 1.0

---

## 1. 技术上下文总结

### 1.1 技术选型

| 类别 | 技术栈 | 说明 |
|------|--------|------|
| **编程语言** | Go >= 1.25.1 | 高性能、简单、跨平台 |
| **Web 框架** | 标准库 `net/http` | 遵循"简单性原则"，不引入第三方框架 |
| **GitHub API** | `google/go-github` + GraphQL v4 | 官方推荐库，支持 REST 和 GraphQL |
| **Markdown 处理** | 标准库 + 简单字符串拼接 | 不引入第三方 Markdown 库，直接拼接字符串 |
| **数据存储** | 无 | 实时获取 API 数据，无需持久化 |
| **命令行解析** | 标准库 `flag` | Go 内置，满足基本需求 |

### 1.2 依赖管理

**核心外部依赖:**
- `github.com/google/go-github/v68` (最新稳定版)
- `github.com/shurcooL/githubv4` (GraphQL 客户端，用于获取 Reactions)

**版本策略:**
- 使用 Go Modules 管理依赖
- 遵循语义化版本控制
- 定期更新依赖，保持安全性

---

## 2. "合宪性"审查

本方案逐条对照 `constitution.md` 原则进行审查：

### 2.1 第一条：简单性原则

| 子条款 | 审查结果 | 说明 |
|--------|----------|------|
| 1.1 (YAGNI) | ✅ 合规 | 仅实现 spec.md 明确要求的功能（Issue/PR/Discussion 导出、Reactions、用户链接） |
| 1.2 (标准库优先) | ✅ 合规 | Web 框架使用 `net/http`，CLI 使用 `flag`，Markdown 使用字符串拼接 |
| 1.3 (反过度工程) | ✅ 合规 | 不引入复杂接口，不使用 ORM/DI 框架，直接使用具体类型 |

**设计决策:**
- 不定义 "Fetcher/Converter" 等抽象接口，直接使用具体函数
- 不使用第三方 Markdown 解析库，通过 `fmt.Sprintf` 直接拼接
- PR Review Comments 直接合并到评论列表，不创建单独的数据结构

### 2.2 第二条：测试先行铁律

| 子条款 | 审查结果 | 说明 |
|--------|----------|------|
| 2.1 (TDD循环) | ✅ 合规 | 所有功能从编写失败测试开始，遵循 Red-Green-Refactor |
| 2.2 (表格驱动) | ✅ 合规 | 单元测试优先采用表格驱动测试风格 |
| 2.3 (拒绝Mocks) | ✅ 合规 | 优先集成测试，使用真实 GitHub API |

**测试策略:**
- `internal/parser`: 纯逻辑，使用表格驱动测试覆盖所有 URL 格式
- `internal/config`: 纯逻辑，单元测试即可
- `internal/github`: 集成测试，使用真实的公开 Issue/PR/Discussion
- `internal/converter`: 表格驱动测试，构造测试数据验证 Markdown 输出
- `internal/cli`: 表格驱动测试，验证参数解析和错误处理

### 2.3 第三条：明确性原则

| 子条款 | 审查结果 | 说明 |
|--------|----------|------|
| 3.1 (错误处理) | ✅ 合规 | 所有错误显式处理，使用 `fmt.Errorf("...: %w", err)` 包装 |
| 3.2 (无全局变量) | ✅ 合规 | 配置通过结构体成员传递，不使用全局变量 |

**错误处理规范:**
```go
// 错误包装示例
if err != nil {
    return fmt.Errorf("failed to parse GitHub URL %q: %w", url, err)
}

// 错误检查示例
if resp.StatusCode == 404 {
    return fmt.Errorf("resource not found: %s", url)
}
```

**配置传递示例:**
```go
// ✅ 正确：通过参数传递
func NewClient(cfg *config.Config) *Client { ... }

// ❌ 错误：全局变量
var githubClient *github.Client
```

---

## 3. 项目结构细化

### 3.1 目录结构

```
issue2md2/
├── cmd/
│   ├── issue2md/
│   │   └── main.go           # CLI 入口
│   └── issue2mdweb/
│       └── main.go           # Web 入口（预留）
├── internal/
│   ├── github/
│   │   ├── client.go         # GitHub API 客户端封装
│   │   └── types.go         # GitHub 数据结构定义
│   ├── parser/
│   │   └── parser.go         # URL 解析逻辑
│   ├── converter/
│   │   ├── converter.go      # Markdown 生成逻辑
│   │   └── frontmatter.go    # Frontmatter 生成
│   ├── cli/
│   │   ├── flags.go          # 命令行标志解析
│   │   ├── help.go           # 帮助信息
│   │   └── error.go          # 错误输出
│   └── config/
│       └── config.go         # 配置管理（环境变量）
├── web/
│   ├── templates/            # Web 模板（预留）
│   └── static/               # 静态资源（预留）
├── specs/
│   └── 001-core-functionality/
│       ├── spec.md
│       ├── api-sketch.md
│       └── plan.md           # 本文档
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 3.2 包职责与依赖关系

#### 3.2.1 internal/config

**职责:**
- 从环境变量读取 `GITHUB_TOKEN`

**依赖:** 无（仅使用标准库 `os`）

**导出:**
```go
// GetGitHubToken 从环境变量读取 GitHub Token
func GetGitHubToken() string
```

#### 3.2.2 internal/parser

**职责:**
- 解析 GitHub URL
- 识别资源类型（Issue/PR/Discussion）
- 提取 owner/repo/number

**依赖:** 无（仅使用标准库 `regexp` 和 `net/url`）

**导出:**
```go
type ResourceType string
type Resource struct {
    Type     ResourceType
    Owner    string
    Repo     string
    Number   int
    Original string
}

func ParseURL(url string) (*Resource, error)
```

#### 3.2.3 internal/github

**职责:**
- 封装 GitHub API 调用
- 获取 Issue/PR/Discussion 完整数据
- 处理 API 认证和限流

**依赖:**
- `internal/config` (获取 token)
- `github.com/google/go-github/v68` (REST API)
- `github.com/shurcooL/githubv4` (GraphQL API，用于 Reactions)

**导出:**
```go
type Issue struct { ... }
type PullRequest struct { ... }
type Discussion struct { ... }
type Comment struct { ... }
type Reactions struct { ... }

type Client struct { ... }

func NewClient(cfg *config.Config) *Client
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
```

#### 3.2.4 internal/converter

**职责:**
- 将 Issue/PR/Discussion 数据转换为 Markdown
- 生成 YAML Frontmatter
- 渲染 Reactions 和用户链接

**依赖:** 无（仅使用标准库 `fmt`, `strings`, `time`）

**导出:**
```go
type Options struct {
    EnableReactions   bool
    EnableUserLinks   bool
}

func DefaultOptions() *Options
func ToMarkdown(issue *Issue, opts *Options) ([]byte, error)
func ToMarkdownPR(pr *PullRequest, opts *Options) ([]byte, error)
func ToMarkdownDiscussion(discussion *Discussion, opts *Options) ([]byte, error)
```

#### 3.2.5 internal/cli

**职责:**
- 解析命令行参数
- 输出帮助信息
- 输出错误信息到 stderr

**依赖:** 无（仅使用标准库 `flag`, `os`, `fmt`）

**导出:**
```go
type Flags struct {
    EnableReactions bool
    EnableUserLinks bool
}

type Args struct {
    URL        string
    OutputFile string
}

func ParseArgs(args []string) (*Flags, *Args, error)
func PrintHelp(w io.Writer)
func PrintError(w io.Writer, err error)
```

### 3.3 依赖图

```
cmd/issue2md/main.go
    ├── internal/cli          (参数解析)
    ├── internal/parser       (URL 解析)
    ├── internal/config       (配置管理)
    └── internal/github       (API 调用)
    └── internal/converter    (Markdown 生成)
```

**关键原则:**
- `internal/parser` 不依赖任何其他包（纯逻辑）
- `internal/converter` 不依赖 `internal/github`（数据解耦）
- 所有依赖通过参数显式注入

---

## 4. 核心数据结构

### 4.1 internal/parser 数据结构

```go
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
```

### 4.2 internal/github 数据结构

```go
// Issue GitHub Issue 数据
type Issue struct {
    Title     string
    Body      string
    Author    string
    AuthorURL string
    CreatedAt time.Time
    Status    string // open, closed
    URL       string
    Reactions *Reactions // 可能为 nil
    Comments  []Comment
}

// PullRequest GitHub Pull Request 数据
type PullRequest struct {
    Title     string
    Body      string
    Author    string
    AuthorURL string
    CreatedAt time.Time
    Status    string // open, closed, merged
    URL       string
    Reactions *Reactions // 可能为 nil
    Comments  []Comment // 包含普通评论和 Review Comments
}

// Discussion GitHub Discussion 数据
type Discussion struct {
    Title     string
    Body      string
    Author    string
    AuthorURL string
    CreatedAt time.Time
    Status    string // open, closed
    URL       string
    Reactions *Reactions // 可能为 nil
    Comments  []Comment
}

// Comment 评论数据
type Comment struct {
    Author     string
    AuthorURL  string
    Body       string
    CreatedAt  time.Time
    Reactions  *Reactions // 可能为 nil
    IsAnswer   bool       // Discussion 特有，标记是否为 Answer
}

// Reactions 反应统计
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

### 4.3 internal/converter 数据结构

```go
// Options Markdown 转换选项
type Options struct {
    EnableReactions bool // 是否启用 Reactions 显示
    EnableUserLinks bool // 是否将用户名渲染为链接
}
```

### 4.4 internal/cli 数据结构

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

---

## 5. 接口设计

### 5.1 internal/parser 接口

```go
// ParseURL 解析 GitHub URL 并返回 Resource 信息
// 如果 URL 格式无效或不是支持的类型，返回错误
//
// 支持的 URL 格式:
//   - https://github.com/owner/repo/issues/{number}
//   - https://github.com/owner/repo/pull/{number}
//   - https://github.com/owner/repo/discussions/{number}
//
// 示例:
//   res, err := parser.ParseURL("https://github.com/owner/repo/issues/123")
//   if err != nil {
//       return err
//   }
//   fmt.Printf("Type: %s, Owner: %s, Repo: %s, Number: %d\n",
//       res.Type, res.Owner, res.Repo, res.Number)
func ParseURL(url string) (*Resource, error)
```

### 5.2 internal/config 接口

```go
// GetGitHubToken 从环境变量读取 GitHub Token
// 如果环境变量未设置，返回空字符串
//
// 环境变量: GITHUB_TOKEN
func GetGitHubToken() string
```

### 5.3 internal/github 接口

```go
// NewClient 创建 GitHub API 客户端
// cfg 从 config 包获取，包含 token 信息
func NewClient(cfg *config.Config) *Client

// FetchIssue 获取指定 Issue 的完整数据
// 包括主楼、所有评论和 Reactions（如果有）
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)

// FetchPullRequest 获取指定 Pull Request 的完整数据
// 包括描述、所有评论、Review Comments 和 Reactions（如果有）
// 注意：不包含 diff 信息和 commits 历史
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)

// FetchDiscussion 获取指定 Discussion 的完整数据
// 包括主楼、所有评论和 Reactions（如果有）
// Answer 评论会标记 IsAnswer = true
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
```

### 5.4 internal/converter 接口

```go
// DefaultOptions 返回默认转换选项
func DefaultOptions() *Options

// ToMarkdown 将 Issue 转换为 Markdown 字符串
// 输出格式符合 spec.md 中定义的标准
func ToMarkdown(issue *Issue, opts *Options) ([]byte, error)

// ToMarkdownPR 将 PullRequest 转换为 Markdown 字符串
// 输出格式符合 spec.md 中定义的标准
func ToMarkdownPR(pr *PullRequest, opts *Options) ([]byte, error)

// ToMarkdownDiscussion 将 Discussion 转换为 Markdown 字符串
// 输出格式符合 spec.md 中定义的标准
// Answer 评论会使用 ✅ 标记
func ToMarkdownDiscussion(discussion *Discussion, opts *Options) ([]byte, error)
```

### 5.5 internal/cli 接口

```go
// ParseArgs 解析命令行参数
// args 通常是 os.Args[1:]
// 返回 Flags 和 Args，如果解析失败返回错误
//
// 支持的标志:
//   -enable-reactions: 启用 Reactions 显示
//   -enable-user-links: 启用用户链接
//
// 参数:
//   url: 必需，GitHub URL
//   output_file: 可选，输出文件路径
//
// 示例:
//   flags, args, err := cli.ParseArgs(os.Args[1:])
//   if err != nil {
//       log.Fatal(err)
//   }
func ParseArgs(args []string) (*Flags, *Args, error)

// PrintHelp 打印使用帮助信息到指定的 io.Writer
func PrintHelp(w io.Writer)

// PrintError 打印错误信息到指定的 io.Writer
func PrintError(w io.Writer, err error)
```

---

## 6. 实现顺序

### Phase 1: 基础设施 (无外部依赖)

**优先级: P0**

1. `internal/config` - 环境变量读取
   - 实现简单，无依赖
   - 为后续包提供配置基础

2. `internal/parser` - URL 解析
   - 纯逻辑，无外部依赖
   - 表格驱动测试覆盖所有 URL 格式

3. `internal/cli` - 命令行解析
   - 纯逻辑，无外部依赖
   - 表格驱动测试覆盖参数组合

### Phase 2: 核心功能 (GitHub API 集成)

**优先级: P0**

4. `internal/github` - GitHub API 客户端
   - 集成 `google/go-github` 和 `githubv4`
   - 实现 Issue/PR/Discussion 获取
   - 集成测试使用真实公开数据

5. `internal/converter` - Markdown 生成
   - 纯逻辑，无外部依赖
   - 实现三种资源类型的转换
   - 表格驱动测试验证输出格式

### Phase 3: 主程序集成

**优先级: P0**

6. `cmd/issue2md/main.go` - CLI 入口
   - 协调所有包
   - 错误处理和输出
   - 集成测试覆盖所有验收标准

### Phase 4: 文档与发布

**优先级: P1**

7. `README.md` - 用户文档
8. `Makefile` - 构建脚本
9. `go.mod` - 依赖管理

---

## 7. 测试策略

### 7.1 单元测试

| 包 | 测试方法 | 测试覆盖率目标 |
|----|----------|----------------|
| `internal/config` | 表格驱动 | 100% |
| `internal/parser` | 表格驱动 | 100% |
| `internal/cli` | 表格驱动 | 100% |
| `internal/converter` | 表格驱动 | 95%+ |
| `internal/github` | 集成测试 | 90%+ |

### 7.2 集成测试

**测试数据:**
- 使用真实的公开 GitHub Issue/PR/Discussion
- 确保 URL 长期有效（选择官方仓库）

**测试覆盖:**
- 所有资源类型（Issue/PR/Discussion）
- 所有错误场景（404, 401, 429）
- 所有 Flag 组合

### 7.3 端到端测试

**测试方法:**
- 使用 `os/exec` 运行编译后的二进制文件
- 验证 stdout 输出符合预期
- 验证文件写入正确
- 验证错误码和错误信息

---

## 8. 错误处理规范

### 8.1 错误类型定义

```go
// 错误常量
const (
    ErrInvalidURL      = "invalid GitHub URL"
    ErrResourceNotFound = "resource not found"
    ErrAuthFailed     = "authentication failed"
    ErrRateLimit      = "rate limit exceeded"
)
```

### 8.2 错误包装模式

```go
// 标准错误包装
if err != nil {
    return fmt.Errorf("failed to fetch issue: %w", err)
}

// 带上下文的错误包装
if err != nil {
    return fmt.Errorf("failed to fetch issue %d from %s/%s: %w", number, owner, repo, err)
}
```

### 8.3 错误响应码

| 场景 | 退出码 | 输出位置 |
|------|--------|----------|
| 参数错误 | 1 | stderr |
| URL 无效 | 1 | stderr |
| 资源不存在 | 1 | stderr |
| 认证失败 | 1 | stderr |
| API 限流 | 1 | stderr |
| 网络错误 | 1 | stderr |
| 成功 | 0 | stdout 或文件 |

---

## 9. 性能考虑

### 9.1 API 调用优化

- 使用 HTTP 连接池（`google/go-github` 默认支持）
- 遵循 GitHub API 限流策略
- 合理设置超时时间（30 秒）

### 9.2 内存管理

- 使用流式输出（不一次性加载所有数据到内存）
- 及时关闭 HTTP 响应体

### 9.3 并发策略

- MVP 阶段不使用并发，保持简单
- 未来可考虑批量处理时使用 goroutine

---

## 10. 安全考虑

### 10.1 Token 安全

- 仅通过环境变量传递 token
- 不在日志中输出 token
- 不在错误信息中暴露 token

### 10.2 输入验证

- 验证 URL 格式，防止注入攻击
- 验证文件路径，防止路径遍历

### 10.3 依赖安全

- 定期更新 `go.mod` 依赖
- 使用 `go vet` 和 `gosec` 进行静态分析

---

## 11. 文档要求

### 11.1 代码注释

- 所有导出函数必须有注释
- 复杂逻辑必须有行内注释
- 使用 `godoc` 风格

### 11.2 用户文档

- README.md: 安装、使用、示例
- CONTRIBUTING.md: 贡献指南
- LICENSE: 开源协议

### 11.3 开发文档

- specs/001-core-functionality/spec.md: 需求规范
- specs/001-core-functionality/plan.md: 技术方案（本文档）
- specs/001-core-functionality/api-sketch.md: API 草稿

---

## 12. 发布计划

### 12.1 版本管理

- 遵循语义化版本控制 (SemVer)
- v0.1.0: MVP 版本
- v0.2.0: Web 界面
- v1.0.0: 稳定版本

### 12.2 分发方式

- GitHub Releases
- Homebrew Formula (macOS/Linux)
- Go install: `go install github.com/owner/issue2md@latest`

---

## 13. 风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| GitHub API 变更 | 高 | 中 | 使用稳定版本 API，及时更新依赖 |
| API 限流 | 中 | 高 | 引导用户使用 token，增加重试逻辑 |
| 跨平台兼容性 | 中 | 低 | 使用 CGO 编译选项，测试多平台 |
| 依赖库漏洞 | 高 | 低 | 定期更新依赖，使用 `go mod tidy` |

---

## 14. 审批与修订

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| 1.0 | 2026-02-08 | Architect | 初始版本 |

---

**附录 A: 参考资料**

- [Go 官方文档](https://golang.org/doc/)
- [Google Go-GitHub 库](https://github.com/google/go-github)
- [GitHub GraphQL API](https://docs.github.com/en/graphql)
- [Conventional Commits](https://www.conventionalcommits.org/)
