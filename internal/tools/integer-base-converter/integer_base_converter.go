// Package radix 实现整数在不同进制（2~64）之间的转换。
package radix

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "integer-base-converter"
	Name        = "整数基转换器"
	Description = "在不同进制之间转换整数（二进制、八进制、十进制、十六进制、Base64 等）"
	Category    = "转换器"
	Icon        = "ArrowsLeftRight"

	minBase = 2
	maxBase = 64
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"integer", "base", "进制", "转换", "binary", "decimal", "hex", "octal", "radix"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// alphabet 为进制数字字母表（与参考项目一致：0-9 a-z A-Z + /）。
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"

// input 是工具的输入结构。
type input struct {
	Value      string `json:"value"`       // 待转换的数字字符串
	FromBase   int    `json:"from_base"`   // 输入进制（2~64）
	CustomBase int    `json:"custom_base"` // 自定义输出进制（2~64）
}

// format 是单个进制结果的标签与值。
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

// Execute 解析输入并返回 6 个目标进制的结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	targets := []struct {
		label string
		base  int
	}{
		{"Binary (2)", 2},
		{"Octal (8)", 8},
		{"Decimal (10)", 10},
		{"Hexadecimal (16)", 16},
		{"Base64 (64)", 64},
		{fmt.Sprintf("Custom (%d)", in.CustomBase), in.CustomBase},
	}

	results := make([]format, 0, len(targets))
	for _, t := range targets {
		v, err := convertBase(in.Value, in.FromBase, t.base)
		if err != nil {
			return "", err
		}
		results = append(results, format{Label: t.label, Value: v})
	}

	out, err := json.Marshal(output{Results: results})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// convertBase 将 value 从 fromBase 进制转换为 toBase 进制。
func convertBase(value string, fromBase, toBase int) (string, error) {
	if fromBase < minBase || fromBase > maxBase {
		return "", fmt.Errorf("输入进制必须在 %d~%d 之间", minBase, maxBase)
	}
	if toBase < minBase || toBase > maxBase {
		return "", fmt.Errorf("输出进制必须在 %d~%d 之间", minBase, maxBase)
	}

	fromAlphabet := alphabet[:fromBase]
	toAlphabet := alphabet[:toBase]

	dec := new(big.Int)
	for _, r := range value {
		idx := strings.IndexRune(fromAlphabet, r)
		if idx < 0 {
			return "", fmt.Errorf("无效数字 %q（进制 %d）", string(r), fromBase)
		}
		dec.Mul(dec, big.NewInt(int64(fromBase)))
		dec.Add(dec, big.NewInt(int64(idx)))
	}

	if dec.Sign() == 0 {
		return "0", nil
	}

	base := big.NewInt(int64(toBase))
	mod := new(big.Int)
	var b []byte
	for dec.Sign() > 0 {
		dec.DivMod(dec, base, mod)
		b = append([]byte{toAlphabet[mod.Int64()]}, b...)
	}
	return string(b), nil
}
