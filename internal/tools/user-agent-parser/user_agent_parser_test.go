package useragentparser

import (
	"encoding/json"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want output
	}{
		{
			name: "Chrome",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
			want: output{
				Browser: nameVersion{Name: "Chrome", Version: "121.0.0"},
				Os:      nameVersion{Name: "Windows", Version: "10"},
				Device:  deviceInfo{Model: "Other"},
			},
		},
		{
			name: "Firefox",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0",
			want: output{
				Browser: nameVersion{Name: "Firefox", Version: "122.0"},
				Os:      nameVersion{Name: "Windows", Version: "10"},
				Device:  deviceInfo{Model: "Other"},
			},
		},
		{
			name: "Safari",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
			want: output{
				Browser: nameVersion{Name: "Safari", Version: "17.2"},
				Os:      nameVersion{Name: "Mac OS X", Version: "10.15.7"},
				Device:  deviceInfo{Vendor: "Apple", Model: "Mac"},
			},
		},
		{
			name: "Edge",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.2210.121",
			want: output{
				Browser: nameVersion{Name: "Edge", Version: "120.0.2210"},
				Os:      nameVersion{Name: "Windows", Version: "10"},
				Device:  deviceInfo{Model: "Other"},
			},
		},
		{
			name: "iPhone",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
			want: output{
				Browser: nameVersion{Name: "Mobile Safari", Version: "17.2"},
				Os:      nameVersion{Name: "iOS", Version: "17.2"},
				Device:  deviceInfo{Vendor: "Apple", Model: "iPhone"},
			},
		},
		{
			name: "未知字符串",
			ua:   "hello world",
			want: output{
				Browser: nameVersion{Name: "Other"},
				Os:      nameVersion{Name: "Other"},
				Device:  deviceInfo{Model: "Other"},
			},
		},
		{
			name: "空字符串",
			ua:   "",
			want: output{},
		},
		{
			name: "仅空白",
			ua:   "   ",
			want: output{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Executor
			out, err := e.Execute(t.Context(), `{"ua":"`+tt.ua+`"}`)
			if err != nil {
				t.Fatalf("Execute 意外错误: %v", err)
			}
			var got output
			if err := json.Unmarshal([]byte(out), &got); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if got != tt.want {
				t.Errorf("Execute 结果 = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestExecuteInvalidInput(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `not-json`); err == nil {
		t.Error("非法输入应报错")
	}
}