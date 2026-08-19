// Package chronometer 实现秒表工具。
// 计时与格式化由前端实时展示（毫秒级 UI），Go 侧仅提供工具注册元数据。
package chronometer

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "chronometer"
	Name        = "秒表"
	Description = "计时器，支持开始、暂停与重置"
	Category    = registry.CategoryMeasurement
	Icon        = "Clock"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"chronometer", "stopwatch", "timer", "time", "lap", "duration", "秒表", "计时", "测速"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// Executor 实现 registry.Executor 接口。秒表为纯前端实时交互工具。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var input map[string]any
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	return `{}`, nil
}
