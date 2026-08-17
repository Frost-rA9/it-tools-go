// Package httpstatuscodes 实现 HTTP 状态码查询工具：内置状态码表，支持按分类浏览与搜索。
// 数据移植自 it-tools（GPLv3）。
package httpstatuscodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "http-status-codes"
	Name        = "HTTP 状态码"
	Description = "查询 HTTP 状态码的说明与分类"
	Category    = "Web"
	Icon        = "Hash"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"http", "status", "codes", "状态码", "状态", "错误码"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// StatusCode 表示一个 HTTP 状态码。
type StatusCode struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// Category 表示一组状态码分类。
type CodeCategory struct {
	Category string       `json:"category"`
	Codes    []StatusCode `json:"codes"`
}

// input 是工具的输入结构。
type input struct {
	Query string `json:"query"` // 搜索关键字；为空返回全部分类
}

// output 是工具的输出结构。
type output struct {
	Results []CodeCategory `json:"results"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out := output{Results: search(in.Query)}
	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// search 在状态码表中检索：查询为空返回全部分类，否则对 code/name/description/category 做大小写不敏感包含匹配。
func search(query string) []CodeCategory {
	q := strings.TrimSpace(query)
	if q == "" {
		return codesByCategories
	}

	q = strings.ToLower(q)
	var matched []StatusCode
	for _, cat := range codesByCategories {
		catName := strings.ToLower(cat.Category)
		if strings.Contains(catName, q) {
			matched = append(matched, cat.Codes...)
			continue
		}
		for _, code := range cat.Codes {
			if matches(code, q) {
				matched = append(matched, code)
			}
		}
	}
	if len(matched) == 0 {
		return nil
	}
	return []CodeCategory{{Category: "Search results", Codes: matched}}
}

// matches 判断状态码是否命中查询词。
func matches(code StatusCode, q string) bool {
	if strings.Contains(strconv.Itoa(code.Code), q) {
		return true
	}
	if strings.Contains(strings.ToLower(code.Name), q) {
		return true
	}
	return strings.Contains(strings.ToLower(code.Description), q)
}