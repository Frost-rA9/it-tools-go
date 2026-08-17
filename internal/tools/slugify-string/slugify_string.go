// Package slugify 实现 Slug 化字符串工具：将文本转为 URL 友好的小写 slug。
// 算法对齐 @sindresorhus/slugify 默认行为（transliterate + decamelize + lowercase）。
package slugify

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "slugify-string"
	Name        = "Slug 化字符串"
	Description = "将字符串转为 URL 友好的小写 slug"
	Category    = "Web"
	Icon        = "LettersCase"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"slugify", "string", "escape", "emoji", "special", "character", "space", "trim", "slug", "打乱", "友好链接"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"` // 待处理文本
}

// output 是工具的输出结构。
type output struct {
	Slug string `json:"slug"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out := output{Slug: slugify(in.Text)}
	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// overridableReplacements 对齐 slugify 的 builtin overridable replacements。
var overridableReplacements = [][2]string{
	{"&", " and "},
	{"🦄", " unicorn "},
	{"♥", " love "},
}

// transliterateMap 覆盖无法通过 NFD 分解的常见非分解字符（其余拉丁重音由 NFD + 去音标处理）。
var transliterateMap = map[rune]string{
	'ß': "ss", 'æ': "ae", 'Æ': "ae", 'œ': "oe", 'Œ': "oe",
	'ø': "o", 'Ø': "o", 'đ': "d", 'Đ': "d", 'ð': "d", 'Ð': "d",
	'þ': "th", 'Þ': "th", 'ł': "l", 'Ł': "l", 'ı': "i",
	'ŋ': "n", 'Ŋ': "n", 'ħ': "h", 'Ħ': "h", 'ŧ': "t", 'Ŧ': "t",
	'ĸ': "k", 'ª': "a", 'º': "o", '©': "c", '®': "r", '™': "tm",
}

// decamelizePatterns 对齐 slugify 的 decamelize 四条正则（RE2 兼容）。
var decamelizePatterns = []struct{ re, repl string }{
	{`([A-Z]{2,})(\d+)`, `$1 $2`},
	{`([a-z\d]+)([A-Z]{2,})`, `$1 $2`},
	{`([a-z\d])([A-Z])`, `$1 $2`},
	{`([A-Z]+)([A-Z][a-rt-z\d]+)`, `$1 $2`},
}

// 预编译正则。
var (
	contractionRe = regexp.MustCompile(`([a-zA-Z\d]+)['\x{2019}]([ts])(\s|$)`)
	nonWordRe     = regexp.MustCompile(`[^a-z0-9]+`)
	collapseSepRe = regexp.MustCompile(`-{2,}`)
)

// slugify 对齐 @sindresorhus/slugify 默认行为。
func slugify(s string) string {
	s = transliterate(s)
	s = decamelize(s)
	s = strings.ToLower(s)
	s = contractionRe.ReplaceAllString(s, `$1$2$3`)
	s = nonWordRe.ReplaceAllString(s, "-")
	s = collapseSepRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// transliterate 音译：应用内置替换 → NFD 规范化 → 去组合音标 → 特殊字符映射。
func transliterate(s string) string {
	for _, rep := range overridableReplacements {
		s = strings.ReplaceAll(s, rep[0], rep[1])
	}
	var sb strings.Builder
	for _, r := range norm.NFD.String(s) {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if repl, ok := transliterateMap[r]; ok {
			sb.WriteString(repl)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// decamelize 分离驼峰单词。
func decamelize(s string) string {
	for _, p := range decamelizePatterns {
		s = regexp.MustCompile(p.re).ReplaceAllString(s, p.repl)
	}
	return s
}