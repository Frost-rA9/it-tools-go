// Package emojipicker 实现 Emoji 选择器（浏览/搜索/复制）。
package emojipicker

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "emoji-picker"
	Name        = "Emoji 选择器"
	Description = "浏览、搜索并复制 Emoji 及其 Unicode 表示"
	Category    = "文本"
	Icon        = "MoodSmile"
)

// Keywords 为搜索关键词。
var Keywords = []string{"emoji", "picker", "unicode", "copy", "paste", "表情", "选择器", "复制"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// 内嵌精简数据（scripts/gen-emojis.mjs 生成，1914 条，9 个分组）。
//
//go:embed data/emojis.json
var emojiData []byte

// emojiEntry 是内嵌数据中的单条 emoji。
type emojiEntry struct {
	Emoji    string   `json:"emoji"`
	Name     string   `json:"name"`
	Group    string   `json:"group"`
	Keywords []string `json:"keywords"`
}

// entries 是启动时解析的全部 emoji（按分组与组内顺序）。
var entries []emojiEntry

func init() {
	if err := json.Unmarshal(emojiData, &entries); err != nil {
		panic(fmt.Sprintf("emoji-picker: 解析内嵌数据失败: %v", err))
	}
}

// input 是工具的输入结构。
type input struct {
	Query string `json:"query"`
}

// emojiOut 是输出中的单条 emoji（含运行时推导的码点与 Unicode 转义）。
type emojiOut struct {
	Emoji      string   `json:"emoji"`
	Name       string   `json:"name"`
	Group      string   `json:"group"`
	Keywords   []string `json:"keywords"`
	Codepoints string   `json:"codepoints"`
	Unicode    string   `json:"unicode"`
}

// output 是工具的输出结构。
type output struct {
	Total  int        `json:"total"`
	Emojis []emojiOut `json:"emojis"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回 emoji 列表：query 为空返回全量（分组序），否则返回加权搜索结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if len(in.Query) > 64 {
		return "", fmt.Errorf("搜索关键词过长（上限 64 字符）")
	}

	res := Search(in.Query)
	out, err := json.Marshal(output{Total: len(res), Emojis: res})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Search 返回匹配 query 的 emoji（含运行时推导字段）。
// 空 query 返回全部（保持数据中的分组顺序）；否则按多字段加权评分：
// 名称命中 3 分、分组命中 2 分、任一关键词命中 1 分，按分数降序、同分保序。
func Search(query string) []emojiOut {
	if query == "" {
		out := make([]emojiOut, 0, len(entries))
		for _, e := range entries {
			out = append(out, toOut(e))
		}
		return out
	}

	lower := strings.ToLower(query)
	type hit struct {
		entry emojiEntry
		score int
	}
	hits := make([]hit, 0)
	for _, e := range entries {
		sc := 0
		if strings.Contains(strings.ToLower(e.Name), lower) {
			sc += 3
		}
		if strings.Contains(strings.ToLower(e.Group), lower) {
			sc += 2
		}
		for _, k := range e.Keywords {
			if strings.Contains(strings.ToLower(k), lower) {
				sc++
				break
			}
		}
		if sc > 0 {
			hits = append(hits, hit{e, sc})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].score > hits[j].score })

	out := make([]emojiOut, 0, len(hits))
	for _, h := range hits {
		out = append(out, toOut(h.entry))
	}
	return out
}

// toOut 由数据条目构造输出（codepoints/unicode 在运行时推导）。
func toOut(e emojiEntry) emojiOut {
	return emojiOut{
		Emoji:      e.Emoji,
		Name:       e.Name,
		Group:      e.Group,
		Keywords:   e.Keywords,
		Codepoints: codePoints(e.Emoji),
		Unicode:    unicodeEscape(e.Emoji),
	}
}

// codePoints 返回首码点的十六进制表示（对齐参考 `0x${codePointAt(0).toString(16)}`）。
func codePoints(s string) string {
	for _, r := range s {
		return fmt.Sprintf("0x%x", r)
	}
	return ""
}

// unicodeEscape 返回 UTF-16 code unit 转义序列（对齐参考 escapeUnicode）：
// 😀 → \ud83d\ude00；星面字符编码为代理对，小写十六进制、4 位补齐。
func unicodeEscape(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r > 0xFFFF {
			r1 := 0xD800 + ((r - 0x10000) >> 10)
			r2 := 0xDC00 + ((r - 0x10000) & 0x3FF)
			fmt.Fprintf(&sb, `\u%04x\u%04x`, r1, r2)
		} else {
			fmt.Fprintf(&sb, `\u%04x`, r)
		}
	}
	return sb.String()
}