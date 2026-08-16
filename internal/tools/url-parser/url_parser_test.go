package urlparser

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    output
		wantErr bool
	}{
		{
			name:  "样例 URL",
			input: `{"url":"https://user:pass@example.com:3000/path?key1=value&key2=value2#the-hash"}`,
			want: output{
				Protocol: "https:",
				Username: "user",
				Password: "pass",
				Hostname: "example.com",
				Port:     "3000",
				Pathname: "/path",
				Search:   "?key1=value&key2=value2",
				Params:   []Param{{Key: "key1", Value: "value"}, {Key: "key2", Value: "value2"}},
			},
		},
		{
			name:  "无查询参数",
			input: `{"url":"http://example.com/path"}`,
			want: output{
				Protocol: "http:",
				Hostname: "example.com",
				Pathname: "/path",
			},
		},
		{
			name:  "无路径",
			input: `{"url":"https://example.com"}`,
			want: output{
				Protocol: "https:",
				Hostname: "example.com",
				Pathname: "",
			},
		},
		{
			name:  "空值参数与加号",
			input: `{"url":"https://example.com/?a=&b=hello+world&b=again"}`,
			want: output{
				Protocol: "https:",
				Hostname: "example.com",
				Pathname: "/",
				Search:   "?a=&b=hello+world&b=again",
				Params:   []Param{{Key: "a", Value: ""}, {Key: "b", Value: "hello world"}, {Key: "b", Value: "again"}},
			},
		},
		{
			name:  "无用户信息",
			input: `{"url":"ftp://files.example.org/pub?x=1"}`,
			want: output{
				Protocol: "ftp:",
				Hostname: "files.example.org",
				Pathname: "/pub",
				Search:   "?x=1",
				Params:   []Param{{Key: "x", Value: "1"}},
			},
		},
		{
			name:    "相对地址非法",
			input:   `{"url":"/path/to/file"}`,
			wantErr: true,
		},
		{
			name:    "无 scheme 非法",
			input:   `{"url":"example.com/path"}`,
			wantErr: true,
		},
		{
			name:    "完全非法",
			input:   `{"url":"://bad"}`,
			wantErr: true,
		},
		{
			name:    "非法输入 JSON",
			input:   `not-json`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Executor
			out, err := e.Execute(t.Context(), tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute 期望错误，得到 %s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute 意外错误: %v", err)
			}
			var got output
			if err := json.Unmarshal([]byte(out), &got); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute 结果 = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []Param
	}{
		{"空", "", nil},
		{"单个", "a=1", []Param{{Key: "a", Value: "1"}}},
		{"重复键", "a=1&a=2", []Param{{Key: "a", Value: "1"}, {Key: "a", Value: "2"}}},
		{"无值键", "a=&b", []Param{{Key: "a", Value: ""}, {Key: "b", Value: ""}}},
		{"加号为空格", "q=hello+world", []Param{{Key: "q", Value: "hello world"}}},
		{"中文转义", "q=%E4%BD%A0%E5%A5%BD", []Param{{Key: "q", Value: "你好"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseParams(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseParams(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}