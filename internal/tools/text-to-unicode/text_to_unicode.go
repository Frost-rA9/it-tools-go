// Package textunicode 实现文本与 Unicode 转义表示之间的转换。
package textunicode

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// 工具元数据。
const (
	ID       = "text-to-unicode"
	Name     = "文本转 Unicode"
	Category = "转换器"
)

// 转换模式。
const (
	ModeTextToUnicode = "text_to_unicode"
	ModeUnicodeToText = "unicode_to_text"
)

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"`
	Mode string `json:"mode"` // text_to_unicode | unicode_to_text
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
	case ModeTextToUnicode:
		result = TextToUnicode(in.Text)
	case ModeUnicodeToText:
		result = UnicodeToText(in.Text)
	default:
		return "", fmt.Errorf("未知模式: %q", in.Mode)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// TextToUnicode 将文本转换为 &#N; 形式的 Unicode 转义字符串。
func TextToUnicode(text string) string {
	var b []byte
	for _, r := range text {
		b = append(b, fmt.Sprintf("&#%d;", r)...)
	}
	return string(b)
}

// unicodeRegexp 匹配 &#N; 形式的 Unicode 转义。
var unicodeRegexp = regexp.MustCompile(`&#(\d+);`)

// UnicodeToText 将 Unicode 转义字符串还原为文本（保留未匹配的部分）。
func UnicodeToText(s string) string {
	return unicodeRegexp.ReplaceAllStringFunc(s, func(m string) string {
		n, err := strconv.ParseInt(m[2:len(m)-1], 10, 32)
		if err != nil {
			return m
		}
		return string(rune(n))
	})
}
