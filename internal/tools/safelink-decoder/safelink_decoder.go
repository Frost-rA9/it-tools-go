// Package safelink 实现 Outlook 安全链接解码器：从 Outlook SafeLink 链接中还原真实 URL。
// 对齐 it-tools 的 safelink-decoder.service。
package safelink

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "safelink-decoder"
	Name        = "Outlook 安全链接解码器"
	Description = "解码 Outlook SafeLink 链接，还原真实 URL"
	Category    = "Web"
	Icon        = "Mailbox"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"outlook", "safelink", "decoder", "安全链接", "解码", "邮件", "链接"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	URL string `json:"url"` // 待解码的 Outlook SafeLink URL
}

// output 是工具的输出结构。
type output struct {
	DecodedURL string `json:"decoded_url"`
}

// safelinkHostSuffix 为 Outlook SafeLink 域名后缀。
const safelinkHostSuffix = ".safelinks.protection.outlook.com"

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	decoded, err := decode(in.URL)
	if err != nil {
		return "", err
	}

	outJSON, err := json.Marshal(output{DecodedURL: decoded})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// decode 校验 SafeLink 域名并提取 url 查询参数。
func decode(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("无效的 URL: %w", err)
	}
	if !strings.Contains(strings.ToLower(u.Host), safelinkHostSuffix) {
		return "", fmt.Errorf("无效的 SafeLink URL：域名不符")
	}
	return u.Query().Get("url"), nil
}