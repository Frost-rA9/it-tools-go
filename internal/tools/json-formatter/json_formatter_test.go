package jsonfmt

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		name      string
		in        string
		wantSub   string // 结果必须包含的子串
		wantLines int
		wantErr   bool
	}{
		{
			name:    "默认缩进2",
			in:      `{"json":"{\"a\":1,\"b\":[1,2]}","indent":"2"}`,
			wantSub: "\"a\": 1",
		},
		{
			name:    "缩进4",
			in:      `{"json":"{\"a\":1}","indent":"4"}`,
			wantSub: "    \"a\": 1",
		},
		{
			name:    "Tab缩进",
			in:      `{"json":"{\"a\":1}","indent":"\t"}`,
			wantSub: "\t\"a\": 1",
		},
		{
			name:      "排序键",
			in:        `{"json":"{\"b\":2,\"a\":1}","sort_keys":true}`,
			wantSub:   "\"a\": 1",
			wantLines: 4,
		},
		{
			name:    "保序",
			in:      `{"json":"{\"b\":2,\"a\":1}","sort_keys":false}`,
			wantSub: "\"b\": 2",
		},
		{name: "空对象", in: `{"json":"{}"}`, wantSub: "{}"},
		{name: "非法缩进", in: `{"json":"{}","indent":"8"}`, wantErr: true},
		{name: "非法JSON", in: `{"json":"{bad}"}`, wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			exec := Executor{}
			out, err := exec.Execute(t.Context(), c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际成功: %s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if !strings.Contains(o.Formatted, c.wantSub) {
				t.Fatalf("结果缺少 %q: %s", c.wantSub, o.Formatted)
			}
			if c.wantLines > 0 && o.LineCount != c.wantLines {
				t.Fatalf("行数不符: got %d, want %d", o.LineCount, c.wantLines)
			}
		})
	}
}

func TestSortKeysOrder(t *testing.T) {
	exec := Executor{}
	// 排序模式：a 在 b 前；保序模式：b 在 a 前
	out, err := exec.Execute(t.Context(), `{"json":"{\"b\":2,\"a\":1}","sort_keys":true}`)
	if err != nil {
		t.Fatal(err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	bi := strings.Index(o.Formatted, "\"b\"")
	ai := strings.Index(o.Formatted, "\"a\"")
	if ai < 0 || bi < 0 || ai > bi {
		t.Fatalf("排序模式应 a 在前: %s", o.Formatted)
	}

	out2, err := exec.Execute(t.Context(), `{"json":"{\"b\":2,\"a\":1}","sort_keys":false}`)
	if err != nil {
		t.Fatal(err)
	}
	var o2 output
	_ = json.Unmarshal([]byte(out2), &o2)
	bi2 := strings.Index(o2.Formatted, "\"b\"")
	ai2 := strings.Index(o2.Formatted, "\"a\"")
	if ai2 < 0 || bi2 < 0 || bi2 > ai2 {
		t.Fatalf("保序模式应 b 在前: %s", o2.Formatted)
	}
}
