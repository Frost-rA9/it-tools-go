// Package jsonmin 实现 JSON 压缩工具：移除 JSON 文本中的空白字符。
package jsonmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

const (
	ID          = "json-minifier"
	Name        = "JSON 压缩"
	Description = "移除 JSON 中的空白，展示压缩前后大小对比"
	Category    = "开发"
	Icon        = "Brackets"
)

var Keywords = []string{"json", "minify", "compress", "compact", "压缩", "减小", "移除空白"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	JSON string `json:"json"` // 待压缩 JSON 文本
}

// output 是工具的输出结构。
type output struct {
	Minified     string  `json:"minified"`      // 压缩结果
	OriginalSize int     `json:"original_size"` // 原始字节数
	MinifiedSize int     `json:"minified_size"` // 压缩后字节数
	Saved        int     `json:"saved"`         // 节省字节数
	SavedPercent float64 `json:"saved_percent"` // 节省百分比
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回压缩结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(in.JSON)); err != nil {
		return "", fmt.Errorf("JSON 解析失败: %w", err)
	}
	minified := buf.String()

	original := len(in.JSON)
	compacted := len(minified)
	saved := original - compacted
	var percent float64
	if original > 0 {
		percent = float64(saved) / float64(original) * 100
	}

	out, err := json.Marshal(output{
		Minified:     minified,
		OriginalSize: original,
		MinifiedSize: compacted,
		Saved:        saved,
		SavedPercent: percent,
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
