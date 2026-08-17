package jwtparser

import (
	"encoding/json"
	"testing"
	"time"
)

// 标准示例 JWT（jwt.io）。
const sampleJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestExecuteSample(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"jwt":"`+sampleJWT+`"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}

	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}

	if len(got.Header) != 2 {
		t.Fatalf("Header claims = %d, want 2", len(got.Header))
	}
	if len(got.Payload) != 3 {
		t.Fatalf("Payload claims = %d, want 3", len(got.Payload))
	}

	// Header：alg 与 typ（按名排序）。
	if got.Header[0].Claim != "alg" || got.Header[0].Value != "HS256" ||
		got.Header[0].ClaimDescription != "Algorithm" ||
		got.Header[0].FriendlyValue != "HMAC using SHA-256" {
		t.Errorf("Header[0] 不符: %+v", got.Header[0])
	}
	if got.Header[1].Claim != "typ" || got.Header[1].Value != "JWT" ||
		got.Header[1].ClaimDescription != "Type" || got.Header[1].FriendlyValue != "" {
		t.Errorf("Header[1] 不符: %+v", got.Header[1])
	}

	// Payload：iat / name / sub（按名排序）。
	if got.Payload[0].Claim != "iat" || got.Payload[0].Value != "1516239022" ||
		got.Payload[0].ClaimDescription != "Issued At" {
		t.Errorf("Payload[0] 不符: %+v", got.Payload[0])
	}
	wantDate := time.Unix(1516239022, 0).Format("2006-01-02 15:04:05")
	if got.Payload[0].FriendlyValue != wantDate {
		t.Errorf("iat friendlyValue = %q, want %q", got.Payload[0].FriendlyValue, wantDate)
	}
	if got.Payload[1].Claim != "name" || got.Payload[1].Value != "John Doe" ||
		got.Payload[1].ClaimDescription != "Full name" {
		t.Errorf("Payload[1] 不符: %+v", got.Payload[1])
	}
	if got.Payload[2].Claim != "sub" || got.Payload[2].Value != "1234567890" ||
		got.Payload[2].ClaimDescription != "Subject" {
		t.Errorf("Payload[2] 不符: %+v", got.Payload[2])
	}
}

func TestExecuteNestedObject(t *testing.T) {
	// payload 含嵌套对象与数组，value 应为缩进 JSON。
	token := "eyJhbGciOiJub25lIn0.eyJvYmplY3QiOnsia2V5IjoiViJ9LCJsaXN0IjpbMSwyXSwibnVsbCI6bnVsbH0."
	var e Executor
	out, err := e.Execute(t.Context(), `{"jwt":"`+token+`"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(got.Payload) != 3 {
		t.Fatalf("Payload claims = %d, want 3", len(got.Payload))
	}
	if got.Payload[0].Claim != "list" || got.Payload[0].Value != "[\n   1,\n   2\n]" {
		t.Errorf("list 值不符: %+v", got.Payload[0])
	}
	if got.Payload[1].Claim != "null" || got.Payload[1].Value != "null" {
		t.Errorf("null 值不符: %+v", got.Payload[1])
	}
	if got.Payload[2].Claim != "object" || got.Payload[2].Value != "{\n   \"key\": \"V\"\n}" {
		t.Errorf("object 值不符: %+v", got.Payload[2])
	}
}

func TestExecuteInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"空字符串", `{"jwt":""}`},
		{"两段", `{"jwt":"a.b"}`},
		{"四段", `{"jwt":"a.b.c.d"}`},
		{"非法 base64", `{"jwt":"!!!.!!!.!!!"}`},
		{"非 JSON payload", `{"jwt":"eyJhbGciOiJub25lIn0.SGVsbG8.sig"}`},
		{"非法输入 JSON", `not-json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Executor
			if _, err := e.Execute(t.Context(), tt.input); err == nil {
				t.Errorf("Execute(%s) 期望错误", tt.input)
			}
		})
	}
}