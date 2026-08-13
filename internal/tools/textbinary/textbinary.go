// Package textbinary 实现文本与 ASCII 二进制表示之间的转换。
package textbinary

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// 工具元数据。
const (
	ID       = "text-to-binary"
	Name     = "文本转 ASCII 二进制"
	Category = "转换器"
)

// 转换模式。
const (
	ModeTextToBinary = "text_to_binary"
	ModeBinaryToText = "binary_to_text"
)

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"`
	Mode string `json:"mode"` // text_to_binary | binary_to_text
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
	switch in.Mode {
	case ModeTextToBinary:
		result = TextToBinary(in.Text)
	case ModeBinaryToText:
		text, err := BinaryToText(in.Text)
		if err != nil {
			return "", err
		}
		result = text
	default:
		return "", fmt.Errorf("未知模式: %q", in.Mode)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// TextToBinary 将文本转换为以空格分隔的 8 位二进制形式。
func TextToBinary(text string) string {
	parts := make([]string, 0, len(text))
	for _, r := range text {
		b := strconv.FormatInt(int64(r), 2)
		if len(b) < 8 {
			b = strings.Repeat("0", 8-len(b)) + b
		}
		parts = append(parts, b)
	}
	return strings.Join(parts, " ")
}

// BinaryToText 将二进制字符串还原为文本（自动剔除空格等非 01 字符），
// 长度不是 8 的倍数时返回错误。
func BinaryToText(binary string) (string, error) {
	clean := strings.Map(func(r rune) rune {
		if r == '0' || r == '1' {
			return r
		}
		return -1
	}, binary)

	if len(clean)%8 != 0 {
		return "", fmt.Errorf("无效的二进制字符串（长度须为 8 的倍数）")
	}

	var sb strings.Builder
	for i := 0; i < len(clean); i += 8 {
		v, err := strconv.ParseUint(clean[i:i+8], 2, 8)
		if err != nil {
			return "", fmt.Errorf("无效的二进制字符串: %w", err)
		}
		sb.WriteRune(rune(v))
	}
	return sb.String(), nil
}
