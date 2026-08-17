package jsondiff

import (
	"encoding/json"
	"testing"
)

func TestExecuteSame(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"left":"{\"a\":1}","right":"{ \"a\": 1 }"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if !got.Same || got.Root != nil {
		t.Errorf("相同 JSON 应 same=true 且 root=null，得到 %+v", got)
	}
}

func TestExecuteAddedRemoved(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"left":"{\"a\":1,\"b\":2}","right":"{\"a\":1,\"c\":3}"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if got.Same || got.Root == nil || got.Root.Status != "children-updated" {
		t.Fatalf("结果不符: %+v", got)
	}
	statuses := map[string]string{}
	for _, c := range got.Root.Children {
		statuses[c.Key.(string)] = c.Status
	}
	if statuses["a"] != "unchanged" || statuses["b"] != "removed" || statuses["c"] != "added" {
		t.Errorf("子节点状态不符: %+v", statuses)
	}
}

func TestExecuteUpdated(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"left":"{\"a\":1}","right":"{\"a\":2}"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if got.Same {
		t.Fatal("应判定为不同")
	}
	child := got.Root.Children[0]
	if child.Status != "updated" || child.OldValue != float64(1) || child.Value != float64(2) {
		t.Errorf("updated 节点不符: %+v", child)
	}
}

func TestExecuteNestedAndArrays(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(),
		`{"left":"{\"obj\":{\"x\":1},\"arr\":[1,2,3]}","right":"{\"obj\":{\"x\":2},\"arr\":[1,9]}"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if got.Same {
		t.Fatal("应判定为不同")
	}
	statuses := map[string]string{}
	for _, c := range got.Root.Children {
		statuses[c.Key.(string)] = c.Status
	}
	if statuses["obj"] != "children-updated" || statuses["arr"] != "children-updated" {
		t.Errorf("子节点状态不符: %+v", statuses)
	}
}

func TestExecuteOnlyShowDifferences(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"left":"{\"a\":1,\"b\":2}","right":"{\"a\":1,\"b\":3}","only_show_differences":true}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(got.Root.Children) != 1 || got.Root.Children[0].Key.(string) != "b" {
		t.Errorf("仅差异过滤不符: %+v", got.Root.Children)
	}
}

func TestExecuteJson5(t *testing.T) {
	// JSON5：注释、单引号、尾逗号。
	var e Executor
	out, err := e.Execute(t.Context(), `{"left":"// comment\n{a:1,}","right":"{a:2}"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if got.Same {
		t.Fatal("应判定为不同")
	}
	if got.Root.Children[0].Status != "updated" {
		t.Errorf("JSON5 解析结果不符: %+v", got.Root.Children[0])
	}
}

func TestExecuteNullVsMissing(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"left":"{\"a\":null}","right":"{}"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if got.Same {
		t.Fatal("应判定为不同")
	}
	if got.Root.Children[0].Status != "removed" {
		t.Errorf("a:null vs 缺失 应 removed，得到 %+v", got.Root.Children[0])
	}
}

func TestExecuteErrors(t *testing.T) {
	var e Executor
	tests := []struct {
		name  string
		input string
	}{
		{"左无效", `{"left":"{bad","right":"{}"}`},
		{"右无效", `{"left":"{}","right":"{bad"}`},
		{"左为空", `{"left":"","right":"{}"}`},
		{"非法输入 JSON", `not-json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := e.Execute(t.Context(), tt.input); err == nil {
				t.Errorf("Execute(%s) 期望错误", tt.input)
			}
		})
	}
}