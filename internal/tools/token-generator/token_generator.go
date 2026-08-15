// Package tokengen 实现随机 Token 生成工具。
package tokengen

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "token-generator"
	Name        = "Token 生成器"
	Description = "生成包含大小写字母、数字与符号的自定义随机 Token"
	Category    = registry.CategoryCrypto
	Icon        = "ArrowsShuffle"
)

// Keywords 为搜索关键词。
var Keywords = []string{"token", "random", "string", "alphanumeric", "symbols", "number", "letters", "lowercase", "uppercase", "password", "令牌", "随机"}

// 默认字符集（与 it-tools 对齐，但修正其遗漏 N/n 的字母表）。
const (
	uppercaseCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseCharset = "abcdefghijklmnopqrstuvwxyz"
	numbersCharset   = "0123456789"
	symbolsCharset   = ".,;:!?./-\"'#{([-|\\@)]=}*+"
)

// lengthMax 与前端滑块上限一致。
const lengthMax = 512

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Length        int    `json:"length"`         // Token 长度（1..512）
	WithUppercase bool   `json:"with_uppercase"` // 包含大写字母
	WithLowercase bool   `json:"with_lowercase"` // 包含小写字母
	WithNumbers   bool   `json:"with_numbers"`   // 包含数字
	WithSymbols   bool   `json:"with_symbols"`   // 包含符号
	Alphabet      string `json:"alphabet"`       // 可选，自定义字符集（优先于开关组合）
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 生成随机 Token 并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	charset, err := buildCharset(in)
	if err != nil {
		return "", err
	}
	if in.Length < 1 || in.Length > lengthMax {
		return "", fmt.Errorf("长度必须在 1..%d 之间（当前 %d）", lengthMax, in.Length)
	}

	result, err := cryptoToken(charset, in.Length)
	if err != nil {
		return "", fmt.Errorf("生成 Token 失败: %w", err)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// buildCharset 根据自定义字符集或开关组合构建候选字符集。
func buildCharset(in input) (string, error) {
	if in.Alphabet != "" {
		return in.Alphabet, nil
	}
	var sb strings.Builder
	if in.WithUppercase {
		sb.WriteString(uppercaseCharset)
	}
	if in.WithLowercase {
		sb.WriteString(lowercaseCharset)
	}
	if in.WithNumbers {
		sb.WriteString(numbersCharset)
	}
	if in.WithSymbols {
		sb.WriteString(symbolsCharset)
	}
	if sb.Len() == 0 {
		return "", fmt.Errorf("至少选择一种字符集（大写/小写/数字/符号）")
	}
	return sb.String(), nil
}

// cryptoToken 使用 crypto/rand 从 charset 中均匀逐位采样，生成长度为 n 的字符串。
func cryptoToken(charset string, n int) (string, error) {
	max := big.NewInt(int64(len(charset)))
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[idx.Int64()])
	}
	return sb.String(), nil
}
