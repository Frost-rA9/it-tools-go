// Package ipv4subnet 实现 IPv4 子网计算器：解析 CIDR/掩码，输出网络地址、掩码、
// 可用主机、广播地址、IP 分类与相邻块，对齐 it-tools 的 ipv4-subnet-calculator。
package ipv4subnet

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
	ID          = "ipv4-subnet-calculator"
	Name        = "IPv4 子网计算器"
	Description = "计算 IPv4 子网信息：网络地址、掩码、可用主机、广播地址与相邻块"
	Category    = registry.CategoryNetwork
	Icon        = "Router"
)

// Keywords 为搜索关键词。
var Keywords = []string{"ipv4", "subnet", "calculator", "mask", "network", "cidr", "netmask", "bitmask", "broadcast", "子网", "掩码"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Address string `json:"address"` // 如 "192.168.0.1/24"、"192.168.0.0/255.255.255.0" 或纯 IP（视为 /32）
}

// output 是工具的输出结构。
type output struct {
	Netmask          string `json:"netmask"`           // base/前缀，如 192.168.0.0/24
	NetworkAddress   string `json:"networkAddress"`    // 网络地址
	NetworkMask      string `json:"networkMask"`       // 子网掩码
	NetworkMaskBin   string `json:"networkMaskBinary"` // 掩码二进制（8 位分组）
	CIDRNotation     string `json:"cidrNotation"`      // /24
	WildcardMask     string `json:"wildcardMask"`      // 通配符掩码
	NetworkSize      uint64 `json:"networkSize"`       // 地址总数
	FirstAddress     string `json:"firstAddress"`      // 首个可用地址
	LastAddress      string `json:"lastAddress"`       // 末个可用地址
	BroadcastAddress string `json:"broadcastAddress"`  // 广播地址
	IPClass          string `json:"ipClass"`           // A-E
	NextBlock        string `json:"nextBlock"`         // 相邻下一块 base/前缀
	PrevBlock        string `json:"prevBlock"`         // 相邻上一块 base/前缀
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 解析输入并返回子网信息 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	out, err := compute(in.Address)
	if err != nil {
		return "", err
	}
	raw, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(raw), nil
}

// compute 解析地址并计算全部子网字段。
func compute(address string) (output, error) {
	ipInt, prefix, err := parseAddress(address)
	if err != nil {
		return output{}, err
	}

	mask := maskFromPrefix(prefix)
	base := ipInt & mask
	hostmask := ^mask
	size := uint64(1) << (32 - prefix)

	base64 := uint64(base)
	broadcast64 := base64 + size - 1
	first, last := base64, broadcast64
	if size > 2 {
		first = base64 + 1
		last = broadcast64 - 1
	}

	next := base64 + size
	if next > 0xFFFFFFFF {
		next = 0xFFFFFFFF
	}
	prev := int64(base64) - int64(size)
	if prev < 0 {
		prev = 0
	}

	return output{
		Netmask:          uintToString(base) + "/" + strconv.Itoa(prefix),
		NetworkAddress:   uintToString(base),
		NetworkMask:      uintToString(mask),
		NetworkMaskBin:   maskBinary(mask),
		CIDRNotation:     "/" + strconv.Itoa(prefix),
		WildcardMask:     uintToString(hostmask),
		NetworkSize:      size,
		FirstAddress:     uintToString(uint32(first)),
		LastAddress:      uintToString(uint32(last)),
		BroadcastAddress: uintToString(uint32(broadcast64)),
		IPClass:          ipClass(base),
		NextBlock:        uintToString(uint32(next)) + "/" + strconv.Itoa(prefix),
		PrevBlock:        uintToString(uint32(prev)) + "/" + strconv.Itoa(prefix),
	}, nil
}

// parseAddress 解析 "ip/prefix"、"ip/掩码" 或纯 "ip"，返回 (ipInt, prefix)。
func parseAddress(address string) (uint32, int, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return 0, 0, fmt.Errorf("地址不能为空")
	}

	ipPart := address
	prefix := 32
	if slash := strings.IndexByte(address, '/'); slash >= 0 {
		ipPart = address[:slash]
		maskPart := strings.TrimSpace(address[slash+1:])
		if maskPart == "" {
			return 0, 0, fmt.Errorf("掩码缺失：%q", address)
		}
		if strings.Contains(maskPart, ".") {
			// 点分掩码 → 前缀长度。
			m, err := parseIPv4(maskPart)
			if err != nil {
				return 0, 0, fmt.Errorf("掩码格式无效: %q", maskPart)
			}
			prefix = prefixFromMask(m)
			if prefix < 0 {
				return 0, 0, fmt.Errorf("掩码 %q 不是连续的 1 序列", maskPart)
			}
		} else {
			n, err := strconv.Atoi(maskPart)
			if err != nil || n < 0 || n > 32 {
				return 0, 0, fmt.Errorf("前缀长度无效: %q（需为 0-32）", maskPart)
			}
			prefix = n
		}
	}

	ipInt, err := parseIPv4(ipPart)
	if err != nil {
		return 0, 0, fmt.Errorf("IPv4 地址无效: %q", ipPart)
	}
	return ipInt, prefix, nil
}

// parseIPv4 将点分十进制 IPv4 解析为 uint32。
func parseIPv4(s string) (uint32, error) {
	parts := strings.Split(s, ".")
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

func maskFromPrefix(prefix int) uint32 {
	if prefix == 0 {
		return 0
	}
	return ^((1 << (32 - prefix)) - 1)
}

func prefixFromMask(mask uint32) int {
	// 掩码须为连续的 1 序列（高位对齐）。
	for n := 0; n <= 32; n++ {
		if mask == maskFromPrefix(n) {
			return n
		}
	}
	return -1
}

// maskBinary 将掩码格式化为 8 位分组的二进制。
func maskBinary(mask uint32) string {
	var parts [4]string
	for i := 0; i < 4; i++ {
		parts[i] = fmt.Sprintf("%08b", byte(mask>>(8*(3-i))))
	}
	return strings.Join(parts[:], ".")
}

// ipClass 按首字节返回 IP 分类（对齐 it-tools getIPClass）。
func ipClass(ip uint32) string {
	first := byte(ip >> 24)
	switch {
	case first < 128:
		return "A"
	case first < 192:
		return "B"
	case first < 224:
		return "C"
	case first < 240:
		return "D"
	default:
		return "E"
	}
}

// uintToString 将 uint32 格式化为点分十进制。
func uintToString(v uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
