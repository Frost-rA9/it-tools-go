package caseconv

import (
	"context"
	"encoding/json"
	"testing"
)

func resultMap(results []format) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Label] = r.Value
	}
	return m
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]string
	}{
		{
			name: "空格分隔小写词",
			text: "lorem ipsum dolor sit amet",
			want: map[string]string{
				"Lowercase":    "lorem ipsum dolor sit amet",
				"Uppercase":    "LOREM IPSUM DOLOR SIT AMET",
				"Camelcase":    "loremIpsumDolorSitAmet",
				"Capitalcase":  "Lorem Ipsum Dolor Sit Amet",
				"Constantcase": "LOREM_IPSUM_DOLOR_SIT_AMET",
				"Dotcase":      "lorem.ipsum.dolor.sit.amet",
				"Headercase":   "Lorem-Ipsum-Dolor-Sit-Amet",
				"Nocase":       "lorem ipsum dolor sit amet",
				"Paramcase":    "lorem-ipsum-dolor-sit-amet",
				"Pascalcase":   "LoremIpsumDolorSitAmet",
				"Pathcase":     "lorem/ipsum/dolor/sit/amet",
				"Sentencecase": "Lorem ipsum dolor sit amet",
				"Snakecase":    "lorem_ipsum_dolor_sit_amet",
				"Mockingcase":  "LoReM IpSuM DoLoR SiT AmEt",
			},
		},
		{
			name: "驼峰输入拆分",
			text: "helloWorld",
			want: map[string]string{
				"Lowercase":   "helloworld",
				"Uppercase":   "HELLOWORLD",
				"Camelcase":   "helloWorld",
				"Pascalcase":  "HelloWorld",
				"Snakecase":   "hello_world",
				"Constantcase": "HELLO_WORLD",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resultMap(Convert(tt.text))
			for label, want := range tt.want {
				if got[label] != want {
					t.Errorf("Convert(%q)[%s] = %q，期望 %q", tt.text, label, got[label], want)
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	in, err := json.Marshal(input{Text: "lorem ipsum dolor sit amet"})
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}

	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}

	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if len(out.Results) != 14 {
		t.Fatalf("结果数量 = %d，期望 14", len(out.Results))
	}

	m := resultMap(out.Results)
	if m["Camelcase"] != "loremIpsumDolorSitAmet" {
		t.Errorf("Camelcase = %q", m["Camelcase"])
	}
	if m["Snakecase"] != "lorem_ipsum_dolor_sit_amet" {
		t.Errorf("Snakecase = %q", m["Snakecase"])
	}
}
