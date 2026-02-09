# Detailed Implementation Tasks - issue2md

**Document ID:** TASKS-001
**Status:** Draft
**Created:** 2026-02-08
**Version:** 1.0

**Legend:**
- `[P]` - 可以并行执行的任务（无依赖关系）
- `→` - 依赖关系（后续任务依赖于前置任务）

---

## Phase 1: Foundation (数据结构定义)

本阶段定义所有核心数据结构和无外部依赖的基础包。

---

### 1.1 internal/config (环境变量读取)

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T001 | [TEST] 创建 config 包测试文件 | `internal/config/config_test.go` | 表格驱动测试 `GetGitHubToken()` 函数，测试环境变量已设置和未设置两种情况 |
| T002 | → [IMP] 实现 config 包 | `internal/config/config.go` | 实现 `GetGitHubToken()` 函数，从环境变量 `GITHUB_TOKEN` 读取 token |

---

### 1.2 internal/parser (URL 解析逻辑)

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T003 | [P] [TEST] 创建 parser 包测试文件 | `internal/parser/parser_test.go` | 表格驱动测试 `ParseURL()` 函数，测试各种 URL 格式（Issue、PR、Discussion）和错误情况 |
| T004 | → [IMP] 实现 parser 包 | `internal/parser/parser.go` | 实现 `Resource` 和 `ResourceType` 类型，以及 `ParseURL()` 函数，支持解析 GitHub URL |

**T003 测试用例（必须覆盖）：**
| 输入 URL | 期望结果 |
|-----------|----------|
| `https://github.com/owner/repo/issues/123` | Type: issue, Owner: owner, Repo: repo, Number: 123 |
| `https://github.com/owner/repo/issues/123#issuecomment-456` | Type: issue, Owner: owner, Repo: repo, Number: 123 |
| `https://github.com/owner/repo/pull/42` | Type: pull_request, Owner: owner, Repo: repo, Number: 42 |
| `https://github.com/owner/repo/pull/42#discussion_r123` | Type: pull_request, Owner: owner, Repo: repo, Number: 42 |
| `https://github.com/owner/repo/discussions/7` | Type: discussion, Owner: owner, Repo: repo, Number: 7 |
| `https://github.com/owner/repo/discussions/7#discussioncomment-890` | Type: discussion, Owner: owner, Repo: repo, Number: 7 |
| `https://example.com/invalid` | Error: "invalid GitHub URL" |
| `https://github.com/owner/repo/tree/main` | Error: "unsupported resource type" |

---

### 1.3 internal/cli (命令行参数解析)

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T005 | [P] [TEST] 创建 cli 包测试文件 | `internal/cli/cli_test.go` | 表格驱动测试 `ParseArgs()` 函数，测试各种参数组合和错误情况 |
| T006 | → [IMP] 实现 cli 包 - 类型定义 | `internal/cli/flags.go` | 实现 `Flags` 和 `Args` 类型 |
| T007 | → [IMP] 实现 cli 包 - 参数解析 | `internal/cli/flags.go` | 实现 `ParseArgs()` 函数，解析命令行参数和标志 |
| T008 | [P] [IMP] 实现 cli 包 - 帮助信息 | `internal/cli/help.go` | 实现 `PrintHelp()` 函数，输出使用帮助 |
| T009 | [P] [IMP] 实现 cli 包 - 错误输出 | `internal/cli/error.go` | 实现 `PrintError()` 函数，输出错误信息到 stderr |

**T005 测试用例（必须覆盖）：**
| 输入参数 | 期望 Flags | 期望 Args | 期望结果 |
|----------|------------|-----------|----------|
| `["https://github.com/owner/repo/issues/123"]` | EnableReactions: false, EnableUserLinks: false | URL: https://github.com/owner/repo/issues/123, OutputFile: "" | 成功 |
| `["-enable-reactions", "https://github.com/owner/repo/issues/123"]` | EnableReactions: true, EnableUserLinks: false | URL: https://github.com/owner/repo/issues/123, OutputFile: "" | 成功 |
| `["-enable-user-links", "https://github.com/owner/repo/issues/123"]` | EnableReactions: false, EnableUserLinks: true | URL: https://github.com/owner/repo/issues/123, OutputFile: "" | 成功 |
| `["-enable-reactions", "-enable-user-links", "https://github.com/owner/repo/issues/123", "output.md"]` | EnableReactions: true, EnableUserLinks: true | URL: https://github.com/owner/repo/issues/123, OutputFile: "output.md" | 成功 |
| `[]` | - | - | Error: "missing required argument: url" |

---

### 1.4 internal/github (数据结构定义)

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T010 | [P] [IMP] 实现 github 包数据结构 | `internal/github/types.go` | 定义 `Issue`, `PullRequest`, `Discussion`, `Comment`, `Reactions` 类型 |

**数据结构要求（参考 plan.md 第 4.2 节）：**
- `Issue`: Title, Body, Author, AuthorURL, CreatedAt, Status, URL, Reactions, Comments
- `PullRequest`: Title, Body, Author, AuthorURL, CreatedAt, Status, URL, Reactions, Comments
- `Discussion`: Title, Body, Author, AuthorURL, CreatedAt, Status, URL, Reactions, Comments
- `Comment`: Author, AuthorURL, Body, CreatedAt, Reactions, IsAnswer
- `Reactions`: ThumbsUp, ThumbsDown, Laugh, Hooray, Confused, Heart, Rocket, Eyes

---

### 1.5 internal/converter (数据结构定义)

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T011 | [P] [IMP] 实现 converter 包类型定义 | `internal/converter/converter.go` | 定义 `Options` 类型，包含 `EnableReactions` 和 `EnableUserLinks` 字段 |

---

### 1.6 go.mod (依赖管理初始化)

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T012 | [P] [IMP] 初始化 Go 模块 | `go.mod` | 创建 `go.mod` 文件，定义模块路径 `github.com/issue2md/issue2md2` 和 Go 版本要求 |
| T013 | → [IMP] 添加 GitHub API 依赖 | `go.mod` | 添加 `github.com/google/go-github/v68` 和 `github.com/shurcooL/githubv4` 依赖 |

---

## Phase 2: GitHub Fetcher (API 交互逻辑，TDD)

本阶段实现 GitHub API 客户端，遵循测试先行原则。

---

### 2.1 GitHub Client 实现

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T014 | [TEST] 创建 GitHub client 测试文件 | `internal/github/client_test.go` | 集成测试，使用真实的公开 Issue/PR/Discussion 数据 |
| T015 | → [IMP] 实现 GitHub Client - NewClient | `internal/github/client.go` | 实现 `NewClient(cfg *config.Config) *Client` 函数，初始化 GitHub API 客户端 |
| T016 | → [IMP] 实现 GitHub Client - FetchIssue | `internal/github/client.go` | 实现 `FetchIssue(owner, repo string, number int) (*Issue, error)` 函数，获取 Issue 数据 |
| T017 | → [IMP] 实现 GitHub Client - FetchPullRequest | `internal/github/client.go` | 实现 `FetchPullRequest(owner, repo string, number int) (*PullRequest, error)` 函数，获取 PR 数据 |
| T018 | → [IMP] 实现 GitHub Client - FetchDiscussion | `internal/github/client.go` | 实现 `FetchDiscussion(owner, repo string, number int) (*Discussion, error)` 函数，获取 Discussion 数据 |

**T014 测试用例（必须覆盖）：**
| 测试类型 | 测试 URL | 期望结果 |
|----------|----------|----------|
| 正常 Issue | `https://github.com/octocat/Hello-World/issues/348` | 成功获取 Issue 数据 |
| 正常 PR | `https://github.com/octocat/Hello-World/pull/348` | 成功获取 PR 数据 |
| 正常 Discussion | `https://github.com/community/community/discussions/12345` | 成功获取 Discussion 数据 |
| 不存在的资源 | `https://github.com/octocat/Hello-World/issues/999999` | 返回 "resource not found" 错误 |
| 私有资源（无 token） | 私有仓库 Issue | 返回认证失败错误或 404 |

**注意：** T014 测试会失败，因为此时 `FetchIssue` 等函数尚未实现。这是 TDD 的正确流程。

---

## Phase 3: Markdown Converter (转换逻辑，TDD)

本阶段实现 Markdown 生成逻辑，遵循测试先行原则。

---

### 3.1 Frontmatter 生成

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T019 | [TEST] 创建 frontmatter 测试文件 | `internal/converter/frontmatter_test.go` | 表格驱动测试 Frontmatter 生成逻辑，验证输出格式正确 |
| T020 | → [IMP] 实现 frontmatter 生成 | `internal/converter/frontmatter.go` | 实现 `generateFrontmatter()` 函数，生成 YAML Frontmatter |

**T019 测试用例（必须覆盖）：**
| 输入数据 | 期望输出 |
|----------|----------|
| Title: "Test Issue", URL: "...", Author: "user", AuthorURL: "...", CreatedAt: time, Status: "open", Type: "issue" | 包含完整的 YAML Frontmatter，字段正确 |
| 包含特殊字符的 Title | 正确转义 YAML 特殊字符 |

---

### 3.2 Reactions 渲染

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T021 | [P] [TEST] 创建 reactions 渲染测试 | `internal/converter/converter_test.go` | 表格驱动测试 Reactions 渲染逻辑 |
| T022 | → [P] [IMP] 实现 reactions 渲染 | `internal/converter/converter.go` | 实现 `renderReactions(r *Reactions) string` 函数，将 Reactions 转换为 emoji + 数量格式 |

**T021 测试用例（必须覆盖）：**
| 输入 Reactions | 期望输出 |
|----------------|----------|
| ThumbsUp: 5, ThumbsDown: 2, Heart: 3 | "👍 5 👎 2 ❤️ 3" |
| 所有字段为 0 | 空字符串 "" |
| nil | 空字符串 "" |

---

### 3.3 用户链接渲染

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T023 | [P] [TEST] 创建用户链接测试 | `internal/converter/converter_test.go` | 表格驱动测试用户名渲染逻辑 |
| T024 | → [P] [IMP] 实现用户链接渲染 | `internal/converter/converter.go` | 实现 `renderUser(username, userURL string, enableLinks bool) string` 函数 |

**T023 测试用例（必须覆盖）：**
| 输入 | EnableLinks | 期望输出 |
|------|-------------|----------|
| username: "octocat", userURL: "https://github.com/octocat" | true | "[@octocat](https://github.com/octocat)" |
| username: "octocat", userURL: "https://github.com/octocat" | false | "@octocat" |

---

### 3.4 Issue 转 Markdown

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T025 | [TEST] 创建 Issue 转 Markdown 测试 | `internal/converter/converter_test.go` | 表格驱动测试 Issue 转 Markdown 逻辑 |
| T026 | → [IMP] 实现 Issue 转 Markdown | `internal/converter/converter.go` | 实现 `ToMarkdown(issue *Issue, opts *Options) ([]byte, error)` 函数 |

**T025 测试用例（必须覆盖）：**
| 输入 Issue | EnableReactions | 期望输出 |
|-------------|-----------------|----------|
| 简单 Issue（无评论） | false | 包含 Frontmatter、标题、正文 |
| 带 Reactions 的 Issue | true | 包含 Reactions 统计 |
| 带评论的 Issue | false | 包含评论列表，按时间正序 |
| 带 Reactions 的评论 | true | 每条评论下方显示 Reactions |

---

### 3.5 PullRequest 转 Markdown

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T027 | [TEST] 创建 PR 转 Markdown 测试 | `internal/converter/converter_test.go` | 表格驱动测试 PR 转 Markdown 逻辑 |
| T028 | → [IMP] 实现 PR 转 Markdown | `internal/converter/converter.go` | 实现 `ToMarkdownPR(pr *PullRequest, opts *Options) ([]byte, error)` 函数 |

**T027 测试用例（必须覆盖）：**
| 输入 PR | 期望输出 |
|---------|----------|
| 简单 PR（无评论） | 包含 Frontmatter、标题、描述，类型为 "pull_request" |
| 带 Review Comments 的 PR | Review Comments 合并到评论列表，不包含代码文件和行号 |
| Merged 状态的 PR | Status 字段为 "merged" |

---

### 3.6 Discussion 转 Markdown

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T029 | [TEST] 创建 Discussion 转 Markdown 测试 | `internal/converter/converter_test.go` | 表格驱动测试 Discussion 转 Markdown 逻辑 |
| T030 | → [IMP] 实现 Discussion 转 Markdown | `internal/converter/converter.go` | 实现 `ToMarkdownDiscussion(discussion *Discussion, opts *Options) ([]byte, error)` 函数 |

**T029 测试用例（必须覆盖）：**
| 输入 Discussion | 期望输出 |
|----------------|----------|
| 简单 Discussion（无评论） | 包含 Frontmatter、标题、正文，类型为 "discussion" |
| 带 Answer 的 Discussion | Answer 评论标记 "✅ **Answer**" |
| 带 Reactions 的 Answer | Answer 下方显示 Reactions 统计 |

---

### 3.7 DefaultOptions

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T031 | [P] [TEST] 创建 DefaultOptions 测试 | `internal/converter/converter_test.go` | 测试默认选项的值 |
| T032 | → [P] [IMP] 实现 DefaultOptions | `internal/converter/converter.go` | 实现 `DefaultOptions() *Options` 函数，返回默认选项（EnableReactions: false, EnableUserLinks: false） |

---

## Phase 4: CLI Assembly (命令行入口集成)

本阶段实现 CLI 主入口，协调所有包的功能。

---

### 4.1 CLI 主程序

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T033 | [IMP] 实现 CLI 主入口 | `cmd/issue2md/main.go` | 实现 main 函数，协调 cli、parser、config、github、converter 包 |

**T033 实现要求（必须遵循）：**
1. 调用 `cli.ParseArgs()` 解析命令行参数
2. 如果解析失败，调用 `cli.PrintError()` 并退出（退出码 1）
3. 如果缺少参数，调用 `cli.PrintHelp()` 并退出（退出码 1）
4. 调用 `config.GetGitHubToken()` 获取 token
5. 调用 `parser.ParseURL()` 解析 URL
6. 根据资源类型调用相应的 `github.Client.Fetch*()` 方法
7. 根据资源类型调用相应的 `converter.ToMarkdown*()` 方法
8. 如果 `OutputFile` 为空，输出到 stdout；否则写入文件
9. 所有错误使用 `fmt.Errorf("...: %w", err)` 包装
10. 输出错误到 stderr，退出码为 1

---

### 4.2 Makefile

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T034 | [P] [IMP] 创建 Makefile | `Makefile` | 定义构建、测试、安装等目标 |

**Makefile 目标（必须包含）：**
```makefile
test        # 运行所有测试
build       # 构建二进制文件
install     # 安装到 $GOPATH/bin
clean       # 清理构建产物
fmt         # 格式化代码
vet         # 运行 go vet
lint        # 运行静态检查
```

---

### 4.3 README.md

| ID | 任务 | 文件 | 描述 |
|----|------|------|------|
| T035 | [P] [IMP] 创建 README.md | `README.md` | 编写用户文档，包括安装、使用、示例 |

**README.md 内容（必须包含）：**
- 项目简介
- 安装方法（go install、Homebrew、源码编译）
- 使用示例
- 命令行参数说明
- 环境变量说明（GITHUB_TOKEN）
- 常见问题

---

## 依赖关系图

```
Phase 1: Foundation
├── T001 → T002 (config)
├── T003 → T004 (parser)
├── T005 → T006 → T007 (cli)
├── T008, T009 [P] (cli)
├── T010 [P] (github types)
├── T011 [P] (converter types)
└── T012 → T013 (go.mod)

Phase 2: GitHub Fetcher
└── T014 → T015 → T016 → T017 → T018 (github client)

Phase 3: Markdown Converter
├── T019 → T020 (frontmatter)
├── T021 → T022 [P] (reactions)
├── T023 → T024 [P] (user links)
├── T025 → T026 (Issue)
├── T027 → T028 (PR)
├── T029 → T030 (Discussion)
└── T031 → T032 [P] (DefaultOptions)

Phase 4: CLI Assembly
├── T033 (main.go)
├── T034 [P] (Makefile)
└── T035 [P] (README.md)
```

---

## 测试覆盖率目标

| 包 | 覆盖率目标 | 测试类型 |
|----|------------|----------|
| internal/config | 100% | 表格驱动单元测试 |
| internal/parser | 100% | 表格驱动单元测试 |
| internal/cli | 100% | 表格驱动单元测试 |
| internal/github | 90%+ | 集成测试（真实 API） |
| internal/converter | 95%+ | 表格驱动单元测试 |

---

## 合宪性检查

| 宪法条款 | 对应任务 | 验证方式 |
|----------|----------|----------|
| 第一条：简单性原则 | 所有实现任务 | 不使用过度抽象，仅实现 Spec 明确要求的功能 |
| 第二条：测试先行铁律 | T001-T031 所有 TEST 任务 | 每个功能先编写测试，再实现 |
| 2.2 (表格驱动) | T003, T005, T019-T032 | 使用表格驱动测试风格 |
| 2.3 (拒绝Mocks) | T014 | 集成测试使用真实 GitHub API |
| 第三条：明确性原则 | 所有实现任务 | 错误使用 `%w` 包装，无全局变量 |

---

## 验收标准对应

| 验收标准 | 对应任务 |
|----------|----------|
| TC01: Issue 导出（基本功能） | T033 + T016 + T026 |
| TC02: Issue 导出（指定输出文件） | T033 |
| TC03: Issue 导出（启用 Reactions） | T022 + T026 + T033 |
| TC04: Issue 导出（启用用户链接） | T024 + T026 + T033 |
| TC05: PR 导出（不含 diff） | T017 + T028 + T033 |
| TC06: Discussion 导出（包含 Answer） | T018 + T030 + T033 |
| TC07: 无效 URL 错误处理 | T004 + T033 |
| TC08: 不存在的资源错误处理 | T014 + T033 |
| TC09: 环境变量认证 | T002 + T015 + T033 |
| TC10: Flag 组合测试 | T007 + T033 |
| TC11: 缺少必需参数 | T007 + T033 |
| TC12: 不支持的 URL 格式 | T004 + T033 |

---

## 执行建议

### 批量执行（并行任务）

**第一批（可并行）：**
- T001, T003, T005, T008, T009, T010, T011, T012, T019, T021, T023, T031, T034, T035

**第二批（依赖第一批）：**
- T002, T004, T006, T013, T020, T022, T024, T032

**第三批（依赖第二批）：**
- T007, T014

**第四批（依赖第三批）：**
- T015, T016, T017, T018, T025, T027, T029

**第五批（依赖第四批）：**
- T026, T028, T030

**第六批（依赖第五批）：**
- T033

---

## 附录：测试数据推荐

**推荐用于集成测试的真实资源：**

| 资源类型 | URL | 说明 |
|----------|------|------|
| Issue | https://github.com/octocat/Hello-World/issues/348 | 官方示例仓库 |
| PR | https://github.com/octocat/Hello-World/pull/348 | 官方示例仓库 |
| Discussion | https://github.com/community/community/discussions/12345 | GitHub 官方社区 |

这些 URL 长期有效且公开访问，适合集成测试。
