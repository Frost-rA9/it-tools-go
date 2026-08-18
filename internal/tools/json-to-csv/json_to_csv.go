// Package jsoncsv 实现 JSON 转 CSV 工具：将 JSON 对象数组转换为 CSV 表格。
package jsoncsv

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

const (
	ID          = "json-to-csv"
	Name        = "JSON 转 CSV"
	Description = "将 JSON 对象数组转换为 CSV 表格"
	Category    = "开发"
	Icon        = "Table"
)

var Keywords = []string{"json", "csv", "convert", "表格", "转换", "对象数组"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	JSON          string `json:"json"`           // 对象数组 JSON
	Delimiter     string `json:"delimiter"`      // "," | ";" | "\t"（默认 ","）
	IncludeHeader *bool  `json:"include_header"` // 是否输出表头（默认 true）
}

// output 是工具的输出结构。
type output struct {
	CSV     string `json:"csv"`
	Rows    int    `json:"rows"`    // 数据行数
	Columns int    `json:"columns"` // 列数
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回 CSV 结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	// 用 UseNumber 保留数字原文，避免 float64 精度损失
	dec := json.NewDecoder(strings.NewReader(in.JSON))
	dec.UseNumber()
	var arr []map[string]any
	if err := dec.Decode(&arr); err != nil {
		return "", fmt.Errorf("JSON 必须是对象数组: %w", err)
	}

	// 收集列键并排序（map 迭代无序，排序保证输出稳定）
	keys := make([]string, 0, 16)
	seen := make(map[string]bool, 16)
	for _, obj := range arr {
		for k := range obj {
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)

	delimiter := in.Delimiter
	if delimiter == "" {
		delimiter = ","
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	switch delimiter {
	case ";":
		w.Comma = ';'
	case "\t":
		w.Comma = '\t'
	default:
		w.Comma = ','
	}

	if in.IncludeHeader == nil || *in.IncludeHeader {
		if len(keys) > 0 {
			if err := w.Write(keys); err != nil {
				return "", fmt.Errorf("写入表头失败: %w", err)
			}
		}
	}
	for _, obj := range arr {
		rec := make([]string, len(keys))
		for i, k := range keys {
			rec[i] = cellString(obj[k])
		}
		if err := w.Write(rec); err != nil {
			return "", fmt.Errorf("写入行失败: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("CSV 写入失败: %w", err)
	}

	out, err := json.Marshal(output{
		CSV:     buf.String(),
		Rows:    len(arr),
		Columns: len(keys),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// cellString 将 JSON 值转换为 CSV 单元格文本。
func cellString(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		// 对象 / 数组 / 混合：紧凑 JSON 表示
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return string(b)
	}
}
