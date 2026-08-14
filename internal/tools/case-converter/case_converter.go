// Package caseconv 实现字符串在多种大小写格式之间的转换。
package caseconv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "case-converter"
	Name        = "大小写转换"
	Description = "在多种大小写格式之间转换字符串（驼峰、蛇形、常量等）"
	Category    = "转换器"
	Icon        = "LetterCaseToggle"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"case", "大小写", "驼峰", "camel", "snake", "upper", "lower"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"`
}

// format 是单个大小写格式的结果。
type format struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// output 是工具的输出结构。
type output struct {
	Results []format `json:"results"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回全部格式的转换结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out, err := json.Marshal(output{Results: Convert(in.Text)})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// wordTransform 对单个词做大小写变换，index 为词在序列中的位置（从 0 开始）。
type wordTransform func(word string, index int) string

// Convert 将输入文本转换为 14 种大小写格式（对齐参考项目 it-tools 的 change-case 语义）。
func Convert(text string) []format {
	words := splitWords(text)

	return []format{
		{"Lowercase", strings.ToLower(text)},
		{"Uppercase", strings.ToUpper(text)},
		{"Camelcase", joinWords(words, "", func(w string, i int) string {
			if i == 0 {
				return strings.ToLower(w)
			}
			return capitalizeWord(w)
		})},
		{"Capitalcase", joinWords(words, " ", capitalizeWordTransform)},
		{"Constantcase", joinWords(words, "_", upperWordTransform)},
		{"Dotcase", joinWords(words, ".", lowerWordTransform)},
		{"Headercase", joinWords(words, "-", capitalizeWordTransform)},
		{"Nocase", joinWords(words, " ", lowerWordTransform)},
		{"Paramcase", joinWords(words, "-", lowerWordTransform)},
		{"Pascalcase", joinWords(words, "", capitalizeWordTransform)},
		{"Pathcase", joinWords(words, "/", lowerWordTransform)},
		{"Sentencecase", joinWords(words, " ", func(w string, i int) string {
			if i == 0 {
				return capitalizeWord(w)
			}
			return strings.ToLower(w)
		})},
		{"Snakecase", joinWords(words, "_", lowerWordTransform)},
		{"Mockingcase", mockingCase(text)},
	}
}

var (
	lowerWordTransform      = func(w string, _ int) string { return strings.ToLower(w) }
	upperWordTransform      = func(w string, _ int) string { return strings.ToUpper(w) }
	capitalizeWordTransform = func(w string, _ int) string { return capitalizeWord(w) }
)

// joinWords 用分隔符连接各词，并按 transform 变换每个词。
func joinWords(words []string, sep string, transform wordTransform) string {
	parts := make([]string, 0, len(words))
	for i, w := range words {
		parts = append(parts, transform(w, i))
	}
	return strings.Join(parts, sep)
}

// capitalizeWord 首字母大写、其余字母小写。
func capitalizeWord(w string) string {
	rs := []rune(w)
	if len(rs) == 0 {
		return w
	}
	rs[0] = unicode.ToUpper(rs[0])
	for i := 1; i < len(rs); i++ {
		rs[i] = unicode.ToLower(rs[i])
	}
	return string(rs)
}

// mockingCase 按字符下标交替大小写（偶数位大写、奇数位小写），对齐参考实现。
func mockingCase(s string) string {
	rs := []rune(s)
	for i := range rs {
		if i%2 == 0 {
			rs[i] = unicode.ToUpper(rs[i])
		} else {
			rs[i] = unicode.ToLower(rs[i])
		}
	}
	return string(rs)
}

// 分隔符：用于在拆分时标记词边界（对应 change-case 的 \0 占位符）。
const sepRune = rune(0)

// isLetter 判断 rune 是否为参考项目 stripRegexp 认可的字母
// （ASCII A-Za-z 及 À-Ö / Ø-ö / ø-ÿ 区段）。
func isLetter(r rune) bool {
	if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' {
		return true
	}
	return r >= 0x00C0 && r <= 0x00D6 || r >= 0x00D8 && r <= 0x00F6 || r >= 0x00F8 && r <= 0x00FF
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }
func isDigit(r rune) bool { return r >= '0' && r <= '9' }

// splitWords 将字符串拆分为词，处理词间分隔符与驼峰/缩写边界（对齐 change-case split）。
func splitWords(s string) []string {
	rs := []rune(s)
	if len(rs) == 0 {
		return nil
	}

	// 第一步：在大小写边界处插入分隔符。
	marked := make([]rune, 0, len(rs)*2)
	for i, r := range rs {
		if i > 0 {
			prev := rs[i-1]
			var next rune
			if i+1 < len(rs) {
				next = rs[i+1]
			}
			if isUpper(r) && (isLower(prev) || isDigit(prev)) {
				marked = append(marked, sepRune)
			} else if isUpper(r) && isUpper(prev) && next != 0 && isLower(next) {
				marked = append(marked, sepRune)
			}
		}
		marked = append(marked, r)
	}

	// 第二步：将非字母连串压缩为单个分隔符。
	cleaned := make([]rune, 0, len(marked))
	prevSep := true
	for _, r := range marked {
		if r == sepRune || !isLetter(r) {
			if !prevSep {
				cleaned = append(cleaned, sepRune)
				prevSep = true
			}
			continue
		}
		cleaned = append(cleaned, r)
		prevSep = false
	}

	parts := strings.Split(string(cleaned), string(sepRune))
	words := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			words = append(words, p)
		}
	}
	return words
}
