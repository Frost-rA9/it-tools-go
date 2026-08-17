// Package useragentparser 实现用户代理分析器工具：解析 UA 的浏览器、操作系统与设备信息。
// 使用 uap-go（uap-core 规则库）；引擎/CPU/设备类型由前端 TS 补充（uap-go 无此数据）。
package useragentparser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ua-parser/uap-go/uaparser"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "user-agent-parser"
	Name        = "用户代理分析器"
	Description = "解析 User-Agent 字符串中的浏览器、操作系统与设备信息"
	Category    = "Web"
	Icon        = "Browser"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"user", "agent", "parser", "browser", "engine", "os", "cpu", "device", "user-agent", "client", "用户代理", "解析"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	UA string `json:"ua"` // 待解析的 User-Agent 字符串
}

// nameVersion 表示一个名称 + 版本。
type nameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// deviceInfo 表示设备信息。
type deviceInfo struct {
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
	Type   string `json:"type"`
}

// output 是工具的输出结构（引擎/CPU/设备类型由前端 TS 补充）。
type output struct {
	Browser nameVersion `json:"browser"`
	Os      nameVersion `json:"os"`
	Device  deviceInfo  `json:"device"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	ua := strings.TrimSpace(in.UA)
	out := output{}
	if ua != "" {
		parser, err := uaparser.New()
		if err != nil {
			return "", fmt.Errorf("初始化解析器失败: %w", err)
		}
		client := parser.Parse(ua)
		out.Browser = nameVersion{
			Name:    client.UserAgent.Family,
			Version: joinVersion(client.UserAgent.Major, client.UserAgent.Minor, client.UserAgent.Patch),
		}
		out.Os = nameVersion{
			Name:    client.Os.Family,
			Version: joinVersion(client.Os.Major, client.Os.Minor, client.Os.Patch),
		}
		out.Device = deviceInfo{
			Vendor: client.Device.Brand,
			Model:  client.Device.Model,
		}
		if out.Device.Model == "" {
			out.Device.Model = client.Device.Family
		}
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// joinVersion 拼接版本段，省略空段（如 "121.0.6167.129" 或仅 "1"）。
func joinVersion(parts ...string) string {
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	return strings.Join(nonEmpty, ".")
}