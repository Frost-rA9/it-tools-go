package loremipsum

import (
	"encoding/json"
	"strings"
	"testing"
)

// marshalInput 将 input 结构编码为 JSON 字符串。
func marshalInput(t *testing.T, in input) string {
	t.Helper()
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}
	return string(raw)
}

// vocabSet 词汇表集合（小写），供词表校验使用。
var vocabSet = func() map[string]bool {
	set := make(map[string]bool, len(vocabulary))
	for _, w := range vocabulary {
		set[w] = true
	}
	return set
}()

// TestExecuteStructure 校验段落数、句数、词数结构与词汇表合法性（固定范围保证确定性）。
func TestExecuteStructure(t *testing.T) {
	exec := Executor{}
	// 2 段 × 3 句 × 10 词，关闭首句固定，可精确拆分断言。
	out, err := exec.Execute(t.Context(), marshalInput(t, input{
		Paragraphs:  2,
		SentenceMin: 3,
		SentenceMax: 3,
		WordMin:     10,
		WordMax:     10,
	}))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}

	paragraphs := strings.Split(o.Text, "\n\n")
	if len(paragraphs) != 2 {
		t.Fatalf("期望 2 段，实际 %d", len(paragraphs))
	}
	for i, para := range paragraphs {
		tokens := strings.Split(para, " ")
		wantTokens := 3 * 10 // 句数 × 词数
		if len(tokens) != wantTokens {
			t.Fatalf("段落 %d 期望 %d 个词，实际 %d", i, wantTokens, len(tokens))
		}
		for j, tok := range tokens {
			stripped := strings.TrimSuffix(tok, ".")
			if !vocabSet[strings.ToLower(stripped)] {
				t.Errorf("段落 %d 词 %d 不在词汇表: %q", i, j, tok)
			}
			if j%10 == 9 {
				if !strings.HasSuffix(tok, ".") {
					t.Errorf("段落 %d 句末词 %q 缺少句号", i, tok)
				}
			} else if strings.Contains(tok, ".") {
				t.Errorf("段落 %d 句内词 %q 含多余标点", i, tok)
			}
			if j%10 == 0 {
				if tok == "" || strings.ToLower(tok[:1]) == tok[:1] {
					t.Errorf("段落 %d 句首词 %q 未大写", i, tok)
				}
			}
		}
	}
}

// TestStartWithLoremIpsum 校验首句固定开关的两个方向。
func TestStartWithLoremIpsum(t *testing.T) {
	exec := Executor{}

	run := func(start bool) string {
		out, err := exec.Execute(t.Context(), marshalInput(t, input{
			Paragraphs:          2,
			SentenceMin:         2,
			SentenceMax:         2,
			WordMin:             5,
			WordMax:             5,
			StartWithLoremIpsum: start,
		}))
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}
		var o output
		if err := json.Unmarshal([]byte(out), &o); err != nil {
			t.Fatalf("解析输出失败: %v", err)
		}
		return o.Text
	}

	if text := run(true); !strings.HasPrefix(text, firstSentence+" ") {
		t.Errorf("开启首句固定时期望以固定句开头，实际: %q", text)
	}
	if text := run(false); strings.HasPrefix(text, firstSentence) {
		t.Errorf("关闭首句固定时不应以固定句开头，实际: %q", text)
	}
}

// TestExecuteAsHTML 校验 HTML 模式下 <p> 包裹与段落连接。
func TestExecuteAsHTML(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), marshalInput(t, input{
		Paragraphs:  2,
		SentenceMin: 2,
		SentenceMax: 2,
		WordMin:     4,
		WordMax:     4,
		AsHTML:      true,
	}))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if !strings.HasPrefix(o.Text, "<p>") || !strings.HasSuffix(o.Text, "</p>") {
		t.Errorf("HTML 模式应整体用 <p>…</p> 包裹，实际: %q", o.Text)
	}
	if n := strings.Count(o.Text, "<p>"); n != 2 {
		t.Errorf("期望 2 个 <p>，实际 %d", n)
	}
	if n := strings.Count(o.Text, "</p>"); n != 2 {
		t.Errorf("期望 2 个 </p>，实际 %d", n)
	}
	if !strings.Contains(o.Text, "</p>\n\n<p>") {
		t.Errorf("段落间应以 </p>\n\n<p> 连接，实际: %q", o.Text)
	}
}

// TestGenerateRandomness 随机性冒烟：连续生成应各不相同。
func TestGenerateRandomness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 20; i++ {
		text := Generate(input{
			Paragraphs:          3,
			SentenceMin:         3,
			SentenceMax:         8,
			WordMin:             8,
			WordMax:             15,
			StartWithLoremIpsum: true,
		})
		if seen[text] {
			t.Fatalf("第 %d 次生成出现重复文本（随机性异常）", i)
		}
		seen[text] = true
	}
}

// TestExecuteErrors 校验非法参数报错。
func TestExecuteErrors(t *testing.T) {
	exec := Executor{}
	tests := []struct {
		name  string
		input input
	}{
		{"段落数 0", input{Paragraphs: 0, SentenceMin: 1, SentenceMax: 1, WordMin: 1, WordMax: 1}},
		{"段落数超上限", input{Paragraphs: 21, SentenceMin: 1, SentenceMax: 1, WordMin: 1, WordMax: 1}},
		{"句数下限 0", input{Paragraphs: 1, SentenceMin: 0, SentenceMax: 1, WordMin: 1, WordMax: 1}},
		{"句数超上限", input{Paragraphs: 1, SentenceMin: 1, SentenceMax: 51, WordMin: 1, WordMax: 1}},
		{"句数 min 大于 max", input{Paragraphs: 1, SentenceMin: 5, SentenceMax: 3, WordMin: 1, WordMax: 1}},
		{"词数下限 0", input{Paragraphs: 1, SentenceMin: 1, SentenceMax: 1, WordMin: 0, WordMax: 1}},
		{"词数超上限", input{Paragraphs: 1, SentenceMin: 1, SentenceMax: 1, WordMin: 1, WordMax: 51}},
		{"词数 min 大于 max", input{Paragraphs: 1, SentenceMin: 1, SentenceMax: 1, WordMin: 4, WordMax: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := exec.Execute(t.Context(), marshalInput(t, tt.input)); err == nil {
				t.Fatalf("期望报错，实际成功")
			}
		})
	}
}