package yamltoml

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    string
		wantErr bool
	}{
		{
			name: "简单映射",
			yaml: "name: John\nage: 30\n",
			want: "age = 30\nname = 'John'",
		},
		{
			name: "嵌套表",
			yaml: "title: Example\nowner:\n  name: Tom\n",
			want: "title = 'Example'\n\n[owner]\nname = 'Tom'",
		},
		{
			name: "空输入",
			yaml: "",
			want: "",
		},
		{
			name: "空白输入",
			yaml: "   \n  ",
			want: "",
		},
		{
			name:    "非法 YAML",
			yaml:    "name: [unclosed\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.yaml)
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

	in, _ := json.Marshal(input{Yaml: "name: John\n"})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if !strings.Contains(out.Result, "name = 'John'") {
		t.Errorf("结果 = %q", out.Result)
	}
}
