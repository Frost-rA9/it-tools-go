package urlcodec

import (
	"encoding/json"
	"testing"
)

func TestEncodeURIComponent(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"空字符串", "", ""},
		{"空格", "hello world", "hello%20world"},
		{"保留字符", "a/b?c=d&e#f", "a%2Fb%3Fc%3Dd%26e%23f"},
		{"未保留字符保持不变", "A-Z_a-z.0-9-~!*'()", "A-Z_a-z.0-9-~!*'()"},
		{"加号", "a+b", "a%2Bb"},
		{"中文", "你好", "%E4%BD%A0%E5%A5%BD"},
		{"emoji", "🎉", "%F0%9F%8E%89"},
		{"混合", "https://example.com/路径?q=你好 world", "https%3A%2F%2Fexample.com%2F%E8%B7%AF%E5%BE%84%3Fq%3D%E4%BD%A0%E5%A5%BD%20world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeURIComponent(tt.in); got != tt.want {
				t.Errorf("encodeURIComponent(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDecodeURIComponent(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"空字符串", "", "", false},
		{"空格编码", "hello%20world", "hello world", false},
		{"加号保持不变", "a+b", "a+b", false},
		{"加号编码", "a%2Bb", "a+b", false},
		{"中文", "%E4%BD%A0%E5%A5%BD", "你好", false},
		{"emoji", "%F0%9F%8E%89", "🎉", false},
		{"已编码 url", "https%3A%2F%2Fexample.com%2F%3Fq%3Dhi", "https://example.com/?q=hi", false},
		{"非法百分号", "%ZZ", "", true},
		{"截断百分号", "%E4", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeURIComponent(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("decodeURIComponent(%q) 期望错误，得到 %q", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeURIComponent(%q) 意外错误: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("decodeURIComponent(%q) = %q, want %q", tt.in, got, tt.want)
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
		{"编码", `{"text":"hello world","action":"encode"}`, `{"result":"hello%20world"}`, false},
		{"解码", `{"text":"hello%20world","action":"decode"}`, `{"result":"hello world"}`, false},
		{"未知操作", `{"text":"x","action":"foo"}`, "", true},
		{"非法输入 JSON", `not-json`, "", true},
		{"解码非法编码", `{"text":"%ZZ","action":"decode"}`, "", true},
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
			// 规范化比较（字段顺序无关）。
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