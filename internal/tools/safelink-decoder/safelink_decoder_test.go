package safelink

import (
	"encoding/json"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "合法 safelink",
			input: `{"url":"https://nam02.safelinks.protection.outlook.com/?url=https%3A%2F%2Fexample.com%2Fpath%3Fq%3D1&data=abc"}`,
			want:  "https://example.com/path?q=1",
		},
		{
			name:  "缺少 url 参数",
			input: `{"url":"https://nam02.safelinks.protection.outlook.com/?data=abc"}`,
			want:  "",
		},
		{
			name:    "非法域名",
			input:   `{"url":"https://example.com/?url=https%3A%2F%2Ffoo.com"}`,
			wantErr: true,
		},
		{
			name:    "非法 URL",
			input:   `{"url":"://bad"}`,
			wantErr: true,
		},
		{
			name:    "非法输入 JSON",
			input:   `not-json`,
			wantErr: true,
		},
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
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if o.DecodedURL != tt.want {
				t.Errorf("decoded_url = %q, want %q", o.DecodedURL, tt.want)
			}
		})
	}
}