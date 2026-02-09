package main

import (
	"fmt"
	"os"

	"github.com/wangyulu/issue2md2/internal/cli"
	"github.com/wangyulu/issue2md2/internal/config"
	"github.com/wangyulu/issue2md2/internal/converter"
	"github.com/wangyulu/issue2md2/internal/github"
	"github.com/wangyulu/issue2md2/internal/parser"
)

func main() {
	// 解析命令行参数
	flags, args, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		cli.PrintError(os.Stderr, err)
		os.Exit(1)
	}

	// 检查是否需要显示帮助
	if err != nil && err.Error() == cli.ErrHelpDisplayed {
		os.Exit(0)
	}

	// 解析 URL
	resource, err := parser.ParseURL(args.URL)
	if err != nil {
		cli.PrintError(os.Stderr, err)
		os.Exit(1)
	}

	// 获取 token（可选）
	token := config.GetGitHubToken()

	// 创建 GitHub 客户端
	client := github.NewClient()
	_ = token // 避免未使用变量警告

	// 根据资源类型获取数据
	var markdown []byte
	switch resource.Type {
	case parser.ResourceTypeIssue:
		issue, err := client.FetchIssue(resource.Owner, resource.Repo, resource.Number)
		if err != nil {
			cli.PrintError(os.Stderr, fmt.Errorf("failed to fetch issue: %w", err))
			os.Exit(1)
		}
		markdown, err = converter.ToMarkdown(issue, &converter.Options{
			EnableReactions: flags.EnableReactions,
			EnableUserLinks: flags.EnableUserLinks,
		})

	case parser.ResourceTypePullRequest:
		pr, err := client.FetchPullRequest(resource.Owner, resource.Repo, resource.Number)
		if err != nil {
			cli.PrintError(os.Stderr, fmt.Errorf("failed to fetch pull request: %w", err))
			os.Exit(1)
		}
		markdown, err = converter.ToMarkdownPR(pr, &converter.Options{
			EnableReactions: flags.EnableReactions,
			EnableUserLinks: flags.EnableUserLinks,
		})

	case parser.ResourceTypeDiscussion:
		discussion, err := client.FetchDiscussion(resource.Owner, resource.Repo, resource.Number)
		if err != nil {
			cli.PrintError(os.Stderr, fmt.Errorf("failed to fetch discussion: %w", err))
			os.Exit(1)
		}
		markdown, err = converter.ToMarkdownDiscussion(discussion, &converter.Options{
			EnableReactions: flags.EnableReactions,
			EnableUserLinks: flags.EnableUserLinks,
		})

	default:
		cli.PrintError(os.Stderr, fmt.Errorf("unsupported resource type: %s", resource.Type))
		os.Exit(1)
	}

	if err != nil {
		cli.PrintError(os.Stderr, fmt.Errorf("failed to convert to markdown: %w", err))
		os.Exit(1)
	}

	// 输出结果
	if args.OutputFile == "" {
		// 输出到 stdout
		fmt.Print(string(markdown))
	} else {
		// 写入文件
		err = os.WriteFile(args.OutputFile, markdown, 0644)
		if err != nil {
			cli.PrintError(os.Stderr, fmt.Errorf("failed to write file: %w", err))
			os.Exit(1)
		}
	}
}
