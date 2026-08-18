// Package jsonfmt 实现 JSON 格式化工具：美化 JSON 文本，支持缩进与按键排序。
package jsonfmt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/registry"
)

const (
	ID          = "json-formatter"
	Name        = "JSON 格式化"
	Description = "美化 JSON，支持缩进与按键排序"
	Category    = "开发"
	Icon        = "Braces"
)

var Keywords = []string{"json", "format", "pretty", "beautify", "indent", "缩进", "美化", "格式化", "排序"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	JSON     string `json:"json"`      // 待格式化 JSON 文本
	Indent   string `json:"indent"`    // 缩进："2" | "4" | "\t"（默认 "2"）
	SortKeys bool   `json:"sort_keys"` // 是否按键排序（false 时保持原顺序）
}

// output 是工具的输出结构。
type output struct {
	Formatted string `json:"formatted"`
	LineCount int    `json:"line_count"`
	CharCount int    `json:"char_count"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回格式化结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	indent := in.Indent
	if indent == "" {
		indent = "2"
	}
	var indentStr string
	switch indent {
	case "4":
		indentStr = "    "
	case "\t":
		indentStr = "\t"
	case "2":
		indentStr = "  "
	default:
		return "", fmt.Errorf("缩进仅支持 2 / 4 / Tab（当前 %q）", indent)
	}

	var formatted []byte
	var err error
	if in.SortKeys {
		// 解析后重新序列化：map 按键排序输出
		var v any
		if err = json.Unmarshal([]byte(in.JSON), &v); err != nil {
			return "", fmt.Errorf("JSON 解析失败: %w", err)
		}
		formatted, err = json.MarshalIndent(v, "", indentStr)
	} else {
		// token 级缩进：保留原始键顺序
		var buf bytes.Buffer
		if err = json.Indent(&buf, []byte(in.JSON), "", indentStr); err != nil {
			return "", fmt.Errorf("JSON 解析失败: %w", err)
		}
		formatted = buf.Bytes()
	}
	if err != nil {
		return "", fmt.Errorf("JSON 格式化失败: %w", err)
	}

	text := string(formatted)
	out, err := json.Marshal(output{
		Formatted: text,
		LineCount: strings.Count(text, "\n") + 1,
		CharCount: len(text),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
