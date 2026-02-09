package cli

import (
	"fmt"
	"os"
	"strings"
)

// 错误常量
const (
	ErrMissingRequiredArg = "missing required argument: %s"
	ErrUnknownFlag       = "unknown flag: %s"
	ErrHelpDisplayed     = "help displayed"
)

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

// ParseArgs 解析命令行参数
// args 通常是 os.Args[1:]
// 返回 Flags 和 Args，如果解析失败返回错误
//
// 支持的标志:
//   -enable-reactions: 启用 Reactions 显示
//   -enable-user-links: 启用用户链接
//
// 位置参数:
//   url: 必需，GitHub URL
//   output_file: 可选，输出文件路径
func ParseArgs(args []string) (*Flags, *Args, error) {
	if len(args) == 0 {
		return nil, nil, fmt.Errorf(ErrMissingRequiredArg, "url")
	}

	flags := &Flags{
		EnableReactions: false,
		EnableUserLinks: false,
	}

	// 解析标志
	remainingArgs := []string{}
	for _, arg := range args {
		switch arg {
		case "-enable-reactions":
			flags.EnableReactions = true
		case "-enable-user-links":
			flags.EnableUserLinks = true
		case "-h":
			PrintHelp(os.Stdout)
			return nil, nil, fmt.Errorf(ErrHelpDisplayed)
		default:
			if strings.HasPrefix(arg, "-") {
				return nil, nil, fmt.Errorf(ErrUnknownFlag, arg)
			}
			remainingArgs = append(remainingArgs, arg)
		}
	}

	// 验证必需参数
	if len(remainingArgs) == 0 {
		return nil, nil, fmt.Errorf(ErrMissingRequiredArg, "url")
	}

	// 解析位置参数
	cliArgs := &Args{
		URL: remainingArgs[0],
	}
	if len(remainingArgs) > 1 {
		cliArgs.OutputFile = remainingArgs[1]
	}

	return flags, cliArgs, nil
}
