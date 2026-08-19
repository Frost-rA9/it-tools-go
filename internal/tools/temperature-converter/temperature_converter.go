// Package temperatureconverter 实现 8 种温标之间的温度换算（以开尔文为中间量）。
//
// 温标：Kelvin、Celsius、Fahrenheit、Rankine、Delisle、Newton、Réaumur、Rømer。
// 公式与参考项目 it-tools（temperature-converter.models.ts，GPLv3）保持一致。
package temperatureconverter

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "temperature-converter"
	Name        = "温度转换器"
	Description = "在 8 种温标之间换算温度"
	Category    = registry.CategoryMeasurement
	Icon        = "Temperature"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"temperature", "celsius", "kelvin", "fahrenheit", "rankine", "delisle", "newton", "reaumur", "romer", "温度", "摄氏", "华氏", "开尔文", "温标"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// scale 描述一种温标及与开尔文的相互换算函数。
type scale struct {
	Scale      string // 温标标识
	Title      string // 中文名称
	Unit       string // 单位符号
	toKelvin   func(float64) float64
	fromKelvin func(float64) float64
}

// scales 全部 8 种温标，按展示顺序排列。
var scales = []scale{
	{"kelvin", "开尔文", "K", func(v float64) float64 { return v }, func(v float64) float64 { return v }},
	{"celsius", "摄氏", "°C", func(v float64) float64 { return v + 273.15 }, func(v float64) float64 { return v - 273.15 }},
	{"fahrenheit", "华氏", "°F", func(v float64) float64 { return (v + 459.67) * (5.0 / 9) }, func(v float64) float64 { return v*(9.0/5) - 459.67 }},
	{"rankine", "兰金", "°R", func(v float64) float64 { return v * (5.0 / 9) }, func(v float64) float64 { return v * (9.0 / 5) }},
	{"delisle", "德利尔", "°De", func(v float64) float64 { return 373.15 - (2.0/3)*v }, func(v float64) float64 { return (3.0 / 2) * (373.15 - v) }},
	{"newton", "牛顿", "°N", func(v float64) float64 { return v*(100.0/33) + 273.15 }, func(v float64) float64 { return (v - 273.15) * (33.0 / 100) }},
	{"reaumur", "列氏", "°Ré", func(v float64) float64 { return v*(5.0/4) + 273.15 }, func(v float64) float64 { return (v - 273.15) * (4.0 / 5) }},
	{"romer", "罗氏", "°Rø", func(v float64) float64 { return (v-7.5)*(40.0/21) + 273.15 }, func(v float64) float64 { return (v-273.15)*(21.0/40) + 7.5 }},
}

// scaleIndex 温标标识 → 下标，便于按 from 定位。
var scaleIndex = func() map[string]int {
	m := make(map[string]int, len(scales))
	for i, s := range scales {
		m[s.Scale] = i
	}
	return m
}()

// input 是工具的输入结构。
type input struct {
	Value float64 `json:"value"` // 输入温度数值
	From  string  `json:"from"`  // 输入温标标识（scale 字段之一）
}

// result 是单个温标的换算结果。
type result struct {
	Scale string  `json:"scale"`
	Title string  `json:"title"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

// output 是工具的输出结构。
type output struct {
	Results []result `json:"results"` // 全部 8 个温标的换算结果
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	idx, ok := scaleIndex[in.From]
	if !ok {
		return "", fmt.Errorf("未知温标: %q", in.From)
	}

	kelvins := scales[idx].toKelvin(in.Value)
	results := make([]result, 0, len(scales))
	for _, s := range scales {
		results = append(results, result{
			Scale: s.Scale,
			Title: s.Title,
			Unit:  s.Unit,
			Value: Round2(s.fromKelvin(kelvins)),
		})
	}

	out, err := json.Marshal(output{Results: results})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Round2 保留 2 位小数（截断式，与参考实现 Math.floor(v*100)/100 语义一致，
// 加微小 epsilon 消除浮点尾差，如 32.00000000000002 不致被截成 31.99）。
func Round2(v float64) float64 {
	return math.Floor(v*100+1e-9) / 100
}
