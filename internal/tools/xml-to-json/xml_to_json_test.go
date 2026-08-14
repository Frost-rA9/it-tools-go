package xmljson

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		want    string
		wantErr bool
	}{
		{
			name: "仅属性",
			xml:  `<a x="1.234" y="It's"/>`,
			want: "{\n  \"a\": {\n    \"_attributes\": {\n      \"x\": \"1.234\",\n      \"y\": \"It's\"\n    }\n  }\n}",
		},
		{
			name: "属性与文本",
			xml:  `<a x="1">text</a>`,
			want: "{\n  \"a\": {\n    \"_attributes\": {\n      \"x\": \"1\"\n    },\n    \"_text\": \"text\"\n  }\n}",
		},
		{
			name: "纯文本",
			xml:  `<a>hello</a>`,
			want: "{\n  \"a\": \"hello\"\n}",
		},
		{
			name: "重复子元素转数组",
			xml:  "<root><item>1</item><item>2</item></root>",
			want: "{\n  \"root\": {\n    \"item\": [\n      \"1\",\n      \"2\"\n    ]\n  }\n}",
		},
		{
			name: "空输入",
			xml:  "",
			want: "",
		},
		{
			name:    "非法 XML",
			xml:    "<a><b></a>",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.xml)
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

	in, _ := json.Marshal(input{Text: `<a x="1"/>`})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if out.Result != "{\n  \"a\": {\n    \"_attributes\": {\n      \"x\": \"1\"\n    }\n  }\n}" {
		t.Errorf("结果 = %q", out.Result)
	}
}
