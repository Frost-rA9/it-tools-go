// Package listconv 实现按配置对换行分隔的列表进行转换。
package listconv

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "list-converter"
	Name        = "列表转换器"
	Description = "对换行分隔的列表进行排序、去重、反转、包装等转换"
	Category    = "转换器"
	Icon        = "List"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"list", "列表", "转换", "sort", "排序", "去重", "反转", "prefix", "suffix", "separator"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text             string `json:"text"`
	LowerCase        bool   `json:"lower_case"`
	TrimItems        bool   `json:"trim_items"`
	ItemPrefix       string `json:"item_prefix"`
	ItemSuffix       string `json:"item_suffix"`
	ListPrefix       string `json:"list_prefix"`
	ListSuffix       string `json:"list_suffix"`
	ReverseList      bool   `json:"reverse_list"`
	SortList         string `json:"sort_list"` // "" | asc | desc
	RemoveDuplicates bool   `json:"remove_duplicates"`
	Separator        string `json:"separator"`
	KeepLineBreaks   bool   `json:"keep_line_breaks"`
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

	out, err := json.Marshal(output{Result: convert(in.Text, in)})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// convert 按配置转换换行分隔的列表（顺序对齐参考项目）。
func convert(text string, opt input) string {
	lineBreak := ""
	if opt.KeepLineBreaks {
		lineBreak = "\n"
	}

	s := text
	if opt.LowerCase {
		s = strings.ToLower(s)
	}

	parts := strings.Split(s, "\n")

	if opt.RemoveDuplicates {
		parts = unique(parts)
	}
	if opt.ReverseList {
		reverse(parts)
	}
	switch opt.SortList {
	case "asc":
		sort.Strings(parts)
	case "desc":
		sort.Sort(sort.Reverse(sort.StringSlice(parts)))
	}
	if opt.TrimItems {
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
	}

	// 移除空项。
	kept := parts[:0]
	for _, p := range parts {
		if p != "" {
			kept = append(kept, p)
		}
	}
	parts = kept

	for i := range parts {
		parts[i] = opt.ItemPrefix + parts[i] + opt.ItemSuffix
	}

	joined := strings.Join(parts, opt.Separator+lineBreak)
	return strings.Join([]string{opt.ListPrefix, joined, opt.ListSuffix}, lineBreak)
}

// unique 保留首次出现的元素（稳定去重）。
func unique(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := items[:0]
	for _, it := range items {
		if _, ok := seen[it]; ok {
			continue
		}
		seen[it] = struct{}{}
		out = append(out, it)
	}
	return out
}

// reverse 原地反转切片。
func reverse(items []string) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
