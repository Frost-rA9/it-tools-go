package qrcodegen

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"text":"https://example.com","foreground":"#000000ff","background":"#ffffffff","level":"medium"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if !strings.HasPrefix(got.DataURL, "data:image/png;base64,") {
		t.Fatalf("Data URL 前缀不符: %q", got.DataURL[:min(len(got.DataURL), 30)])
	}
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(got.DataURL, "data:image/png;base64,"))
	if err != nil || len(b) < 8 || string(b[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatalf("输出不是合法 PNG: %v", err)
	}
}

func TestEmptyText(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"text":""}`)
	if err != nil || out != `{"data_url":""}` {
		t.Errorf("空文本结果 = %s, err=%v", out, err)
	}
}

func TestInvalidLevel(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `{"text":"x","level":"invalid"}`); err == nil {
		t.Error("非法纠错级别应报错")
	}
}
