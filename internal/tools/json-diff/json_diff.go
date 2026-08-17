// Package jsondiff 实现 JSON 差异比较工具：递归对比两个 JSON，生成差异树。
// 解析用 JSON5（对齐 it-tools），diff 语义对齐 it-tools 的 json-diff.models。
package jsondiff

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/swaggest/assertjson/json5"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "json-diff"
	Name        = "JSON 差异比较"
	Description = "递归比较两个 JSON，标记新增、删除、修改与不变"
	Category    = "Web"
	Icon        = "GitCompare"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"json", "diff", "compare", "difference", "object", "data", "差异", "比较", "对比"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Left                string `json:"left"`                  // 左侧 JSON（JSON5）
	Right               string `json:"right"`                 // 右侧 JSON（JSON5）
	OnlyShowDifferences bool   `json:"only_show_differences"` // 是否仅显示差异
}

// Difference 表示差异树中的一个节点。
type Difference struct {
	Key      any          `json:"key"`
	Type     string       `json:"type"`     // object | array | value
	Status   string       `json:"status"`   // added | removed | updated | unchanged | children-updated
	OldValue any          `json:"oldValue,omitempty"`
	Value    any          `json:"value,omitempty"`
	Children []Difference `json:"children,omitempty"`
}

// output 是工具的输出结构。
type output struct {
	Same bool        `json:"same"` // 两侧 JSON 是否完全相同
	Root *Difference `json:"root"` // 差异树根节点（same 时为 null）
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	left, err := parse(in.Left)
	if err != nil {
		return "", fmt.Errorf("左侧 JSON 无效: %w", err)
	}
	right, err := parse(in.Right)
	if err != nil {
		return "", fmt.Errorf("右侧 JSON 无效: %w", err)
	}

	out := output{
		Same: deepEqual(left, right),
	}
	if !out.Same {
		root := diff(value{value: left, present: true}, value{value: right, present: true}, "", in.OnlyShowDifferences)
		out.Root = &root
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// value 包装一个可能缺失的 JSON 值（区分「键不存在」与「值为 null」）。
type value struct {
	value   any
	present bool
}

// parse 用 JSON5 解析输入。
func parse(s string) (any, error) {
	var v any
	if err := json5.Unmarshal([]byte(s), &v); err != nil {
		return nil, err
	}
	return v, nil
}

// deepEqual 判断两个 JSON 值是否深度相等（对齐 lodash isEqual）。
func deepEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// getType 返回值的类型：object | array | value。
func getType(v value) string {
	if !v.present || v.value == nil {
		return "value"
	}
	switch v.value.(type) {
	case map[string]any:
		return "object"
	case []any:
		return "array"
	}
	return "value"
}

// getStatus 计算两个值的差异状态（对齐 it-tools 的 getStatus）。
func getStatus(a, b value) string {
	if !a.present {
		return "added"
	}
	if !b.present {
		return "removed"
	}
	if deepEqual(a.value, b.value) {
		return "unchanged"
	}
	if getType(a) == "object" && getType(b) == "object" || getType(a) == "array" && getType(b) == "array" {
		return "children-updated"
	}
	return "updated"
}

// diff 递归计算差异树（对齐 it-tools 的 diff）。
func diff(a, b value, key any, onlyShow bool) Difference {
	aType, bType := getType(a), getType(b)
	if aType == "array" && bType == "array" {
		return Difference{
			Key:      key,
			Type:     "array",
			Children: diffArrays(a, b, onlyShow),
			OldValue: a.value,
			Value:    b.value,
			Status:   getStatus(a, b),
		}
	}
	if aType == "object" && bType == "object" {
		return Difference{
			Key:      key,
			Type:     "object",
			Children: diffObjects(a, b, onlyShow),
			OldValue: a.value,
			Value:    b.value,
			Status:   getStatus(a, b),
		}
	}
	return Difference{
		Key:      key,
		Type:     "value",
		OldValue: a.value,
		Value:    b.value,
		Status:   getStatus(a, b),
	}
}

// diffObjects 对比对象键的并集。
func diffObjects(a, b value, onlyShow bool) []Difference {
	keys := map[string]struct{}{}
	if a.present {
		for k := range a.value.(map[string]any) {
			keys[k] = struct{}{}
		}
	}
	if b.present {
		for k := range b.value.(map[string]any) {
			keys[k] = struct{}{}
		}
	}
	result := make([]Difference, 0, len(keys))
	for k := range keys {
		d := diff(valueAt(a, k), valueAt(b, k), k, onlyShow)
		if !onlyShow || d.Status != "unchanged" {
			result = append(result, d)
		}
	}
	return result
}

// diffArrays 按索引对比数组（对齐 it-tools 的 diffArrays）。
func diffArrays(a, b value, onlyShow bool) []Difference {
	aLen, bLen := 0, 0
	if a.present {
		aLen = len(a.value.([]any))
	}
	if b.present {
		bLen = len(b.value.([]any))
	}
	maxLen := aLen
	if bLen > maxLen {
		maxLen = bLen
	}
	result := make([]Difference, 0, maxLen)
	for i := 0; i < maxLen; i++ {
		d := diff(valueAt(a, i), valueAt(b, i), i, onlyShow)
		if !onlyShow || d.Status != "unchanged" {
			result = append(result, d)
		}
	}
	return result
}

// valueAt 取对象键或数组索引处的值，不存在时 present=false。
func valueAt(v value, k any) value {
	if !v.present {
		return value{}
	}
	switch t := v.value.(type) {
	case map[string]any:
		val, ok := t[k.(string)]
		return value{value: val, present: ok}
	case []any:
		idx, ok := k.(int)
		if !ok || idx < 0 || idx >= len(t) {
			return value{}
		}
		return value{value: t[idx], present: true}
	}
	return value{}
}