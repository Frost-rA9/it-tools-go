package htmlentities

import (
	"encoding/json"
	"testing"
)

func TestEscape(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"空字符串", "", ""},
		{"无特殊字符", "plain text", "plain text"},
		{"五个特殊字符", `&<>"'`, "&amp;&lt;&gt;&quot;&#39;"},
		{"标题样例", "<title>IT Tool</title>", "&lt;title&gt;IT Tool&lt;/title&gt;"},
		{"中文", "你好", "你好"},
		{"混合", `a < b > c & d "e" 'f'`, `a &lt; b &gt; c &amp; d &quot;e&quot; &#39;f&#39;`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escape.Replace(tt.in); got != tt.want {
				t.Errorf("escape(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"空字符串", "", ""},
		{"无实体", "plain text", "plain text"},
		{"五个实体", "&amp;&lt;&gt;&quot;&#39;", `&<>"'`},
		{"标题样例", "&lt;title&gt;IT Tool&lt;/title&gt;", "<title>IT Tool</title>"},
		{"单遍不递归", "&amp;lt;", "&lt;"},
		{"未知实体保持", "&nbsp;&copy;", "&nbsp;&copy;"},
		{"数字实体保持", "&#65;", "&#65;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unescape.Replace(tt.in); got != tt.want {
				t.Errorf("unescape(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"转义", `{"text":"<b>hi</b>","action":"escape"}`, `{"result":"&lt;b&gt;hi&lt;/b&gt;"}`, false},
		{"反转义", `{"text":"&lt;b&gt;hi&lt;/b&gt;","action":"unescape"}`, `{"result":"<b>hi</b>"}`, false},
		{"未知操作", `{"text":"x","action":"foo"}`, "", true},
		{"非法输入 JSON", `not-json`, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Executor
			out, err := e.Execute(t.Context(), tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute 期望错误，得到 %s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute 意外错误: %v", err)
			}
			var got, want map[string]string
			if err := json.Unmarshal([]byte(out), &got); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
				t.Fatalf("解析期望输出失败: %v", err)
			}
			if got["result"] != want["result"] {
				t.Errorf("Execute 结果 = %s, want %s", out, tt.want)
			}
		})
	}
}