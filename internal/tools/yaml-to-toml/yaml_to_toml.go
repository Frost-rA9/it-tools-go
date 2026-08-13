// Package yamltoml 实现将 YAML 文本转换为 TOML。
package yamltoml

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// 工具元数据。
const (
	ID       = "yaml-to-toml"
	Name     = "YAML 转 TOML"
	Category = "转换器"
)

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

// Convert 将 YAML 文本转换为 TOML。
func Convert(y string) (string, error) {
	if strings.TrimSpace(y) == "" {
		return "", nil
	}

	var v interface{}
	if err := yaml.Unmarshal([]byte(y), &v); err != nil {
		return "", fmt.Errorf("无效的 YAML: %w", err)
	}
	if v == nil {
		return "", nil
	}

	b, err := toml.Marshal(normalize(v))
	if err != nil {
		return "", fmt.Errorf("序列化 TOML 失败: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

// normalize 将 YAML 解码结果递归归一化为 TOML 友好类型
// （字符串化 map key，nil 值转为空字符串）。
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
	case nil:
		return ""
	default:
		return v
	}
}
