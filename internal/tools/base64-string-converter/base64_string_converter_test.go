package base64string

import (
	"context"
	"encoding/json"
	"testing"
)

func TestExecute(t *testing.T) {
	e := Executor{}

	tests := []struct {
		name    string
		input   input
		want    string
		wantErr bool
	}{
		{
			name:  "编码",
			input: input{Text: "hello world", Action: "encode"},
			want:  "aGVsbG8gd29ybGQ=",
		},
		{
			name:  "解码",
			input: input{Text: "aGVsbG8gd29ybGQ=", Action: "decode"},
			want:  "hello world",
		},
		{
			name:  "中文编码",
			input: input{Text: "你好", Action: "encode"},
			want:  "5L2g5aW9",
		},
		{
			name:  "空字符串编码",
			input: input{Text: "", Action: "encode"},
			want:  "",
		},
		{
			name:  "URL-safe 编码",
			input: input{Text: "hello?world", Action: "encode", URLSafe: true},
			want:  "aGVsbG8_d29ybGQ=",
		},
		{
			name:  "URL-safe 解码",
			input: input{Text: "aGVsbG8_d29ybGQ=", Action: "decode", URLSafe: true},
			want:  "hello?world",
		},
		{
			name:    "非法 Base64 解码",
			input:   input{Text: "!!!not-base64!!!", Action: "decode"},
			wantErr: true,
		},
		{
			name:    "未知操作",
			input:   input{Text: "abc", Action: "foo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("序列化输入失败: %v", err)
			}

			outJSON, err := e.Execute(context.Background(), string(in))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望出错，但成功了（输出 %q）", outJSON)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute 返回错误: %v", err)
			}

			var out output
			if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
				t.Fatalf("反序列化输出失败: %v", err)
			}
			if out.Result != tt.want {
				t.Errorf("结果 = %q，期望 %q", out.Result, tt.want)
			}
		})
	}
}
