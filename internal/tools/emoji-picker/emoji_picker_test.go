package emojipicker

import (
	"encoding/json"
	"strings"
	"testing"
)

// knownGroups 数据文件中的全部分组（对齐 unicode-emoji-json 0.9.0）。
var knownGroups = []string{
	"Smileys & Emotion", "People & Body", "Animals & Nature", "Food & Drink",
	"Travel & Places", "Activities", "Objects", "Symbols", "Flags",
}

func TestDataIntegrity(t *testing.T) {
	if len(entries) < 1800 {
		t.Fatalf("条目数过少: %d", len(entries))
	}
	seen := make(map[string]bool, len(entries))
	groups := make(map[string]bool)
	for _, e := range entries {
		if seen[e.Emoji] {
			t.Errorf("重复 emoji: %q", e.Emoji)
		}
		seen[e.Emoji] = true
		if e.Name == "" || e.Group == "" {
			t.Errorf("字段缺失: %+v", e)
		}
		groups[e.Group] = true
	}
	if len(groups) != len(knownGroups) {
		t.Errorf("分组数异常: %d（期望 %d）", len(groups), len(knownGroups))
	}
	for _, g := range knownGroups {
		if !groups[g] {
			t.Errorf("缺少分组: %s", g)
		}
	}
}

func TestSearchName(t *testing.T) {
	res := Search("Grinning face")
	if len(res) == 0 {
		t.Fatal("搜索 'Grinning face' 无结果")
	}
	if res[0].Emoji != "😀" {
		t.Errorf("名称完全匹配应置顶，实际第一条 %q", res[0].Emoji)
	}
}

func TestSearchKeywordOnly(t *testing.T) {
	// "welsh" 仅出现在 Flag Wales 的关键词中（名称/分组不含）。
	res := Search("welsh")
	if len(res) == 0 {
		t.Fatal("关键词搜索 'welsh' 无结果")
	}
	if res[0].Name != "Flag Wales" {
		t.Errorf("纯关键词命中，期望 Flag Wales，实际 %q", res[0].Name)
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	upper := Search("SMILE")
	lower := Search("smile")
	if len(upper) == 0 || len(upper) != len(lower) {
		t.Fatalf("大小写结果不一致: %d vs %d", len(upper), len(lower))
	}
	for i := range upper {
		if upper[i].Emoji != lower[i].Emoji {
			t.Fatalf("第 %d 条大小写结果不同: %q vs %q", i, upper[i].Emoji, lower[i].Emoji)
		}
	}
}

func TestSearchEmptyReturnsAll(t *testing.T) {
	res := Search("")
	if len(res) != len(entries) {
		t.Errorf("空 query 应返回全量: %d vs %d", len(res), len(entries))
	}
	// 全量按分组顺序：首条应为 Smileys & Emotion 分组，末条应为 Flags。
	if res[0].Group != "Smileys & Emotion" {
		t.Errorf("首条分组异常: %q", res[0].Group)
	}
	if res[len(res)-1].Group != "Flags" {
		t.Errorf("末条分组异常: %q", res[len(res)-1].Group)
	}
}

func TestSearchNoResult(t *testing.T) {
	if res := Search("zzzqqq"); len(res) != 0 {
		t.Errorf("期望无结果，实际 %d 条", len(res))
	}
}

func TestUnicodeAux(t *testing.T) {
	tests := []struct {
		emoji      string
		codePoints string
		unicode    string
	}{
		{"😀", "0x1f600", `\ud83d\ude00`},
		{"A", "0x41", `\u0041`},
		{"中", "0x4e2d", `\u4e2d`},
		{"👨‍👩‍👧", "0x1f468", `\ud83d\udc68\u200d\ud83d\udc69\u200d\ud83d\udc67`},
	}
	for _, tt := range tests {
		if got := codePoints(tt.emoji); got != tt.codePoints {
			t.Errorf("codePoints(%s) = %q，期望 %q", tt.emoji, got, tt.codePoints)
		}
		if got := unicodeEscape(tt.emoji); got != tt.unicode {
			t.Errorf("unicodeEscape(%s) = %q，期望 %q", tt.emoji, got, tt.unicode)
		}
	}
}

func TestExecuteQueryTooLong(t *testing.T) {
	exec := Executor{}
	in, _ := json.Marshal(input{Query: strings.Repeat("a", 65)})
	if _, err := exec.Execute(t.Context(), string(in)); err == nil {
		t.Fatal("超长 query 应报错")
	}
	// 边界：64 字符允许
	in64, _ := json.Marshal(input{Query: strings.Repeat("a", 64)})
	if _, err := exec.Execute(t.Context(), string(in64)); err != nil {
		t.Fatalf("64 字符 query 不应报错: %v", err)
	}
}