package uuidgen

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func TestExecuteV4(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Version: "v4", Count: 5})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("v4 生成失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	lines := strings.Split(o.Result, "\n")
	if len(lines) != 5 {
		t.Fatalf("期望 5 行，实际 %d", len(lines))
	}
	seen := map[string]bool{}
	for _, line := range lines {
		if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(line) {
			t.Errorf("v4 格式不正确: %q", line)
		}
		if seen[line] {
			t.Errorf("v4 出现重复: %q", line)
		}
		seen[line] = true
	}
}

func TestExecuteNIL(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Version: "NIL", Count: 1})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatal(err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	if o.Result != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("NIL 结果: %q", o.Result)
	}
}

func TestExecuteV1(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Version: "v1", Count: 2})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("v1 生成失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	for _, line := range strings.Split(o.Result, "\n") {
		if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(line) {
			t.Errorf("v1 格式不正确: %q", line)
		}
	}
}

func TestExecuteV3V5Deterministic(t *testing.T) {
	exec := Executor{}
	ns := "6ba7b811-9dad-11d1-80b4-00c04fd430c8" // URL 命名空间
	for _, ver := range []string{"v3", "v5"} {
		raw, _ := json.Marshal(input{Version: ver, Count: 2, Namespace: ns, Name: "it-tools"})
		out1, err := exec.Execute(t.Context(), string(raw))
		if err != nil {
			t.Fatalf("%s 生成失败: %v", ver, err)
		}
		out2, _ := exec.Execute(t.Context(), string(raw))
		var o1, o2 output
		_ = json.Unmarshal([]byte(out1), &o1)
		_ = json.Unmarshal([]byte(out2), &o2)
		if o1.Result != o2.Result {
			t.Errorf("%s 应确定性生成相同结果", ver)
		}
		lines := strings.Split(o1.Result, "\n")
		if len(lines) != 2 || lines[0] != lines[1] {
			t.Errorf("%s 同参应生成相同 UUID", ver)
		}
	}
}

func TestExecuteV5KnownVector(t *testing.T) {
	// v5(DNS 命名空间, "www.example.com") 的已知标准向量。
	exec := Executor{}
	raw, _ := json.Marshal(input{Version: "v5", Count: 1, Namespace: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", Name: "www.example.com"})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatal(err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	if o.Result != "2ed6657d-e927-568b-95e1-2665a8aea6a2" {
		t.Errorf("v5 向量不符: %q", o.Result)
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name  string
		input input
	}{
		{"未知版本", input{Version: "v2", Count: 1}},
		{"数量 0", input{Version: "v4", Count: 0}},
		{"数量超限", input{Version: "v4", Count: 51}},
		{"v3 非法命名空间", input{Version: "v3", Count: 1, Namespace: "bad", Name: "x"}},
		{"v3 缺名称", input{Version: "v3", Count: 1, Namespace: "6ba7b811-9dad-11d1-80b4-00c04fd430c8"}},
	}
	for _, c := range cases {
		raw, _ := json.Marshal(c.input)
		if _, err := exec.Execute(t.Context(), string(raw)); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}
