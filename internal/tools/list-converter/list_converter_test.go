package listconv

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name string
		text string
		opt  input
		want string
	}{
		{
			name: "默认分隔",
			text: "b\na\nc",
			opt:  input{Separator: ", "},
			want: "b, a, c",
		},
		{
			name: "升序",
			text: "b\na\nc",
			opt:  input{Separator: ", ", SortList: "asc"},
			want: "a, b, c",
		},
		{
			name: "降序",
			text: "b\na\nc",
			opt:  input{Separator: ", ", SortList: "desc"},
			want: "c, b, a",
		},
		{
			name: "去重",
			text: "a\nb\na\nc",
			opt:  input{Separator: ", ", RemoveDuplicates: true},
			want: "a, b, c",
		},
		{
			name: "反转",
			text: "a\nb\nc",
			opt:  input{Separator: ", ", ReverseList: true},
			want: "c, b, a",
		},
		{
			name: "转小写",
			text: "A\nBb\nC",
			opt:  input{Separator: ", ", LowerCase: true},
			want: "a, bb, c",
		},
		{
			name: "去空格",
			text: "  a  \n b \n c",
			opt:  input{Separator: ", ", TrimItems: true},
			want: "a, b, c",
		},
		{
			name: "去空格后移除空项",
			text: "a\n   \nb",
			opt:  input{Separator: ", ", TrimItems: true},
			want: "a, b",
		},
		{
			name: "条目前后缀",
			text: "a\nb",
			opt:  input{Separator: ", ", ItemPrefix: "[", ItemSuffix: "]"},
			want: "[a], [b]",
		},
		{
			name: "列表前后缀",
			text: "a\nb",
			opt:  input{Separator: ", ", ListPrefix: "start", ListSuffix: "end"},
			want: "starta, bend",
		},
		{
			name: "保留换行",
			text: "a\nb\nc",
			opt:  input{Separator: ",", KeepLineBreaks: true},
			want: "\na,\nb,\nc\n",
		},
		{
			name: "空输入",
			text: "",
			opt:  input{Separator: ", "},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert(tt.text, tt.opt); got != tt.want {
				t.Errorf("convert(%q) = %q，期望 %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	in, _ := json.Marshal(input{Text: "b\na\nc", Separator: ", ", SortList: "asc"})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if out.Result != "a, b, c" {
		t.Errorf("结果 = %q，期望 %q", out.Result, "a, b, c")
	}
}
