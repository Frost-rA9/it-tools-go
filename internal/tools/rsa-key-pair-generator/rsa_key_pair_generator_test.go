package rsakeypair

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"
)

func TestExecuteGenerate(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Bits: 1024})
	out, err := exec.Execute(t.Context(), string(raw))
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}

	// 公钥 PEM 应为 PKCS#1 RSA PUBLIC KEY。
	if !strings.Contains(o.PublicKeyPem, "-----BEGIN RSA PUBLIC KEY-----") {
		t.Errorf("公钥头不正确:\n%s", o.PublicKeyPem)
	}
	if !strings.Contains(o.PrivateKeyPem, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Errorf("私钥头不正确:\n%s", o.PrivateKeyPem)
	}

	// 私钥可解析为 PKCS#1。
	privBlock, _ := pem.Decode([]byte(o.PrivateKeyPem))
	if privBlock == nil {
		t.Fatal("私钥 PEM 解析失败")
	}
	priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		t.Fatalf("私钥解析失败: %v", err)
	}
	if priv.N.BitLen() != 1024 {
		t.Errorf("密钥位数 %d != 1024", priv.N.BitLen())
	}

	// 公钥可解析且与私钥匹配。
	pubBlock, _ := pem.Decode([]byte(o.PublicKeyPem))
	if pubBlock == nil {
		t.Fatal("公钥 PEM 解析失败")
	}
	pub, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		t.Fatalf("公钥解析失败: %v", err)
	}
	if pub.N.Cmp(priv.N) != 0 {
		t.Error("公钥与私钥不匹配")
	}
}

func TestSignVerifyRoundTrip(t *testing.T) {
	exec := Executor{}
	raw, _ := json.Marshal(input{Bits: 1024})
	out, _ := exec.Execute(t.Context(), string(raw))
	var o output
	_ = json.Unmarshal([]byte(out), &o)

	privBlock, _ := pem.Decode([]byte(o.PrivateKeyPem))
	priv, _ := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	sum := sha256.Sum256([]byte("hello world"))
	digest := sum[:]

	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}
	if err := rsa.VerifyPKCS1v15(&priv.PublicKey, crypto.SHA256, digest, sig); err != nil {
		t.Fatalf("验签失败: %v", err)
	}
}

func TestBitsValidation(t *testing.T) {
	exec := Executor{}
	for _, bits := range []int{0, 128, 255, 1023, 16392, 2052} {
		raw, _ := json.Marshal(input{Bits: bits})
		if _, err := exec.Execute(t.Context(), string(raw)); err == nil {
			t.Errorf("bits=%d 应报错", bits)
		}
	}
}
