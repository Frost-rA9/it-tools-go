package tomljson

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		toml    string
		want    string
		wantErr bool
	}{
		{
			name: "简单",
			toml: "name = 'John'\nage = 30\n",
			want: "{\n   \"age\": 30,\n   \"name\": \"John\"\n}",
		},
		{
			name: "嵌套表",
			toml: "title = 'Example'\n[owner]\nname = 'Tom'\n",
			want: "{\n   \"owner\": {\n      \"name\": \"Tom\"\n   },\n   \"title\": \"Example\"\n}",
		},
		{
			name: "数组",
			toml: "nums = [1, 2, 3]\n",
			want: "{\n   \"nums\": [\n      1,\n      2,\n      3\n   ]\n}",
		},
		{
			name: "空输入",
			toml: "",
			want: "",
		},
		{
			name:    "非法 TOML",
			toml:    "name = [unclosed\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.toml)
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

	in, _ := json.Marshal(input{Text: "name = 'John'\n"})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if out.Result != "{\n   \"name\": \"John\"\n}" {
		t.Errorf("结果 = %q", out.Result)
	}
}
