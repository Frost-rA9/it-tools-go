// Package textdiff 实现文本比较（行级差异 + 行内精炼高亮）。
package textdiff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "text-diff"
	Name        = "文本比较"
	Description = "对比两段文本的差异，高亮新增与删除"
	Category    = "文本"
	Icon        = "FileDiff"
)

// Keywords 为搜索关键词。
var Keywords = []string{"text", "diff", "compare", "string", "差异", "对比", "比较", "文本差异"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// 输入限制（防滥用，保证 DP 内存可控）。
const (
	maxChars = 100000 // 每侧最大字符数（rune）
	maxLines = 3000   // 每侧最大行数（3000×3000 DP ≈ 36MB）
)

// input 是工具的输入结构。
type input struct {
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
}

// seg 是行内片段：Text 为文本，Changed 标记是否为差异部分。
type seg struct {
	Text    string `json:"t"`
	Changed bool   `json:"c"`
}

// row 是输出中的一行对比。
// Type: equal（两栏同内容）| delete（仅旧侧）| insert（仅新侧）。
type row struct {
	Type        string `json:"type"`
	OldNo       int    `json:"old_no,omitempty"` // 旧侧行号（从 1 开始，0 表示无）
	NewNo       int    `json:"new_no,omitempty"`
	Old         string `json:"old,omitempty"`          // equal 行的旧侧内容
	New         string `json:"new,omitempty"`          // equal 行的新侧内容
	OldSegments []seg  `json:"old_segments,omitempty"` // delete 行的行内片段
	NewSegments []seg  `json:"new_segments,omitempty"` // insert 行的行内片段
}

// stats 汇总差异统计。
type stats struct {
	OldLines int `json:"old_lines"`
	NewLines int `json:"new_lines"`
	Removed  int `json:"removed"`
	Added    int `json:"added"`
	Changed  int `json:"changed"` // 配对的修改行对数
}

// output 是工具的输出结构。
type output struct {
	Equal bool  `json:"equal"`
	Stats stats `json:"stats"`
	Rows  []row `json:"rows"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回行级差异结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if utf8.RuneCountInString(in.OldText) > maxChars || utf8.RuneCountInString(in.NewText) > maxChars {
		return "", fmt.Errorf("文本过长（每侧上限 %d 字符）", maxChars)
	}

	oldLines := splitLines(in.OldText)
	newLines := splitLines(in.NewText)
	if len(oldLines) > maxLines || len(newLines) > maxLines {
		return "", fmt.Errorf("文本行数过多（每侧上限 %d 行）", maxLines)
	}

	out, err := json.Marshal(build(oldLines, newLines))
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// splitLines 按 \n 拆行：忽略末尾多余尾行；行内容去掉 \r（容忍 CRLF）。
// 空文本返回空切片。
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	if len(lines) > 1 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for i, l := range lines {
		lines[i] = strings.TrimSuffix(l, "\r")
	}
	return lines
}

// build 由旧/新行组装对比行列表与统计。
// 连续的非相等操作组成一个变更块，块内删除/新增行按序配对做行内精炼；
// 多余的删除/新增不做行内精炼（整行高亮）。
func build(oldLines, newLines []string) output {
	ops := lcsDiff(oldLines, newLines)

	out := output{Stats: stats{OldLines: len(oldLines), NewLines: len(newLines)}}
	ai, bj := 0, 0 // 旧/新行游标
	for idx := 0; idx < len(ops); {
		op := ops[idx]
		if op.kind == '=' {
			out.Rows = append(out.Rows, row{
				Type: "equal", OldNo: ai + 1, NewNo: bj + 1,
				Old: oldLines[ai], New: newLines[bj],
			})
			ai++
			bj++
			idx++
			continue
		}

		// 收集变更块内的删除行与新增行（含各自行号）
		type line struct {
			text string
			no   int
		}
		var dels, ins []line
		for idx < len(ops) && ops[idx].kind != '=' {
			switch ops[idx].kind {
			case '-':
				dels = append(dels, line{oldLines[ai], ai + 1})
				ai++
			case '+':
				ins = append(ins, line{newLines[bj], bj + 1})
				bj++
			}
			idx++
		}

		out.Stats.Removed += len(dels)
		out.Stats.Added += len(ins)

		// 配对做行内精炼
		pairN := min(len(dels), len(ins))
		for k := 0; k < pairN; k++ {
			os, ns := splitInline(dels[k].text, ins[k].text)
			out.Rows = append(out.Rows,
				row{Type: "delete", OldNo: dels[k].no, OldSegments: os},
				row{Type: "insert", NewNo: ins[k].no, NewSegments: ns},
			)
			out.Stats.Changed++
		}
		// 多余的删除/新增：整行高亮
		for k := pairN; k < len(dels); k++ {
			out.Rows = append(out.Rows, row{Type: "delete", OldNo: dels[k].no,
				OldSegments: []seg{{Text: dels[k].text, Changed: true}}})
		}
		for k := pairN; k < len(ins); k++ {
			out.Rows = append(out.Rows, row{Type: "insert", NewNo: ins[k].no,
				NewSegments: []seg{{Text: ins[k].text, Changed: true}}})
		}
	}

	out.Equal = out.Stats.Removed == 0 && out.Stats.Added == 0
	return out
}

// splitInline 对配对行的内容做行内精炼：剥离公共前缀与公共后缀，
// 中间部分标记为差异片段（按 rune 处理，中文安全）。
// 返回旧行与新行的片段序列。
func splitInline(o, n string) (os, ns []seg) {
	or, nr := []rune(o), []rune(n)

	p := 0
	for p < len(or) && p < len(nr) && or[p] == nr[p] {
		p++
	}
	s := 0
	maxS := len(or) - p
	if len(nr)-p < maxS {
		maxS = len(nr) - p
	}
	for s < maxS && or[len(or)-1-s] == nr[len(nr)-1-s] {
		s++
	}

	oldMid := string(or[p : len(or)-s])
	newMid := string(nr[p : len(nr)-s])
	if oldMid == "" && newMid == "" {
		// 内容完全相同（防御分支）
		return []seg{{Text: o}}, []seg{{Text: n}}
	}

	if p > 0 {
		os = append(os, seg{Text: string(or[:p])})
		ns = append(ns, seg{Text: string(nr[:p])})
	}
	if oldMid != "" {
		os = append(os, seg{Text: oldMid, Changed: true})
	}
	if newMid != "" {
		ns = append(ns, seg{Text: newMid, Changed: true})
	}
	if s > 0 {
		os = append(os, seg{Text: string(or[len(or)-s:])})
		ns = append(ns, seg{Text: string(nr[len(nr)-s:])})
	}
	return os, ns
}