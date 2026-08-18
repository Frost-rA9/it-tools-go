package jsontoml

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    string
		wantErr bool
	}{
		{
			name: "简单映射",
			json: `{"name":"John","age":30}`,
			want: "age = 30\nname = 'John'",
		},
		{
			name: "嵌套表",
			json: `{"server":{"host":"localhost","port":8080}}`,
			want: "[server]\nhost = 'localhost'\nport = 8080",
		},
		{
			name: "数组",
			json: `{"tags":["a","b"]}`,
			want: "tags = ['a', 'b']",
		},
		{
			name: "浮点数",
			json: `{"pi":3.14}`,
			want: "pi = 3.14",
		},
		{
			name: "null 值转为空字符串",
			json: `{"a":null}`,
			want: "a = ''",
		},
		{
			name: "空输入",
			json: "",
			want: "",
		},
		{
			name:    "非法 JSON",
			json:    `{"a":`,
			wantErr: true,
		},
		{
			name:    "根节点数组",
			json:    `[1,2]`,
			wantErr: true,
		},
		{
			name:    "根节点标量",
			json:    `"hello"`,
			wantErr: true,
		},
		{
			name:    "多余内容",
			json:    `{} {}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.json)
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

	in, _ := json.Marshal(input{Text: `{"name":"John"}`})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if out.Result != "name = 'John'" {
		t.Errorf("结果 = %q", out.Result)
	}
}
