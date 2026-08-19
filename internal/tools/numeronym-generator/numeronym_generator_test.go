package numeronymgen

import (
	"encoding/json"
	"testing"
)

// marshalInput 将 input 结构编码为 JSON 字符串。
func marshalInput(t *testing.T, in input) string {
	t.Helper()
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}
	return string(raw)
}

func TestExecuteNumeronym(t *testing.T) {
	exec := Executor{}
	tests := []struct {
		name string
		word string
		want string
	}{
		{"经典 i18n", "internationalization", "i18n"},
		{"a11y", "accessibility", "a11y"},
		{"l10n", "localization", "l10n"},
		{"长度 4", "abcd", "a2d"},
		{"长度 3 原样", "abc", "abc"},
		{"长度 2 原样", "ab", "ab"},
		{"单字符原样", "a", "a"},
		{"中文词（rune 计数）", "国际化的", "国2的"},
		{"多词串", "hello world", "h9d"},
		{"带符号长串", "test!", "t3!"},
		{"空输入", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Execute(t.Context(), marshalInput(t, input{Word: tt.word}))
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if o.Numeronym != tt.want {
				t.Errorf("期望 %q，实际 %q", tt.want, o.Numeronym)
			}
		})
	}
}