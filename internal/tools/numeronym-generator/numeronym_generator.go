// Package numeronymgen 实现数字名称生成器（如 internationalization → i18n）。
package numeronymgen

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "numeronym-generator"
	Name        = "数字名称生成器"
	Description = "生成数字名称缩写（如 internationalization → i18n）"
	Category    = "文本"
	Icon        = "AB"
)

// Keywords 为搜索关键词。
var Keywords = []string{"numeronym", "generator", "abbreviation", "i18n", "a11y", "l10n", "数字名称", "数字缩写"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Word string `json:"word"`
}

// output 是工具的输出结构。
type output struct {
	Numeronym string `json:"numeronym"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回数字名称缩写。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out, err := json.Marshal(output{Numeronym: Generate(in.Word)})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Generate 生成数字名称缩写（对齐参考项目 it-tools）：
// 长度 ≤ 3 返回原词；否则返回 首字符 + (长度-2) + 末字符。
// 以 Unicode 码点（rune）计数。
func Generate(word string) string {
	rs := []rune(word)
	if len(rs) <= 3 {
		return word
	}
	return string(rs[0]) + strconv.Itoa(len(rs)-2) + string(rs[len(rs)-1])
}