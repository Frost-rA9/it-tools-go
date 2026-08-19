package ipv4rangeexpander

import (
	"encoding/json"
	"testing"
)

func TestCalculateKnownVector(t *testing.T) {
	out, err := calculate("192.168.1.1", "192.168.6.255")
	if err != nil {
		t.Fatalf("calculate 失败: %v", err)
	}
	want := output{
		OldSize:  1535,
		NewStart: "192.168.0.0",
		NewEnd:   "192.168.7.255",
		NewCIDR:  "192.168.0.0/21",
		NewSize:  2048,
	}
	if out != want {
		t.Errorf("已知向量不符:\n  期望 %+v\n  实际 %+v", want, out)
	}
}

func TestCalculateCases(t *testing.T) {
	cases := []struct {
		name       string
		start, end string
		newStart   string
		newEnd     string
		newCidr    string
		oldSize    uint64
		newSize    uint64
	}{
		{"单地址", "10.0.0.5", "10.0.0.5", "10.0.0.5", "10.0.0.5", "10.0.0.5/32", 1, 1},
		{"整块", "192.168.0.0", "192.168.0.255", "192.168.0.0", "192.168.0.255", "192.168.0.0/24", 256, 256},
		{"跨 /24", "192.168.0.10", "192.168.1.10", "192.168.0.0", "192.168.1.255", "192.168.0.0/23", 257, 512},
		{"全范围", "0.0.0.0", "255.255.255.255", "0.0.0.0", "255.255.255.255", "0.0.0.0/0", 4294967296, 4294967296},
		{"非对齐内退", "10.0.0.100", "10.0.0.200", "10.0.0.0", "10.0.0.255", "10.0.0.0/24", 101, 256},
	}
	for _, c := range cases {
		out, err := calculate(c.start, c.end)
		if err != nil {
			t.Errorf("%s: 计算失败: %v", c.name, err)
			continue
		}
		if out.NewStart != c.newStart || out.NewEnd != c.newEnd || out.NewCIDR != c.newCidr ||
			out.OldSize != c.oldSize || out.NewSize != c.newSize {
			t.Errorf("%s: 不符 newStart=%s newEnd=%s cidr=%s old=%d new=%d（实际 %s/%s/%s/%d/%d）",
				c.name, c.newStart, c.newEnd, c.newCidr, c.oldSize, c.newSize,
				out.NewStart, out.NewEnd, out.NewCIDR, out.OldSize, out.NewSize)
		}
	}
}

func TestCalculateEndBelowStart(t *testing.T) {
	if _, err := calculate("192.168.6.255", "192.168.1.1"); err == nil {
		t.Errorf("end < start: 期望报错")
	}
}

func TestCalculateErrors(t *testing.T) {
	cases := []struct{ start, end string }{
		{"", "1.2.3.4"},
		{"1.2.3.4", ""},
		{"999.1.1.1", "1.2.3.4"},
		{"1.2.3.4", "1.2.3"},
		{"1.2.3.4", "1.2.3.x"},
	}
	for _, c := range cases {
		if _, err := calculate(c.start, c.end); err == nil {
			t.Errorf("start=%q end=%q: 期望报错", c.start, c.end)
		}
	}
}

func TestExecuteEndToEnd(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{StartIP: "192.168.1.1", EndIP: "192.168.6.255"})
	outJSON, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("输出解析失败: %v", err)
	}
	if out.NewCIDR != "192.168.0.0/21" {
		t.Errorf("Execute 输出不符: %+v", out)
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	for _, in := range []string{
		`{bad`,
		`{"startIp":"1.2.3.4"}`,
		`{"startIp":"2.2.2.2","endIp":"1.1.1.1"}`,
		`{"startIp":"bad","endIp":"1.1.1.1"}`,
	} {
		if _, err := exec.Execute(t.Context(), in); err == nil {
			t.Errorf("输入 %q: 期望报错", in)
		}
	}
}

func TestToolMetadata(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != "网络" || meta.Icon != "ArrowsMaximize" {
		t.Errorf("元数据不一致: %+v", meta)
	}
	if meta.Name == "" || meta.Description == "" || len(meta.Keywords) == 0 {
		t.Errorf("元数据缺字段: %+v", meta)
	}
}
