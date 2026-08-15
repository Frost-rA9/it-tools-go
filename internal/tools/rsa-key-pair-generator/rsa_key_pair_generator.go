// Package rsakeypair 实现 RSA 密钥对生成工具（对齐 it-tools rsa-key-pair-generator）。
package rsakeypair

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "rsa-key-pair-generator"
	Name        = "RSA 密钥对生成器"
	Description = "生成新的随机 RSA 私钥与公钥 PEM 证书密钥对"
	Category    = registry.CategoryCrypto
	Icon        = "Certificate"
)

// Keywords 为搜索关键词。
var Keywords = []string{"rsa", "key", "pair", "generator", "public", "private", "secret", "ssh", "pem", "密钥对"}

// 位数边界。
const (
	bitsMin = 256
	bitsMax = 16384
)

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Bits int `json:"bits"` // 密钥位数（256..16384，8 的倍数）
}

// output 是工具的输出结构。
type output struct {
	PublicKeyPem  string `json:"public_key_pem"`
	PrivateKeyPem string `json:"private_key_pem"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 生成 RSA 密钥对并返回 PKCS#1 PEM 结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if in.Bits < bitsMin || in.Bits > bitsMax || in.Bits%8 != 0 {
		return "", fmt.Errorf("位数必须在 %d..%d 且为 8 的倍数（当前 %d）", bitsMin, bitsMax, in.Bits)
	}

	priv, err := rsa.GenerateKey(rand.Reader, in.Bits)
	if err != nil {
		return "", fmt.Errorf("生成 RSA 密钥失败: %w", err)
	}

	// PKCS#1 PEM，与 node-forge 的 publicKeyToPem/privateKeyToPem 格式一致。
	privatePem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	publicPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&priv.PublicKey)})

	out, err := json.Marshal(output{
		PublicKeyPem:  string(publicPem),
		PrivateKeyPem: string(privatePem),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
