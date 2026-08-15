package ulidgen

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func TestNewULIDFormat(t *testing.T) {
	id, err := newULID()
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	if len(id) != 26 {
		t.Errorf("ULID 长度应为 26，实际 %d", len(id))
	}
	if !regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`).MatchString(id) {
		t.Errorf("ULID 含非法字符: %q", id)
	}
}

func TestEncodeULIDKnownVector(t *testing.T) {
	// 参考实现（oklog/ulid 与 npm ulid 一致）向量：时间戳 1469918176385、随机数全零。
	got := encodeULID(1469918176385, [10]byte{})
	if got != "01ARYZ6S410000000000000000" {
		t.Errorf("参考向量不符: got %q", got)
	}
}

// ulidTimestamp 从 ULID 字符串解码 48 位时间戳（前 10 字符的 50 位去掉末尾 2 个随机位）。
func ulidTimestamp(t *testing.T, id string) uint64 {
	t.Helper()
	var ts uint64
	for i := 0; i < 10; i++ {
		idx := strings.IndexByte(crockford, id[i])
		if idx < 0 {
			t.Fatalf("非法字符: %q", id[i])
		}
		ts = ts*32 + uint64(idx)
	}
	return ts >> 2
}

func TestNewULIDMonotonicTimestamp(t *testing.T) {
	// 时间戳部分随调用保持单调不减（48 位毫秒时间戳）。
	prev := uint64(0)
	for i := 0; i < 200; i++ {
		id, err := newULID()
		if err != nil {
			t.Fatal(err)
		}
		ts := ulidTimestamp(t, id)
		if ts < prev {
			t.Errorf("ULID 时间戳回退: %d < %d", ts, prev)
		}
		prev = ts
	}
}

func TestNewULIDUnique(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		id, err := newULID()
		if err != nil {
			t.Fatal(err)
		}
		if seen[id] {
			t.Fatalf("出现重复 ULID: %q", id)
		}
		seen[id] = true
	}
}

func TestExecuteRaw(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Count: 3, Format: FormatRaw})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("raw 生成失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	lines := strings.Split(o.Result, "\n")
	if len(lines) != 3 {
		t.Fatalf("期望 3 行，实际 %d", len(lines))
	}
	for _, line := range lines {
		if len(line) != 26 {
			t.Errorf("ULID 长度不正确: %q", line)
		}
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Count: 2, Format: FormatJSON})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("json 生成失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	var ids []string
	if err := json.Unmarshal([]byte(o.Result), &ids); err != nil {
		t.Fatalf("结果不是有效 JSON 数组: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("JSON 数组长度应为 2，实际 %d", len(ids))
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name  string
		input input
	}{
		{"数量 0", input{Count: 0, Format: FormatRaw}},
		{"数量超限", input{Count: 101, Format: FormatRaw}},
		{"未知格式", input{Count: 1, Format: "bad"}},
	}
	for _, c := range cases {
		raw, _ := json.Marshal(c.input)
		if _, err := exec.Execute(t.Context(), string(raw)); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}
