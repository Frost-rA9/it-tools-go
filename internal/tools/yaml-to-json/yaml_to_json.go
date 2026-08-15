// Package yamljson 实现将 YAML 文本转换为 JSON。
package yamljson

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
	"gopkg.in/yaml.v3"
)

// 工具元数据。
const (
	ID          = "yaml-to-json"
	Name        = "YAML 转 JSON"
	Description = "将 YAML 转换为 JSON"
	Category    = "转换器"
	Icon        = "AlignJustified"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"yaml", "json", "转换", "convert"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Yaml string `json:"yaml"`
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

	result, err := Convert(in.Yaml)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Convert 将 YAML 文本转换为 JSON（3 空格缩进，对齐参考项目）。
func Convert(y string) (string, error) {
	var v interface{}
	if err := yaml.Unmarshal([]byte(y), &v); err != nil {
		return "", fmt.Errorf("无效的 YAML: %w", err)
	}
	if v == nil {
		return "", nil
	}

	b, err := json.MarshalIndent(normalize(v), "", "   ")
	if err != nil {
		return "", fmt.Errorf("序列化 JSON 失败: %w", err)
	}
	return string(b), nil
}

// normalize 将 YAML 解码结果递归归一化为 JSON 友好类型（字符串化 map key）。
func normalize(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			out[k] = normalize(val)
		}
		return out
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			out[fmt.Sprint(k)] = normalize(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, val := range t {
			out[i] = normalize(val)
		}
		return out
	default:
		return v
	}
}
