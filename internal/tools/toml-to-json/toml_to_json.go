// Package tomljson 实现将 TOML 文本转换为 JSON。
package tomljson

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "toml-to-json"
	Name        = "TOML 转 JSON"
	Description = "将 TOML 转换为 JSON"
	Category    = "转换器"
	Icon        = "Braces"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"toml", "json", "转换", "convert", "transform"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

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
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	result, err := Convert(in.Text)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Convert 将 TOML 文本转换为 JSON（3 空格缩进，对齐参考项目）。
func Convert(t string) (string, error) {
	if strings.TrimSpace(t) == "" {
		return "", nil
	}

	var v map[string]interface{}
	if err := toml.Unmarshal([]byte(t), &v); err != nil {
		return "", fmt.Errorf("无效的 TOML: %w", err)
	}

	b, err := json.MarshalIndent(v, "", "   ")
	if err != nil {
		return "", fmt.Errorf("序列化 JSON 失败: %w", err)
	}
	return string(b), nil
}
