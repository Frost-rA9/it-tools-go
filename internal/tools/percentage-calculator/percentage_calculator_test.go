package percentagecalculator

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"it-tools-go/internal/registry"
)

// run 辅助执行器，返回解析后的 output。
func run(t *testing.T, in input) output {
	t.Helper()
	raw, err := json.Marshal(in)
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
	return out
}

func ptr(f float64) *float64 { return &f }

func TestPercentOf(t *testing.T) {
	cases := []struct {
		name string
		x, y float64
		want string
	}{
		{"整数", 50, 200, "100"},
		{"小数", 12.5, 80, "10"},
		{"负数", -50, 200, "-100"},
		{"X 为零", 0, 100, "0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := run(t, input{Mode: ModePercentOf, X: ptr(c.x), Y: ptr(c.y)})
			if out.Result != c.want {
				t.Errorf("percent_of(%v, %v) = %q, want %q", c.x, c.y, out.Result, c.want)
			}
		})
	}
}

func TestWhatPercent(t *testing.T) {
	cases := []struct {
		name string
		x, y float64
		want string
	}{
		{"整数", 25, 200, "12.5"},
		{"整百", 100, 100, "100"},
		{"小数", 7.5, 30, "25"},
		{"除零返回空", 10, 0, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := run(t, input{Mode: ModeWhatPercent, X: ptr(c.x), Y: ptr(c.y)})
			if out.Result != c.want {
				t.Errorf("what_percent(%v, %v) = %q, want %q", c.x, c.y, out.Result, c.want)
			}
		})
	}
}

func TestChange(t *testing.T) {
	cases := []struct {
		name string
		x, y float64
		want string
	}{
		{"增长一倍", 50, 100, "100"},
		{"增长两成", 10, 12, "20"},
		{"下降一半", 100, 50, "-50"},
		{"不变", 30, 30, "0"},
		{"除零返回空", 0, 10, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := run(t, input{Mode: ModeChange, X: ptr(c.x), Y: ptr(c.y)})
			if out.Result != c.want {
				t.Errorf("change(%v, %v) = %q, want %q", c.x, c.y, out.Result, c.want)
			}
		})
	}
}

func TestMissingFields(t *testing.T) {
	// 未填写的输入（nil）应返回空结果而不是错误。
	cases := []input{
		{Mode: ModePercentOf, X: nil, Y: ptr(100)},
		{Mode: ModePercentOf, X: ptr(10), Y: nil},
		{Mode: ModePercentOf, X: nil, Y: nil},
	}
	for i, in := range cases {
		out := run(t, in)
		if out.Result != "" {
			t.Errorf("case %d: result = %q, want empty", i, out.Result)
		}
	}
}

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != registry.CategoryMath || meta.Icon != "Percentage" {
		t.Errorf("元数据不符: %+v", meta)
	}
	if !reflect.DeepEqual(meta.Keywords, Keywords) {
		t.Errorf("Keywords 不符: %v", meta.Keywords)
	}
}
