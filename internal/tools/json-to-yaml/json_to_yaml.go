// Package jsontoyaml 实现将 JSON 文本转换为 YAML。
package jsontoyaml

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"it-tools-go/internal/registry"
	"gopkg.in/yaml.v3"
)

// 工具元数据。
const (
	ID          = "json-to-yaml"
	Name        = "JSON 转 YAML"
	Description = "将 JSON 转换为 YAML"
	Category    = "转换器"
	Icon        = "Braces"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"json", "yaml", "转换", "convert"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"`
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

	result, err := Convert(in.Text)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Convert 将 JSON 文本转换为 YAML。
func Convert(j string) (string, error) {
	if strings.TrimSpace(j) == "" {
		return "", nil
	}

	var v interface{}
	dec := json.NewDecoder(strings.NewReader(j))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return "", fmt.Errorf("无效的 JSON: %w", err)
	}
	var extra interface{}
	if err := dec.Decode(&extra); err != io.EOF {
		return "", fmt.Errorf("无效的 JSON: 存在多余内容")
	}

	b, err := yaml.Marshal(normalize(v))
	if err != nil {
		return "", fmt.Errorf("序列化 YAML 失败: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

// normalize 将 JSON 解码结果（UseNumber）递归归一化为 YAML 友好类型：
// json.Number 按格式转 int64（整数形式）或 float64。
func normalize(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			out[k] = normalize(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, val := range t {
			out[i] = normalize(val)
		}
		return out
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return i
		}
		if f, err := t.Float64(); err == nil {
			return f
		}
		return t.String()
	default:
		return v
	}
}
