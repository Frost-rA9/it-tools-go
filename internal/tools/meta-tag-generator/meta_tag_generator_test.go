package metatag

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"siteName", "site_name"},
		{"HTTPServer", "http_server"},
		{"image:alt", "image:alt"},
		{"article:publishedTime", "article:published_time"},
		{"type", "type"},
		{"sIPAddress", "s_ip_address"},
	}
	for _, tt := range tests {
		if got := toSnakeCase(tt.in); got != tt.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFlattenMetadata(t *testing.T) {
	got := flattenMetadata(map[string]any{
		"type":  "website",
		"title": "Hello",
		"image": "", // 空字符串应跳过
		"tags":  []any{"a", "b"},
	}, "og")
	want := []metaFlat{
		{Key: "og:tags", Value: "a"},
		{Key: "og:tags", Value: "b"},
		{Key: "og:title", Value: "Hello"},
		{Key: "og:type", Value: "website"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("flatten = %+v, want %+v", got, want)
	}
}

func TestPickTwitterCompatible(t *testing.T) {
	existing := []metaFlat{
		{Key: "og:title", Value: "T"},
		{Key: "og:description", Value: "D"},
		{Key: "og:image", Value: "img.png"},
	}
	// twitter 已有 twitter:title，则仅补 description 与 image。
	twitter := []metaFlat{{Key: "twitter:title", Value: "TW"}}
	got := pickTwitterCompatible(existing, twitter)
	want := []metaFlat{
		{Key: "twitter:description", Value: "D"},
		{Key: "twitter:image", Value: "img.png"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pick = %+v, want %+v", got, want)
	}
}

func TestGenerateMeta(t *testing.T) {
	got := generateMeta(
		map[string]any{"type": "website", "title": "Hello"},
		map[string]any{"card": "summary", "site": "@x"},
		true,
	)
	want := `<!-- og meta -->
<meta property="og:title" value="Hello" />
<meta property="og:type" value="website" />

<!-- twitter meta -->
<meta name="twitter:card" value="summary" />
<meta name="twitter:site" value="@x" />
<meta name="twitter:title" value="Hello" />`
	if got != want {
		t.Errorf("generateMeta 输出不符:\n%s\n--- want ---\n%s", got, want)
	}
}

func TestExecuteSchemas(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"action":"schemas"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(o.Schemas.Base) != 3 {
		t.Errorf("base 区块 = %d, want 3", len(o.Schemas.Base))
	}
	if len(o.Schemas.Types) != 11 {
		t.Errorf("types 数量 = %d, want 11", len(o.Schemas.Types))
	}
	// 校验 website 的 type 选择器含分组选项。
	general := o.Schemas.Base[0]
	if general.Name != "General information" {
		t.Errorf("首个 base 区块 = %q", general.Name)
	}
	var hasGroup bool
	for _, e := range general.Elements {
		if e.Key == "type" {
			for _, o := range e.Options {
				if o.Type == "group" {
					hasGroup = true
				}
			}
		}
	}
	if !hasGroup {
		t.Error("type 选择器应包含分组选项")
	}
}

func TestExecuteMeta(t *testing.T) {
	var e Executor
	input := `{"action":"meta","metadata":{"type":"article","title":"T","article:tag":["a","b"],"image":"","twitter:card":"summary_large_image"}}`
	out, err := e.Execute(t.Context(), input)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if !strings.Contains(o.Meta, `<meta property="og:article:tag" value="a" />`) {
		t.Errorf("缺少数组多标签: %s", o.Meta)
	}
	if !strings.Contains(o.Meta, `<meta property="og:article:tag" value="b" />`) {
		t.Errorf("缺少数组多标签 b: %s", o.Meta)
	}
	if strings.Contains(o.Meta, `og:image"`) {
		t.Errorf("空 image 不应出现: %s", o.Meta)
	}
	if !strings.Contains(o.Meta, `<meta name="twitter:title" value="T" />`) {
		t.Errorf("缺少 twitter 兼容字段: %s", o.Meta)
	}
}

func TestExecuteErrors(t *testing.T) {
	var e Executor
	tests := []struct {
		name  string
		input string
	}{
		{"未知操作", `{"action":"foo"}`},
		{"缺少 metadata", `{"action":"meta"}`},
		{"非法输入", `not-json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := e.Execute(t.Context(), tt.input); err == nil {
				t.Errorf("Execute(%s) 期望错误", tt.input)
			}
		})
	}
}