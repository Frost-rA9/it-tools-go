package tokengen

import (
	"encoding/json"
	"regexp"
	"strings"
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

func TestExecuteToken(t *testing.T) {
	exec := Executor{}

	tests := []struct {
		name    string
		input   input
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name:  "默认全开关，长度 64",
			input: input{Length: 64, WithUppercase: true, WithLowercase: true, WithNumbers: true},
			check: func(t *testing.T, r string) {
				if len(r) != 64 {
					t.Errorf("期望长度 64，实际 %d", len(r))
				}
				if ok := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(r); !ok {
					t.Errorf("结果包含非法字符: %q", r)
				}
			},
		},
		{
			name:  "仅大写",
			input: input{Length: 32, WithUppercase: true},
			check: func(t *testing.T, r string) {
				if ok := regexp.MustCompile(`^[A-Z]+$`).MatchString(r); !ok {
					t.Errorf("仅大写模式出现非大写字符: %q", r)
				}
			},
		},
		{
			name:  "仅数字",
			input: input{Length: 16, WithNumbers: true},
			check: func(t *testing.T, r string) {
				if ok := regexp.MustCompile(`^[0-9]+$`).MatchString(r); !ok {
					t.Errorf("仅数字模式出现非数字字符: %q", r)
				}
			},
		},
		{
			name:  "全开关含符号",
			input: input{Length: 128, WithUppercase: true, WithLowercase: true, WithNumbers: true, WithSymbols: true},
			check: func(t *testing.T, r string) {
				if len(r) != 128 {
					t.Errorf("期望长度 128，实际 %d", len(r))
				}
				if !strings.ContainsAny(r, symbolsCharset) {
					t.Errorf("开启符号却未出现任何符号: %q", r)
				}
			},
		},
		{
			name:  "自定义字符集优先",
			input: input{Length: 20, WithUppercase: true, WithLowercase: false, WithNumbers: false, Alphabet: "xy"},
			check: func(t *testing.T, r string) {
				if ok := regexp.MustCompile(`^[xy]+$`).MatchString(r); !ok {
					t.Errorf("自定义字符集出现非法字符: %q", r)
				}
			},
		},
		{
			name:  "最大长度 512",
			input: input{Length: 512, WithUppercase: true, WithLowercase: true, WithNumbers: true},
			check: func(t *testing.T, r string) {
				if len(r) != 512 {
					t.Errorf("期望长度 512，实际 %d", len(r))
				}
			},
		},
		{
			name:    "长度 0 报错",
			input:   input{Length: 0, WithUppercase: true},
			wantErr: true,
		},
		{
			name:    "长度超上限报错",
			input:   input{Length: 513, WithUppercase: true},
			wantErr: true,
		},
		{
			name:    "无任何字符集报错",
			input:   input{Length: 16},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Execute(t.Context(), marshalInput(t, tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际成功: %s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			tt.check(t, o.Result)
		})
	}
}

func TestCryptoTokenRandomness(t *testing.T) {
	// 连续生成应各不相同（随机性冒烟测试）。
	seen := make(map[string]bool)
	for i := 0; i < 20; i++ {
		tok, err := cryptoToken("abcdefghijklmnopqrstuvwxyz", 16)
		if err != nil {
			t.Fatalf("生成失败: %v", err)
		}
		if seen[tok] {
			t.Fatalf("出现重复 Token: %q", tok)
		}
		seen[tok] = true
	}
}
