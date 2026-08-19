// Package textstatistics 实现文本统计（字符/单词/行数/字节大小）。
package textstatistics

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "text-statistics"
	Name        = "文本统计"
	Description = "统计文本的字符、单词、行数与字节大小"
	Category    = "文本"
	Icon        = "FileText"
)

// Keywords 为搜索关键词。
var Keywords = []string{"text", "statistics", "length", "characters", "count", "bytes", "文本", "统计", "字数", "字符", "行数", "字节"}

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
	Characters int    `json:"characters"` // 字符数（对齐 JS String.length，UTF-16 code unit 计数）
	Words      int    `json:"words"`      // 单词数（按空白切分）
	Lines      int    `json:"lines"`      // 行数（按 \n / \r\n / \r 切分）
	Bytes      int    `json:"bytes"`      // 字节数（UTF-8 编码）
	SizeText   string `json:"size_text"`  // 字节大小的人类可读格式（对齐 it-tools formatBytes）
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回统计结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	st := CountStats(in.Text)
	out, err := json.Marshal(output{
		Characters: st.Characters,
		Words:      st.Words,
		Lines:      st.Lines,
		Bytes:      st.Bytes,
		SizeText:   formatBytes(st.Bytes),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Stat 保存文本统计结果。
type Stat struct {
	Characters int
	Words      int
	Lines      int
	Bytes      int
}

// CountStats 统计输入文本的各项指标：
//   - Characters：对齐参考项目 JS `text.length` 的 UTF-16 code unit 计数（星面字符记 2）。
//   - Words：按 Unicode 空白切分词数（strings.Fields 语义，首尾多余空白不产生空词）。
//   - Lines：按 \r\n / \r / \n 切分行数；空文本为 0（对齐参考项目 split 语义）。
//   - Bytes：UTF-8 编码字节数（对齐 TextEncoder().encode 长度）。
func CountStats(text string) Stat {
	chars := 0
	for _, r := range []rune(text) {
		if r > 0xFFFF {
			chars += 2 // 代理对占 2 个 UTF-16 code unit
		} else {
			chars++
		}
	}

	words := 0
	if text != "" {
		words = len(strings.Fields(text))
	}

	lines := 0
	if text != "" {
		lines = len(lineBreakRe.Split(text, -1))
	}

	return Stat{
		Characters: chars,
		Words:      words,
		Lines:      lines,
		Bytes:      len([]byte(text)),
	}
}

// lineBreakRe 对齐参考项目 JS `/\r\n|\r|\n/` 的行分隔符。
var lineBreakRe = regexp.MustCompile(`\r\n|\r|\n`)

// sizes 对齐参考项目 formatBytes 的单位表。
var sizes = []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

// formatBytes 将字节数格式化为人类可读大小（1024 进制，最多 2 位小数去尾零）。
func formatBytes(bytes int) string {
	if bytes == 0 {
		return "0 Bytes"
	}
	i := int(math.Log(float64(bytes)) / math.Log(1024))
	if i >= len(sizes) {
		i = len(sizes) - 1
	}
	v := float64(bytes) / math.Pow(1024, float64(i))
	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	return s + " " + sizes[i]
}