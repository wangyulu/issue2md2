package cli

import (
	"fmt"
	"io"
)

// PrintHelp 打印使用帮助信息到指定的 io.Writer
func PrintHelp(w io.Writer) {
	fmt.Fprintln(w, "Usage: issue2md [flags] <url> [output_file]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  -enable-reactions")
	fmt.Fprintln(w, "        Enable reactions display (default: false)")
	fmt.Fprintln(w, "  -enable-user-links")
	fmt.Fprintln(w, "        Render usernames as links to GitHub profiles (default: false)")
	fmt.Fprintln(w, "  -h")
	fmt.Fprintln(w, "        Show this help message")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Environment Variables:")
	fmt.Fprintln(w, "  GITHUB_TOKEN")
	fmt.Fprintln(w, "        GitHub personal access token (optional, for private repos)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  issue2md https://github.com/owner/repo/issues/123")
	fmt.Fprintln(w, "  issue2md -enable-reactions https://github.com/owner/repo/issues/123 output.md")
	fmt.Fprintln(w, "  GITHUB_TOKEN=ghp_xxx issue2md https://github.com/owner/private-repo/issues/1")
}
