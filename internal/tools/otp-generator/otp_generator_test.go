package otpgen

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

// rfcSecret 是 RFC 4226 附录 D 测试用的 secret：ASCII "12345678901234567890" 的 base32 编码。
const rfcSecret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

// exec 便捷执行工具并解析输出。
func exec(t *testing.T, input string) output {
	t.Helper()
	var e Executor
	out, err := e.Execute(t.Context(), input)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	return o
}

func TestHOTPRFCVectors(t *testing.T) {
	// RFC 4226 附录 D：counter 0..3 → 755224 / 287082 / 359152 / 969429。
	want := map[int]string{0: "755224", 1: "287082", 2: "359152", 3: "969429"}
	for c, code := range want {
		o := exec(t, `{"action":"codes","secret":"`+rfcSecret+`","now":`+strconv.FormatInt(int64(c)*30000, 10)+`}`)
		if o.Current != code {
			t.Errorf("counter=%d current = %q, want %q", c, o.Current, code)
		}
		if o.Counter != int64(c) {
			t.Errorf("counter 字段 = %d, want %d", o.Counter, c)
		}
	}
}

func TestTOTP(t *testing.T) {
	// epoch 59 → counter 1 → current 287082（RFC 6238 SHA1）。
	o := exec(t, `{"action":"codes","secret":"`+rfcSecret+`","now":59000}`)
	if o.Epoch != 59 || o.Counter != 1 {
		t.Fatalf("epoch/counter = %d/%d, want 59/1", o.Epoch, o.Counter)
	}
	if o.Current != "287082" {
		t.Errorf("current = %q, want 287082", o.Current)
	}
	if o.Previous != "755224" || o.Next != "359152" {
		t.Errorf("previous/next = %q/%q, want 755224/359152", o.Previous, o.Next)
	}
	if o.NextIn != 1 {
		t.Errorf("next_in = %d, want 1", o.NextIn)
	}
	if o.SecretHex != "3132333435363738393031323334353637383930" {
		t.Errorf("secret_hex = %q", o.SecretHex)
	}
}

func TestGenerateSecret(t *testing.T) {
	o := exec(t, `{"action":"generate"}`)
	if len(o.Secret) != 16 {
		t.Fatalf("secret 长度 = %d, want 16", len(o.Secret))
	}
	for _, c := range o.Secret {
		if !strings.ContainsRune(base32Alphabet, c) {
			t.Errorf("secret 含非法字符 %q", c)
		}
	}
}

func TestBuildURI(t *testing.T) {
	o := exec(t, `{"action":"uri","secret":"JBSWY3DPEHPK3PXP"}`)
	want := "otpauth://totp/IT-Tools:demo-user?algorithm=SHA1&digits=6&issuer=IT-Tools&period=30&secret=JBSWY3DPEHPK3PXP"
	if o.URI != want {
		t.Errorf("uri = %q, want %q", o.URI, want)
	}

	// 自定义 app/account（含特殊字符应被转义）。
	o = exec(t, `{"action":"uri","secret":"JBSWY3DPEHPK3PXP","app":"My App","account":"user@example.com"}`)
	if !strings.Contains(o.URI, "My%20App:user%40example.com") {
		t.Errorf("自定义 app/account 转义不符: %q", o.URI)
	}
}

func TestBuildQR(t *testing.T) {
	o := exec(t, `{"action":"qr","secret":"JBSWY3DPEHPK3PXP"}`)
	if !strings.HasPrefix(o.QRDataURL, "data:image/png;base64,") {
		t.Fatalf("qr 前缀不符: %.40q", o.QRDataURL)
	}
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(o.QRDataURL, "data:image/png;base64,"))
	if err != nil {
		t.Fatalf("base64 解码失败: %v", err)
	}
	if len(b) < 8 || string(b[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatalf("不是合法的 PNG: %x", b[:min(len(b), 8)])
	}
}

func TestErrors(t *testing.T) {
	var e Executor
	tests := []struct {
		name  string
		input string
	}{
		{"非法 base32", `{"action":"codes","secret":"!!!!"}`},
		{"空 secret", `{"action":"codes","secret":""}`},
		{"未知操作", `{"action":"foo","secret":"JBSWY3DPEHPK3PXP"}`},
		{"非法输入 JSON", `not-json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := e.Execute(t.Context(), tt.input); err == nil {
				t.Errorf("Execute(%s) 期望错误", tt.input)
			}
		})
	}
}
