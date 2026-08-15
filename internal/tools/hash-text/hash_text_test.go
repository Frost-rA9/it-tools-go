package hashtext

import (
	"encoding/json"
	"os"
	"testing"
)

// cryptoJSVectors 从 crypto-js 生成的参考向量（testdata/crypto-js-vectors.json）。
func loadVectors(t *testing.T) map[string]map[string]map[string]string {
	t.Helper()
	data, err := os.ReadFile("testdata/crypto-js-vectors.json")
	if err != nil {
		t.Fatalf("读取测试向量失败: %v", err)
	}
	var vectors map[string]map[string]map[string]string
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatalf("解析测试向量失败: %v", err)
	}
	return vectors
}

// TestExecuteAgainstCryptoJS 用 crypto-js 参考向量校验 Execute 的完整输出。
func TestExecuteAgainstCryptoJS(t *testing.T) {
	vectors := loadVectors(t)
	exec := Executor{}

	// 仅校验 Hex 编码下全部算法（其余编码在 TestEncodeVectors 单独校验）。
	for text, algos := range vectors {
		raw, err := json.Marshal(input{Text: text, Encoding: "Hex"})
		if err != nil {
			t.Fatal(err)
		}
		out, err := exec.Execute(t.Context(), string(raw))
		if err != nil {
			t.Fatalf("Execute(%q) 失败: %v", text, err)
		}
		var o output
		if err := json.Unmarshal([]byte(out), &o); err != nil {
			t.Fatalf("解析输出失败: %v", err)
		}
		if len(o.Results) != len(algoOrder) {
			t.Fatalf("结果数量 %d != %d", len(o.Results), len(algoOrder))
		}
		for i, r := range o.Results {
			want := algos[r.Algo]["hex"]
			if r.Digest != want {
				t.Errorf("text=%q algo=%s:\n got %s\nwant %s", text, r.Algo, r.Digest, want)
			}
			if r.Algo != algoOrder[i] {
				t.Errorf("算法顺序错位: %s != %s", r.Algo, algoOrder[i])
			}
		}
	}
}

// TestEncodeVectors 校验 4 种编码与 crypto-js 一致。
func TestEncodeVectors(t *testing.T) {
	vectors := loadVectors(t)
	encodings := []string{"hex", "base64", "base64url", "bin"}
	for text, algos := range vectors {
		for algo, vals := range algos {
			d, err := digest(algo, []byte(text))
			if err != nil {
				t.Fatalf("digest(%s) 失败: %v", algo, err)
			}
			for _, enc := range encodings {
				got, err := encode(uppercase(enc), d)
				if err != nil {
					t.Fatalf("encode(%s) 失败: %v", enc, err)
				}
				if got != vals[enc] {
					t.Errorf("text=%q algo=%s enc=%s:\n got %s\nwant %s", text, algo, enc, got, vals[enc])
				}
			}
		}
	}
}

// TestExecuteErrors 校验错误输入。
func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name  string
		input string
	}{
		{"非法 JSON", `{`},
		{"未知编码", `{"text":"hello","encoding":"Xx"}`},
	}
	for _, c := range cases {
		if _, err := exec.Execute(t.Context(), c.input); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}

// uppercase 将小写编码名转为执行时的编码名。
func uppercase(s string) string {
	if s == "base64url" {
		return "Base64url"
	}
	if s == "base64" {
		return "Base64"
	}
	if s == "bin" {
		return "Bin"
	}
	return "Hex"
}
