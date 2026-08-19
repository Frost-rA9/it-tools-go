package ipv6ula

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

// 固定时间戳，便于确定性断言。
const fixedTS = int64(1700000000000)

// 已知向量：SHA1("170000000000020:37:06:12:34:56") = 2d9080a90af4be113e109924d4dd2fb729eb4229，
// 低 40 bits = b729eb4229 → 前缀 fdb7:29eb:4229。
var knownOutput = output{
	ULA:                "fdb7:29eb:4229::/48",
	FirstRoutableBlock: "fdb7:29eb:4229:0::/64",
	LastRoutableBlock:  "fdb7:29eb:4229:ffff::/64",
}

func TestGenerateKnownVector(t *testing.T) {
	got, err := generate("20:37:06:12:34:56", fixedTS)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	if got != knownOutput {
		t.Errorf("已知向量不符:\n  期望 %+v\n  实际 %+v", knownOutput, got)
	}
}

func TestGenerateFormat(t *testing.T) {
	ulaRe := regexp.MustCompile(`^fd[0-9a-f]{2}:[0-9a-f]{4}:[0-9a-f]{4}::/48$`)
	blockRe := regexp.MustCompile(`^fd[0-9a-f]{2}:[0-9a-f]{4}:[0-9a-f]{4}:(0|ffff)::/64$`)

	cases := []struct {
		name string
		mac  string
	}{
		{"冒号分隔", "20:37:06:12:34:56"},
		{"短横线分隔", "20-37-06-12-34-56"},
		{"大写", "20:37:06:12:34:56"},
		{"小写", "20:37:06:12:34:56"},
		{"三字节", "00:11:22"},
		{"六字节", "00:11:22:33:44:55"},
	}
	for _, c := range cases {
		out, err := generate(c.mac, fixedTS)
		if err != nil {
			t.Errorf("%s: 生成失败: %v", c.name, err)
			continue
		}
		if !ulaRe.MatchString(out.ULA) {
			t.Errorf("%s: ULA 格式不正确: %q", c.name, out.ULA)
		}
		for _, block := range []string{out.FirstRoutableBlock, out.LastRoutableBlock} {
			if !blockRe.MatchString(block) {
				t.Errorf("%s: /64 块格式不正确: %q", c.name, block)
			}
		}
		// 三项应共享同一 /48 前缀。
		prefix := out.ULA[:len(out.ULA)-len("::/48")]
		if !strings.HasPrefix(out.FirstRoutableBlock, prefix+":") || !strings.HasPrefix(out.LastRoutableBlock, prefix+":") {
			t.Errorf("%s: 三项前缀不一致: %q / %q / %q", c.name, out.ULA, out.FirstRoutableBlock, out.LastRoutableBlock)
		}
	}
}

func TestGenerateDeterministic(t *testing.T) {
	a, err := generate("20:37:06:12:34:56", fixedTS)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := generate("20:37:06:12:34:56", fixedTS)
	if a != b {
		t.Errorf("同 ts + 同 MAC 应确定性生成相同结果:\n  %+v\n  %+v", a, b)
	}
}

func TestGenerateDiffersAcrossTimestamps(t *testing.T) {
	a, err := generate("20:37:06:12:34:56", fixedTS)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := generate("20:37:06:12:34:56", fixedTS+1)
	if a.ULA == b.ULA {
		t.Errorf("不同 ts 应生成不同 ULA（碰撞概率 2^-40）: %q", a.ULA)
	}
}

func TestGenerateErrors(t *testing.T) {
	cases := []struct {
		name string
		mac  string
	}{
		{"空 MAC", ""},
		{"缺组", "20:37"},
		{"半字节", "20:3"},
		{"单字节组", "2:37:06:12:34:56"},
		{"非法字符", "zz:37:06:12:34:56"},
		{"多余分隔符", "20::37:06:12:34:56"},
		{"多余字节", "20:37:06:12:34:56:78"},
		{"纯数字", "203706123456"},
	}
	for _, c := range cases {
		if _, err := generate(c.mac, fixedTS); err == nil {
			t.Errorf("%s: 期望报错，实际成功（MAC=%q）", c.name, c.mac)
		}
	}
}

func TestExecuteEndToEnd(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{MACAddress: "20:37:06:12:34:56"})
	outJSON, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("输出解析失败: %v", err)
	}
	if !regexp.MustCompile(`^fd[0-9a-f]{2}:[0-9a-f]{4}:[0-9a-f]{4}::/48$`).MatchString(out.ULA) {
		t.Errorf("Execute 输出格式不正确: %q", out.ULA)
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name  string
		input string
	}{
		{"非法 JSON", `{bad`},
		{"空 MAC", `{"macAddress":""}`},
		{"非法 MAC", `{"macAddress":"not-a-mac"}`},
	}
	for _, c := range cases {
		if _, err := exec.Execute(t.Context(), c.input); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}

func TestToolMetadata(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != "网络" || meta.Icon != "BuildingFactory" {
		t.Errorf("元数据不一致: %+v", meta)
	}
	if meta.Name == "" || meta.Description == "" || len(meta.Keywords) == 0 {
		t.Errorf("元数据缺字段: %+v", meta)
	}
}
