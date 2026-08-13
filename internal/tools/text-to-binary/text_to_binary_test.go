package textbinary

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTextToBinary(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"Hi", "Hi", "01001000 01101001"},
		{"空格", "a b", "01100001 00100000 01100010"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextToBinary(tt.text); got != tt.want {
				t.Errorf("TextToBinary(%q) = %q，期望 %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestBinaryToText(t *testing.T) {
	tests := []struct {
		name    string
		binary  string
		want    string
		wantErr bool
	}{
		{"带空格", "01001000 01101001", "Hi", false},
		{"无空格", "0100100001101001", "Hi", false},
		{"单字符", "01001000", "H", false},
		{"空字符串", "", "", false},
		{"非8的倍数", "0100100", "", true},
		{"含非法字符", "0100100x", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BinaryToText(tt.binary)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望出错，但成功（%q）", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("返回错误: %v", err)
			}
			if got != tt.want {
				t.Errorf("BinaryToText(%q) = %q，期望 %q", tt.binary, got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	t.Run("文本转二进制", func(t *testing.T) {
		in, _ := json.Marshal(input{Text: "Hi", Mode: ModeTextToBinary})
		outJSON, err := e.Execute(context.Background(), string(in))
		if err != nil {
			t.Fatalf("Execute 返回错误: %v", err)
		}
		var out output
		if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
			t.Fatalf("反序列化输出失败: %v", err)
		}
		if out.Result != "01001000 01101001" {
			t.Errorf("结果 = %q", out.Result)
		}
	})

	t.Run("二进制转文本", func(t *testing.T) {
		in, _ := json.Marshal(input{Text: "01001000 01101001", Mode: ModeBinaryToText})
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

	t.Run("非法二进制报错", func(t *testing.T) {
		in, _ := json.Marshal(input{Text: "0100100", Mode: ModeBinaryToText})
		if _, err := e.Execute(context.Background(), string(in)); err == nil {
			t.Fatal("期望出错，但成功")
		}
	})
}
