package encryption

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

// cryptoJSVectors 为 crypto-js 生成的加密参考向量。
type cryptoJSVectors struct {
	Msg    string            `json:"msg"`
	Pass   string            `json:"pass"`
	Salt   string            `json:"saltHex"`
	Vector map[string]string `json:"vectors"`
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

func mustSalt(t *testing.T, hexStr string) []byte {
	t.Helper()
	s, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// TestEncryptAgainstCryptoJS 用固定盐校验加密输出与 crypto-js 完全一致。
func TestEncryptAgainstCryptoJS(t *testing.T) {
	v := loadVectors(t)
	salt := mustSalt(t, v.Salt)
	for algo, want := range v.Vector {
		spec := algos[algo]
		got, err := encryptWithSalt([]byte(v.Msg), []byte(v.Pass), algo, spec, salt)
		if err != nil {
			t.Fatalf("%s 加密失败: %v", algo, err)
		}
		if got != want {
			t.Errorf("%s:\n got %s\nwant %s", algo, got, want)
		}
	}
}

// TestDecryptCryptoJSVectors 用 crypto-js 生成的密文解密并比对明文。
func TestDecryptCryptoJSVectors(t *testing.T) {
	v := loadVectors(t)
	for algo, ct := range v.Vector {
		spec := algos[algo]
		got, err := decrypt(ct, v.Pass, algo, spec)
		if err != nil {
			t.Fatalf("%s 解密失败: %v", algo, err)
		}
		if got != v.Msg {
			t.Errorf("%s 解密结果: %q != %q", algo, got, v.Msg)
		}
	}
}

// TestRoundTrip 校验各算法加密→解密回环，覆盖边界长度与多字节文本。
func TestRoundTrip(t *testing.T) {
	messages := []string{
		"",
		"a",
		"1234567890123456",        // 16 字节
		"123456789012345678901234", // 24 字节
		"Hello, 世界! 你好世界 abcdefghijklmnopqrstuvwxyz",
	}
	for algo := range algos {
		spec := algos[algo]
		for _, msg := range messages {
			ct, err := encrypt(msg, "my secret key", algo, spec)
			if err != nil {
				t.Fatalf("%s 加密失败: %v", algo, err)
			}
			if ct == "" {
				t.Fatalf("%s 加密结果为空", algo)
			}
			pt, err := decrypt(ct, "my secret key", algo, spec)
			if err != nil {
				t.Fatalf("%s 解密失败: %v", algo, err)
			}
			if pt != msg {
				t.Errorf("%s 回环失败: %q != %q", algo, pt, msg)
			}
		}
	}
}

// TestDecryptErrors 校验错误用例。
func TestDecryptErrors(t *testing.T) {
	v := loadVectors(t)
	spec := algos["AES"]
	cases := []struct {
		name string
		ct   string
	}{
		{"无效 Base64", "!!!not-base64!!!"},
		{"缺少 Salted__", "QUJDREVGRw=="},
		{"密文过短", "U2FsdGVkX1"},
		{"密钥错误", v.Vector["AES"]},
	}
	for _, c := range cases {
		secret := "my secret key"
		if c.name == "密钥错误" {
			secret = "wrong key"
		}
		if _, err := decrypt(c.ct, secret, "AES", spec); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}

// TestExecute 校验 Execute 入口的完整协议。
func TestExecute(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Mode: ModeEncrypt, Text: "hi", Secret: "secret", Algo: "AES"})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if o.Result == "" {
		t.Fatal("加密结果为空")
	}

	raw, _ = json.Marshal(input{Mode: ModeDecrypt, Text: o.Result, Secret: "secret", Algo: "AES"})
	out, err = exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if o.Result != "hi" {
		t.Errorf("解密结果: %q", o.Result)
	}

	if _, err := exec.Execute(t.Context(), `{"mode":"bad","text":"","secret":"","algo":"AES"}`); err == nil {
		t.Error("未知模式应报错")
	}
	if _, err := exec.Execute(t.Context(), `{"mode":"encrypt","text":"hi","secret":"s","algo":"X"}`); err == nil {
		t.Error("未知算法应报错")
	}
}
