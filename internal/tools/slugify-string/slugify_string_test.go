package slugify

import (
	"encoding/json"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"My file path", "my-file-path"},
		{"Déjà Vu!", "deja-vu"},
		{"Larry the 🦆 Duck", "larry-the-duck"},
		{"Larry the 🦄 Duck", "larry-the-unicorn-duck"},
		{"BAR and baz", "bar-and-baz"},
		{"fooBar", "foo-bar"},
		{"HTTPServer", "http-server"},
		{"I love APIs", "i-love-apis"},
		{"don't", "dont"},
		{"it's a test", "its-a-test"},
		{"Über café", "uber-cafe"},
		{"你好世界", ""},
		{"straße", "strasse"},
		{"æther", "aether"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := slugify(tt.in); got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"text":"My File path"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if o.Slug != "my-file-path" {
		t.Errorf("slug = %q, want %q", o.Slug, "my-file-path")
	}
}

func TestExecuteInvalidInput(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `not-json`); err == nil {
		t.Error("非法输入应报错")
	}
}