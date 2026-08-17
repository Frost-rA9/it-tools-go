// Package jwtparser 实现 JWT 解析器工具：解码 JWT 的 header 与 payload，展示各 claim。
// 行为对齐 it-tools（基于 jwt-decode），仅解码不验签。
package jwtparser

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "jwt-parser"
	Name        = "JWT 解析器"
	Description = "解码 JWT 的 Header 与 Payload，展示每个 claim"
	Category    = "Web"
	Icon        = "Key"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"jwt", "parser", "decode", "typ", "alg", "iss", "sub", "aud", "exp", "nbf", "iat", "jti", "json", "web", "token", "解析", "解码"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	JWT string `json:"jwt"` // 待解析的 JWT 字符串
}

// claim 表示一个解析出的 claim。
type claim struct {
	Claim            string `json:"claim"`
	Value            string `json:"value"`
	ClaimDescription string `json:"claimDescription,omitempty"`
	FriendlyValue    string `json:"friendlyValue,omitempty"`
}

// output 是工具的输出结构。
type output struct {
	Header  []claim `json:"header"`
	Payload []claim `json:"payload"`
}

// claimDescriptions 为 IANA 注册的 JWT claim 描述（子集，覆盖常见 claim）。
var claimDescriptions = map[string]string{
	"typ": "Type", "alg": "Algorithm", "iss": "Issuer", "sub": "Subject",
	"aud": "Audience", "exp": "Expiration Time", "nbf": "Not Before", "iat": "Issued At",
	"jti": "JWT ID", "name": "Full name", "given_name": "Given name(s) or first name(s)",
	"family_name": "Surname(s) or last name(s)", "middle_name": "Middle name(s)",
	"nickname": "Casual name", "preferred_username": "Shorthand name by which the End-User wishes to be referred to",
	"profile": "Profile page URL", "picture": "Profile picture URL", "website": "Web page or blog URL",
	"email": "Preferred e-mail address", "email_verified": "True if the e-mail address has been verified; otherwise false",
	"gender": "Gender", "birthdate": "Birthday", "zoneinfo": "Time zone", "locale": "Locale",
	"phone_number": "Preferred telephone number", "phone_number_verified": "True if the phone number has been verified; otherwise false",
	"address": "Preferred postal address", "updated_at": "Time the information was last updated",
	"azp": "Authorized party - the party to which the ID Token was issued", "nonce": "Value used to associate a Client session with an ID Token",
	"auth_time": "Time when the authentication occurred", "at_hash": "Access Token hash value",
	"c_hash": "Code hash value", "acr": "Authentication Context Class Reference", "amr": "Authentication Methods References",
	"sub_jwk": "Public key used to check the signature of an ID Token", "cnf": "Confirmation",
	"sid": "Session ID", "vot": "Vector of Trust value", "vtm": "Vector of Trust trustmark URL",
	"act": "Actor", "scope": "Scope Values", "client_id": "Client Identifier",
	"may_act": "Authorized Actor - the party that is authorized to become the actor",
	"jcard": "jCard data", "at_use_nbr": "Number of API requests for which the access token can be used",
	"roles": "Roles", "groups": "Groups", "entitlements": "Entitlements",
	"token_introspection": "Token introspection response",
}

// algorithmDescriptions 为 RFC 7518 §3.1 定义的 JWS 算法说明。
var algorithmDescriptions = map[string]string{
	"HS256": "HMAC using SHA-256", "HS384": "HMAC using SHA-384", "HS512": "HMAC using SHA-512",
	"RS256": "RSASSA-PKCS1-v1_5 using SHA-256", "RS384": "RSASSA-PKCS1-v1_5 using SHA-384", "RS512": "RSASSA-PKCS1-v1_5 using SHA-512",
	"ES256": "ECDSA using P-256 and SHA-256", "ES384": "ECDSA using P-384 and SHA-384", "ES512": "ECDSA using P-521 and SHA-512",
	"PS256": "RSASSA-PSS using SHA-256 and MGF1 with SHA-256", "PS384": "RSASSA-PSS using SHA-384 and MGF1 with SHA-384",
	"PS512": "RSASSA-PSS using SHA-512 and MGF1 with SHA-512", "none": "No digital signature or MAC performed",
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	header, payload, err := decode(in.JWT)
	if err != nil {
		return "", fmt.Errorf("无效的 JWT: %w", err)
	}

	out := output{
		Header:  parseClaims(header),
		Payload: parseClaims(payload),
	}
	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// decode 将 JWT 拆为三段，base64url 解码 header 与 payload 并校验为 JSON 对象。
func decode(jwt string) (header, payload map[string]any, err error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("格式应为 header.payload.signature")
	}

	header = map[string]any{}
	if err := decodeSegment(parts[0], &header); err != nil {
		return nil, nil, fmt.Errorf("Header 解码失败: %w", err)
	}
	payload = map[string]any{}
	if err := decodeSegment(parts[1], &payload); err != nil {
		return nil, nil, fmt.Errorf("Payload 解码失败: %w", err)
	}
	return header, payload, nil
}

// decodeSegment base64url 解码一段并反序列化为 JSON。
func decodeSegment(seg string, v any) error {
	data, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		// 兼容带 padding 的编码。
		data, err = base64.URLEncoding.DecodeString(seg)
	}
	if err != nil {
		return fmt.Errorf("base64url 解码失败: %w", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("JSON 解析失败: %w", err)
	}
	return nil
}

// parseClaims 将 claim 对象转为有序 claim 列表（对齐 it-tools 的 _.map）。
func parseClaims(obj map[string]any) []claim {
	// 以固定顺序遍历：为保证输出稳定，按 claim 名排序。
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	claims := make([]claim, 0, len(obj))
	for _, k := range keys {
		v := obj[k]
		c := claim{
			Claim:            k,
			Value:            formatValue(v),
			ClaimDescription: claimDescriptions[k],
			FriendlyValue:    friendlyValue(k, v),
		}
		claims = append(claims, c)
	}
	return claims
}

// formatValue 对齐 it-tools：对象/数组用 3 空格缩进 JSON，其余转字符串。
func formatValue(v any) string {
	switch t := v.(type) {
	case nil:
		return "null"
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		b, _ := json.Marshal(t)
		return string(b)
	default:
		b, _ := json.MarshalIndent(v, "", "   ")
		return string(b)
	}
}

// friendlyValue 对齐 it-tools：exp/nbf/iat 显示本地时间，alg 显示算法说明。
func friendlyValue(claim string, v any) string {
	switch claim {
	case "exp", "nbf", "iat":
		if n, ok := toNumber(v); ok {
			return time.Unix(n, 0).Format("2006-01-02 15:04:05")
		}
	case "alg":
		if s, ok := v.(string); ok {
			return algorithmDescriptions[s]
		}
	}
	return ""
}

// toNumber 将 JSON 数值（float64）转为 int64。
func toNumber(v any) (int64, bool) {
	f, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return int64(f), true
}