// Package encryption 实现文本加密/解密工具，与 it-tools（crypto-js）的
// AES/TripleDES/Rabbit/RC4 passphrase 模式字节级兼容。
package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "encryption"
	Name        = "加密 / 解密"
	Description = "使用 AES、TripleDES、Rabbit 或 RC4 算法加密文本与解密密文"
	Category    = registry.CategoryCrypto
	Icon        = "Lock"
)

// Keywords 为搜索关键词。
var Keywords = []string{"cypher", "encipher", "text", "AES", "TripleDES", "Rabbit", "RC4", "加密", "解密"}

// saltedMagic 为 OpenSSL 格式的盐前缀。
var saltedMagic = []byte("Salted__")

// algoSpec 描述加密算法的参数。
type algoSpec struct {
	keyLen int  // 派生密钥字节数
	ivLen  int  // 派生 IV 字节数
	block  bool // true=分组密码（PKCS7），false=流密码（无填充）
}

// algos 与 crypto-js 各算法的 keySize/ivSize 对齐。
var algos = map[string]algoSpec{
	"AES":       {keyLen: 32, ivLen: 16, block: true},
	"TripleDES": {keyLen: 24, ivLen: 8, block: true},
	"Rabbit":    {keyLen: 16, ivLen: 8, block: false},
	"RC4":       {keyLen: 32, ivLen: 0, block: false},
}

// 转换模式。
const (
	ModeEncrypt = "encrypt"
	ModeDecrypt = "decrypt"
)

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Mode   string `json:"mode"`   // encrypt | decrypt
	Text   string `json:"text"`   // 待加密文本或密文（OpenSSL 格式 Base64）
	Secret string `json:"secret"` // 密钥口令
	Algo   string `json:"algo"`   // AES | TripleDES | Rabbit | RC4
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 加密或解密文本并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	spec, ok := algos[in.Algo]
	if !ok {
		return "", fmt.Errorf("未知算法: %q（仅支持 AES/TripleDES/Rabbit/RC4）", in.Algo)
	}

	var result string
	var err error
	switch in.Mode {
	case ModeEncrypt:
		result, err = encrypt(in.Text, in.Secret, in.Algo, spec)
	case ModeDecrypt:
		result, err = decrypt(in.Text, in.Secret, in.Algo, spec)
	default:
		return "", fmt.Errorf("未知模式: %q（仅支持 encrypt/decrypt）", in.Mode)
	}
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// encrypt 使用随机盐加密文本，输出 OpenSSL 格式（Base64("Salted__"+盐+密文)）。
func encrypt(text, secret, algo string, spec algoSpec) (string, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}
	return encryptWithSalt([]byte(text), []byte(secret), algo, spec, salt)
}

// encryptWithSalt 使用固定盐加密（供测试与确定性向量对比）。
func encryptWithSalt(plain, secret []byte, algo string, spec algoSpec, salt []byte) (string, error) {
	ct, err := cryptoRun(plain, secret, salt, algo, spec, true)
	if err != nil {
		return "", err
	}
	payload := append([]byte{}, saltedMagic...)
	payload = append(payload, salt...)
	payload = append(payload, ct...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

// decrypt 解析 OpenSSL 格式并解密密文。
func decrypt(text, secret, algo string, spec algoSpec) (string, error) {
	payload, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("密文不是有效的 Base64: %w", err)
	}
	if len(payload) < len(saltedMagic)+8 {
		return "", fmt.Errorf("密文格式不正确（长度过短）")
	}
	if string(payload[:len(saltedMagic)]) != string(saltedMagic) {
		return "", fmt.Errorf("密文格式不正确（缺少 Salted__ 前缀）")
	}
	salt := payload[len(saltedMagic) : len(saltedMagic)+8]
	ct := payload[len(saltedMagic)+8:]
	pt, err := cryptoRun(ct, []byte(secret), salt, algo, spec, false)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

// cryptoRun 执行核心加密/解密：EVP 派生密钥与 IV，然后按算法处理。
func cryptoRun(data, secret, salt []byte, algo string, spec algoSpec, encrypt bool) ([]byte, error) {
	derived := evpBytesToKey(secret, salt, spec.keyLen+spec.ivLen)
	key := derived[:spec.keyLen]
	var iv []byte
	if spec.ivLen > 0 {
		iv = derived[spec.keyLen : spec.keyLen+spec.ivLen]
	}

	switch {
	case spec.block:
		return runBlock(data, key, iv, algo, encrypt)
	case algo == "Rabbit":
		return runRabbit(data, key, iv)
	case algo == "RC4":
		return runRC4(data, key)
	}
	return nil, fmt.Errorf("未知算法: %q", algo)
}

// runBlock 处理 AES / TripleDES（CBC + PKCS7）。
func runBlock(data, key, iv []byte, algo string, encrypt bool) ([]byte, error) {
	var block cipher.Block
	var err error
	if algo == "TripleDES" {
		block, err = des.NewTripleDESCipher(key)
	} else {
		block, err = aes.NewCipher(key)
	}
	if err != nil {
		return nil, fmt.Errorf("初始化 %s 失败: %w", algo, err)
	}

	payload := data
	if encrypt {
		payload = pkcs7Pad(data, block.BlockSize())
		mode := cipher.NewCBCEncrypter(block, iv)
		out := make([]byte, len(payload))
		mode.CryptBlocks(out, payload)
		return out, nil
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	out, err = pkcs7Unpad(out, block.BlockSize())
	if err != nil {
		return nil, fmt.Errorf("解密失败（密钥错误或密文损坏）: %w", err)
	}
	return out, nil
}

// runRC4 处理 RC4 流密码（无填充）。
func runRC4(data, key []byte) ([]byte, error) {
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("初始化 RC4 失败: %w", err)
	}
	out := make([]byte, len(data))
	c.XORKeyStream(out, data)
	return out, nil
}

// runRabbit 处理 Rabbit 流密码（无填充）。
func runRabbit(data, key, iv []byte) ([]byte, error) {
	r, err := newRabbitCipher(key, iv)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	r.XORKeyStream(out, data)
	return out, nil
}

// pkcs7Pad 按 blockSize 进行 PKCS#7 填充。
func pkcs7Pad(data []byte, blockSize int) []byte {
	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	for i := range pad {
		pad[i] = byte(n)
	}
	return append(data, pad...)
}

// pkcs7Unpad 移除 PKCS#7 填充并校验。
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("密文长度非法")
	}
	n := int(data[len(data)-1])
	if n == 0 || n > blockSize {
		return nil, fmt.Errorf("填充非法")
	}
	for _, b := range data[len(data)-n:] {
		if int(b) != n {
			return nil, fmt.Errorf("填充校验失败")
		}
	}
	return data[:len(data)-n], nil
}

// evpBytesToKey 复刻 crypto-js EvpKDF（MD5、1 迭代）：
// D1 = MD5(password+salt)；Di = MD5(D(i-1)+password+salt)，截取前 total 字节。
func evpBytesToKey(password, salt []byte, total int) []byte {
	out := make([]byte, 0, total)
	var prev []byte
	for len(out) < total {
		h := md5.New()
		if len(prev) > 0 {
			h.Write(prev)
		}
		h.Write(password)
		h.Write(salt)
		prev = h.Sum(nil)
		out = append(out, prev...)
	}
	return out[:total]
}

// 编译期断言：Rabbit 实现 cipher.Stream 接口。
var _ cipher.Stream = (*rabbitStream)(nil)
