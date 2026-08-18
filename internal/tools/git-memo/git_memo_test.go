package gitmemo

import "testing"

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{}`)
	if err != nil || out != `{}` {
		t.Errorf("Execute() = %q, %v", out, err)
	}
}

func TestExecuteInvalidInput(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `not-json`); err == nil {
		t.Error("非法输入应报错")
	}
}
