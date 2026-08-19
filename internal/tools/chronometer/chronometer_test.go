package chronometer

import (
	"context"
	"testing"

	"it-tools-go/internal/registry"
)

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Name != "秒表" || meta.Category != registry.CategoryMeasurement || meta.Icon != "Clock" {
		t.Errorf("元数据不符: %+v", meta)
	}
}

func TestExecutorReturnsEmpty(t *testing.T) {
	out, err := (Executor{}).Execute(context.Background(), `{}`)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	if out != `{}` {
		t.Errorf("输出 = %q, want {}", out)
	}
}

func TestExecuteInvalidJSON(t *testing.T) {
	if _, err := (Executor{}).Execute(context.Background(), `not json`); err == nil {
		t.Error("非法 JSON 应返回错误")
	}
}
