// Package ipv4rangeexpander 实现 IPv4 范围扩展器：计算覆盖指定起始与结束地址的
// 最小 CIDR 块，对齐 it-tools 的 ipv4-range-expander。
package ipv4rangeexpander

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
	ID          = "ipv4-range-expander"
	Name        = "IPv4 范围扩展器"
	Description = "计算覆盖指定起始与结束地址的最小 CIDR 块"
	Category    = registry.CategoryNetwork
	Icon        = "ArrowsMaximize"
)

// Keywords 为搜索关键词。
var Keywords = []string{"ipv4", "range", "expander", "subnet", "cidr", "creator", "地址", "范围", "子网"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	StartIP string `json:"startIp"` // 起始地址，如 "192.168.1.1"
	EndIP   string `json:"endIp"`   // 结束地址，如 "192.168.6.255"
}

// output 是工具的输出结构。
type output struct {
	OldSize  uint64 `json:"oldSize"`  // 原范围内地址数
	NewStart string `json:"newStart"` // 覆盖 CIDR 的起始地址
	NewEnd   string `json:"newEnd"`   // 覆盖 CIDR 的结束地址
	NewCIDR  string `json:"newCidr"`  // 覆盖 CIDR，如 192.168.0.0/21
	NewSize  uint64 `json:"newSize"`  // 覆盖 CIDR 的地址数
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 解析输入并返回覆盖 CIDR 结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	out, err := calculate(in.StartIP, in.EndIP)
	if err != nil {
		return "", err
	}
	raw, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(raw), nil
}

// calculate 计算覆盖 [startIp, endIp] 的最小 CIDR 块。
func calculate(startIP, endIP string) (output, error) {
	start, err := parseIPv4(startIP)
	if err != nil {
		return output{}, fmt.Errorf("起始地址无效: %q", startIP)
	}
	end, err := parseIPv4(endIP)
	if err != nil {
		return output{}, fmt.Errorf("结束地址无效: %q", endIP)
	}
	if end < start {
		return output{}, fmt.Errorf("结束地址 %s 低于起始地址 %s，请交换两者", endIP, startIP)
	}

	// 最小覆盖掩码：找 start 与 end 最高不同位。
	mask := 32
	if x := start ^ end; x != 0 {
		mask = 31 - highestBit(x)
	}
	hostBits := 32 - mask
	hostMask := uint32(0)
	if hostBits > 0 {
		hostMask = (uint32(1) << hostBits) - 1
	}
	newStart := start &^ hostMask
	newEnd := newStart | hostMask

	oldSize := uint64(end) - uint64(start) + 1
	newSize := uint64(1) << hostBits

	return output{
		OldSize:  oldSize,
		NewStart: uintToString(newStart),
		NewEnd:   uintToString(newEnd),
		NewCIDR:  uintToString(newStart) + "/" + strconv.Itoa(mask),
		NewSize:  newSize,
	}, nil
}

// highestBit 返回 v 的最高置位位位置（0 为最低位）。
func highestBit(v uint32) int {
	p := 0
	for v > 1 {
		v >>= 1
		p++
	}
	return p
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

// uintToString 将 uint32 格式化为点分十进制。
func uintToString(v uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
