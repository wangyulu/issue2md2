package config

import (
	"os"
)

// GetGitHubToken 从环境变量读取 GitHub Token
// 如果环境变量未设置，返回空字符串
func GetGitHubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}
