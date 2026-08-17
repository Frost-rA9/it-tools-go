// Package deviceinformation 实现设备信息工具。
// 该工具的屏幕/浏览器/系统信息均来自浏览器 API（window.screen / navigator），Go 后端无法访问，
// 因此由前端直接读取展示，本包仅提供占位注册以保证工具出现在列表中。
package deviceinformation

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "device-information"
	Name        = "设备信息"
	Description = "查看屏幕、窗口与浏览器环境信息"
	Category    = "Web"
	Icon        = "DeviceDesktop"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"device", "information", "screen", "pixel", "ratio", "status", "data", "computer", "size", "user", "agent", "设备", "信息", "屏幕"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
// 该工具数据由前端浏览器 API 直接读取，此处返回空对象。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in map[string]any
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	out, err := json.Marshal(map[string]any{})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}