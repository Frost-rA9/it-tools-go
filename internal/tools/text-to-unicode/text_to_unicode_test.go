package textunicode

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTextToUnicode(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"Hi", "Hi", "&#72;&#105;"},
		{"中文", "中", "&#20013;"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextToUnicode(tt.text); got != tt.want {
				t.Errorf("TextToUnicode(%q) = %q，期望 %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestUnicodeToText(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"Hi", "&#72;&#105;", "Hi"},
		{"中文", "&#20013;", "中"},
		{"含普通文本", "hello &#72;&#105; world", "hello Hi world"},
		{"空字符串", "", ""},
		{"无转义", "plain text", "plain text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnicodeToText(tt.s); got != tt.want {
				t.Errorf("UnicodeToText(%q) = %q，期望 %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	t.Run("文本转Unicode", func(t *testing.T) {
		in, _ := json.Marshal(input{Text: "Hi", Mode: ModeTextToUnicode})
		outJSON, err := e.Execute(context.Background(), string(in))
		if err != nil {
			t.Fatalf("Execute 返回错误: %v", err)
		}
		var out output
		json.Unmarshal([]byte(outJSON), &out)
		if out.Result != "&#72;&#105;" {
			t.Errorf("结果 = %q", out.Result)
		}
	})

	t.Run("Unicode转文本", func(t *testing.T) {
		in, _ := json.Marshal(input{Text: "&#72;&#105;", Mode: ModeUnicodeToText})
		outJSON, err := e.Execute(context.Background(), string(in))
		if err != nil {
			t.Fatalf("Execute 返回错误: %v", err)
		}
		var out output
		json.Unmarshal([]byte(outJSON), &out)
		if out.Result != "Hi" {
			t.Errorf("结果 = %q", out.Result)
		}
	})
}
