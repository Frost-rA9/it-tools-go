package basicauth

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{"user:pass", "user", "pass"},
		{"admin:admin", "admin", "admin"},
		{"空输入", "", ""},
		{"含特殊字符", "a:b", "c/d"},
		{"含中文", "用户", "密码"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := "Authorization: Basic " +
				base64.StdEncoding.EncodeToString([]byte(tt.username+":"+tt.password))
			var e Executor
			out, err := e.Execute(t.Context(), `{"username":"`+tt.username+`","password":"`+tt.password+`"}`)
			if err != nil {
				t.Fatalf("Execute 意外错误: %v", err)
			}
			var got output
			if err := json.Unmarshal([]byte(out), &got); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if got.Header != want {
				t.Errorf("Header = %q, want %q", got.Header, want)
			}
		})
	}
}

func TestExecuteInvalidInput(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `not-json`); err == nil {
		t.Error("非法输入应报错")
	}
}