// Package gitmemo 实现 Git 备忘录工具。
// 内容由前端静态展示，Go 侧仅提供工具注册元数据。
package gitmemo

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

const (
	ID          = "git-memo"
	Name        = "Git 备忘录"
	Description = "常用 Git 命令速查表"
	Category    = "开发"
	Icon        = "BrandGit"
)

var Keywords = []string{"git", "push", "force", "pull", "commit", "amend", "rebase", "merge", "reset", "soft", "hard", "lease", "备忘录", "速查"}

func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// Executor 实现 registry.Executor 接口。Git 备忘录为纯静态前端内容。
type Executor struct{}

func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var input map[string]any
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	return `{}`, nil
}
