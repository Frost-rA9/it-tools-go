// Package ipv6ula 实现 IPv6 ULA 生成器：按 RFC 4193 方法 1（时间戳 + MAC → SHA1 → 低 40 bits）
// 生成唯一的本地 IPv6 地址（ULA）前缀，对齐 it-tools 的 ipv6-ula-generator。
package ipv6ula

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "ipv6-ula-generator"
	Name        = "IPv6 ULA 生成器"
	Description = "生成 RFC 4193 唯一的本地 IPv6 地址（ULA），基于时间戳与 MAC 地址"
	Category    = registry.CategoryNetwork
	Icon        = "BuildingFactory"
)

// Keywords 为搜索关键词。
var Keywords = []string{"ipv6", "ula", "generator", "rfc4193", "network", "private", "本地地址", "唯一本地", "生成器"}

// macRe 对齐 it-tools macAddressValidation：2-5 组“XX:/-”加末组 XX（共 3-6 字节）。
var macRe = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){2,5}([0-9A-Fa-f]{2})$`)

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	MACAddress string `json:"macAddress"` // MAC 地址，如 "20:37:06:12:34:56"
}

// output 是工具的输出结构。
type output struct {
	ULA                string `json:"ula"`                // fdXX:XXXX:XXXX::/48
	FirstRoutableBlock string `json:"firstRoutableBlock"` // fdXX:XXXX:XXXX:0::/64
	LastRoutableBlock  string `json:"lastRoutableBlock"`  // fdXX:XXXX:XXXX:ffff::/64
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 解析输入并返回生成的 ULA 前缀 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out, err := generate(in.MACAddress, time.Now().UnixMilli())
	if err != nil {
		return "", err
	}

	raw, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(raw), nil
}

// generate 校验 MAC 并按 ts 生成 ULA /48 前缀与两个可路由 /64 块。
// 算法：SHA1(ts + mac) 的十六进制取低 40 bits 作 Global ID，
// 前缀拆分为 fd + 2hex + ":" + 4hex + ":" + 4hex（与 it-tools 完全一致）。
func generate(mac string, ts int64) (output, error) {
	if !macRe.MatchString(mac) {
		return output{}, fmt.Errorf("MAC 地址格式无效: %q（期望如 d2:5f:61:07:3d:63）", mac)
	}

	sum := sha1.Sum([]byte(fmt.Sprintf("%d%s", ts, mac)))
	hex40 := fmt.Sprintf("%x", sum)[30:] // 后 10 个 hex 字符 = 低 40 bits

	prefix := "fd" + hex40[0:2] + ":" + hex40[2:6] + ":" + hex40[6:]
	return output{
		ULA:                prefix + "::/48",
		FirstRoutableBlock: prefix + ":0::/64",
		LastRoutableBlock:  prefix + ":ffff::/64",
	}, nil
}
