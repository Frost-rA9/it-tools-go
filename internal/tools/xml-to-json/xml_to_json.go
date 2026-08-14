// Package xmljson 实现将 XML 文本转换为 JSON（对齐 it-tools 的 _attributes/_text 约定）。
package xmljson

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
	ID          = "xml-to-json"
	Name        = "XML 转 JSON"
	Description = "将 XML 转换为 JSON"
	Category    = "转换器"
	Icon        = "Braces"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"xml", "json", "转换", "convert"}

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

// Convert 将 XML 文本转换为 JSON（2 空格缩进，对齐参考项目）。
func Convert(x string) (string, error) {
	if strings.TrimSpace(x) == "" {
		return "", nil
	}

	mv, err := mxj.NewMapXml([]byte(x))
	if err != nil {
		return "", fmt.Errorf("无效的 XML: %w", err)
	}

	b, err := json.MarshalIndent(toXmljs(map[string]interface{}(mv)), "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化 JSON 失败: %w", err)
	}
	return string(b), nil
}

// toXmljs 将 mxj 约定（属性 `-x`、文本 `#text`）递归转换为
// it-tools/xml-js compact 约定（属性 `_attributes`、文本 `_text`）。
func toXmljs(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		attrs := make(map[string]interface{})
		hasAttr := false
		var text interface{}
		hasText := false
		children := make(map[string]interface{})
		for k, val := range t {
			switch {
			case strings.HasPrefix(k, "-"):
				hasAttr = true
				attrs[k[1:]] = val
			case k == "#text":
				hasText = true
				text = val
			default:
				children[k] = toXmljs(val)
			}
		}
		if !hasAttr && !hasText {
			return children
		}
		out := make(map[string]interface{})
		if hasAttr {
			out["_attributes"] = attrs
		}
		if hasText {
			out["_text"] = text
		}
		for k, val := range children {
			out[k] = val
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, val := range t {
			out[i] = toXmljs(val)
		}
		return out
	default:
		return v
	}
}
