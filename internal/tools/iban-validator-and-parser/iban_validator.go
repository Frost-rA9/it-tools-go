// Package ibanvalidator 实现 IBAN 验证器和解析器。
//
// 自研零依赖实现：MOD-97 校验 + 内嵌各国 IBAN 规格表（iban_data.go），
// 错误语义对齐参考项目 it-tools（ibantools），输出中文文案。
package ibanvalidator

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "iban-validator-and-parser"
	Name        = "IBAN 验证器和解析器"
	Description = "验证 IBAN 有效性并解析国家、BBAN 等信息"
	Category    = registry.CategoryData
	Icon        = "BuildingBank"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"iban", "validator", "and", "parser", "bic", "bank", "IBAN", "银行", "账号", "验证", "解析"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Iban string `json:"iban"` // 待校验的 IBAN（可含空格/连字符）
}

// output 是工具的输出结构。
type output struct {
	Valid          bool     `json:"valid"`          // 是否有效
	Errors         []string `json:"errors"`         // 错误列表（中文）
	QRIban         bool     `json:"qrIban"`         // 是否为瑞士 QR-IBAN
	CountryCode    string   `json:"countryCode"`    // 国家代码（ISO2）
	BBAN           string   `json:"bban"`           // BBAN 部分
	FriendlyFormat string   `json:"friendlyFormat"` // 每 4 位空格分组的友好格式
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out := Validate(in.Iban)
	res, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(res), nil
}

// Validate 规范化并校验 IBAN，返回完整解析结果。
func Validate(raw string) output {
	iban := normalize(raw)
	out := output{}

	if iban == "" {
		out.Errors = append(out.Errors, "未提供 IBAN")
		return out
	}
	out.CountryCode = iban[:2]
	out.BBAN = safeBban(iban)
	out.FriendlyFormat = FriendlyFormat(iban)

	// 国家代码：前两位应为字母。
	if len(iban) < 4 || !isAlpha(iban[0]) || !isAlpha(iban[1]) {
		out.Errors = append(out.Errors, "缺少 IBAN 国家代码")
	}

	// 通用长度范围（ISO 13616：15-34 位）。
	if len(iban) < 15 || len(iban) > 34 {
		out.Errors = append(out.Errors, "IBAN 长度错误（应在 15-34 位之间）")
	}

	// 校验位（第 3-4 位）必须为数字。
	if len(iban) >= 4 {
		if !isAllDigits(iban[2:4]) {
			out.Errors = append(out.Errors, "校验位不是数字")
		}
	}

	// 国家规格表校验（BBAN 长度与结构）。
	if sp, ok := specIndex[out.CountryCode]; ok && len(iban) >= 4 {
		bban := iban[4:]
		if len(bban) != sp.BLength {
			out.Errors = append(out.Errors, "BBAN 长度错误（应为 "+fmt.Sprint(sp.BLength)+" 位）")
		}
		if !sp.bbanRegex.MatchString(bban) {
			out.Errors = append(out.Errors, "BBAN 格式错误")
		}
		if len(iban) != sp.ILength {
			out.Errors = append(out.Errors, "IBAN 长度错误（应为 "+fmt.Sprint(sp.ILength)+" 位）")
		}
	}

	// MOD-97 校验：先确保字符合法。
	if containsOnlyAlphaNum(iban) {
		if mod97(iban) != 1 {
			out.Errors = append(out.Errors, "IBAN 校验和不通过")
		}
	} else if len(iban) >= 2 && isAlpha(iban[0]) && isAlpha(iban[1]) {
		out.Errors = append(out.Errors, "IBAN 包含非法字符")
	}

	out.Valid = len(out.Errors) == 0
	out.QRIban = len(iban) == 21 && strings.HasPrefix(iban, "CH") && mod97(iban) == 1 && isQRAccount(iban)
	return out
}

// normalize 规范化输入：大写并去除空格与连字符。
func normalize(iban string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(iban), " ", ""), "-", ""))
}

// mod97 计算 IBAN 的 MOD 97 余数（字母转数字 A=10..Z=35，前 4 位移至末尾）。
// 调用方需确保 iban 仅含字母数字且长度 >= 4；有效 IBAN 满足 mod97 == 1。
func mod97(iban string) int {
	s := iban[4:] + iban[:4]
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			b.WriteByte(byte('1' + (c-'A')/10))
			b.WriteByte(byte('0' + (c-'A')%10))
		} else {
			b.WriteByte(c)
		}
	}
	digits := b.String()
	rem := 0
	for i := 0; i < len(digits); i++ {
		rem = (rem*10 + int(digits[i]-'0')) % 97
	}
	return rem
}

// isQRAccount 判断瑞士 IBAN 的账号是否为 QR-IBAN 账号
// （BBAN 第 6-10 位，即 12 位账号的前 5 位，为 30000 或 31000-31999）。
func isQRAccount(iban string) bool {
	if len(iban) != 21 {
		return false
	}
	acctPrefix := iban[4+5 : 4+5+5] // BBAN[5:10]
	return acctPrefix == "30000" || (acctPrefix >= "31000" && acctPrefix <= "31999")
}

// FriendlyFormat 将 IBAN 格式化为每 4 位空格分组。
func FriendlyFormat(iban string) string {
	if iban == "" {
		return ""
	}
	var b strings.Builder
	for i, c := range iban {
		if i > 0 && i%4 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(c)
	}
	return b.String()
}

func safeBban(iban string) string {
	if len(iban) >= 4 {
		return iban[4:]
	}
	return ""
}

func isAlpha(c byte) bool { return c >= 'A' && c <= 'Z' }

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return len(s) > 0
}

func containsOnlyAlphaNum(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= '0' && c <= '9' || c >= 'A' && c <= 'Z') {
			return false
		}
	}
	return len(s) > 0
}

// spec 是单个国家的 IBAN 规格。
type spec struct {
	Code      string // ISO2 国家代码
	ILength   int    // IBAN 总长度
	BLength   int    // BBAN 长度
	bbanRegex *regexp.Regexp
}

// specIndex 国家代码 → 规格。
var specIndex map[string]spec

func init() {
	specIndex = make(map[string]spec, len(ibanSpecs))
	for _, s := range ibanSpecs {
		specIndex[s.Code] = s
	}
}
