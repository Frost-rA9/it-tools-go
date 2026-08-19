// Package loremipsum 实现 Lorem ipsum 占位文本生成器。
package loremipsum

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"unicode"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "lorem-ipsum-generator"
	Name        = "Lorem ipsum 生成器"
	Description = "生成 Lorem ipsum 占位文本，可自定义段落与句子数量"
	Category    = "文本"
	Icon        = "AlignJustified"
)

// Keywords 为搜索关键词。
var Keywords = []string{"lorem", "ipsum", "placeholder", "filler", "random", "generator", "占位文本", "占位符", "随机文本"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// 参数上限（对齐参考项目 it-tools 的滑块范围）。
const (
	maxParagraphs    = 20
	maxSentenceCount = 50
	maxWordCount     = 50
)

// input 是工具的输入结构。
type input struct {
	Paragraphs          int  `json:"paragraphs"`
	SentenceMin         int  `json:"sentence_min"`
	SentenceMax         int  `json:"sentence_max"`
	WordMin             int  `json:"word_min"`
	WordMax             int  `json:"word_max"`
	StartWithLoremIpsum bool `json:"start_with_lorem_ipsum"`
	AsHTML              bool `json:"as_html"`
}

// output 是工具的输出结构。
type output struct {
	Text string `json:"text"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回生成的占位文本。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if err := validate(in); err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Text: Generate(in)})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// validate 校验输入参数的范围与合法性。
func validate(in input) error {
	if in.Paragraphs < 1 || in.Paragraphs > maxParagraphs {
		return fmt.Errorf("段落数必须在 1 到 %d 之间", maxParagraphs)
	}
	if in.SentenceMin < 1 || in.SentenceMax > maxSentenceCount || in.SentenceMin > in.SentenceMax {
		return fmt.Errorf("每段句数范围必须在 1 到 %d 之间且最小值不大于最大值", maxSentenceCount)
	}
	if in.WordMin < 1 || in.WordMax > maxWordCount || in.WordMin > in.WordMax {
		return fmt.Errorf("每句词数范围必须在 1 到 %d 之间且最小值不大于最大值", maxWordCount)
	}
	return nil
}

// firstSentence 是经典 Lorem ipsum 开篇句（对齐参考项目）。
const firstSentence = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."

// Generate 按参数生成占位文本，每次调用结果随机：
// 每段句数在 [SentenceMin, SentenceMax] 内随机，每句词数在 [WordMin, WordMax] 内随机。
func Generate(in input) string {
	paragraphs := make([]string, 0, in.Paragraphs)
	for p := 0; p < in.Paragraphs; p++ {
		sentenceCount := rand.IntN(in.SentenceMax-in.SentenceMin+1) + in.SentenceMin
		sentences := make([]string, 0, sentenceCount)
		for s := 0; s < sentenceCount; s++ {
			wordCount := rand.IntN(in.WordMax-in.WordMin+1) + in.WordMin
			sentences = append(sentences, generateSentence(wordCount))
		}
		if in.StartWithLoremIpsum && p == 0 {
			sentences[0] = firstSentence
		}
		paragraphs = append(paragraphs, strings.Join(sentences, " "))
	}

	if in.AsHTML {
		return "<p>" + strings.Join(paragraphs, "</p>\n\n<p>") + "</p>"
	}
	return strings.Join(paragraphs, "\n\n")
}

// generateSentence 生成一句：从词汇表随机取 wordCount 个词，首字母大写、句号结尾。
func generateSentence(wordCount int) string {
	words := make([]string, 0, wordCount)
	for i := 0; i < wordCount; i++ {
		words = append(words, vocabulary[rand.IntN(len(vocabulary))])
	}
	sentence := strings.Join(words, " ")
	rs := []rune(sentence)
	rs[0] = unicode.ToUpper(rs[0])
	return string(rs) + "."
}

// vocabulary 拉丁占位词表（对齐参考项目 it-tools 的 lorem-ipsum-generator.service.ts）。
var vocabulary = []string{
	"a", "ac", "accumsan", "ad", "adipiscing", "aenean", "aliquam", "aliquet", "amet", "ante",
	"aptent", "arcu", "at", "auctor", "bibendum", "blandit", "class", "commodo", "condimentum", "congue",
	"consectetur", "consequat", "conubia", "convallis", "cras", "cubilia", "cum", "curabitur", "curae", "dapibus",
	"diam", "dictum", "dictumst", "dignissim", "dolor", "donec", "dui", "duis", "egestas", "eget",
	"eleifend", "elementum", "elit", "enim", "erat", "eros", "est", "et", "etiam", "eu",
	"euismod", "facilisi", "faucibus", "felis", "fermentum", "feugiat", "fringilla", "fusce", "gravida", "habitant",
	"habitasse", "hac", "hendrerit", "himenaeos", "iaculis", "id", "imperdiet", "in", "inceptos", "integer",
	"interdum", "ipsum", "justo", "lacinia", "lacus", "laoreet", "lectus", "leo", "ligula", "litora",
	"lobortis", "lorem", "luctus", "maecenas", "magna", "magnis", "malesuada", "massa", "mattis", "mauris",
	"metus", "mi", "molestie", "mollis", "montes", "morbi", "mus", "nam", "nascetur", "natoque",
	"nec", "neque", "netus", "nisi", "nisl", "non", "nostra", "nulla", "nullam", "nunc",
	"odio", "orci", "ornare", "parturient", "pellentesque", "penatibus", "per", "pharetra", "phasellus", "placerat",
	"platea", "porta", "porttitor", "posuere", "potenti", "praesent", "pretium", "primis", "proin", "pulvinar",
	"purus", "quam", "quis", "quisque", "rhoncus", "ridiculus", "risus", "rutrum", "sagittis", "sapien",
	"scelerisque", "sed", "sem", "semper", "senectus", "sit", "sociis", "sociosqu", "sodales", "sollicitudin",
	"suscipit", "suspendisse", "taciti", "tellus", "tempor", "tempus", "tincidunt", "torquent", "tortor", "turpis",
	"ullamcorper", "ultrices", "ultricies", "urna", "varius", "vehicula", "vel", "velit", "venenatis", "vestibulum",
	"vitae", "vivamus", "viverra", "volutpat", "vulputate",
}