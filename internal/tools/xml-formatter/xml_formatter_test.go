package xmlfmt

import (
	"encoding/json"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "嵌套元素",
			in:   `<root><a x="1"><b>text</b></a><c/></root>`,
			want: "<root>\n  <a x=\"1\">\n    <b>text</b>\n  </a>\n  <c/>\n</root>\n",
		},
		{
			name: "已格式化",
			in:   "<root>\n  <a>1</a>\n</root>",
			want: "<root>\n  <a>1</a>\n</root>\n",
		},
		{
			name: "注释",
			in:   `<root><!-- 注释 --><a>1</a></root>`,
			want: "<root>\n  <!-- 注释 -->\n  <a>1</a>\n</root>\n",
		},
		{
			name: "空元素",
			in:   `<root><a/></root>`,
			want: "<root>\n  <a/>\n</root>\n",
		},
		{name: "非法XML", in: `<root><a></root>`, want: ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want == "" {
				if _, err := Format(c.in, "  "); err == nil {
					t.Fatalf("期望报错: %q", c.in)
				}
				return
			}
			got, err := Format(c.in, "  ")
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("格式化不符:\ngot:\n%s\nwant:\n%s", got, c.want)
			}
		})
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"xml":"<root><a x=\"1\"><b>t</b></a></root>"}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if o.Formatted == "" || o.LineCount != 5 {
		t.Fatalf("输出不符: %+v", o)
	}

	// Tab 缩进
	out2, err := exec.Execute(t.Context(), `{"xml":"<r><a/></r>","indent":"\t"}`)
	if err != nil {
		t.Fatal(err)
	}
	var o2 output
	_ = json.Unmarshal([]byte(out2), &o2)
	if o2.Formatted != "<r>\n\t<a/>\n</r>\n" {
		t.Fatalf("Tab 缩进不符: %q", o2.Formatted)
	}

	if _, err := exec.Execute(t.Context(), `{"xml":"<bad>","indent":"8"}`); err == nil {
		t.Fatal("非法缩进应报错")
	}
}
