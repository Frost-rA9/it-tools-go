package macaddressgen

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

var macRe = regexp.MustCompile(`^(([0-9a-fA-F]{2})[-:.]?){5}[0-9a-fA-F]{2}$`)

// segments 返回 MAC 中全部字节 hex 段。
func segments(mac string) []string {
	return regexp.MustCompile(`[0-9a-fA-F]{2}`).FindAllString(mac, -1)
}

func TestGenerateKnownPrefix(t *testing.T) {
	macs, err := generate(1, "64:16:7F", ":", "upper")
	if err != nil {
		t.Fatalf("generate 失败: %v", err)
	}
	if len(macs) != 1 {
		t.Fatalf("期望 1 条，实际 %d", len(macs))
	}
	got := macs[0]
	if !macRe.MatchString(got) {
		t.Errorf("格式不正确: %q", got)
	}
	if len(segments(got)) != 6 {
		t.Errorf("应为 6 字节: %q", got)
	}
	if !regexp.MustCompile(`^64:16:7F:`).MatchString(got) {
		t.Errorf("应保留前缀 64:16:7F: %q", got)
	}
}

func TestGenerateFullPrefixDeterministic(t *testing.T) {
	// 6 字节前缀 → 完全确定，count 1。
	macs, err := generate(1, "64:16:7F:AA:BB:CC", ":", "upper")
	if err != nil {
		t.Fatalf("generate 失败: %v", err)
	}
	if macs[0] != "64:16:7F:AA:BB:CC" {
		t.Errorf("6 字节前缀应完全确定: %q", macs[0])
	}
}

func TestSplitPrefix(t *testing.T) {
	cases := []struct {
		prefix string
		want   []string
	}{
		{"64:16:7F", []string{"64", "16", "7F"}},
		{"64167F", []string{"64", "16", "7F"}},
		{"64-16-7F", []string{"64", "16", "7F"}},
		{"64.16.7F", []string{"64", "16", "7F"}},
		{"64:16-7F", []string{"64", "16", "7F"}}, // 混用分隔符
		{"6416F", []string{"64", "16", "0F"}},    // 奇数 hex 末段补零
		{"6", []string{"06"}},
		{"", nil},
	}
	for _, c := range cases {
		got, err := splitPrefix(c.prefix)
		if err != nil {
			t.Errorf("prefix=%q: 拆分失败: %v", c.prefix, err)
			continue
		}
		if len(got) != len(c.want) {
			t.Errorf("prefix=%q: 期望 %v 实际 %v", c.prefix, c.want, got)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("prefix=%q: 段[%d] 期望 %s 实际 %s", c.prefix, i, c.want[i], got[i])
			}
		}
	}
}

func TestCaseAndSeparator(t *testing.T) {
	// lower + 短横线。
	macs, err := generate(1, "64:16:7F", "-", "lower")
	if err != nil {
		t.Fatal(err)
	}
	if macs[0] != strings.ToLower(macs[0]) {
		t.Errorf("lower 应全小写: %q", macs[0])
	}
	if len(segments(macs[0])) != 6 {
		t.Errorf("分隔符 - 解析错误: %q", macs[0])
	}

	// upper + 点。
	macs, _ = generate(1, "64:16:7F", ".", "upper")
	if macs[0] != strings.ToUpper(macs[0]) {
		t.Errorf("upper 应全大写: %q", macs[0])
	}
	if len(segments(macs[0])) != 6 {
		t.Errorf("分隔符 . 解析错误: %q", macs[0])
	}

	// 无分隔符 → 12 位连续 hex。
	macs, _ = generate(1, "64167F", "", "upper")
	if len(macs[0]) != 12 {
		t.Errorf("无分隔符应 12 hex: %q", macs[0])
	}
}

func TestRandomNoPrefixFormat(t *testing.T) {
	macs, err := generate(10, "", ":", "upper")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, m := range macs {
		if !macRe.MatchString(m) {
			t.Errorf("随机 MAC 格式不正确: %q", m)
		}
		if seen[m] {
			t.Errorf("出现重复: %q", m)
		}
		seen[m] = true
	}
}

func TestGenerateErrors(t *testing.T) {
	cases := []struct {
		name                string
		count               int
		prefix, sep, csMode string
	}{
		{"数量 0", 0, "", ":", "upper"},
		{"数量超限", 101, "", ":", "upper"},
		{"前缀 7 字节", 1, "64:16:7F:AA:BB:CC:DD", ":", "upper"},
		{"前缀非法字符", 1, "zz:16:7F", ":", "upper"},
		{"前缀段超长", 1, "6416:7F", ":", "upper"},
		{"分隔符非法", 1, "", "_", "upper"},
		{"大小写非法", 1, "", ":", "mixed"},
	}
	for _, c := range cases {
		if _, err := generate(c.count, c.prefix, c.sep, c.csMode); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}

func TestExecuteEndToEnd(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Count: 2, Prefix: "64:16:7F", Separator: ":", Case: "upper"})
	outJSON, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("输出解析失败: %v", err)
	}
	lines := strings.Split(out.MACAddresses, "\n")
	if len(lines) != 2 {
		t.Fatalf("期望 2 行，实际 %d", len(lines))
	}
	for _, l := range lines {
		if !regexp.MustCompile(`^64:16:7F:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$`).MatchString(l) {
			t.Errorf("端到端输出不符: %q", l)
		}
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	for _, in := range []string{
		`{bad`,
		`{"count":0}`,
		`{"prefix":"zz:11:22"}`,
		`{"prefix":"aa:bb:cc:dd:ee:ff:00"}`,
	} {
		if _, err := exec.Execute(t.Context(), in); err == nil {
			t.Errorf("输入 %q: 期望报错", in)
		}
	}
}

func TestToolMetadata(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != "网络" || meta.Icon != "Devices" {
		t.Errorf("元数据不一致: %+v", meta)
	}
	if meta.Name == "" || meta.Description == "" || len(meta.Keywords) == 0 {
		t.Errorf("元数据缺字段: %+v", meta)
	}
}
