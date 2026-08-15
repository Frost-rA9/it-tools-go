package bcrypt

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Mode: ModeHash, Text: "secret", Cost: 10})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("hash 失败: %v", err)
	}
	var o hashOutput
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^\$2a\$10\$[./A-Za-z0-9]{53}$`).MatchString(o.Result) {
		t.Errorf("哈希格式不正确: %q", o.Result)
	}
}

func TestHashRandomSalt(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Mode: ModeHash, Text: "secret", Cost: 10})
	out1, _ := exec.Execute(t.Context(), string(raw))
	out2, _ := exec.Execute(t.Context(), string(raw))
	var h1, h2 hashOutput
	_ = json.Unmarshal([]byte(out1), &h1)
	_ = json.Unmarshal([]byte(out2), &h2)
	if h1.Result == h2.Result {
		t.Errorf("两次哈希应因随机盐而不同")
	}
}

func TestCompare(t *testing.T) {
	exec := Executor{}

	// 先生成一个有效哈希
	raw, _ := json.Marshal(input{Mode: ModeHash, Text: "correct horse", Cost: 10})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatal(err)
	}
	var ho hashOutput
	_ = json.Unmarshal([]byte(out), &ho)

	cases := []struct {
		name string
		text string
		hash string
		want bool
	}{
		{"匹配", "correct horse", ho.Result, true},
		{"不匹配", "wrong", ho.Result, false},
		{"非法哈希", "correct horse", "not-a-hash", false},
		{"空哈希", "correct horse", "", false},
	}
	for _, c := range cases {
		raw, _ := json.Marshal(input{Mode: ModeCompare, Text: c.text, Hash: c.hash})
		out, err := exec.Execute(t.Context(), string(raw))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		var co compareOutput
		if err := json.Unmarshal([]byte(out), &co); err != nil {
			t.Fatal(err)
		}
		if co.Match != c.want {
			t.Errorf("%s: match=%v, want %v", c.name, co.Match, c.want)
		}
	}
}

func TestHashCostValidation(t *testing.T) {
	exec := Executor{}
	for _, cost := range []int{3, 32, 100} {
		raw, _ := json.Marshal(input{Mode: ModeHash, Text: "x", Cost: cost})
		if _, err := exec.Execute(t.Context(), string(raw)); err == nil {
			t.Errorf("cost=%d 应报错", cost)
		}
	}

	// 不传 cost 时用默认值 10
	raw, _ := json.Marshal(input{Mode: ModeHash, Text: "x"})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("默认 cost 失败: %v", err)
	}
	var o hashOutput
	_ = json.Unmarshal([]byte(out), &o)
	if !strings.HasPrefix(o.Result, "$2a$10$") {
		t.Errorf("默认 cost 应生成 $2a$10$ 哈希: %q", o.Result)
	}
}

func TestInvalidMode(t *testing.T) {
	exec := Executor{}
	if _, err := exec.Execute(t.Context(), `{"mode":"bad","text":"x"}`); err == nil {
		t.Error("未知模式应报错")
	}
	if _, err := exec.Execute(t.Context(), `{`); err == nil {
		t.Error("非法 JSON 应报错")
	}
}
