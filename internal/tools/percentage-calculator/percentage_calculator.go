// Package percentagecalculator 实现百分比计算：求比例数值、占比与增减幅度。
package percentagecalculator

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "percentage-calculator"
	Name        = "百分比计算器"
	Description = "计算百分比数值、占比与增减幅度"
	Category    = registry.CategoryMath
	Icon        = "Percentage"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"percentage", "percent", "百分比", "占比", "增减", "计算", "%"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// 计算模式。
const (
	ModePercentOf   = "percent_of"   // X% of Y：X 的百分之 Y 是多少
	ModeWhatPercent = "what_percent" // X 是 Y 的百分之几
	ModeChange      = "change"       // Y 相对 X 的增减百分比
)

// input 是工具的输入结构。X/Y 用指针以区分「未填写」与「0」。
type input struct {
	Mode string   `json:"mode"` // percent_of | what_percent | change
	X    *float64 `json:"x"`    // 第一个数值
	Y    *float64 `json:"y"`    // 第二个数值（change 模式下为 To 值）
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"` // 计算结果（无效输入时为空字符串）
}

// 各模式的计算公式。
var calcFuncs = map[string]func(x, y float64) float64{
	ModePercentOf:   func(x, y float64) float64 { return x / 100 * y },
	ModeWhatPercent: func(x, y float64) float64 { return 100 * x / y },
	ModeChange:      func(x, y float64) float64 { return (y - x) / x * 100 },
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	result := ""
	if in.X != nil && in.Y != nil {
		calc, ok := calcFuncs[in.Mode]
		if !ok {
			return "", fmt.Errorf("未知模式: %q", in.Mode)
		}
		v := calc(*in.X, *in.Y)
		// 除零等导致非有限结果时返回空（如 what_percent 的 Y=0、change 的 X=0）。
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			result = strconv.FormatFloat(v, 'f', -1, 64)
		}
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
