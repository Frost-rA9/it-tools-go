package jsonxml

import (
	"context"
	"encoding/json"
	"strings"
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
			name: "属性",
			json: `{"a":{"_attributes":{"x":"1.234","y":"It's"}}}`,
			want: `<a x="1.234" y="It's"/>`,
		},
		{
			name: "属性与文本",
			json: `{"a":{"_attributes":{"x":"1"},"_text":"text"}}`,
			want: `<a x="1">text</a>`,
		},
		{
			name: "纯文本",
			json: `{"a":"hello"}`,
			want: `<a>hello</a>`,
		},
		{
			name: "嵌套",
			json: `{"root":{"b":"2"}}`,
			want: `<root><b>2</b></root>`,
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
			name:    "根节点非对象",
			json:    `[1,2]`,
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

	in, _ := json.Marshal(input{Text: `{"a":{"_attributes":{"x":"1"}}}`})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if !strings.Contains(out.Result, `x="1"`) {
		t.Errorf("结果 = %q", out.Result)
	}
}
