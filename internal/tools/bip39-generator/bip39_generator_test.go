package bip39gen

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
)

func TestEntropyToMnemonicOfficialVector(t *testing.T) {
	// BIP39/Trezor 官方向量：全零 128 位熵。
	entropy, _ := hex.DecodeString("00000000000000000000000000000000")
	got, err := entropyToMnemonic(entropy, languages["English"])
	if err != nil {
		t.Fatalf("编码失败: %v", err)
	}
	want := strings.Repeat("abandon ", 11) + "about"
	if got != want {
		t.Errorf("官方向量不符:\n got %q\nwant %q", got, want)
	}
}

func TestMnemonicToEntropyOfficialVector(t *testing.T) {
	mnemonic := strings.Repeat("abandon ", 11) + "about"
	got, err := mnemonicToEntropy(mnemonic, languages["English"])
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}
	want, _ := hex.DecodeString("00000000000000000000000000000000")
	if hex.EncodeToString(got) != hex.EncodeToString(want) {
		t.Errorf("官方向量不符: got %x", got)
	}
}

func TestRoundTripAllLanguages(t *testing.T) {
	entropy, err := hex.DecodeString("7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f")
	if err != nil {
		t.Fatal(err)
	}
	for lang := range languages {
		mnemonic, err := entropyToMnemonic(entropy, languages[lang])
		if err != nil {
			t.Fatalf("%s 编码失败: %v", lang, err)
		}
		back, err := mnemonicToEntropy(mnemonic, languages[lang])
		if err != nil {
			t.Fatalf("%s 解码失败: %v", lang, err)
		}
		if hex.EncodeToString(back) != hex.EncodeToString(entropy) {
			t.Errorf("%s 回环失败: %x", lang, back)
		}
	}
}

func TestExecuteGenerate(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Mode: ModeGenerate, Language: "English"})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("generate 失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	if len(o.Result) != 32 || !isHex(o.Result) {
		t.Errorf("generate 应返回 32 位十六进制熵: %q", o.Result)
	}
}

func TestExecuteEntropyToMnemonic(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Mode: ModeEntropyToMnemonic, Language: "Chinese simplified", Text: "00000000000000000000000000000000"})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("entropy-to-mnemonic 失败: %v", err)
	}
	var o output
	_ = json.Unmarshal([]byte(out), &o)
	if len(strings.Fields(o.Result)) != 12 {
		t.Errorf("应生成 12 个助记词: %q", o.Result)
	}
}

func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	cases := []struct {
		name  string
		input input
	}{
		{"未知语言", input{Mode: ModeGenerate, Language: "Klingon"}},
		{"熵非法 hex", input{Mode: ModeEntropyToMnemonic, Language: "English", Text: "zz"}},
		{"熵位数不足", input{Mode: ModeEntropyToMnemonic, Language: "English", Text: "aabb"}},
		{"助记词校验失败", input{Mode: ModeMnemonicToEntropy, Language: "English", Text: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon"}},
		{"助记词含未知词", input{Mode: ModeMnemonicToEntropy, Language: "English", Text: "zzzz abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon"}},
		{"未知模式", input{Mode: "bad", Language: "English"}},
	}
	for _, c := range cases {
		raw, _ := json.Marshal(c.input)
		if _, err := exec.Execute(t.Context(), string(raw)); err == nil {
			t.Errorf("%s: 期望报错", c.name)
		}
	}
}

func isHex(s string) bool {
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return false
		}
	}
	return true
}
