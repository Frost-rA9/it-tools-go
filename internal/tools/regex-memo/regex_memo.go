// Package regexmemo 实现正则表达式速查表工具。
// 内容由前端静态展示，Go 侧仅提供工具注册元数据。
package regexmemo

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

const (
	ID          = "regex-memo"
	Name        = "正则表达式速查表"
	Description = "常用正则表达式语法速查"
	Category    = "开发"
	Icon        = "BrandJavascript"
)

var Keywords = []string{"regex", "regular", "expression", "正则", "速查", "cheatsheet", "memo", "模式"}

func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// Executor 实现 registry.Executor 接口。速查表为纯静态前端内容。
type Executor struct{}

func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var input map[string]any
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	return `{}`, nil
}
