package httpstatuscodes

import (
	"encoding/json"
	"testing"
)

func TestSearchEmptyReturnsAll(t *testing.T) {
	res := search("")
	if len(res) != 5 {
		t.Fatalf("分类数 = %d, want 5", len(res))
	}
	total := 0
	for _, cat := range res {
		total += len(cat.Codes)
	}
	if total != 63 {
		t.Errorf("状态码总数 = %d, want 63", total)
	}
}

func TestSearchByCode(t *testing.T) {
	res := search("404")
	if len(res) != 1 || res[0].Category != "Search results" || len(res[0].Codes) != 1 {
		t.Fatalf("搜索结果不符: %+v", res)
	}
	code := res[0].Codes[0]
	if code.Code != 404 || code.Name != "Not Found" {
		t.Errorf("命中不符: %+v", code)
	}
}

func TestSearchByName(t *testing.T) {
	res := search("teapot")
	if len(res) != 1 || len(res[0].Codes) != 1 {
		t.Fatalf("搜索结果不符: %+v", res)
	}
	if res[0].Codes[0].Code != 418 {
		t.Errorf("命中不符: %+v", res[0].Codes[0])
	}
}

func TestSearchByDescription(t *testing.T) {
	res := search("coffee")
	if len(res) != 1 || len(res[0].Codes) != 1 {
		t.Fatalf("搜索结果不符: %+v", res)
	}
	if res[0].Codes[0].Code != 418 {
		t.Errorf("命中不符: %+v", res[0].Codes[0])
	}
}

func TestSearchByCategory(t *testing.T) {
	res := search("redirection")
	if len(res) != 1 || len(res[0].Codes) != 9 {
		t.Fatalf("搜索结果不符: %+v", res)
	}
	if res[0].Category != "Search results" {
		t.Errorf("分类不符: %+v", res[0])
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	res := search("CREATED")
	if len(res) != 1 || len(res[0].Codes) != 1 || res[0].Codes[0].Code != 201 {
		t.Errorf("大小写不敏感匹配不符: %+v", res)
	}
}

func TestSearchNoMatch(t *testing.T) {
	if res := search("zzzzz"); res != nil {
		t.Errorf("期望无结果，得到 %+v", res)
	}
}

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"query":"404"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(got.Results) != 1 || got.Results[0].Codes[0].Code != 404 {
		t.Errorf("Execute 结果不符: %+v", got)
	}

	if _, err := e.Execute(t.Context(), `not-json`); err == nil {
		t.Error("非法输入应报错")
	}
}