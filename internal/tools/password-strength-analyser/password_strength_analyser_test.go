package pwdstrength

import (
	"encoding/json"
	"math"
	"testing"
)

func TestExecuteVectors(t *testing.T) {
	// JS 参考结果（与 it-tools service.ts 完全一致的移植）。
	cases := []struct {
		password string
		duration string
		entropy  float64
		charset  int
		length   int
		score    float64
	}{
		{"", "Instantly", 0, 0, 0, 0},
		{"abc", "Instantly", 14.101319154423276, 26, 3, 0.11016655589393184},
		{"Abc123!", "18 hours, 47 seconds", 45.882121961743465, 94, 7, 0.3584540778261208},
		{"correct horse battery staple", "7.54e+29 millennia, 3 centuries", 164.02346786357202, 58, 28, 1},
		{"P@ssw0rd1234567890!", "978,639,208,454,308,400 millennia, 8 centuries", 124.53718818187511, 94, 19, 0.9729467826708993},
	}

	for _, c := range cases {
		raw, _ := json.Marshal(input{Password: c.password})
		out, err := (Executor{}).Execute(t.Context(), string(raw))
		if err != nil {
			t.Fatalf("%q: %v", c.password, err)
		}
		var o output
		if err := json.Unmarshal([]byte(out), &o); err != nil {
			t.Fatal(err)
		}
		if o.CrackDuration != c.duration {
			t.Errorf("%q: 时长 %q != %q", c.password, o.CrackDuration, c.duration)
		}
		if o.Entropy != c.entropy {
			t.Errorf("%q: 熵 %v != %v", c.password, o.Entropy, c.entropy)
		}
		if o.CharsetLength != c.charset || o.PasswordLength != c.length {
			t.Errorf("%q: charset/len %d/%d != %d/%d", c.password, o.CharsetLength, o.PasswordLength, c.charset, c.length)
		}
		if o.Score != c.score {
			t.Errorf("%q: 得分 %v != %v", c.password, o.Score, c.score)
		}
	}
}

func TestCharsetLength(t *testing.T) {
	cases := []struct {
		password string
		want     int
	}{
		{"", 0},
		{"abcdef", 26},
		{"ABCDEF", 26},
		{"123456", 10},
		{"!@#$", 32},
		{"aA1!", 94},
		{"a_b", 58}, // 下划线算特殊字符（\W|_）
	}
	for _, c := range cases {
		if got := charsetLength(c.password); got != c.want {
			t.Errorf("%q: charset=%d != %d", c.password, got, c.want)
		}
	}
}

func TestPrettifyEdgeCases(t *testing.T) {
	if got := prettifyQuantity(math.Inf(1)); got != "Infinity" {
		t.Errorf("Infinity: %q", got)
	}
	if got := groupThousands(978639208454308400); got != "978,639,208,454,308,400" {
		t.Errorf("groupThousands: %q", got)
	}
}
