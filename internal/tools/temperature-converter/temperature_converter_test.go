package temperatureconverter

import (
	"context"
	"encoding/json"
	"math"
	"testing"

	"it-tools-go/internal/registry"
)

// convert 辅助：给定源温标与数值，返回全部换算结果（scale→value）。
func convert(t *testing.T, from string, value float64) map[string]float64 {
	t.Helper()
	raw, err := json.Marshal(input{Value: value, From: from})
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}
	outStr, err := (Executor{}).Execute(context.Background(), string(raw))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	m := make(map[string]float64, len(out.Results))
	for _, r := range out.Results {
		m[r.Scale] = r.Value
	}
	return m
}

func near(a, b float64) bool { return math.Abs(a-b) < 0.01 }

func TestCelsiusBaseline(t *testing.T) {
	m := convert(t, "celsius", 0)
	if !near(m["kelvin"], 273.15) {
		t.Errorf("0°C → K = %v, want 273.15", m["kelvin"])
	}
	if !near(m["fahrenheit"], 32) {
		t.Errorf("0°C → °F = %v, want 32", m["fahrenheit"])
	}
	if m["celsius"] != 0 {
		t.Errorf("0°C → °C = %v, want 0", m["celsius"])
	}
	if m["rankine"] != 491.67 {
		t.Errorf("0°C → °R = %v, want 491.67", m["rankine"])
	}
}

func TestBoilingAndBody(t *testing.T) {
	// 100°C = 212°F = 373.15K
	m := convert(t, "celsius", 100)
	if !near(m["fahrenheit"], 212) {
		t.Errorf("100°C → °F = %v, want 212", m["fahrenheit"])
	}
	if !near(m["kelvin"], 373.15) {
		t.Errorf("100°C → K = %v, want 373.15", m["kelvin"])
	}

	// 人体体温 37°C ≈ 98.6°F
	m = convert(t, "celsius", 37)
	if !near(m["fahrenheit"], 98.6) {
		t.Errorf("37°C → °F = %v, want 98.6", m["fahrenheit"])
	}
}

func TestKelvinBaseline(t *testing.T) {
	m := convert(t, "kelvin", 0)
	// 绝对零度：-273.15°C、-459.67°F、0 R、559.73°De
	if !near(m["celsius"], -273.15) {
		t.Errorf("0K → °C = %v, want -273.15", m["celsius"])
	}
	if !near(m["fahrenheit"], -459.67) {
		t.Errorf("0K → °F = %v, want -459.67", m["fahrenheit"])
	}
	if m["kelvin"] != 0 {
		t.Errorf("0K → K = %v, want 0", m["kelvin"])
	}
	if !near(m["delisle"], 559.73) {
		t.Errorf("0K → °De = %v, want 559.73", m["delisle"])
	}
}

func TestFahrenheitRoundtrip(t *testing.T) {
	m := convert(t, "fahrenheit", 212)
	if !near(m["kelvin"], 373.15) {
		t.Errorf("212°F → K = %v, want 373.15", m["kelvin"])
	}
}

func TestLessCommonScales(t *testing.T) {
	// 参考 models.ts 的基准关系：
	// Rømer 7.5° = 273.15K（冰点）；Newton 33° = 373.15K（沸点）
	m := convert(t, "romer", 7.5)
	if !near(m["kelvin"], 273.15) {
		t.Errorf("7.5°Rø → K = %v, want 273.15", m["kelvin"])
	}
	m = convert(t, "newton", 33)
	if !near(m["kelvin"], 373.15) {
		t.Errorf("33°N → K = %v, want 373.15", m["kelvin"])
	}

	// 德利尔：100°C = 0°De（沸点为 0）
	m = convert(t, "delisle", 0)
	if !near(m["kelvin"], 373.15) {
		t.Errorf("0°De → K = %v, want 373.15", m["kelvin"])
	}
}

func TestResultMetadata(t *testing.T) {
	outStr, err := (Executor{}).Execute(context.Background(), `{"value":0,"from":"kelvin"}`)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(out.Results) != 8 {
		t.Fatalf("结果数 = %d, want 8", len(out.Results))
	}
	wantScales := []string{"kelvin", "celsius", "fahrenheit", "rankine", "delisle", "newton", "reaumur", "romer"}
	for i, r := range out.Results {
		if r.Scale != wantScales[i] {
			t.Errorf("results[%d].scale = %q, want %q", i, r.Scale, wantScales[i])
		}
		if r.Title == "" || r.Unit == "" {
			t.Errorf("results[%d] 缺少 title/unit: %+v", i, r)
		}
	}
}

func TestUnknownScale(t *testing.T) {
	if _, err := (Executor{}).Execute(context.Background(), `{"value":10,"from":"lux"}`); err == nil {
		t.Error("未知温标应返回错误")
	}
}

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != registry.CategoryMeasurement || meta.Icon != "Temperature" {
		t.Errorf("元数据不符: %+v", meta)
	}
}
