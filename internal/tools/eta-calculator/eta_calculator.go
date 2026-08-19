// Package etacalculator 实现 ETA 计算：根据单位数量与消耗速率推算总时长与结束时间。
package etacalculator

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "eta-calculator"
	Name        = "ETA 计算器"
	Description = "根据单位数量与消耗速率推算总时长与结束时间"
	Category    = registry.CategoryMath
	Icon        = "Hourglass"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"eta", "calculator", "estimate", "estimated", "time", "arrival", "duration", "预计", "耗时", "估算"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	UnitCount              float64 `json:"unitCount"`              // 待消耗的单位总数（如 500 个盘子）
	UnitPerTimeSpan        float64 `json:"unitPerTimeSpan"`        // 每个时间段消耗的单位数（如 3 个盘子）
	TimeSpan               float64 `json:"timeSpan"`               // 时间段数值（如 5）
	TimeSpanUnitMultiplier float64 `json:"timeSpanUnitMultiplier"` // 时间段单位换算毫秒（秒 1000 / 分 60000 / 时 3600000 / 天 86400000）
	StartedAtMs            int64   `json:"startedAtMs"`            // 消耗开始时间（Unix 毫秒时间戳）
}

// output 是工具的输出结构。
type output struct {
	DurationMs   float64 `json:"durationMs"`   // 总时长（毫秒）
	DurationText string  `json:"durationText"` // 中文总时长描述（如 "5 小时 10 分"）
	EndAtMs      int64   `json:"endAtMs"`      // 结束时间（Unix 毫秒时间戳）
	EndAtText    string  `json:"endAtText"`    // 中文结束时间描述（相对 + 绝对）
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out := output{}
	valid := in.UnitCount > 0 && in.UnitPerTimeSpan > 0 && in.TimeSpan > 0 && in.TimeSpanUnitMultiplier > 0
	if valid {
		timeSpanMs := in.TimeSpan * in.TimeSpanUnitMultiplier
		durationMs := in.UnitCount / (in.UnitPerTimeSpan / timeSpanMs)
		if isFinite(durationMs) && durationMs >= 0 {
			endAtMs := in.StartedAtMs + int64(durationMs)
			out = output{
				DurationMs:   durationMs,
				DurationText: DurationText(durationMs),
				EndAtMs:      endAtMs,
				EndAtText:    EndAtText(durationMs, in.StartedAtMs, endAtMs),
			}
		}
	}

	res, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(res), nil
}

func isFinite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// DurationText 将毫秒时长格式化为中文描述，0 值分段省略，全零时输出 "0 毫秒"。
func DurationText(ms float64) string {
	total := int64(ms)
	hours := total / (3600 * 1000)
	mins := (total / (60 * 1000)) % 60
	secs := (total / 1000) % 60
	millis := total % 1000

	parts := []struct {
		n    int64
		unit string
	}{
		{hours, "小时"}, {mins, "分"}, {secs, "秒"}, {millis, "毫秒"},
	}
	var b []byte
	for _, p := range parts {
		if p.n > 0 {
			if len(b) > 0 {
				b = append(b, ' ')
			}
			b = strconv.AppendInt(b, p.n, 10)
			b = append(b, ' ')
			b = append(b, []byte(p.unit)...)
		}
	}
	if len(b) == 0 {
		return "0 毫秒"
	}
	return string(b)
}

// EndAtText 生成结束时间的中文描述：相对当前时间的近似描述 + 绝对时间。
func EndAtText(durationMs float64, startedAtMs, endAtMs int64) string {
	rel := time.Duration(endAtMs-time.Now().UnixMilli()) * time.Millisecond
	relText := RelativeText(rel)
	absText := time.UnixMilli(endAtMs).Format("2006-01-02 15:04")
	if relText == "" {
		return absText
	}
	return relText + "（" + absText + "）"
}

// RelativeText 将相对时长（可为负）转换为中文近似描述；接近零时返回空串。
func RelativeText(d time.Duration) string {
	if d < 0 {
		d = -d
		neg := formatAgo(d)
		if neg == "" {
			return ""
		}
		return "已结束 " + neg
	}
	return formatAgo(d)
}

// formatAgo 将非负时长格式化为"约 X 天后/小时后/分钟后/秒后"。
func formatAgo(d time.Duration) string {
	switch {
	case d < 5*time.Second:
		return ""
	case d < time.Minute:
		return fmt.Sprintf("约 %d 秒后", int((d+500*time.Millisecond)/time.Second))
	case d < time.Hour:
		return fmt.Sprintf("约 %d 分钟后", int((d+30*time.Second)/time.Minute))
	case d < 24*time.Hour:
		return fmt.Sprintf("约 %d 小时后", int((d+30*time.Minute)/time.Hour))
	default:
		return fmt.Sprintf("约 %d 天后", int((d+12*time.Hour)/(24*time.Hour)))
	}
}
