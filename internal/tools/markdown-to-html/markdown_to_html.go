// Package markdownhtml 实现将 Markdown 文本转换为 HTML。
package markdownhtml

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/yuin/goldmark"
	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "markdown-to-html"
	Name        = "Markdown 转 HTML"
	Description = "将 Markdown 转换为 HTML"
	Category    = "转换器"
	Icon        = "Markdown"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"markdown", "html", "转换", "convert", "md"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// md 为 goldmark 解析器实例。
var md = goldmark.New()

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"`
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(ctx context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	var buf bytes.Buffer
	if err := md.Convert([]byte(in.Text), &buf); err != nil {
		return "", fmt.Errorf("转换 Markdown 失败: %w", err)
	}

	out, err := json.Marshal(output{Result: buf.String()})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
