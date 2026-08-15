// Package hmacgen 实现 HMAC 生成工具（对齐 it-tools hmac-generator）。
package hmacgen

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "hmac-generator"
	Name        = "HMAC 生成器"
	Description = "使用密钥与哈希函数计算基于哈希的消息认证码（HMAC）"
	Category    = registry.CategoryCrypto
	Icon        = "Key"
)

// Keywords 为搜索关键词。
var Keywords = []string{"hmac", "generator", "MD5", "SHA1", "SHA256", "SHA224", "SHA512", "SHA384", "SHA3", "RIPEMD160", "密钥", "认证码"}

// algoOrder 与 it-tools 的展示顺序一致。
var algoOrder = []string{"MD5", "RIPEMD160", "SHA1", "SHA3", "SHA224", "SHA256", "SHA384", "SHA512"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text     string `json:"text"`     // 待计算文本
	Secret   string `json:"secret"`   // 密钥
	Algo     string `json:"algo"`     // 哈希函数
	Encoding string `json:"encoding"` // Hex | Base64 | Base64url | Bin
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 计算 HMAC 并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	newHash, err := hashFunc(in.Algo)
	if err != nil {
		return "", err
	}

	mac := hmac.New(newHash, []byte(in.Secret))
	mac.Write([]byte(in.Text))
	digest := mac.Sum(nil)

	encoded, err := encode(in.Encoding, digest)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Result: encoded})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// hashFunc 返回指定算法的哈希构造函数。
func hashFunc(algo string) (func() hash.Hash, error) {
	switch algo {
	case "MD5":
		return md5.New, nil
	case "RIPEMD160":
		return ripemd160.New, nil
	case "SHA1":
		return sha1.New, nil
	case "SHA3":
		// crypto-js 的 SHA3 为 legacy Keccak（0x01 填充），与 hash-text 一致。
		return sha3.NewLegacyKeccak512, nil
	case "SHA224":
		return sha256.New224, nil
	case "SHA256":
		return sha256.New, nil
	case "SHA384":
		return sha512.New384, nil
	case "SHA512":
		return sha512.New, nil
	default:
		return nil, fmt.Errorf("未知哈希函数: %q", algo)
	}
}

// encode 按编码输出摘要字符串。
func encode(enc string, d []byte) (string, error) {
	switch enc {
	case "Hex":
		return hex.EncodeToString(d), nil
	case "Base64":
		return base64.StdEncoding.EncodeToString(d), nil
	case "Base64url":
		return base64.RawURLEncoding.EncodeToString(d), nil
	case "Bin":
		return hexToBin(d), nil
	default:
		return "", fmt.Errorf("未知摘要编码: %q（仅支持 Hex/Base64/Base64url/Bin）", enc)
	}
}

// hexToBin 对齐 crypto-js convertHexToBin：每个十六进制字符 → 4 位二进制。
func hexToBin(d []byte) string {
	var b strings.Builder
	b.Grow(len(d) * 8)
	for _, by := range d {
		b.WriteString(nibbleBin[by>>4])
		b.WriteString(nibbleBin[by&0x0f])
	}
	return b.String()
}

// nibbleBin 为十六进制半字节的 4 位二进制表示。
var nibbleBin = [16]string{"0000", "0001", "0010", "0011", "0100", "0101", "0110", "0111", "1000", "1001", "1010", "1011", "1100", "1101", "1110", "1111"}
