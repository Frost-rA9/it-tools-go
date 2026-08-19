package mathevaluator

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"it-tools-go/internal/registry"
)

func TestExecuteSuccess(t *testing.T) {
	outStr, err := (Executor{}).Execute(context.Background(), `{"expression":"2+3*4"}`)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if out.Result != "14" || out.Error != "" {
		t.Errorf("输出不符: %+v", out)
	}
}

func TestExecuteEvalError(t *testing.T) {
	outStr, err := (Executor{}).Execute(context.Background(), `{"expression":"sqrt(-1)"}`)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if out.Result != "" || len(out.Error) == 0 {
		t.Errorf("输出不符: %+v", out)
	}
}

func TestExecuteEmpty(t *testing.T) {
	for _, expr := range []string{"", "   "} {
		outStr, err := (Executor{}).Execute(context.Background(), `{"expression":"`+expr+`"}`)
		if err != nil {
			t.Fatalf("Execute(%q) 返回错误: %v", expr, err)
		}
		var out output
		if err := json.Unmarshal([]byte(outStr), &out); err != nil {
			t.Fatalf("解析输出失败: %v", err)
		}
		if out.Result != "" || out.Error != "" {
			t.Errorf("空表达式应返回空结果: %+v", out)
		}
	}
}

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != registry.CategoryMath || meta.Icon != "Math" {
		t.Errorf("元数据不符: %+v", meta)
	}
	if !reflect.DeepEqual(meta.Keywords, Keywords) {
		t.Errorf("Keywords 不符: %v", meta.Keywords)
	}
}
