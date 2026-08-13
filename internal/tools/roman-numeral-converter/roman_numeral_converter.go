// Package romannumeral 实现阿拉伯数字与罗马数字之间的转换。
package romannumeral

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 工具元数据。
const (
	ID       = "roman-numeral-converter"
	Name     = "罗马数字转换器"
	Category = "转换器"

	MinArabic = 1
	MaxArabic = 3999
)

// 转换模式。
const (
	ModeArabicToRoman = "arabic_to_roman"
	ModeRomanToArabic = "roman_to_arabic"
)

// input 是工具的输入结构。
type input struct {
	Value string `json:"value"` // 阿拉伯数字或罗马数字字符串
	Mode  string `json:"mode"`  // arabic_to_roman | roman_to_arabic
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"` // 转换结果（无效输入时为空字符串）
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
	case ModeArabicToRoman:
		n, err := strconv.Atoi(strings.TrimSpace(in.Value))
		if err != nil {
			return "", fmt.Errorf("无效的数字: %q", in.Value)
		}
		result = ArabicToRoman(n)
	case ModeRomanToArabic:
		n, ok := RomanToArabic(strings.ToUpper(strings.TrimSpace(in.Value)))
		if ok {
			result = strconv.Itoa(n)
		}
	default:
		return "", fmt.Errorf("未知模式: %q", in.Mode)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// 罗马数字符号表，按数值降序排列（映射关系有唯一最优拆分）。
var romanTable = []struct {
	value int
	glyph string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

// ArabicToRoman 将阿拉伯数字转换为罗马数字，超出 [1, 3999] 范围返回空字符串。
func ArabicToRoman(num int) string {
	if num < MinArabic || num > MaxArabic {
		return ""
	}
	var b strings.Builder
	for _, r := range romanTable {
		for num >= r.value {
			b.WriteString(r.glyph)
			num -= r.value
		}
	}
	return b.String()
}

// romanRegexp 校验标准罗马数字（匹配 1~3999 的合法写法）。
var romanRegexp = regexp.MustCompile(`^M{0,3}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$`)

// isValidRomanNumber 判断是否为合法的罗马数字（空字符串视为非法）。
func isValidRomanNumber(s string) bool {
	return s != "" && romanRegexp.MatchString(s)
}

// romanValues 罗马数字字符对应的数值。
var romanValues = map[byte]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}

// RomanToArabic 将罗马数字转换为阿拉伯数字，非法输入返回 ok=false。
func RomanToArabic(s string) (int, bool) {
	if !isValidRomanNumber(s) {
		return 0, false
	}
	total := 0
	for i := 0; i < len(s); i++ {
		cur := romanValues[s[i]]
		if i+1 < len(s) && romanValues[s[i+1]] > cur {
			total -= cur
		} else {
			total += cur
		}
	}
	return total, true
}
