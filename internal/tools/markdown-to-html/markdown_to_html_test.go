package markdownhtml

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	e := Executor{}

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "标题",
			text: "# Hello",
			want: "<h1>Hello</h1>\n",
		},
		{
			name: "段落",
			text: "Hello **world**",
			want: "<p>Hello <strong>world</strong></p>\n",
		},
		{
			name: "列表",
			text: "- a\n- b",
			want: "<ul>\n<li>a</li>\n<li>b</li>\n</ul>\n",
		},
		{
			name: "代码块",
			text: "```go\nfmt.Println(\"hi\")\n```",
			want: "<pre><code class=\"language-go\">fmt.Println(&quot;hi&quot;)\n</code></pre>\n",
		},
		{
			name: "空输入",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, _ := json.Marshal(input{Text: tt.text})
			outJSON, err := e.Execute(context.Background(), string(in))
			if err != nil {
				t.Fatalf("Execute 返回错误: %v", err)
			}
			var out output
			if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
				t.Fatalf("反序列化输出失败: %v", err)
			}
			if !strings.Contains(out.Result, tt.want) {
				t.Errorf("结果 = %q，期望包含 %q", out.Result, tt.want)
			}
		})
	}
}
