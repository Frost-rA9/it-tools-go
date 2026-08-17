// Package otpgen 实现 OTP 代码生成器：TOTP（HMAC-SHA1、6 位、30 秒周期）以及 otpauth URI 与二维码。
// 逻辑对齐 it-tools 的 otp.service（RFC 4226/6238）。
package otpgen

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "otp-generator"
	Name        = "OTP 代码生成器"
	Description = "生成 TOTP 一次性验证码（HMAC-SHA1、6 位、30 秒）并提供 otpauth URI 与二维码"
	Category    = "Web"
	Icon        = "DeviceMobile"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"otp", "code", "generator", "validator", "one", "time", "password", "authentication", "MFA", "mobile", "device", "security", "TOTP", "Time", "HMAC", "验证码", "动态口令"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Action  string `json:"action"`  // generate | codes | uri | qr
	Secret  string `json:"secret"`  // base32 secret
	Now     *int64 `json:"now"`     // 可选，毫秒时间戳（默认当前时间）
	App     string `json:"app"`     // 可选，URI 中的应用名（默认 IT-Tools）
	Account string `json:"account"` // 可选，URI 中的账户名（默认 demo-user）
}

// output 是各 action 的输出结构。
type output struct {
	Secret    string `json:"secret,omitempty"`
	Previous  string `json:"previous,omitempty"`
	Current   string `json:"current,omitempty"`
	Next      string `json:"next,omitempty"`
	Epoch     int64  `json:"epoch"`
	Counter   int64  `json:"counter"`
	SecretHex string `json:"secret_hex,omitempty"`
	NextIn    int64  `json:"next_in"`
	URI       string `json:"uri,omitempty"`
	QRDataURL string `json:"qr_data_url,omitempty"`
}

const (
	timeStepSec  = 30
	digits       = 6
	secretLength = 16
)

const base32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	var out output
	var err error
	switch in.Action {
	case "generate":
		out = output{Secret: generateSecret(secretLength)}
	case "codes":
		out, err = computeCodes(in.Secret, in.Now)
	case "uri":
		out, err = buildURI(in)
	case "qr":
		out, err = buildQR(in)
	default:
		return "", fmt.Errorf("未知操作: %q（仅支持 generate/codes/uri/qr）", in.Action)
	}
	if err != nil {
		return "", err
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// generateSecret 生成给定长度的随机 base32 字符串（RFC 4648 字母表）。
func generateSecret(length int) string {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		panic("crypto/rand 读取失败: " + err.Error())
	}
	var sb strings.Builder
	for _, b := range buf {
		sb.WriteByte(base32Alphabet[int(b)%len(base32Alphabet)])
	}
	return sb.String()
}

// decodeSecret 规范化并解码 base32 secret（大写、去除 padding，兼容空格与短杠）。
func decodeSecret(secret string) ([]byte, error) {
	normalized := strings.NewReplacer(" ", "", "-", "").Replace(strings.ToUpper(secret))
	normalized = strings.TrimRight(normalized, "=")
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	key, err := decoder.DecodeString(normalized)
	if err != nil {
		return nil, fmt.Errorf("无效的 base32 secret: %w", err)
	}
	if len(key) == 0 {
		return nil, fmt.Errorf("无效的 base32 secret: 为空")
	}
	return key, nil
}

// computeCodes 计算 TOTP 上一/当前/下一验证码及相关信息。
func computeCodes(secret string, nowMs *int64) (output, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return output{}, err
	}
	now := time.Now().UnixMilli()
	if nowMs != nil {
		now = *nowMs
	}
	epoch := now / 1000
	counter := epoch / timeStepSec

	return output{
		Previous:  hotp(key, counter-1),
		Current:   hotp(key, counter),
		Next:      hotp(key, counter+1),
		Epoch:     epoch,
		Counter:   counter,
		SecretHex: hex.EncodeToString(key),
		NextIn:    timeStepSec - epoch%timeStepSec,
	}, nil
}

// hotp 计算基于计数器的一次性密码（RFC 4226，HMAC-SHA1 动态截断）。
func hotp(key []byte, counter int64) string {
	msg := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		msg[i] = byte(counter)
		counter >>= 8
	}

	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	digest := mac.Sum(nil)

	offset := digest[19] & 0x0F
	binary := int64(digest[offset]&0x7F)<<24 |
		int64(digest[offset+1]&0xFF)<<16 |
		int64(digest[offset+2]&0xFF)<<8 |
		int64(digest[offset+3]&0xFF)

	return fmt.Sprintf("%0*d", digits, binary%1_000_000)
}

// buildURI 构建 otpauth://totp/ URI（对齐 it-tools 的 buildKeyUri）。
func buildURI(in input) (output, error) {
	if _, err := decodeSecret(in.Secret); err != nil {
		return output{}, err
	}
	app := in.App
	if app == "" {
		app = "IT-Tools"
	}
	account := in.Account
	if account == "" {
		account = "demo-user"
	}

	params := url.Values{}
	params.Set("issuer", app)
	params.Set("secret", strings.ToUpper(strings.TrimRight(in.Secret, "=")))
	params.Set("algorithm", "SHA1")
	params.Set("digits", "6")
	params.Set("period", "30")

	label := encodeComponent(app) + ":" + encodeComponent(account)
	uri := "otpauth://totp/" + label + "?" + params.Encode()
	return output{URI: uri}, nil
}

// encodeComponent 对齐 JS encodeURIComponent：仅保留未保留字符，其余按 UTF-8 字节百分号转义。
func encodeComponent(s string) string {
	const unreserved = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.!~*'()"
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if strings.IndexByte(unreserved, b) >= 0 {
			sb.WriteByte(b)
		} else {
			fmt.Fprintf(&sb, "%%%02X", b)
		}
	}
	return sb.String()
}

// writeCloser 包装 io.Writer 以满足 standard.NewWith 的 io.WriteCloser 参数。
type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }

// buildQR 生成 otpauth URI 的 PNG 二维码并返回 base64 data URL。
func buildQR(in input) (output, error) {
	uriOut, err := buildURI(in)
	if err != nil {
		return output{}, err
	}

	qrc, err := qrcode.NewWith(uriOut.URI, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	if err != nil {
		return output{}, fmt.Errorf("生成二维码失败: %w", err)
	}

	var buf bytes.Buffer
	w := standard.NewWithWriter(
		&writeCloser{&buf},
		standard.WithQRWidth(8),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	)
	if err := qrc.Save(w); err != nil {
		return output{}, fmt.Errorf("渲染二维码失败: %w", err)
	}

	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	return output{QRDataURL: dataURL}, nil
}