// Package urlcodec 实现 URL 字符串编码/解码工具，语义对齐 JS 的 encodeURIComponent/decodeURIComponent。
package urlcodec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "url-encoder-decoder"
	Name        = "URL 编码/解码"
	Description = "在 URL 字符串与其百分号编码形式之间进行转换"
	Category    = "Web"
	Icon        = "Link"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"url", "uri", "encode", "decode", "编码", "解码", "percent-encoding", "percent-encode"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text   string `json:"text"`   // 待处理文本
	Action string `json:"action"` // encode | decode
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	var result string
	switch in.Action {
	case "encode":
		result = encodeURIComponent(in.Text)
	case "decode":
		decoded, err := decodeURIComponent(in.Text)
		if err != nil {
			return "", fmt.Errorf("URL 解码失败: %w", err)
		}
		result = decoded
	default:
		return "", fmt.Errorf("未知操作: %q（仅支持 encode/decode）", in.Action)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// encodeURIComponent 对齐 JS encodeURIComponent：
// 除 A-Za-z0-9-_.!~*'() 外的每个 UTF-8 字节均转义为 %XX（空格为 %20）。
func encodeURIComponent(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		if isURIComponentUnreserved(b) {
			sb.WriteByte(b)
		} else {
			fmt.Fprintf(&sb, "%%%02X", b)
		}
	}
	return sb.String()
}

// isURIComponentUnreserved 判断字节是否属于 encodeURIComponent 的未转义集合。
func isURIComponentUnreserved(b byte) bool {
	switch {
	case 'a' <= b && b <= 'z', 'A' <= b && b <= 'Z', '0' <= b && b <= '9':
		return true
	case b == '-' || b == '_' || b == '.' || b == '!' || b == '~' || b == '*' || b == '\'' || b == '(' || b == ')':
		return true
	}
	return false
}

// decodeURIComponent 对齐 JS decodeURIComponent：解码 %XX 序列，
// 但字面 '+' 保持原样（Go 的 url.QueryUnescape 会把 '+' 当作空格，故先转义为 %2B）。
// 非法百分号序列或非合法 UTF-8 结果返回错误，对应 JS 的 URIError。
func decodeURIComponent(s string) (string, error) {
	escaped := strings.ReplaceAll(s, "+", "%2B")
	decoded, err := url.QueryUnescape(escaped)
	if err != nil {
		return "", err
	}
	if !utf8.ValidString(decoded) {
		return "", fmt.Errorf("非法的 UTF-8 序列")
	}
	return decoded, nil
}