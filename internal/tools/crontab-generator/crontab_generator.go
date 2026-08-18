// Package crontab 实现 Crontab 表达式生成器：表单生成 cron 表达式，或直接解析表达式为可读描述。
package crontab

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"it-tools-go/internal/registry"
)

const (
	ID          = "crontab-generator"
	Name        = "Crontab 表达式生成器"
	Description = "表单生成 cron 表达式，或输入表达式实时解析为可读描述"
	Category    = "开发"
	Icon        = "Clock"
)

var Keywords = []string{"cron", "crontab", "cronjob", "定时任务", "表达式", "调度", "schedule", "定时"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// fieldInput 描述单个 cron 字段的选择。
type fieldInput struct {
	Type   string `json:"type"`   // every | step | range | list
	From   int    `json:"from"`   // range 起点
	To     int    `json:"to"`     // range 终点
	Step   int    `json:"step"`   // step 步长
	Values []int  `json:"values"` // list 取值
}

// input 是工具的输入结构。
type input struct {
	Action     string     `json:"action"`     // generate（表单，默认）| parse（直接解析表达式）
	Expression string     `json:"expression"` // parse 模式下的 cron 表达式
	Minute     fieldInput `json:"minute"`
	Hour       fieldInput `json:"hour"`
	DayOfMonth fieldInput `json:"day_of_month"`
	Month      fieldInput `json:"month"`
	DayOfWeek  fieldInput `json:"day_of_week"`
}

// output 是工具的输出结构。
type output struct {
	Expression string `json:"expression"` // 标准 5 段 cron 表达式
	Summary    string `json:"summary"`    // 人类可读描述
}

// fieldMeta 描述字段的取值边界与中文单位，用于校验与描述生成。
type fieldMeta struct {
	name   string // 中文名
	min    int    // 允许的最小值
	max    int    // 允许的最大值
	unit   string // 描述单位，如 "分钟" / "小时"
	suffix string // 描述后缀，如 "执行"
}

var fieldsOrder = []struct {
	key  string
	meta fieldMeta
}{
	{key: "minute", meta: fieldMeta{name: "分钟", min: 0, max: 59, unit: "分钟"}},
	{key: "hour", meta: fieldMeta{name: "小时", min: 0, max: 23, unit: "小时"}},
	{key: "day_of_month", meta: fieldMeta{name: "日", min: 1, max: 31, unit: "天"}},
	{key: "month", meta: fieldMeta{name: "月", min: 1, max: 12, unit: "月"}},
	{key: "day_of_week", meta: fieldMeta{name: "星期", min: 0, max: 7, unit: "周"}},
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入：generate 由表单字段生成表达式，parse 直接解析表达式为描述。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	if in.Action == "parse" {
		desc, err := parseCron(in.Expression)
		if err != nil {
			return "", err
		}
		out, err := json.Marshal(output{Expression: in.Expression, Summary: desc})
		if err != nil {
			return "", fmt.Errorf("序列化输出失败: %w", err)
		}
		return string(out), nil
	}

	expr, summary, err := build(in)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Expression: expr, Summary: summary})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// build 生成 5 段表达式与描述文本。
func build(in input) (string, string, error) {
	fields := map[string]fieldInput{
		"minute":       in.Minute,
		"hour":         in.Hour,
		"day_of_month": in.DayOfMonth,
		"month":        in.Month,
		"day_of_week":  in.DayOfWeek,
	}

	parts := make([]string, 0, 5)
	descs := make([]string, 0, 5)
	for _, f := range fieldsOrder {
		part, desc, err := renderField(fields[f.key], f.meta)
		if err != nil {
			return "", "", fmt.Errorf("%s字段: %w", f.meta.name, err)
		}
		parts = append(parts, part)
		descs = append(descs, desc)
	}

	// 描述主体：各字段片段逗号拼接
	summary := strings.Join(descs, "，")
	return strings.Join(parts, " "), summary, nil
}

// renderField 将单个字段输入渲染为表达式片段与描述片段。
func renderField(f fieldInput, m fieldMeta) (string, string, error) {
	switch f.Type {
	case "", "every":
		return "*", "每" + m.unit, nil
	case "step":
		step := f.Step
		if step < 1 {
			return "", "", fmt.Errorf("步长必须 ≥ 1（当前 %d）", step)
		}
		return fmt.Sprintf("*/%d", step), fmt.Sprintf("每 %d %s", step, m.unit), nil
	case "range":
		if f.From < m.min || f.To > m.max || f.From > f.To {
			return "", "", fmt.Errorf("范围 %d-%d 超出 %d-%d", f.From, f.To, m.min, m.max)
		}
		return fmt.Sprintf("%d-%d", f.From, f.To), fmt.Sprintf("%d-%d %s", f.From, f.To, m.unit), nil
	case "list":
		if len(f.Values) == 0 {
			return "", "", fmt.Errorf("列表不能为空")
		}
		vals := make([]int, len(f.Values))
		copy(vals, f.Values)
		for _, v := range vals {
			if v < m.min || v > m.max {
				return "", "", fmt.Errorf("取值 %d 超出 %d-%d", v, m.min, m.max)
			}
		}
		sort.Ints(vals)
		parts := make([]string, len(vals))
		for i, v := range vals {
			parts[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(parts, ","), fmt.Sprintf("%s %s", strings.Join(parts, "、"), m.unit), nil
	default:
		return "", "", fmt.Errorf("未知模式 %q（仅支持 every/step/range/list）", f.Type)
	}
}
