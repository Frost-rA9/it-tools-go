package baseconv

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConvertBase(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fromBase int
		toBase   int
		want     string
		wantErr  bool
	}{
		{"十进制转二进制", "42", 10, 2, "101010", false},
		{"十进制转十六进制", "42", 10, 16, "2a", false},
		{"十六进制转十进制", "2a", 16, 10, "42", false},
		{"十六进制转十进制2", "ff", 16, 10, "255", false},
		{"十进制转Base64", "42", 10, 64, "G", false},
		{"Base64转十进制", "G", 64, 10, "42", false},
		{"零", "0", 10, 2, "0", false},
		{"空字符串", "", 10, 2, "0", false},
		{"大数", "18446744073709551616", 10, 2, "1" + repeat("0", 64), false},
		{"非法数字", "2", 2, 10, "", true},
		{"非法数字2", "g", 16, 10, "", true},
		{"输入进制越界", "1", 1, 10, "", true},
		{"输出进制越界", "1", 10, 65, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertBase(tt.value, tt.fromBase, tt.toBase)
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
				t.Errorf("结果 = %q，期望 %q", got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	in, err := json.Marshal(input{Value: "42", FromBase: 10, CustomBase: 42})
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}

	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}

	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if len(out.Results) != 6 {
		t.Fatalf("结果数量 = %d，期望 6", len(out.Results))
	}

	m := make(map[string]string)
	for _, r := range out.Results {
		m[r.Label] = r.Value
	}
	if m["Binary (2)"] != "101010" {
		t.Errorf("Binary (2) = %q", m["Binary (2)"])
	}
	if m["Hexadecimal (16)"] != "2a" {
		t.Errorf("Hexadecimal (16) = %q", m["Hexadecimal (16)"])
	}
	if m["Base64 (64)"] != "G" {
		t.Errorf("Base64 (64) = %q", m["Base64 (64)"])
	}
	if m["Custom (42)"] != "10" {
		t.Errorf("Custom (42) = %q", m["Custom (42)"])
	}

	t.Run("非法数字报错", func(t *testing.T) {
		in, _ := json.Marshal(input{Value: "2", FromBase: 2, CustomBase: 42})
		if _, err := e.Execute(context.Background(), string(in)); err == nil {
			t.Fatal("期望出错，但成功")
		}
	})
}

func repeat(s string, n int) string {
	b := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		b = append(b, s...)
	}
	return string(b)
}
