package ipv4subnet

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestComputeKnownVector(t *testing.T) {
	out, err := compute("192.168.0.1/24")
	if err != nil {
		t.Fatalf("compute 失败: %v", err)
	}
	want := output{
		Netmask:          "192.168.0.0/24",
		NetworkAddress:   "192.168.0.0",
		NetworkMask:      "255.255.255.0",
		NetworkMaskBin:   "11111111.11111111.11111111.00000000",
		CIDRNotation:     "/24",
		WildcardMask:     "0.0.0.255",
		NetworkSize:      256,
		FirstAddress:     "192.168.0.1",
		LastAddress:      "192.168.0.254",
		BroadcastAddress: "192.168.0.255",
		IPClass:          "C",
		NextBlock:        "192.168.1.0/24",
		PrevBlock:        "192.167.255.0/24",
	}
	if out != want {
		t.Errorf("已知向量不符:\n  期望 %+v\n  实际 %+v", want, out)
	}
}

func TestComputeBoundaries(t *testing.T) {
	cases := []struct {
		name      string
		address   string
		network   string
		first     string
		last      string
		broadcast string
		size      uint64
	}{
		{"/8", "10.1.2.3/8", "10.0.0.0", "10.0.0.1", "10.255.255.254", "10.255.255.255", 1 << 24},
		{"/31", "192.168.0.0/31", "192.168.0.0", "192.168.0.0", "192.168.0.1", "192.168.0.1", 2},
		{"/32", "192.168.0.5/32", "192.168.0.5", "192.168.0.5", "192.168.0.5", "192.168.0.5", 1},
		{"/0", "0.0.0.0/0", "0.0.0.0", "0.0.0.1", "255.255.255.254", "255.255.255.255", uint64(1) << 32},
	}
	for _, c := range cases {
		out, err := compute(c.address)
		if err != nil {
			t.Errorf("%s: compute 失败: %v", c.name, err)
			continue
		}
		if out.NetworkAddress != c.network || out.FirstAddress != c.first ||
			out.LastAddress != c.last || out.BroadcastAddress != c.broadcast || out.NetworkSize != c.size {
			t.Errorf("%s: 不符 network=%s first=%s last=%s bc=%s size=%d（实际 %s/%s/%s/%s/%d）",
				c.name, c.network, c.first, c.last, c.broadcast, c.size,
				out.NetworkAddress, out.FirstAddress, out.LastAddress, out.BroadcastAddress, out.NetworkSize)
		}
	}
}

func TestParseForms(t *testing.T) {
	// 掩码形式与纯 IP 形式。
	if out, err := compute("192.168.0.1/255.255.255.0"); err != nil || out.CIDRNotation != "/24" {
		t.Errorf("掩码形式失败: %+v, err=%v", out, err)
	}
	if out, err := compute("192.168.0.1"); err != nil || out.CIDRNotation != "/32" {
		t.Errorf("纯 IP 应视为 /32: %+v, err=%v", out, err)
	}
	// 带空格。
	if _, err := compute(" 10.0.0.0/24 "); err != nil {
		t.Errorf("带空格应通过: %v", err)
	}
	// 非 8 位边界掩码 255.255.255.248 → /29。
	if out, err := compute("10.0.0.1/255.255.255.248"); err != nil || out.CIDRNotation != "/29" {
		t.Errorf("掩码 /29 失败: %+v, err=%v", out, err)
	}
}

func TestIPClass(t *testing.T) {
	cases := []struct{ ip, class string }{
		{"1.2.3.4", "A"},
		{"127.0.0.1", "A"},
		{"128.0.0.1", "B"},
		{"191.255.255.255", "B"},
		{"192.168.0.1", "C"},
		{"223.255.255.255", "C"},
		{"224.0.0.1", "D"},
		{"239.255.255.255", "D"},
		{"240.0.0.1", "E"},
		{"255.255.255.255", "E"},
	}
	for _, c := range cases {
		ip, err := parseIPv4(c.ip)
		if err != nil {
			t.Fatal(err)
		}
		if got := ipClass(ip); got != c.class {
			t.Errorf("ip=%s 期望 %s 实际 %s", c.ip, c.class, got)
		}
	}
}

func TestComputeErrors(t *testing.T) {
	cases := []string{
		"",
		"192.168.1.1/33",
		"192.168.1.1/-1",
		"192.168.1.1/abc",
		"192.168.1.1/255.0.255.0", // 非连续掩码
		"999.1.1.1",
		"1.2.3", // 段数不足
		"1.2.3.4.5",
		"1.2.3.x",
		"/24", // 缺 IP
	}
	for _, c := range cases {
		if _, err := compute(c); err == nil {
			t.Errorf("%q: 期望报错，实际成功", c)
		}
	}
}

func TestNeighborClamp(t *testing.T) {
	// prev 越界 clamp 到 0。
	if out, err := compute("10.0.0.1/24"); err != nil || !strings.HasPrefix(out.PrevBlock, "9.255.255.0/24") {
		// 10.0.0.0 - 256 = 9.255.255.0
		t.Errorf("prev clamp 失败: %q err=%v", out.PrevBlock, err)
	}
	// next 越界 clamp。
	if out, err := compute("255.255.255.254/32"); err != nil || out.NextBlock != "255.255.255.255/32" {
		t.Errorf("next clamp 失败: %q err=%v", out.NextBlock, err)
	}
}

func TestExecuteEndToEnd(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Address: "192.168.0.1/24"})
	outJSON, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("输出解析失败: %v", err)
	}
	if out.NetworkAddress != "192.168.0.0" || out.NetworkMask != "255.255.255.0" {
		t.Errorf("Execute 输出不符: %+v", out)
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	for _, in := range []string{`{bad`, `{"address":""}`, `{"address":"1.1.1.1/99"}`} {
		if _, err := exec.Execute(t.Context(), in); err == nil {
			t.Errorf("输入 %q: 期望报错", in)
		}
	}
}

func TestToolMetadata(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != "网络" || meta.Icon != "Router" {
		t.Errorf("元数据不一致: %+v", meta)
	}
	if meta.Name == "" || meta.Description == "" || len(meta.Keywords) == 0 {
		t.Errorf("元数据缺字段: %+v", meta)
	}
}
