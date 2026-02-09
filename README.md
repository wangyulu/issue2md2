# issue2md

GitHub Issue/PR/Discussion to Markdown Converter

一个简单、快速的命令行工具，用于将 GitHub Issue、Pull Request 或 Discussion 导出为 Markdown 文件。

## 功能特性

- ✅ 支持 Issue、Pull Request 和 Discussion
- ✅ 输出标准 GitHub Flavored Markdown 格式
- ✅ 包含完整的 YAML Frontmatter
- ✅ 可选：显示 Reactions 统计
- ✅ 可选：用户名渲染为 GitHub 链接
- ✅ 按时间正序排列所有评论
- ✅ 支持公开仓库和私有仓库（需认证）
- ✅ 轻量级：仅使用必要的 GitHub API 客户端库

## 安装

### 使用 go install

```bash
go install github.com/wangyulu/issue2md2@latest
```

### 从源码编译

```bash
git clone https://github.com/wangyulu/issue2md2.git
cd issue2md2
make build
```

编译后的二进制文件位于当前目录 `issue2md`。

## 使用方法

### 基本用法

```bash
# 输出到 stdout（默认）
./issue2md https://github.com/owner/repo/issues/123

# 输出到文件
./issue2md https://github.com/owner/repo/issues/123 issue.md

# 启用 Reactions
./issue2md -enable-reactions https://github.com/owner/repo/issues/123

# 启用用户链接
./issue2md -enable-user-links https://github.com/owner/repo/issues/123

# 所有选项组合
./issue2md -enable-reactions -enable-user-links https://github.com/owner/repo/issues/123 output.md
```

### 支持的资源类型

| 资源类型 | URL 示例 |
|---------|---------|
| Issue | `https://github.com/owner/repo/issues/123` |
| Pull Request | `https://github.com/owner/repo/pull/42` |
| Discussion | `https://github.com/owner/repo/discussions/7` |

### 命令行参数

```bash
issue2md [flags] <url> [output_file]
```

**Flags:**

| Flag | 说明 | 默认值 |
|------|------|--------|
| `-enable-reactions` | 启用 Reactions 显示 | `false` |
| `-enable-user-links` | 渲染用户名为 GitHub 链接 | `false` |
| `-h` | 显示帮助信息 | - |

**位置参数:**

| 参数 | 说明 | 是否必需 |
|------|------|----------|
| `<url>` | GitHub Issue/PR/Discussion 的完整 URL | 必需 |
| `[output_file]` | 输出文件路径，省略则输出到 stdout | 可选 |

## 环境变量

### GITHUB_TOKEN

GitHub Personal Access Token，用于访问私有仓库或提高 API 限流上限。

**设置方法:**
```bash
export GITHUB_TOKEN=ghp_xxx
```

**注意:**
- 仅支持通过环境变量传递 token，不支持 `--token` 参数（避免在 Shell 历史中泄露）
- 未设置 token 时，只能访问公开仓库（60 次/小时）
- 已认证：5000 次/小时
- 获取 Personal Access Token: https://github.com/settings/tokens

## 输出格式

### Frontmatter

所有输出都包含 YAML Frontmatter：

```yaml
---
title: "Issue Title"
url: "https://github.com/owner/repo/issues/123"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "issue"
---
```

### 内容结构

1. Frontmatter
2. 主楼（标题 + 正文 + 可选 reactions）
3. 评论列表（按时间正序，扁平化展示）

### 特殊标记

- **Reactions**: 当启用时，显示为 `👍 5 👎 2 ❤️ 3`
- **用户链接**: 当启用时，用户名显示为 `[@octocat](https://github.com/octocat)`
- **Discussion Answer**: Answer 评论标记为 `✅ **Answer**`

## 示例输出

### Issue

```markdown
---
title: "Testing comments"
url: "https://github.com/octocat/Hello-World/issues/348"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2017-05-22T18:47:38Z"
status: "open"
type: "issue"
---

# Testing comments

Let's add some, shall we?

---

## Comments

### @bgammill commented at 2017-05-22T21:00:09Z

Here is a shiny new comment.

### @operate2v commented at 2017-05-23T00:00:27Z

A shiny new comment! :tada:
```

### Pull Request

```markdown
---
title: "PR Title Example"
url: "https://github.com/owner/repo/pull/42"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "pull_request"
---

# PR Title Example

This is the PR description.

---

## Comments

### @reviewer1 reviewed at 2024-01-03T14:30:00Z

This looks good, but consider handling edge cases.

### @octocat commented at 2024-01-04T09:00:00Z

Good point, I'll update it.
```

### Discussion

```markdown
---
title: "Discussion Title"
url: "https://github.com/owner/repo/discussions/7"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "discussion"
---

# Discussion Title

This is the discussion body.

---

## Comments

### @octocat commented at 2024-01-02T10:00:00Z

Initial question about something.

### @expert commented at 2024-01-03T14:30:00Z

> This is the discussion body.

Here is a detailed answer with explanation.

✅ **Answer**
```

## 错误处理

| 场景 | 错误信息 | 退出码 |
|------|----------|--------|
| 缺少必需参数 | `missing required argument: url` | 1 |
| URL 格式错误 | `invalid GitHub URL: {url}` | 1 |
| 不支持的资源类型 | `unsupported resource type: {url}` | 1 |
| 资源不存在 | `resource not found: {url}` | 1 |
| 认证失败 | `authentication failed: check GITHUB_TOKEN` | 1 |
| API 限流 | `rate limit exceeded: {GitHub API message}` | 1 |
| 网络错误 | `failed to fetch data: {error}` | 1 |

## 常见问题

### Q: 如何访问私有仓库？

A: 设置 `GITHUB_TOKEN` 环境变量：

```bash
export GITHUB_TOKEN=ghp_xxx
./issue2md https://github.com/owner/private-repo/issues/1
```

### Q: API 限流怎么办？

A: 设置 `GITHUB_TOKEN` 可以大幅提高限流上限：
- 未认证: 60 次/小时
- 已认证: 5000 次/小时

### Q: 为什么 PR 没有 diff 信息？

A: 工具的设计目标是归档"讨论过程"，而不是代码变更。如果需要 diff 信息，建议直接使用 GitHub 的导出功能。

### Q: Discussion 的 Answer 如何识别？

A: 被标记为 Answer 的评论会显示 `✅ **Answer**` 标记。

### Q: 为什么图片链接保持原样而不下载？

A: 为了保持 Markdown 文件的简洁性和可移植性，图片链接保持原样，不下载到本地。

## 开发

### 运行测试

```bash
make test
```

### 构建二进制文件

```bash
make build
```

### 安装到 $GOPATH/bin

```bash
make install
```

### 格式化代码

```bash
make fmt
```

### 运行静态检查

```bash
make vet
```

## 项目结构

```
issue2md2/
├── cmd/
│   └── issue2md/          # CLI 入口
├── internal/
│   ├── config/             # 环境变量配置
│   ├── parser/             # URL 解析
│   ├── github/             # GitHub API 客户端
│   ├── converter/          # Markdown 生成
│   └── cli/               # 命令行接口
├── specs/                 # 技术规范
├── Makefile
├── go.mod
└── README.md
```

## 技术栈

- Go 1.25.1
- `github.com/google/go-github/v68` - GitHub REST API 客户端
- `github.com/shurcooL/githubv4` - GitHub GraphQL API 客户端

## 设计原则

本项目严格遵循以下设计原则：

1. **简单性原则**：遵循 Go 语言"少即是多"哲学，只实现明确要求的功能
2. **测试先行**：所有功能从编写测试开始，使用表格驱动测试风格
3. **明确性**：所有错误显式处理，使用 `fmt.Errorf("...: %w", err)` 包装
4. **无全局变量**：所有依赖通过函数参数或结构体成员显式注入

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

- 项目主页: https://github.com/wangyulu/issue2md2
- 问题反馈: https://github.com/wangyulu/issue2md2/issues