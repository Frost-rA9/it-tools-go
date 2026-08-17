// Package keycodeinfo 实现按键码信息工具。
// 该工具需监听浏览器键盘事件（document keydown）以获取 event.key/keyCode/code/location/修饰键，
// Go 后端无法捕获，因此由前端直接读取展示，本包仅提供占位注册以保证工具出现在列表中。
package keycodeinfo

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "keycode-info"
	Name        = "按键码信息"
	Description = "显示键盘按键的 key、keyCode、code、location 与修饰键信息"
	Category    = "Web"
	Icon        = "Keyboard"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"keycode", "info", "code", "key", "javascript", "event", "keycodes", "which", "keyboard", "press", "modifier", "alt", "ctrl", "meta", "shift", "按键", "键码", "键盘"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
// 该工具数据由前端键盘事件直接读取，此处返回空对象。
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