package etacalculator

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"it-tools-go/internal/registry"
)

func TestDurationText(t *testing.T) {
	cases := []struct {
		name string
		ms   float64
		want string
	}{
		{"完整分段", 5*3600*1000 + 10*60*1000 + 30*1000 + 500, "5 小时 10 分 30 秒 500 毫秒"},
		{"省略零值", 65 * 1000, "1 分 5 秒"},
		{"纯小时", 2 * 3600 * 1000, "2 小时"},
		{"纯毫秒", 999, "999 毫秒"},
		{"超过一天", 26*3600*1000 + 5*60*1000, "26 小时 5 分"},
		{"全零", 0, "0 毫秒"},
		{"小数截断", 1500.7, "1 秒 500 毫秒"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := DurationText(c.ms); got != c.want {
				t.Errorf("DurationText(%v) = %q, want %q", c.ms, got, c.want)
			}
		})
	}
}

func TestRelativeText(t *testing.T) {
	cases := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"接近现在", 3 * time.Second, ""},
		{"几秒", 4 * time.Second, ""},
		{"一分钟", 60 * time.Second, "约 1 分钟后"},
		{"几分钟", 5 * time.Minute, "约 5 分钟后"},
		{"几小时", 3 * time.Hour, "约 3 小时后"},
		{"一天多", 30 * time.Hour, "约 1 天后"},
		{"已结束分钟", -5 * time.Minute, "已结束 约 5 分钟后"},
		{"已结束刚结束", -3 * time.Second, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := RelativeText(c.d); got != c.want {
				t.Errorf("RelativeText(%v) = %q, want %q", c.d, got, c.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	// 参考示例：3 分钟洗 5 个盘子，500 个盘子需 5 小时。
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.Local).UnixMilli()
	raw := `{"unitCount":500,"unitPerTimeSpan":5,"timeSpan":3,"timeSpanUnitMultiplier":60000,"startedAtMs":` + strconv.FormatInt(start, 10) + `}`
	outStr, err := (Executor{}).Execute(context.Background(), raw)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	// 500 / (5 / 180000) = 18,000,000 ms = 5 小时。
	if out.DurationMs != 18_000_000 {
		t.Errorf("DurationMs = %v, want 18000000", out.DurationMs)
	}
	if out.DurationText != "5 小时" {
		t.Errorf("DurationText = %q, want %q", out.DurationText, "5 小时")
	}
	if out.EndAtMs != start+18_000_000 {
		t.Errorf("EndAtMs = %v, want %v", out.EndAtMs, start+18_000_000)
	}
}

func TestExecuteInvalid(t *testing.T) {
	cases := []string{
		`{"unitCount":500,"unitPerTimeSpan":0,"timeSpan":3,"timeSpanUnitMultiplier":60000,"startedAtMs":0}`,
		`{"unitCount":500,"unitPerTimeSpan":5,"timeSpan":0,"timeSpanUnitMultiplier":60000,"startedAtMs":0}`,
		`{"unitCount":0,"unitPerTimeSpan":5,"timeSpan":3,"timeSpanUnitMultiplier":60000,"startedAtMs":0}`,
		`{"unitCount":500,"unitPerTimeSpan":5,"timeSpan":3,"timeSpanUnitMultiplier":0,"startedAtMs":0}`,
	}
	for i, raw := range cases {
		outStr, err := (Executor{}).Execute(context.Background(), raw)
		if err != nil {
			t.Fatalf("case %d: Execute 返回错误: %v", i, err)
		}
		var out output
		if err := json.Unmarshal([]byte(outStr), &out); err != nil {
			t.Fatalf("case %d: 解析输出失败: %v", i, err)
		}
		if out.DurationText != "" || out.EndAtMs != 0 {
			t.Errorf("case %d: 无效输入应返回空结果, got %+v", i, out)
		}
	}
}

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != registry.CategoryMath || meta.Icon != "Hourglass" {
		t.Errorf("元数据不符: %+v", meta)
	}
}
