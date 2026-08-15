// Package bip39gen 实现 BIP39 助记词生成工具（对齐 it-tools bip39-generator）。
package bip39gen

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39/wordlists"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "bip39-generator"
	Name        = "BIP39 密码生成器"
	Description = "从随机或指定的熵生成 BIP39 助记词，或从助记词还原熵"
	Category    = registry.CategoryCrypto
	Icon        = "ListNumbers"
)

// Keywords 为搜索关键词。
var Keywords = []string{"BIP39", "passphrase", "generator", "mnemonic", "entropy", "助记词", "熵"}

// 转换模式。
const (
	ModeGenerate          = "generate"
	ModeEntropyToMnemonic = "entropy-to-mnemonic"
	ModeMnemonicToEntropy = "mnemonic-to-entropy"
)

// defaultEntropyBits 为随机熵的位数。
const defaultEntropyBits = 128

// languages 语言标签 → 字词表。
var languages = map[string][]string{
	"English":             wordlists.English,
	"Chinese simplified":  wordlists.ChineseSimplified,
	"Chinese traditional": wordlists.ChineseTraditional,
	"Czech":               wordlists.Czech,
	"French":              wordlists.French,
	"Italian":             wordlists.Italian,
	"Japanese":            wordlists.Japanese,
	"Korean":              wordlists.Korean,
	"Portuguese":          portugueseWordlist,
	"Spanish":             wordlists.Spanish,
}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Mode     string `json:"mode"`     // generate | entropy-to-mnemonic | mnemonic-to-entropy
	Language string `json:"language"` // 语言标签
	Text     string `json:"text"`     // 熵（hex）或助记词
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理 BIP39 转换并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	words, ok := languages[in.Language]
	if !ok {
		return "", fmt.Errorf("未知语言: %q", in.Language)
	}

	var result string
	switch in.Mode {
	case ModeGenerate:
		entropy, err := generateEntropy(defaultEntropyBits)
		if err != nil {
			return "", err
		}
		result = hex.EncodeToString(entropy)
	case ModeEntropyToMnemonic:
		entropy, err := hex.DecodeString(strings.TrimSpace(in.Text))
		if err != nil {
			return "", fmt.Errorf("熵不是有效的十六进制字符串: %w", err)
		}
		mnemonic, err := entropyToMnemonic(entropy, words)
		if err != nil {
			return "", err
		}
		result = mnemonic
	case ModeMnemonicToEntropy:
		entropy, err := mnemonicToEntropy(strings.TrimSpace(in.Text), words)
		if err != nil {
			return "", err
		}
		result = hex.EncodeToString(entropy)
	default:
		return "", fmt.Errorf("未知模式: %q（仅支持 generate/entropy-to-mnemonic/mnemonic-to-entropy）", in.Mode)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// generateEntropy 生成 bitSize 位随机熵（128..256，32 的倍数）。
func generateEntropy(bitSize int) ([]byte, error) {
	if bitSize < 128 || bitSize > 256 || bitSize%32 != 0 {
		return nil, fmt.Errorf("熵位数必须在 128..256 且为 32 的倍数（当前 %d）", bitSize)
	}
	entropy := make([]byte, bitSize/8)
	if _, err := rand.Read(entropy); err != nil {
		return nil, fmt.Errorf("生成随机熵失败: %w", err)
	}
	return entropy, nil
}

// entropyToMnemonic 将熵编码为助记词（BIP39：SHA256 校验和 + 11 位分组索引字词表）。
func entropyToMnemonic(entropy []byte, words []string) (string, error) {
	entBits := len(entropy) * 8
	if entBits < 128 || entBits > 256 || entBits%32 != 0 {
		return "", fmt.Errorf("熵位数必须在 128..256 且为 32 的倍数（当前 %d）", entBits)
	}
	if len(words) != 2048 {
		return "", fmt.Errorf("字词表长度必须为 2048（当前 %d）", len(words))
	}

	checksum := sha256.Sum256(entropy)
	checkBits := entBits / 32
	totalBits := entBits + checkBits

	// 位流 = 熵 + 校验和的前 checkBits 位。
	bits := make([]bool, 0, totalBits)
	for _, b := range entropy {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (b>>i)&1 == 1)
		}
	}
	for i := 0; i < checkBits; i++ {
		bits = append(bits, (checksum[i/8]>>(7-i%8))&1 == 1)
	}

	parts := make([]string, 0, totalBits/11)
	for i := 0; i < totalBits; i += 11 {
		var idx int
		for j := 0; j < 11; j++ {
			if bits[i+j] {
				idx |= 1 << (10 - j)
			}
		}
		parts = append(parts, words[idx])
	}
	return strings.Join(parts, " "), nil
}

// mnemonicToEntropy 将助记词解码回熵并校验校验和。
func mnemonicToEntropy(mnemonic string, words []string) ([]byte, error) {
	parts := strings.Fields(mnemonic)
	n := len(parts)
	if n < 12 || n > 24 || n%3 != 0 {
		return nil, fmt.Errorf("助记词词数必须为 12..24 且为 3 的倍数（当前 %d）", n)
	}
	if len(words) != 2048 {
		return nil, fmt.Errorf("字词表长度必须为 2048（当前 %d）", len(words))
	}

	index := make(map[string]int, len(words))
	for i, w := range words {
		index[w] = i
	}

	entBits := n * 32 / 3 // 词数 * 11 - 校验位后熵位数；实际 = 32 * (n/3) * ... 用标准公式
	checkBits := entBits / 32
	totalBits := entBits + checkBits

	// 位流。
	bits := make([]bool, 0, totalBits)
	for _, p := range parts {
		idx, ok := index[p]
		if !ok {
			return nil, fmt.Errorf("助记词包含不在字词表中的词: %q", p)
		}
		for i := 10; i >= 0; i-- {
			bits = append(bits, (idx>>i)&1 == 1)
		}
	}

	entropyBytes := make([]byte, entBits/8)
	for i := 0; i < entBits; i++ {
		if bits[i] {
			entropyBytes[i/8] |= 0x80 >> (i % 8)
		}
	}

	checksum := sha256.Sum256(entropyBytes)
	wantCheck := 0
	for i := 0; i < checkBits; i++ {
		if (checksum[i/8]>>(7-i%8))&1 == 1 {
			wantCheck |= 1 << (checkBits - 1 - i)
		}
	}
	gotCheck := 0
	for i := 0; i < checkBits; i++ {
		if bits[entBits+i] {
			gotCheck |= 1 << (checkBits - 1 - i)
		}
	}
	if gotCheck != wantCheck {
		return nil, fmt.Errorf("校验和不匹配（助记词无效）")
	}
	return entropyBytes, nil
}
