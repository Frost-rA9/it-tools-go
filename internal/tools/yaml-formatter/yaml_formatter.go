// Package yamlfmt 实现 YAML 格式化工具：基于 yaml.v3 的 Node 规范化重编码（保留注释）。
package yamlfmt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	"it-tools-go/internal/registry"
)

const (
	ID          = "yaml-formatter"
	Name        = "YAML 格式化"
	Description = "美化 YAML：规范化缩进，保留注释"
	Category    = "开发"
	Icon        = "FileText"
)

var Keywords = []string{"yaml", "yml", "format", "pretty", "美化", "格式化", "缩进", "indent"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	YAML   string `json:"yaml"`   // 待格式化 YAML
	Indent string `json:"indent"` // "2" | "4"（默认 "2"）
}

// output 是工具的输出结构。
type output struct {
	Formatted string `json:"formatted"`
	LineCount int    `json:"line_count"`
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
	switch indent {
	case "2", "4":
	default:
		return "", fmt.Errorf("缩进仅支持 2 / 4（当前 %q）", indent)
	}

	formatted, err := Format(in.YAML, indent)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{
		Formatted: formatted,
		LineCount: strings.Count(strings.TrimRight(formatted, "\n"), "\n") + 1,
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Format 解析 YAML 到 Node 后按指定缩进重新编码（保留注释）。
func Format(yamlText, indent string) (string, error) {
	if strings.TrimSpace(yamlText) == "" {
		return "", fmt.Errorf("YAML 为空")
	}

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(yamlText), &node); err != nil {
		return "", fmt.Errorf("YAML 解析失败: %w", err)
	}

	indentN := 2
	if indent == "4" {
		indentN = 4
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(indentN)
	if err := enc.Encode(&node); err != nil {
		return "", fmt.Errorf("YAML 编码失败: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", fmt.Errorf("YAML 编码失败: %w", err)
	}
	return buf.String(), nil
}
