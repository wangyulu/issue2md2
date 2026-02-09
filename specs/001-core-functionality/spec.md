# Spec #001: Core Functionality - GitHub to Markdown Converter

**Status:** Draft
**Created:** 2026-02-08
**Priority:** P0

---

## 1. 用户故事 (User Stories)

### 1.1 当前 MVP: CLI 工具

> 作为一个开发者，我希望能够将 GitHub Issue/PR/Discussion 的内容导出为 Markdown 文件，以便于归档、文档编写或离线阅读。

**Acceptance Criteria:**
- 支持从 GitHub URL 导出 Issue、PR、Discussion
- 输出符合 GitHub Flavored Markdown 标准的格式
- 支持通过环境变量配置认证
- 默认输出到 stdout，可重定向到文件

### 1.2 未来: Web 界面 (Future User Story)

> 作为一个用户，我希望通过 Web 界面输入 GitHub URL 并直接下载 Markdown 文件，无需安装命令行工具。

**Implementation Note:** 此功能暂不实现，但需要预留架构扩展点。

---

## 2. 功能性需求 (Functional Requirements)

### 2.1 支持的资源类型

| 资源类型 | URL 模式 | 提取内容 |
|---------|---------|---------|
| Issue | `https://github.com/{owner}/{repo}/issues/{number}` | 标题、作者、创建时间、状态、主楼、所有评论 |
| Pull Request | `https://github.com/{owner}/{repo}/pull/{number}` | 标题、作者、创建时间、状态、描述、评论、Review评论 |
| Discussion | `https://github.com/{owner}/{repo}/discussions/{number}` | 标题、作者、创建时间、状态、主楼、所有评论 |

**PR 特殊处理:**
- 不提取 diff 信息
- 不提取 commits 历史
- Review Comments 不保留代码文件和行号上下文

### 2.2 命令行接口

**命令格式:**
```bash
issue2md [flags] <url> [output_file]
```

**参数说明:**
- `<url>`: 必需参数，GitHub Issue/PR/Discussion 的完整 URL
- `[output_file]`: 可选参数，指定输出文件路径。若省略则输出到 stdout

**Flags:**
| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-enable-reactions` | bool | false | 是否在主楼和评论下方显示 Reactions 统计 |
| `-enable-user-links` | bool | false | 是否将用户名渲染为指向其 GitHub 主页的链接 |

**认证配置:**
- 仅通过环境变量 `GITHUB_TOKEN` 配置 Personal Access Token
- **严禁** 提供 `--token` 参数，避免在 Shell 历史中泄露密钥
- 未配置 token 时，仅能访问公开仓库（受 GitHub API 限流限制）

### 2.3 Markdown 输出格式

#### 2.3.1 Frontmatter (YAML)

必须包含以下元数据：

```yaml
---
title: "Issue Title"
url: "https://github.com/owner/repo/issues/123"
author: "username"
author_url: "https://github.com/username"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "issue"
---
```

#### 2.3.2 内容结构

**总体顺序:**
1. Frontmatter
2. 主楼（标题 + 正文 + 可选 reactions）
3. 评论列表（按时间正序，扁平化展示）

**评论渲染规则:**
- Issue/Discussion: 主楼 → 评论1 → 评论2 → ...
- PR: 描述 → 评论1 → 评论2 → Review Comments → ...
- Discussion Answer: 使用引用块或 Emoji (✅) 显著标记

**格式化规则:**
- 代码块、表情符号保留 GitHub 原始格式
- 图片直接保留原始链接，不下载到本地
- 评论排序统一按时间**正序**
- PR Review Comments 不按 Review 分组，与其他评论一起按时间线平铺展示

### 2.4 URL 解析逻辑

必须正确识别以下 URL 格式：

```go
// Issue
https://github.com/owner/repo/issues/123
https://github.com/owner/repo/issues/123#issuecomment-456

// Pull Request
https://github.com/owner/repo/pull/42
https://github.com/owner/repo/pull/42#discussion_r123

// Discussion
https://github.com/owner/repo/discussions/7
https://github.com/owner/repo/discussions/7#discussioncomment-890
```

**解析错误处理:**
- URL 格式不符合上述模式 → 返回明确的错误信息
- 无法解析 owner/repo → 返回明确的错误信息

---

## 3. 非功能性需求 (Non-Functional Requirements)

### 3.1 架构设计原则

**简单性 (Simplicity):**
- 遵循 Go 语言"少即是多"哲学
- 优先使用标准库（`net/http`, `encoding/json` 等）
- 避免不必要的抽象和接口

**解耦 (Decoupling):**
- URL 解析逻辑与 GitHub API 调用逻辑分离
- Markdown 生成逻辑与数据获取逻辑分离
- 为未来 Web 版本预留数据层复用点

**依赖管理:**
- 仅引入必需的外部依赖
- 优先使用 Go 官方推荐的 GitHub API 客户端库（如 `google/go-github`）

### 3.2 错误处理

**错误传播规范:**
- 所有错误必须显式处理
- 错误传递必须使用 `fmt.Errorf("...: %w", err)` 进行包装

**错误场景处理:**

| 场景 | 错误信息 | 行为 |
|------|---------|------|
| URL 格式错误 | `invalid GitHub URL: {url}` | 输出到 stderr，退出码 1 |
| 资源不存在 (404) | `resource not found: {url}` | 输出到 stderr，退出码 1 |
| 认证失败 (401) | `authentication failed: check GITHUB_TOKEN` | 输出到 stderr，退出码 1 |
| API 限流 (429) | `rate limit exceeded: {GitHub API message}` | 输出到 stderr，退出码 1 |
| 网络错误 | `failed to fetch data: {error}` | 输出到 stderr，退出码 1 |

### 3.3 测试要求

**测试策略:**
- 优先编写**表格驱动测试 (Table-Driven Tests)**
- 核心逻辑使用单元测试覆盖
- 集成测试使用真实的 GitHub API（或可配置的 mock）
- 测试数据使用真实的 GitHub Issue/PR/Discussion 作为案例

---

## 4. 验收标准 (Acceptance Criteria)

### 4.1 功能测试 Case

#### TC01: Issue 导出 (基本功能)
```bash
issue2md https://github.com/octocat/Hello-World/issues/348
```
**Expected:**
- 正确输出到 stdout
- 包含完整的 Frontmatter
- 包含主楼和所有评论
- 按时间正序排列

#### TC02: Issue 导出 (指定输出文件)
```bash
issue2md https://github.com/octocat/Hello-World/issues/348 issue.md
```
**Expected:**
- 正确写入 `issue.md`
- 文件内容与 stdout 输出一致

#### TC03: Issue 导出 (启用 Reactions)
```bash
issue2md -enable-reactions https://github.com/octocat/Hello-World/issues/348
```
**Expected:**
- 主楼和评论下方显示 Reactions 统计
- 格式示例: `👍 5 👎 2`

#### TC04: Issue 导出 (启用用户链接)
```bash
issue2md -enable-user-links https://github.com/octocat/Hello-World/issues/348
```
**Expected:**
- 用户名渲染为 Markdown 链接
- 示例: `[@octocat](https://github.com/octocat)`

#### TC05: PR 导出 (不含 diff)
```bash
issue2md https://github.com/octocat/Hello-World/pull/348
```
**Expected:**
- 不包含 diff 内容
- 包含 Review Comments（不含代码文件和行号）

#### TC06: Discussion 导出 (包含 Answer)
```bash
issue2md https://github.com/community/community/discussions/12345
```
**Expected:**
- Answer 评论有显著标记（✅ 或引用块）

#### TC07: 无效 URL 错误处理
```bash
issue2md https://example.com/invalid
```
**Expected:**
- 输出错误信息到 stderr
- 退出码为 1

#### TC08: 不存在的资源错误处理
```bash
issue2md https://github.com/octocat/Hello-World/issues/999999
```
**Expected:**
- 输出 "resource not found" 错误到 stderr
- 退出码为 1

#### TC09: 环境变量认证
```bash
export GITHUB_TOKEN=ghp_xxx
issue2md https://github.com/owner/private-repo/issues/1
```
**Expected:**
- 成功导出私有仓库 Issue
- 不支持通过 `--token` 参数传递

#### TC10: Flag 组合测试
```bash
issue2md -enable-reactions -enable-user-links https://github.com/octocat/Hello-World/issues/348 output.md
```
**Expected:**
- 两个 Flag 同时生效
- 输出到文件 `output.md`

### 4.2 负面测试 Case

#### TC11: 缺少必需参数
```bash
issue2md
```
**Expected:**
- 输出使用帮助信息
- 退出码为 1

#### TC12: 不支持的 URL 格式
```bash
issue2md https://github.com/owner/repo/tree/main
```
**Expected:**
- 输出 "unsupported resource type" 错误
- 退出码为 1

---

## 5. 输出格式示例

### 5.1 Issue 输出示例

```markdown
---
title: "Issue Title Example"
url: "https://github.com/owner/repo/issues/123"
author: "octocat"
author_url: "https://github.com/octocat"
created_at: "2024-01-01T12:00:00Z"
status: "open"
type: "issue"
---

# Issue Title Example

This is the issue body content.

Some **bold** and *italic* text.

```go
func main() {
    fmt.Println("Hello, World!")
}
```

An image: ![Alt text](https://github.com/owner/repo/assets/12345/image.png)

## Reactions

👍 5 👎 2 ❤️ 3

---

## Comments

### @octocat commented at 2024-01-02T10:00:00Z

This is the first comment.

Some code example:

```javascript
console.log("test");
```

👍 2

### @defunkt commented at 2024-01-03T14:30:00Z

Replying to the first comment.

👍 1 👎 1
```

### 5.2 PR 输出示例

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

## Reactions

👍 3

---

## Comments

### @octocat commented at 2024-01-02T10:00:00Z

LGTM! Ready to merge.

### @reviewer1 reviewed at 2024-01-03T14:30:00Z

This looks good, but consider handling edge cases.

👍 1

### @octocat commented at 2024-01-04T09:00:00Z

Good point, I'll update it.
```

### 5.3 Discussion 输出示例

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

### @octocat commented at 2024-01-04T09:00:00Z

Thank you! This helps a lot.
```

---

## 6. 技术债务与未来改进

### 6.1 已知限制
- 不支持自定义 Markdown 模板
- 不支持批量导出
- 不支持私有仓库（除非配置 token）

### 6.2 未来增强功能
- 批量导出多个 Issue/PR
- 自定义 Markdown 模板支持
- Web 界面实现
- 支持更多平台（GitLab, Bitbucket 等）
- 导出为 PDF/HTML 格式
- 搜索和过滤功能

---

## 7. 实现计划

### Phase 1: 核心功能 (MVP)
- [ ] URL 解析逻辑
- [ ] GitHub API 客户端封装
- [ ] 数据结构定义
- [ ] Markdown 生成逻辑
- [ ] CLI 参数解析
- [ ] 单元测试
- [ ] 集成测试

### Phase 2: 增强功能
- [ ] Reactions 支持
- [ ] 用户链接支持
- [ ] Discussion Answer 标记
- [ ] 错误处理完善

### Phase 3: 文档与发布
- [ ] README 文档
- [ ] 示例代码
- [ ] GitHub Release
- [ ] Homebrew Formula（可选）
