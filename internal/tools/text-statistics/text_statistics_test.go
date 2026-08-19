package textstatistics

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

func TestExecuteStatistics(t *testing.T) {
	exec := Executor{}
	tests := []struct {
		name    string
		text    string
		want    Stat
		wantStr string
	}{
		{
			name:    "空文本全为 0",
			text:    "",
			want:    Stat{0, 0, 0, 0},
			wantStr: "0 Bytes",
		},
		{
			name:    "纯 ASCII",
			text:    "hello world",
			want:    Stat{Characters: 11, Words: 2, Lines: 1, Bytes: 11},
			wantStr: "11 Bytes",
		},
		{
			name:    "中文（UTF-16 每字 1，UTF-8 每字 3 字节）",
			text:    "你好世界",
			want:    Stat{Characters: 4, Words: 1, Lines: 1, Bytes: 12},
			wantStr: "12 Bytes",
		},
		{
			name:    "emoji 星面字符记 2 个 code unit",
			text:    "a😀",
			want:    Stat{Characters: 3, Words: 1, Lines: 1, Bytes: 5},
			wantStr: "5 Bytes",
		},
		{
			name:    "LF 换行",
			text:    "a\nb",
			want:    Stat{Characters: 3, Words: 2, Lines: 2, Bytes: 3},
			wantStr: "3 Bytes",
		},
		{
			name:    "CRLF 与 CR 换行",
			text:    "a\r\nb\rc",
			want:    Stat{Characters: 6, Words: 3, Lines: 3, Bytes: 6},
			wantStr: "6 Bytes",
		},
		{
			name:    "末尾换行也算一行",
			text:    "a\n",
			want:    Stat{Characters: 2, Words: 1, Lines: 2, Bytes: 2},
			wantStr: "2 Bytes",
		},
		{
			name:    "首尾空白不产生空词",
			text:    "  a   b  ",
			want:    Stat{Characters: 9, Words: 2, Lines: 1, Bytes: 9},
			wantStr: "9 Bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Execute(t.Context(), marshalInput(t, input{Text: tt.text}))
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if o.Characters != tt.want.Characters || o.Words != tt.want.Words ||
				o.Lines != tt.want.Lines || o.Bytes != tt.want.Bytes {
				t.Errorf("期望 %+v，实际 characters=%d words=%d lines=%d bytes=%d",
					tt.want, o.Characters, o.Words, o.Lines, o.Bytes)
			}
			if o.SizeText != tt.wantStr {
				t.Errorf("SizeText 期望 %q，实际 %q", tt.wantStr, o.SizeText)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes int
		want  string
	}{
		{0, "0 Bytes"},
		{1, "1 Bytes"},
		{1023, "1023 Bytes"},
		{1024, "1 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1 MB"},
		{2*1024*1024 + 512*1024, "2.5 MB"},
		{2000000000, "1.86 GB"},
	}
	for _, tt := range tests {
		if got := formatBytes(tt.bytes); got != tt.want {
			t.Errorf("formatBytes(%d) 期望 %q，实际 %q", tt.bytes, tt.want, got)
		}
	}
}