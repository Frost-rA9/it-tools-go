package ibanvalidator

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"it-tools-go/internal/registry"
)

// computeChecksum 依据 MOD-97 计算 IBAN 校验位（测试辅助，用于构造合法 IBAN）。
// 校验位 = 98 - mod(bban + 国家 + "00")；国家+BBAN 的长度决定 IBAN 总长度。
func computeChecksum(cc, bban string) string {
	rem := modDigits(bban + cc + "00")
	cs := 98 - rem
	return string([]byte{byte('0' + cs/10), byte('0' + cs%10)})
}

// modDigits 将字母数字串（字母转 A=10..Z=35）计算 MOD 97 余数。
func modDigits(s string) int {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			v := int(c-'A') + 10
			b.WriteByte(byte('0' + v/10))
			b.WriteByte(byte('0' + v%10))
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

func TestValidExamples(t *testing.T) {
	cases := []struct {
		name string
		iban string
		cc   string
	}{
		{"法国", "FR7630006000011234567890189", "FR"},
		{"德国", "DE89370400440532013000", "DE"},
		{"英国", "GB29NWBK60161331926819", "GB"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := Validate(c.iban)
			if !out.Valid {
				t.Errorf("应有效，结果: %+v", out)
			}
			if out.CountryCode != c.cc {
				t.Errorf("国家 = %q, want %q", out.CountryCode, c.cc)
			}
			if out.BBAN == "" || out.FriendlyFormat == "" {
				t.Error("缺少 BBAN/友好格式")
			}
		})
	}
}

func TestNormalization(t *testing.T) {
	// 带空格与小写应同样识别。
	out := Validate(" fr76 3000 6000 0112 3456 7890 189 ")
	if !out.Valid {
		t.Errorf("规范化后应有效: %+v", out)
	}
	out = Validate("GB29-NWBK-6016-1331-9268-19")
	if !out.Valid {
		t.Errorf("连字符规范化后应有效: %+v", out)
	}
}

func TestInvalidChecksum(t *testing.T) {
	// 把 FR 示例的校验位 76 改 77。
	out := Validate("FR7730006000011234567890189")
	if out.Valid {
		t.Error("校验位错误应无效")
	}
	if !strings.Contains(strings.Join(out.Errors, " "), "校验和") {
		t.Errorf("应有校验和错误: %+v", out.Errors)
	}
}

func TestEmpty(t *testing.T) {
	out := Validate("")
	if out.Valid {
		t.Error("空输入应无效")
	}
	if len(out.Errors) == 0 {
		t.Error("空输入应有错误提示")
	}
}

func TestTooShort(t *testing.T) {
	out := Validate("DE89")
	if out.Valid {
		t.Error("过短应无效")
	}
}

func TestChecksumNotNumber(t *testing.T) {
	out := Validate("FRAB30006000011234567890189")
	if out.Valid {
		t.Error("校验位非数字应无效")
	}
	if !strings.Contains(strings.Join(out.Errors, " "), "校验位") {
		t.Errorf("应有校验位错误: %+v", out.Errors)
	}
}

func TestWrongBbanLength(t *testing.T) {
	// FR 但 BBAN 长度错误（去掉 2 位）。
	out := Validate("FR7630006000011234567890")
	if out.Valid {
		t.Error("BBAN 长度错误应无效")
	}
	if !strings.Contains(strings.Join(out.Errors, " "), "BBAN") {
		t.Errorf("应有 BBAN 错误: %+v", out.Errors)
	}
}

func TestWrongBbanFormat(t *testing.T) {
	// GB 的 BBAN 首 4 位应为字母，构造数字开头且长度合法，再修正校验位。
	raw := "GB" + computeChecksum("GB", "123460161331926819") + "123460161331926819"
	if len(raw) != 22 {
		t.Fatalf("长度 = %d, want 22", len(raw))
	}
	out := Validate(raw)
	if !strings.Contains(strings.Join(out.Errors, " "), "BBAN 格式错误") {
		t.Errorf("应有 BBAN 格式错误: %+v", out.Errors)
	}
}

func TestQRIBAN(t *testing.T) {
	// 构造一个合法的瑞士 QR-IBAN：银行代码 00762（5 位），账号 300000000000（12 位，前缀 30000）。
	bban := "00762" + "300000000000"
	if len(bban) != 17 {
		t.Fatalf("BBAN 长度 = %d, want 17", len(bban))
	}
	iban := "CH" + computeChecksum("CH", bban) + bban
	out := Validate(iban)
	if !out.Valid {
		t.Fatalf("QR-IBAN 应有效: %+v", out.Errors)
	}
	if !out.QRIban {
		t.Error("应识别为 QR-IBAN")
	}

	// 非 QR：账号前缀非 3xxxx。
	bban2 := "00762" + "123450000000"
	iban2 := "CH" + computeChecksum("CH", bban2) + bban2
	out2 := Validate(iban2)
	if out2.QRIban {
		t.Error("普通瑞士账号不应判定为 QR-IBAN")
	}
}

func TestUnknownCountry(t *testing.T) {
	// 未知国家 XX：长度合法（16）且 MOD-97 通过时应视为有效（无规格表仅做通用校验）。
	ibaan := "XX" + computeChecksum("XX", "123456789012") + "123456789012"
	out := Validate(ibaan)
	if len(ibaan) != 16 {
		t.Fatalf("长度 = %d", len(ibaan))
	}
	if !out.Valid {
		t.Errorf("未知国家但长度与校验和正确应有效: %+v", out.Errors)
	}
}

func TestFriendlyFormat(t *testing.T) {
	if got := FriendlyFormat("FR7630006000011234567890189"); got != "FR76 3000 6000 0112 3456 7890 189" {
		t.Errorf("友好格式 = %q", got)
	}
	if FriendlyFormat("") != "" {
		t.Error("空输入友好格式应为空")
	}
}

func TestExecute(t *testing.T) {
	outStr, err := (Executor{}).Execute(context.Background(), `{"iban":"DE89370400440532013000"}`)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if !out.Valid || out.CountryCode != "DE" {
		t.Errorf("输出不符: %+v", out)
	}
}

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != registry.CategoryData || meta.Icon != "BuildingBank" {
		t.Errorf("元数据不符: %+v", meta)
	}
}
