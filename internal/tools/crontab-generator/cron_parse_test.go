package crontab

import "testing"

func TestParseCron(t *testing.T) {
	cases := []struct {
		expr string
		want string
	}{
		{"* * * * *", "每分钟"},
		{"*/15 * * * *", "每 15 分钟"},
		{"40 * * * *", "每小时第 40 分钟"},
		{"0 9-17 * * *", "第 0 分钟，9 至 17 点"},
		{"0 9-17 * * 1-5", "第 0 分钟，9 至 17 点，周一至周五"},
		{"0 0 1 1 *", "0 时 0 分，每月 1 日，1 月"},
		{"0 0 * * 0", "0 时 0 分，周日"},
		{"30 4 1,15 * *", "4 时 30 分，每月 1、15 日"},
		{"0 0 * * mon-fri", "0 时 0 分，周一至周五"},
		{"0 0 1 jan *", "0 时 0 分，每月 1 日，1 月"},
		{"@daily", "每天 0 时 0 分执行"},
		{"@yearly", "每年 1 月 1 日 0 时 0 分执行"},
		{"@reboot", "系统启动时执行"},
		{"0 40 * * * *", "第 0 秒，每小时第 40 分钟"},
		{"*/10 0 1 * *", "每 10 分钟，第 0 点，每月 1 日"},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			got, err := parseCron(c.expr)
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("描述不符: got %q, want %q", got, c.want)
			}
		})
	}
}

func TestParseCronErrors(t *testing.T) {
	bad := []string{
		"",
		"* * *",
		"61 * * * *",
		"* 24 * * *",
		"* * 0 * *",
		"* * * 13 *",
		"* * * * 8",
		"abc * * * *",
		"*/0 * * * *",
		"10-5 * * * *",
		"@invalid",
	}
	for _, expr := range bad {
		t.Run(expr, func(t *testing.T) {
			if _, err := parseCron(expr); err == nil {
				t.Fatalf("期望报错: %q", expr)
			}
		})
	}
}

func TestExecuteParseAction(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"action":"parse","expression":"*/5 * * * *"}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	if want := `{"expression":"*/5 * * * *","summary":"每 5 分钟"}`; out != want {
		t.Fatalf("输出不符: got %s, want %s", out, want)
	}

	if _, err := exec.Execute(t.Context(), `{"action":"parse","expression":"61 * * * *"}`); err == nil {
		t.Fatal("非法表达式应报错")
	}
}
