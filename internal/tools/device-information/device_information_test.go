package deviceinformation

import "testing"

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"x":1}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	if out != "{}" {
		t.Errorf("输出应为空对象，得到 %s", out)
	}
}

func TestExecuteInvalidInput(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `not-json`); err == nil {
		t.Error("非法输入应报错")
	}
}