package crontab

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	cases := []struct {
		name     string
		in       input
		wantExpr string
		wantErr  bool
	}{
		{name: "默认每分钟", in: input{}, wantExpr: "* * * * *", wantErr: false},
		{
			name:     "每15分钟",
			in:       input{Minute: fieldInput{Type: "step", Step: 15}},
			wantExpr: "*/15 * * * *", wantErr: false,
		},
		{
			name:     "小时范围",
			in:       input{Minute: fieldInput{Type: "list", Values: []int{0}}, Hour: fieldInput{Type: "range", From: 9, To: 17}},
			wantExpr: "0 9-17 * * *", wantErr: false,
		},
		{
			name: "列表",
			in: input{
				Minute: fieldInput{Type: "list", Values: []int{30, 0, 15}},
			},
			wantExpr: "0,15,30 * * * *", wantErr: false,
		},
		{
			name: "完整组合",
			in: input{
				Minute:     fieldInput{Type: "range", From: 0, To: 30},
				Hour:       fieldInput{Type: "step", Step: 2},
				DayOfMonth: fieldInput{Type: "list", Values: []int{1, 15}},
				Month:      fieldInput{Type: "every"},
				DayOfWeek:  fieldInput{Type: "every"},
			},
			wantExpr: "0-30 */2 1,15 * *", wantErr: false,
		},
		{name: "步长0报错", in: input{Minute: fieldInput{Type: "step", Step: 0}}, wantErr: true},
		{name: "小时越界报错", in: input{Hour: fieldInput{Type: "range", From: 0, To: 24}}, wantErr: true},
		{name: "范围倒置报错", in: input{Hour: fieldInput{Type: "range", From: 17, To: 9}}, wantErr: true},
		{name: "空列表报错", in: input{Minute: fieldInput{Type: "list"}}, wantErr: true},
		{name: "未知模式报错", in: input{Minute: fieldInput{Type: "foo"}}, wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr, summary, err := build(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际成功: %q", expr)
				}
				return
			}
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if expr != c.wantExpr {
				t.Fatalf("表达式不符: got %q, want %q", expr, c.wantExpr)
			}
			if summary == "" {
				t.Fatal("描述为空")
			}
		})
	}
}

func TestRenderField(t *testing.T) {
	m := fieldMeta{name: "分钟", min: 0, max: 59, unit: "分钟"}
	expr, desc, err := renderField(fieldInput{Type: "step", Step: 5}, m)
	if err != nil {
		t.Fatalf("意外错误: %v", err)
	}
	if expr != "*/5" || desc != "每 5 分钟" {
		t.Fatalf("step 渲染不符: expr=%q desc=%q", expr, desc)
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"minute":{"type":"every"},"hour":{"type":"every"},"day_of_month":{"type":"every"},"month":{"type":"every"},"day_of_week":{"type":"every"}}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if o.Expression != "* * * * *" {
		t.Fatalf("表达式不符: %q", o.Expression)
	}
	if reflect.DeepEqual(o.Summary, "") {
		t.Fatal("描述为空")
	}
}
