package yamlfmt

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		indent string
		want   string
	}{
		{
			name:   "映射与序列",
			in:     "a: 1\nb:\n- x\n- y\nc: z\n",
			indent: "2",
			want:   "a: 1\nb:\n  - x\n  - y\nc: z\n",
		},
		{
			name:   "嵌套",
			in:     "server:\n  host: localhost\n  port: 8080\n",
			indent: "2",
			want:   "server:\n  host: localhost\n  port: 8080\n",
		},
		{
			name:   "注释保留",
			in:     "# 顶部注释\na: 1 # 行注释\n",
			indent: "2",
			want:   "# 顶部注释\na: 1 # 行注释\n",
		},
		{
			name:   "缩进4",
			in:     "list:\n- 1\n- 2\n",
			indent: "4",
			want:   "list:\n    - 1\n    - 2\n",
		},
		{
			name:   "标量",
			in:     "hello: world\n",
			indent: "2",
			want:   "hello: world\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Format(c.in, c.indent)
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("格式化不符:\ngot:\n%s\nwant:\n%s", got, c.want)
			}
		})
	}
}

func TestFormatErrors(t *testing.T) {
	bad := []string{"", "   ", "a: [1, 2"} // 空与非法
	for _, in := range bad {
		if _, err := Format(in, "2"); err == nil {
			t.Fatalf("期望报错: %q", in)
		}
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"yaml":"a: 1\nb:\n- x\n"}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(o.Formatted, "  - x") || o.LineCount < 3 {
		t.Fatalf("输出不符: %+v", o)
	}

	if _, err := exec.Execute(t.Context(), `{"yaml":"a: 1","indent":"8"}`); err == nil {
		t.Fatal("非法缩进应报错")
	}
}
