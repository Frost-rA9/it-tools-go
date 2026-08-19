// Package macaddressgen 实现 MAC 地址生成器：按前缀、大小写与分隔符随机生成
// MAC 地址，对齐 it-tools 的 mac-address-generator。
package macaddressgen

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "mac-address-generator"
	Name        = "MAC 地址生成器"
	Description = "按前缀、大小写与分隔符随机生成 MAC 地址"
	Category    = registry.CategoryNetwork
	Icon        = "Devices"
)

// Keywords 为搜索关键词。
var Keywords = []string{"mac", "address", "generator", "random", "prefix", "mac 地址", "生成器"}

// 数量上限（对齐 it-tools）。
const countMax = 100

var (
	hexRe   = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	hasSep  = regexp.MustCompile(`[^0-9a-fA-F]`)
	sepRe   = regexp.MustCompile(`[^0-9a-fA-F]+`)
	validRe = regexp.MustCompile(`^[0-9a-fA-F:\-\. ]+$`)
)

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Count     int    `json:"count"`     // 生成数量（1..100）
	Prefix    string `json:"prefix"`    // 可选前缀（hex/分隔符混用，最多 6 字节）
	Separator string `json:"separator"` // ":" | "-" | "." | ""（无）
	Case      string `json:"case"`      // "upper" | "lower"
}

// output 是工具的输出结构。
type output struct {
	MACAddresses string `json:"macAddresses"` // 多行连接
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 解析输入并返回生成的 MAC 列表 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	macs, err := generate(in.Count, in.Prefix, in.Separator, in.Case)
	if err != nil {
		return "", err
	}
	out := output{MACAddresses: strings.Join(macs, "\n")}
	raw, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(raw), nil
}

// generate 生成 count 条 MAC 地址，返回列表。
func generate(count int, prefix, separator, caseMode string) ([]string, error) {
	if count < 1 || count > countMax {
		return nil, fmt.Errorf("数量必须在 1..%d 之间（当前 %d）", countMax, count)
	}
	if separator != ":" && separator != "-" && separator != "." && separator != "" {
		return nil, fmt.Errorf("分隔符无效: %q（仅支持 : - . 或无）", separator)
	}
	if caseMode != "" && caseMode != "upper" && caseMode != "lower" {
		return nil, fmt.Errorf("大小写模式无效: %q（仅支持 upper/lower）", caseMode)
	}

	prefixBytes, err := splitPrefix(prefix)
	if err != nil {
		return nil, err
	}

	macs := make([]string, count)
	for i := 0; i < count; i++ {
		mac, err := randomMAC(prefixBytes, separator, caseMode)
		if err != nil {
			return nil, err
		}
		macs[i] = mac
	}
	return macs, nil
}

// splitPrefix 拆分前缀为字节 hex 列表（对齐 it-tools splitPrefix + padStart）。
// 纯 hex 按每 2 字符切；含分隔符按非 hex 切；每字节补零到 2 位；最多 6 字节。
func splitPrefix(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if !validRe.MatchString(raw) {
		return nil, fmt.Errorf("前缀只能包含十六进制字符或分隔符 : / - / . / 空格")
	}

	var segs []string
	if hasSep.MatchString(raw) {
		segs = sepRe.Split(raw, -1)
	} else {
		for i := 0; i < len(raw); i += 2 {
			end := i + 2
			if end > len(raw) {
				end = len(raw)
			}
			segs = append(segs, raw[i:end])
		}
	}

	bytes := make([]string, 0, 6)
	for _, s := range segs {
		if s == "" {
			continue
		}
		if len(s) > 2 || !hexRe.MatchString(s) {
			return nil, fmt.Errorf("前缀段 %q 无效（每段需为 1-2 位十六进制）", s)
		}
		// 补零到 2 位，保留原始大小写（对齐 it-tools padStart）。
		if len(s) == 1 {
			s = "0" + s
		}
		bytes = append(bytes, s)
	}
	if len(bytes) > 6 {
		return nil, fmt.Errorf("前缀最多 6 字节（当前 %d 字节）", len(bytes))
	}
	return bytes, nil
}

// randomMAC 由前缀 + 随机补足字节构建一条 MAC 地址，应用分隔符与大小写。
func randomMAC(prefixBytes []string, separator, caseMode string) (string, error) {
	bytes := make([]string, 6)
	n := copy(bytes, prefixBytes)

	rest := make([]byte, 6-n)
	if _, err := crand.Read(rest); err != nil {
		return "", fmt.Errorf("读取安全随机数失败: %w", err)
	}
	for i, b := range rest {
		bytes[n+i] = fmt.Sprintf("%02x", b)
	}

	mac := strings.Join(bytes, separator)
	switch caseMode {
	case "upper":
		mac = strings.ToUpper(mac)
	case "lower":
		mac = strings.ToLower(mac)
	}
	return mac, nil
}
