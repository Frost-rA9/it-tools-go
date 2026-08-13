package romannumeral

import (
	"context"
	"encoding/json"
	"testing"
)

func TestArabicToRoman(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want string
	}{
		{"最小值", 1, "I"},
		{"常规", 42, "XLII"},
		{"含减法规则", 4, "IV"},
		{"含减法规则2", 9, "IX"},
		{"复杂", 1999, "MCMXCIX"},
		{"复杂2", 3549, "MMMDXLIX"},
		{"最大值", 3999, "MMMCMXCIX"},
		{"小于下界", 0, ""},
		{"负数", -5, ""},
		{"大于上界", 4000, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArabicToRoman(tt.num); got != tt.want {
				t.Errorf("ArabicToRoman(%d) = %q，期望 %q", tt.num, got, tt.want)
			}
		})
	}
}

func TestRomanToArabic(t *testing.T) {
	tests := []struct {
		name  string
		roman string
		want  int
		wantOK bool
	}{
		{"最小值", "I", 1, true},
		{"常规", "XLII", 42, true},
		{"减法规则", "IV", 4, true},
		{"复杂", "MCMXCIX", 1999, true},
		{"最大值", "MMMCMXCIX", 3999, true},
		{"非法-四个连续", "IIII", 0, false},
		{"非法-超上限", "MMMM", 0, false},
		{"非法-字母", "abc", 0, false},
		{"空字符串", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := RomanToArabic(tt.roman)
			if ok != tt.wantOK || got != tt.want {
				t.Errorf("RomanToArabic(%q) = (%d, %v)，期望 (%d, %v)", tt.roman, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	tests := []struct {
		name    string
		input   input
		want    string
		wantErr bool
	}{
		{name: "阿拉伯转罗马", input: input{Value: "42", Mode: ModeArabicToRoman}, want: "XLII"},
		{name: "罗马转阿拉伯", input: input{Value: "XLII", Mode: ModeRomanToArabic}, want: "42"},
		{name: "罗马小写转阿拉伯", input: input{Value: "xlii", Mode: ModeRomanToArabic}, want: "42"},
		{name: "阿拉伯越界返回空", input: input{Value: "4000", Mode: ModeArabicToRoman}, want: ""},
		{name: "罗马非法返回空", input: input{Value: "IIII", Mode: ModeRomanToArabic}, want: ""},
		{name: "阿拉伯非法输入报错", input: input{Value: "abc", Mode: ModeArabicToRoman}, wantErr: true},
		{name: "未知模式报错", input: input{Value: "1", Mode: "foo"}, wantErr: true},
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
