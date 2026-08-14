// Package jsonxml 实现将 JSON 文本转换为 XML（对齐 it-tools 的 _attributes/_text 约定）。
package jsonxml

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clbanning/mxj/v2"
	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "json-to-xml"
	Name        = "JSON 转 XML"
	Description = "将 JSON 转换为 XML"
	Category    = "转换器"
	Icon        = "Braces"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"json", "xml", "转换", "convert"}

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

// Convert 将 JSON 文本转换为 XML。
func Convert(j string) (string, error) {
	if strings.TrimSpace(j) == "" {
		return "", nil
	}

	var v interface{}
	if err := json.Unmarshal([]byte(j), &v); err != nil {
		return "", fmt.Errorf("无效的 JSON: %w", err)
	}

	v = toMxj(v)
	m, ok := v.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("JSON 根节点必须是对象")
	}

	var xmlBytes []byte
	var err error
	if len(m) == 1 {
		for root, val := range m {
			switch val.(type) {
			case map[string]interface{}:
				xmlBytes, err = mxj.Map(val.(map[string]interface{})).Xml(root)
			case []interface{}:
				xmlBytes, err = mxj.Map(m).Xml("root")
			default:
				xmlBytes, err = mxj.Map(map[string]interface{}{"#text": val}).Xml(root)
			}
		}
	} else {
		xmlBytes, err = mxj.Map(m).Xml("root")
	}
	if err != nil {
		return "", fmt.Errorf("序列化 XML 失败: %w", err)
	}
	return strings.TrimSpace(string(xmlBytes)), nil
}

// toMxj 将 it-tools/xml-js compact 约定（属性 `_attributes`、文本 `_text`）
// 递归转换为 mxj 约定（属性 `-x`、文本 `#text`）。
func toMxj(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{})
		if attrs, ok := t["_attributes"].(map[string]interface{}); ok {
			for k, val := range attrs {
				out["-"+k] = val
			}
		}
		if text, ok := t["_text"]; ok {
			out["#text"] = text
		}
		for k, val := range t {
			if k == "_attributes" || k == "_text" {
				continue
			}
			out[k] = toMxj(val)
		}
		return out
	case []interface{}:
		for i := range t {
			t[i] = toMxj(t[i])
		}
		return t
	default:
		return v
	}
}
