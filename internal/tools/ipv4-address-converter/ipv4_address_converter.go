// Package ipv4addressconv 实现 IPv4 地址转换器：将 IPv4 地址转换为十进制、十六进制、
// 二进制与 IPv6 映射形式，对齐 it-tools 的 ipv4-address-converter。
package ipv4addressconv

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "ipv4-address-converter"
	Name        = "IPv4 地址转换器"
	Description = "将 IPv4 地址转换为十进制、十六进制、二进制与 IPv6 映射形式"
	Category    = registry.CategoryNetwork
	Icon        = "Binary"
)

// Keywords 为搜索关键词。
var Keywords = []string{"ipv4", "address", "converter", "decimal", "hexadecimal", "binary", "ipv6", "地址", "转换"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	IP string `json:"ip"` // IPv4 地址，如 "192.168.1.1"
}

// output 是工具的输出结构。
type output struct {
	Decimal     uint32 `json:"decimal"`     // 十进制整数
	Hexadecimal string `json:"hexadecimal"` // 大写十六进制（8 位）
	Binary      string `json:"binary"`      // 二进制（32 位）
	IPv6        string `json:"ipv6"`        // IPv4 映射 IPv6（完整 8 组）
	IPv6Short   string `json:"ipv6Short"`   // IPv4 映射 IPv6（压缩形式）
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 解析输入并返回各进制转换结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	ip, err := parseIPv4(in.IP)
	if err != nil {
		return "", fmt.Errorf("IPv4 地址无效: %q", in.IP)
	}

	full, short := ipv6Mapped(ip)
	out := output{
		Decimal:     ip,
		Hexadecimal: fmt.Sprintf("%08X", ip),
		Binary:      fmt.Sprintf("%032b", ip),
		IPv6:        full,
		IPv6Short:   short,
	}
	raw, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(raw), nil
}

// ipv6Mapped 返回 IPv4 的 IPv4 映射 IPv6 地址（完整与压缩形式）。
func ipv6Mapped(ip uint32) (full, short string) {
	// 末两组 hextet：每 2 个 octet 一组，小写 hex。
	last := fmt.Sprintf("%02x%02x:%02x%02x", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	return "0000:0000:0000:0000:0000:ffff:" + last, "::ffff:" + last
}

// parseIPv4 将点分十进制 IPv4 解析为 uint32。
func parseIPv4(s string) (uint32, error) {
	parts := strings.Split(strings.TrimSpace(s), ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("需为 4 段地址")
	}
	var v uint32
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 || n > 255 {
			return 0, fmt.Errorf("段 %q 不是 0-255 的整数", p)
		}
		v |= uint32(n) << (8 * (3 - i))
	}
	return v, nil
}
