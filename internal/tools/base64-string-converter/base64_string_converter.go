// Package base64string 实现 Base64 字符串编码/解码工具。
package base64string

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Tool 描述工具的元数据。
const (
	ID       = "base64-string-converter"
	Name     = "Base64 字符串编码/解码"
	Category = "转换器"
)

// input 是工具的输入结构。
type input struct {
	Text    string `json:"text"`     // 待处理文本
	Action  string `json:"action"`   // encode | decode
	URLSafe bool   `json:"url_safe"` // 是否使用 URL-safe 编码（- 与 _ 代替 + 与 /）
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

	var result string
	switch in.Action {
	case "encode":
		if in.URLSafe {
			result = base64.URLEncoding.EncodeToString([]byte(in.Text))
		} else {
			result = base64.StdEncoding.EncodeToString([]byte(in.Text))
		}
	case "decode":
		decoded, err := decode(in.Text, in.URLSafe)
		if err != nil {
			return "", fmt.Errorf("Base64 解码失败: %w", err)
		}
		result = string(decoded)
	default:
		return "", fmt.Errorf("未知操作: %q（仅支持 encode/decode）", in.Action)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// decode 解码 Base64 文本，urlSafe 为真时使用 URL-safe 字符集，并兼容标准字符。
func decode(s string, urlSafe bool) ([]byte, error) {
	if urlSafe {
		return base64.URLEncoding.DecodeString(s)
	}
	return base64.StdEncoding.DecodeString(s)
}
