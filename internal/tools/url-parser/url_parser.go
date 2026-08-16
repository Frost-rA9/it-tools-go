// Package urlparser 实现 URL 分析器工具，解析 URL 的协议、用户信息、主机、端口、路径与查询参数。
// 行为对齐浏览器 new URL()（要求绝对地址）。
package urlparser

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
	ID          = "url-parser"
	Name        = "URL 分析器"
	Description = "解析 URL 的协议、用户名、密码、主机名、端口、路径与查询参数"
	Category    = "Web"
	Icon        = "Unlink"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"url", "parser", "protocol", "origin", "params", "port", "username", "password", "href", "解析", "分析"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	URL string `json:"url"` // 待解析的 URL
}

// Param 表示一个查询参数键值对。
type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// output 是工具的输出结构，字段对齐 it-tools 的属性列表（不含 hash/origin/href）。
type output struct {
	Protocol string  `json:"protocol"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Hostname string  `json:"hostname"`
	Port     string  `json:"port"`
	Pathname string  `json:"pathname"`
	Search   string  `json:"search"`
	Params   []Param `json:"params"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	u, err := parse(in.URL)
	if err != nil {
		return "", fmt.Errorf("无效的 URL: %w", err)
	}

	username := ""
	password := ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}

	out := output{
		Protocol: strings.ToLower(u.Scheme) + ":",
		Username: username,
		Password: password,
		Hostname: strings.ToLower(u.Hostname()),
		Port:     u.Port(),
		Pathname: u.EscapedPath(),
		Search:   "",
		Params:   parseParams(u.RawQuery),
	}
	if u.RawQuery != "" {
		out.Search = "?" + u.RawQuery
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// parse 解析 URL 并校验其为绝对地址（对齐 new URL()：必须含 scheme 与 host）。
func parse(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() || u.Host == "" {
		return nil, fmt.Errorf("必须是绝对 URL")
	}
	return u, nil
}

// parseParams 有序解析查询参数（保留重复键、空值；'+' 解码为空格，对齐 URLSearchParams）。
func parseParams(rawQuery string) []Param {
	if rawQuery == "" {
		return nil
	}
	var params []Param
	for _, pair := range strings.Split(rawQuery, "&") {
		if pair == "" {
			continue
		}
		key, value := pair, ""
		if idx := strings.IndexByte(pair, '='); idx >= 0 {
			key, value = pair[:idx], pair[idx+1:]
		}
		k, err1 := url.QueryUnescape(key)
		v, err2 := url.QueryUnescape(value)
		if err1 != nil || err2 != nil {
			// 非法转义：保留原始片段，避免解析失败。
			params = append(params, Param{Key: key, Value: value})
			continue
		}
		params = append(params, Param{Key: k, Value: v})
	}
	return params
}