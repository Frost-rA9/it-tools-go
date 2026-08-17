package mimetypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestExecuteList(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"list":true}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(o.All) != len(mimeToExtensions) {
		t.Errorf("all 数量 = %d, want %d", len(o.All), len(mimeToExtensions))
	}
	if len(o.All) < 500 {
		t.Errorf("all 数量过少: %d", len(o.All))
	}
	// 数据已排序。
	for i := 1; i < len(o.All); i++ {
		if o.All[i-1].MimeType > o.All[i].MimeType {
			t.Fatal("MIME 类型未按字典序排序")
		}
	}
}

func TestExecuteByMime(t *testing.T) {
	var e Executor
	tests := []struct {
		mime string
		want []string
	}{
		{"application/pdf", []string{"pdf"}},
		{"image/jpeg", []string{"jpg", "jpeg", "jpe"}},
		{"text/plain", []string{"txt", "text", "conf", "def", "list", "log", "in", "ini"}},
	}
	for _, tt := range tests {
		out, err := e.Execute(t.Context(), `{"mime_type":"`+tt.mime+`"}`)
		if err != nil {
			t.Fatalf("Execute 意外错误: %v", err)
		}
		var o output
		if err := json.Unmarshal([]byte(out), &o); err != nil {
			t.Fatalf("解析输出失败: %v", err)
		}
		if !reflect.DeepEqual(o.Extensions, tt.want) {
			t.Errorf("by_mime(%s) = %v, want %v", tt.mime, o.Extensions, tt.want)
		}
	}
}

func TestExecuteByExtension(t *testing.T) {
	var e Executor
	tests := []struct {
		ext  string
		want string
	}{
		{"pdf", "application/pdf"},
		{"jpg", "image/jpeg"},
		{"jpeg", "image/jpeg"},
		{".PNG", "image/png"}, // 前导点与大小写
	}
	for _, tt := range tests {
		out, err := e.Execute(t.Context(), `{"extension":"`+tt.ext+`"}`)
		if err != nil {
			t.Fatalf("Execute 意外错误: %v", err)
		}
		var o output
		if err := json.Unmarshal([]byte(out), &o); err != nil {
			t.Fatalf("解析输出失败: %v", err)
		}
		if len(o.MimeTypes) != 1 || o.MimeTypes[0] != tt.want {
			t.Errorf("by_extension(%q) = %v, want [%s]", tt.ext, o.MimeTypes, tt.want)
		}
	}
}

func TestExecuteUnknown(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"mime_type":"application/unknown-xyz"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	if out != `{}` {
		t.Errorf("未知 mime 输出 = %s", out)
	}

	out, err = e.Execute(t.Context(), `{"extension":"zzzz"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	if out != `{}` {
		t.Errorf("未知扩展名输出 = %s", out)
	}
}

func TestExecuteErrors(t *testing.T) {
	var e Executor
	tests := []struct {
		name  string
		input string
	}{
		{"无查询参数", `{}`},
		{"非法输入", `not-json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := e.Execute(t.Context(), tt.input); err == nil {
				t.Errorf("Execute(%s) 期望错误", tt.input)
			}
		})
	}
}