// Package hashtext 实现文本哈希工具（对齐 it-tools hash-text）。
package hashtext

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "hash-text"
	Name        = "哈希文本"
	Description = "使用 MD5、SHA1、SHA224、SHA256、SHA384、SHA512、SHA3 或 RIPEMD160 对文本进行哈希"
	Category    = registry.CategoryCrypto
	Icon        = "EyeOff"
)

// Keywords 为搜索关键词。
var Keywords = []string{"hash", "digest", "crypto", "security", "text", "MD5", "SHA1", "SHA256", "SHA224", "SHA512", "SHA384", "SHA3", "RIPEMD160", "哈希", "摘要"}

// 默认摘要编码。
const defaultEncoding = "Hex"

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// algoOrder 与 it-tools 的展示顺序一致。
var algoOrder = []string{"MD5", "SHA1", "SHA256", "SHA224", "SHA512", "SHA384", "SHA3", "RIPEMD160"}

// input 是工具的输入结构。
type input struct {
	Text     string `json:"text"`     // 待哈希文本
	Encoding string `json:"encoding"` // Hex | Base64 | Base64url | Bin
}

// digestResult 描述单个算法的摘要结果。
type digestResult struct {
	Algo   string `json:"algo"`
	Digest string `json:"digest"`
}

// output 是工具的输出结构。
type output struct {
	Results []digestResult `json:"results"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 计算全部算法的摘要并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if in.Encoding == "" {
		in.Encoding = defaultEncoding
	}
	if !validEncoding(in.Encoding) {
		return "", fmt.Errorf("未知摘要编码: %q（仅支持 Hex/Base64/Base64url/Bin）", in.Encoding)
	}

	results := make([]digestResult, 0, len(algoOrder))
	for _, algo := range algoOrder {
		d, err := digest(algo, []byte(in.Text))
		if err != nil {
			return "", err
		}
		encoded, err := encode(in.Encoding, d)
		if err != nil {
			return "", err
		}
		results = append(results, digestResult{Algo: algo, Digest: encoded})
	}

	out, err := json.Marshal(output{Results: results})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// digest 计算指定算法的原始摘要字节。
func digest(algo string, data []byte) ([]byte, error) {
	switch algo {
	case "MD5":
		sum := md5.Sum(data)
		return sum[:], nil
	case "SHA1":
		sum := sha1.Sum(data)
		return sum[:], nil
	case "SHA224":
		sum := sha256.Sum224(data)
		return sum[:], nil
	case "SHA256":
		sum := sha256.Sum256(data)
		return sum[:], nil
	case "SHA384":
		sum := sha512.Sum384(data)
		return sum[:], nil
	case "SHA512":
		sum := sha512.Sum512(data)
		return sum[:], nil
	case "SHA3":
		// crypto-js 的 SHA3 为 legacy Keccak（0x01 填充），非 FIPS-202。
		h := sha3.NewLegacyKeccak512()
		h.Write(data)
		return h.Sum(nil), nil
	case "RIPEMD160":
		h := ripemd160.New()
		h.Write(data)
		return h.Sum(nil), nil
	default:
		return nil, fmt.Errorf("未知算法: %q", algo)
	}
}

func validEncoding(enc string) bool {
	switch enc {
	case "Hex", "Base64", "Base64url", "Bin":
		return true
	}
	return false
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
		return "", fmt.Errorf("未知摘要编码: %q", enc)
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
