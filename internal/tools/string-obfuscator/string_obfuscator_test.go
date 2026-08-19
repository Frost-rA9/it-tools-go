package stringobfuscator

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

func TestExecuteObfuscate(t *testing.T) {
	exec := Executor{}
	tests := []struct {
		name    string
		input   input
		want    string
		wantErr bool
	}{
		{
			name:  "默认 keepFirst=4 keepLast=0（对齐参考默认）",
			input: input{Text: "hello world very long string", KeepFirst: 4},
			want:  "hell************************",
		},
		{
			name:  "keepLast 生效",
			input: input{Text: "hello world", KeepFirst: 4, KeepLast: 3},
			want:  "hell****rld",
		},
		{
			name:  "keepSpace 保留空格",
			input: input{Text: "hello world", KeepFirst: 4, KeepSpace: true},
			want:  "hell* *****",
		},
		{
			name:  "keepSpace 关闭时空格也被遮蔽",
			input: input{Text: "hello world", KeepFirst: 4, KeepSpace: false},
			want:  "hell*******",
		},
		{
			name:  "keepFirst+keepLast 覆盖全文时全部保留",
			input: input{Text: "hello world", KeepFirst: 11},
			want:  "hello world",
		},
		{
			name:  "keepFirst+keepLast 重叠时全部保留",
			input: input{Text: "hello", KeepFirst: 4, KeepLast: 4},
			want:  "hello",
		},
		{
			name:  "中文混淆（rune 计数，空格保留）",
			input: input{Text: "你好世界 мир", KeepFirst: 2, KeepSpace: true},
			want:  "你好** ***",
		},
		{
			name:  "自定义替换字符",
			input: input{Text: "hello world", KeepFirst: 4, ReplacementChar: "#"},
			want:  "hell#######",
		},
		{
			name:  "负数钳制为 0",
			input: input{Text: "abcd", KeepFirst: -1, KeepLast: -1},
			want:  "****",
		},
		{
			name:  "空输入",
			input: input{Text: ""},
			want:  "",
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
			if o.Result != tt.want {
				t.Errorf("期望 %q，实际 %q", tt.want, o.Result)
			}
		})
	}
}