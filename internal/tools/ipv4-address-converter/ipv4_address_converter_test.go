package ipv4addressconv

import (
	"encoding/json"
	"testing"
)

func TestExecuteKnownVector(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{IP: "192.168.1.1"})
	outJSON, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("输出解析失败: %v", err)
	}
	want := output{
		Decimal:     3232235777,
		Hexadecimal: "C0A80101",
		Binary:      "11000000101010000000000100000001",
		IPv6:        "0000:0000:0000:0000:0000:ffff:c0a8:0101",
		IPv6Short:   "::ffff:c0a8:0101",
	}
	if out != want {
		t.Errorf("已知向量不符:\n  期望 %+v\n  实际 %+v", want, out)
	}
}

func TestExecuteEdges(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name string
		ip   string
		dec  uint32
		hex  string
		bin  string
	}{
		{"0.0.0.0", "0.0.0.0", 0, "00000000", "00000000000000000000000000000000"},
		{"255.255.255.255", "255.255.255.255", 4294967295, "FFFFFFFF", "11111111111111111111111111111111"},
		{"10.0.0.1", "10.0.0.1", 167772161, "0A000001", "00001010000000000000000000000001"},
	}
	for _, c := range cases {
		raw, _ := json.Marshal(input{IP: c.ip})
		outJSON, err := exec.Execute(t.Context(), string(raw))
		if err != nil {
			t.Errorf("%s: 生成失败: %v", c.name, err)
			continue
		}
		var out output
		_ = json.Unmarshal([]byte(outJSON), &out)
		if out.Decimal != c.dec || out.Hexadecimal != c.hex || out.Binary != c.bin {
			t.Errorf("%s: 不符 dec=%d hex=%s bin=%s（实际 %d/%s/%s）",
				c.name, c.dec, c.hex, c.bin, out.Decimal, out.Hexadecimal, out.Binary)
		}
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []input{
		{IP: ""},
		{IP: "999.1.1.1"},
		{IP: "1.2.3"},
		{IP: "1.2.3.4.5"},
		{IP: "1.2.3.x"},
		{IP: "192.168.1"},
	}
	for _, in := range cases {
		raw, _ := json.Marshal(in)
		if _, err := exec.Execute(t.Context(), string(raw)); err == nil {
			t.Errorf("输入 %+v: 期望报错", in)
		}
	}
	// 非法 JSON。
	if _, err := exec.Execute(t.Context(), `{bad`); err == nil {
		t.Errorf("非法 JSON: 期望报错")
	}
}

func TestIP6Mapped(t *testing.T) {
	// 独立验证 IPv6 映射组格式。
	full, short := ipv6Mapped(3232235777) // 192.168.1.1
	if full != "0000:0000:0000:0000:0000:ffff:c0a8:0101" || short != "::ffff:c0a8:0101" {
		t.Errorf("IPv6 映射不符: full=%q short=%q", full, short)
	}
}

func TestToolMetadata(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != "网络" || meta.Icon != "Binary" {
		t.Errorf("元数据不一致: %+v", meta)
	}
	if meta.Name == "" || meta.Description == "" || len(meta.Keywords) == 0 {
		t.Errorf("元数据缺字段: %+v", meta)
	}
}
