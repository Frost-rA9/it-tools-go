package regextester

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		text    string
		flags   Flags
		want    []MatchResult
		wantErr bool
	}{
		{
			name:  "基本全局匹配",
			regex: `\d+`,
			text:  "a1b22c333",
			flags: Flags{G: true},
			want: []MatchResult{
				{Index: 1, Value: "1"},
				{Index: 3, Value: "22"},
				{Index: 6, Value: "333"},
			},
		},
		{
			name:  "非全局仅第一个",
			regex: `\d+`,
			text:  "a1b22c333",
			want:  []MatchResult{{Index: 1, Value: "1"}},
		},
		{
			name:  "忽略大小写",
			regex: `abc`,
			text:  "xABCy",
			flags: Flags{G: true, I: true},
			want:  []MatchResult{{Index: 1, Value: "ABC"}},
		},
		{
			name:  "多行模式",
			regex: `^foo`,
			text:  "bar\nfoo\nbaz",
			flags: Flags{G: true, M: true},
			want:  []MatchResult{{Index: 4, Value: "foo"}},
		},
		{
			name:  "单行模式点匹配换行",
			regex: `a.b`,
			text:  "a\nb",
			flags: Flags{G: true, S: true},
			want:  []MatchResult{{Index: 0, Value: "a\nb"}},
		},
		{
			name:  "编号捕获组",
			regex: `(\w+)@(\w+)`,
			text:  "u@d",
			flags: Flags{G: true},
			want: []MatchResult{
				{
					Index: 0, Value: "u@d",
					Captures: []GroupCapture{
						{Name: "1", Value: "u", Start: 0, End: 1},
						{Name: "2", Value: "d", Start: 2, End: 3},
					},
				},
			},
		},
		{
			name:  "命名捕获组",
			regex: `(?P<user>\w+)@(?P<domain>\w+)`,
			text:  "u@d",
			flags: Flags{G: true},
			want: []MatchResult{
				{
					Index: 0, Value: "u@d",
					Groups: []GroupCapture{
						{Name: "user", Value: "u", Start: 0, End: 1},
						{Name: "domain", Value: "d", Start: 2, End: 3},
					},
				},
			},
		},
		{
			name:  "未参与匹配的组被跳过",
			regex: `(a)|(b)`,
			text:  "a",
			flags: Flags{G: true},
			want: []MatchResult{
				{Index: 0, Value: "a", Captures: []GroupCapture{{Name: "1", Value: "a", Start: 0, End: 1}}},
			},
		},
		{
			name:  "无匹配",
			regex: `\d+`,
			text:  "abc",
			flags: Flags{G: true},
			want:  nil,
		},
		{
			name: "空正则",
			regex: "",
			text:  "abc",
			want:  nil,
		},
		{
			name: "空文本",
			regex: `a`,
			text:  "",
			want:  nil,
		},
		{
			name:    "非法正则",
			regex:   `(unclosed`,
			text:    "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.regex, tt.text, tt.flags)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望出错，但成功（%v）", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("返回错误: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("匹配数 = %d，期望 %d（%+v）", len(got), len(tt.want), got)
			}
			for i := range tt.want {
				w := tt.want[i]
				g := got[i]
				if g.Index != w.Index || g.Value != w.Value {
					t.Errorf("第 %d 个匹配 = %+v，期望 %+v", i, g, w)
				}
				if len(g.Captures) != len(w.Captures) || len(g.Groups) != len(w.Groups) {
					t.Errorf("第 %d 个匹配组 = %+v，期望 %+v", i, g, w)
					continue
				}
				for j := range w.Captures {
					if g.Captures[j] != w.Captures[j] {
						t.Errorf("第 %d 个匹配 captures[%d] = %+v，期望 %+v", i, j, g.Captures[j], w.Captures[j])
					}
				}
				for j := range w.Groups {
					if g.Groups[j] != w.Groups[j] {
						t.Errorf("第 %d 个匹配 groups[%d] = %+v，期望 %+v", i, j, g.Groups[j], w.Groups[j])
					}
				}
			}
		})
	}
}

func TestConvertZeroLength(t *testing.T) {
	// 零长匹配不应 panic（Go FindAll 自带 advance 处理）。
	_, err := Convert(`a*`, "bbb", Flags{G: true})
	if err != nil {
		t.Fatalf("返回错误: %v", err)
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	in, _ := json.Marshal(input{Regex: `(\w+)`, Text: "hello", Flags: Flags{G: true}})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if len(out.Matches) != 1 || out.Matches[0].Value != "hello" {
		t.Errorf("结果 = %+v", out.Matches)
	}
}
