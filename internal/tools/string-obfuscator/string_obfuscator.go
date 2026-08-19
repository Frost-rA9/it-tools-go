// Package stringobfuscator 实现字符串混淆器（遮蔽中间部分）。
package stringobfuscator

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "string-obfuscator"
	Name        = "字符串混淆器"
	Description = "用替换字符遮蔽字符串的中间部分，可保留开头与结尾的字符数"
	Category    = "文本"
	Icon        = "EyeOff"
)

// Keywords 为搜索关键词。
var Keywords = []string{"string", "obfuscator", "secret", "token", "hide", "obscure", "mask", "masking", "混淆", "遮蔽", "打码", "隐藏"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text            string `json:"text"`
	KeepFirst       int    `json:"keep_first"`
	KeepLast        int    `json:"keep_last"`
	KeepSpace       bool   `json:"keep_space"`
	ReplacementChar string `json:"replacement_char"`
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回混淆结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out, err := json.Marshal(output{
		Result: Obfuscate(in.Text, in.KeepFirst, in.KeepLast, in.KeepSpace, in.ReplacementChar),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Obfuscate 按字符下标遮蔽字符串：
//   - keepSpace 且字符为 ASCII 空格时保留（对齐参考项目）；
//   - 下标在 [0, keepFirst) 或 [len-keepLast, len) 内保留原字符；
//   - 其余字符替换为 replacementChar（负数钳制为 0）。
//
// 以 Unicode 码点（rune）为单位处理，对 BMP 字符与参考项目 JS split('') 语义一致。
func Obfuscate(text string, keepFirst, keepLast int, keepSpace bool, replacementChar string) string {
	if keepFirst < 0 {
		keepFirst = 0
	}
	if keepLast < 0 {
		keepLast = 0
	}

	repl := rune('*')
	if rs := []rune(replacementChar); len(rs) > 0 {
		repl = rs[0]
	}

	rs := []rune(text)
	out := make([]rune, 0, len(rs))
	for i, r := range rs {
		if keepSpace && r == ' ' {
			out = append(out, r)
			continue
		}
		if i < keepFirst || i >= len(rs)-keepLast {
			out = append(out, r)
		} else {
			out = append(out, repl)
		}
	}
	return string(out)
}