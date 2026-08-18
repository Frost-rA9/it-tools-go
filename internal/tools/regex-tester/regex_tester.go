// Package regextester 实现正则表达式测试器：对输入文本执行正则匹配并返回匹配详情。
package regextester

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "regex-tester"
	Name        = "Regex 测试器"
	Description = "用示例文本测试正则表达式，展示匹配与捕获组"
	Category    = "开发"
	Icon        = "Language"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"regex", "tester", "regular", "expression", "正则", "匹配", "捕获"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// Flags 是正则匹配选项。
type Flags struct {
	G bool `json:"g"` // 全局（返回全部匹配；false 仅第一个）
	I bool `json:"i"` // 忽略大小写
	M bool `json:"m"` // 多行（^ $ 匹配行边界）
	S bool `json:"s"` // 单行（. 匹配换行）
}

// input 是工具的输入结构。
type input struct {
	Regex string `json:"regex"`
	Text  string `json:"text"`
	Flags Flags  `json:"flags"`
}

// GroupCapture 描述一个捕获组（编号或命名）。
type GroupCapture struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// MatchResult 描述一次匹配。
type MatchResult struct {
	Index    int            `json:"index"`
	Value    string         `json:"value"`
	Captures []GroupCapture `json:"captures"` // 编号捕获组（1..n）
	Groups   []GroupCapture `json:"groups"`   // 命名捕获组
}

// output 是工具的输出结构。
type output struct {
	Matches []MatchResult `json:"matches"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	matches, err := Convert(in.Regex, in.Text, in.Flags)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Matches: matches})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Convert 对 text 执行正则匹配。flags.g=false 时仅返回第一个匹配（对齐 it-tools）。
func Convert(regex, text string, flags Flags) ([]MatchResult, error) {
	if regex == "" || text == "" {
		return nil, nil
	}

	pattern := regex
	var sb strings.Builder
	if flags.I {
		sb.WriteString("(?i)")
	}
	if flags.M {
		sb.WriteString("(?m)")
	}
	if flags.S {
		sb.WriteString("(?s)")
	}
	sb.WriteString(pattern)

	re, err := regexp.Compile(sb.String())
	if err != nil {
		return nil, fmt.Errorf("无效的正则表达式: %w", err)
	}

	n := -1
	if !flags.G {
		n = 1
	}
	locs := re.FindAllStringSubmatchIndex(text, n)
	names := re.SubexpNames() // [0] 为整体匹配，名称空。

	results := make([]MatchResult, 0, len(locs))
	for _, loc := range locs {
		m := MatchResult{
			Index: loc[0],
			Value: text[loc[0]:loc[1]],
		}
		for i := 1; i < len(names); i++ {
			s, e := loc[2*i], loc[2*i+1]
			if s < 0 || e < 0 {
				continue // 该组未参与本次匹配
			}
			gc := GroupCapture{Name: strconv.Itoa(i), Value: text[s:e], Start: s, End: e}
			if names[i] != "" {
				gc.Name = names[i]
				m.Groups = append(m.Groups, gc)
			} else {
				m.Captures = append(m.Captures, gc)
			}
		}
		results = append(results, m)
	}
	return results, nil
}
