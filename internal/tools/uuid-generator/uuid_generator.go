// Package uuidgen 实现 UUID 生成工具（对齐 it-tools uuid-generator）。
package uuidgen

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "uuid-generator"
	Name        = "UUID 生成器"
	Description = "生成 NIL、v1、v3、v4 或 v5 版本的 UUID（128 位通用唯一标识符）"
	Category    = registry.CategoryCrypto
	Icon        = "Fingerprint"
)

// Keywords 为搜索关键词。
var Keywords = []string{"uuid", "v4", "v1", "v3", "v5", "nil", "random", "id", "identifier", "unique", "token", "string", "标识"}

// 数量上限。
const countMax = 50

// 支持生成 UUID 的版本。
var versions = []string{"NIL", "v1", "v3", "v4", "v5"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Version   string `json:"version"`   // NIL | v1 | v3 | v4 | v5
	Count     int    `json:"count"`     // 生成数量（1..50）
	Namespace string `json:"namespace"` // v3/v5 的命名空间 UUID
	Name      string `json:"name"`      // v3/v5 的名称
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 生成 UUID 列表并返回结果 JSON（换行连接）。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if in.Count < 1 || in.Count > countMax {
		return "", fmt.Errorf("数量必须在 1..%d 之间（当前 %d）", countMax, in.Count)
	}

	generator, err := newGenerator(in)
	if err != nil {
		return "", err
	}

	ids := make([]string, 0, in.Count)
	for i := 0; i < in.Count; i++ {
		id, err := generator(i)
		if err != nil {
			return "", err
		}
		ids = append(ids, id)
	}

	out, err := json.Marshal(output{Result: strings.Join(ids, "\n")})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// generator 返回按索引生成 UUID 字符串的函数。
func newGenerator(in input) (func(int) (string, error), error) {
	switch in.Version {
	case "NIL":
		return func(int) (string, error) { return uuid.Nil.String(), nil }, nil
	case "v4":
		return func(int) (string, error) { return uuid.New().String(), nil }, nil
	case "v1":
		return func(int) (string, error) {
			id, err := uuid.NewUUID()
			if err != nil {
				return "", fmt.Errorf("生成 v1 UUID 失败: %w", err)
			}
			return id.String(), nil
		}, nil
	case "v3", "v5":
		ns, err := uuid.Parse(in.Namespace)
		if err != nil {
			return nil, fmt.Errorf("命名空间不是有效的 UUID: %w", err)
		}
		if in.Name == "" {
			return nil, fmt.Errorf("v%s 模式需要提供名称（name）", in.Version)
		}
		data := []byte(in.Name)
		if in.Version == "v3" {
			return func(int) (string, error) { return uuid.NewMD5(ns, data).String(), nil }, nil
		}
		return func(int) (string, error) { return uuid.NewSHA1(ns, data).String(), nil }, nil
	default:
		return nil, fmt.Errorf("未知版本: %q（仅支持 NIL/v1/v3/v4/v5）", in.Version)
	}
}
