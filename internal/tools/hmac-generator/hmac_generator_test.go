package hmacgen

import (
	"encoding/json"
	"os"
	"testing"
)

type cryptoJSVectors struct {
	Text    string                       `json:"text"`
	Secret  string                       `json:"secret"`
	Vectors map[string]map[string]string `json:"vectors"`
}

func loadVectors(t *testing.T) cryptoJSVectors {
	t.Helper()
	data, err := os.ReadFile("testdata/crypto-js-vectors.json")
	if err != nil {
		t.Fatalf("读取测试向量失败: %v", err)
	}
	var v cryptoJSVectors
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("解析测试向量失败: %v", err)
	}
	return v
}

// TestExecuteAgainstCryptoJS 用 crypto-js 参考向量校验全部算法与编码。
func TestExecuteAgainstCryptoJS(t *testing.T) {
	v := loadVectors(t)
	encodings := []string{"hex", "base64", "base64url", "bin"}
	for algo, vals := range v.Vectors {
		for _, enc := range encodings {
			raw, _ := json.Marshal(input{Text: v.Text, Secret: v.Secret, Algo: algo, Encoding: toEncodingName(enc)})
			out, err := (Executor{}).Execute(t.Context(), string(raw))
			if err != nil {
				t.Fatalf("%s/%s 失败: %v", algo, enc, err)
			}
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatal(err)
			}
			if o.Result != vals[enc] {
				t.Errorf("%s/%s:\n got %s\nwant %s", algo, enc, o.Result, vals[enc])
			}
		}
	}
}

func toEncodingName(s string) string {
	switch s {
	case "base64url":
		return "Base64url"
	case "base64":
		return "Base64"
	case "bin":
		return "Bin"
	default:
		return "Hex"
	}
}

func TestRFC4231HMACSHA256(t *testing.T) {
	// RFC 4231 Test Case 1：key=0x0b×20, data="Hi There"。
	raw, _ := json.Marshal(input{
		Text:     "Hi There",
		Secret:   string([]byte{0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b}),
		Algo:     "SHA256",
		Encoding: "Hex",
	})
	out, err := (Executor{}).Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	const want = "b0344c61d8db38535ca8afceaf0bf12b881dc200c9833da726e9376c2e32cff7"
	if o.Result != want {
		t.Errorf("RFC 4231 向量不符:\n got %s\nwant %s", o.Result, want)
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name  string
		input string
	}{
		{"未知算法", `{"text":"a","secret":"b","algo":"XX","encoding":"Hex"}`},
		{"未知编码", `{"text":"a","secret":"b","algo":"SHA256","encoding":"Xx"}`},
		{"非法 JSON", `{`},
	}
	for _, c := range cases {
		if _, err := exec.Execute(t.Context(), c.input); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}
