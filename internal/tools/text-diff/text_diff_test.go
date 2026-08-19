package textdiff

import (
	"encoding/json"
	"strings"
	"testing"
)

// run 执行工具并解析输出。
func run(t *testing.T, oldText, newText string) output {
	t.Helper()
	in, err := json.Marshal(input{OldText: oldText, NewText: newText})
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}
	raw, err := (Executor{}).Execute(t.Context(), string(in))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(raw), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	return o
}

func TestIdenticalText(t *testing.T) {
	o := run(t, "a\nb\nc", "a\nb\nc")
	if !o.Equal {
		t.Error("相等文本应 equal=true")
	}
	if len(o.Rows) != 3 {
		t.Fatalf("期望 3 行，实际 %d", len(o.Rows))
	}
	for i, r := range o.Rows {
		if r.Type != "equal" || r.OldNo != i+1 || r.NewNo != i+1 {
			t.Errorf("第 %d 行应为 equal 且行号 %d，实际 %+v", i, i+1, r)
		}
	}
	if o.Stats.Removed != 0 || o.Stats.Added != 0 || o.Stats.Changed != 0 {
		t.Errorf("相等文本统计应为 0，实际 %+v", o.Stats)
	}
}

func TestAllInsert(t *testing.T) {
	o := run(t, "", "a\nb")
	if o.Equal {
		t.Error("全新增应 equal=false")
	}
	if len(o.Rows) != 2 {
		t.Fatalf("期望 2 行，实际 %d", len(o.Rows))
	}
	for i, r := range o.Rows {
		if r.Type != "insert" || r.NewNo != i+1 {
			t.Errorf("第 %d 行应为 insert，实际 %+v", i, r)
		}
	}
	if o.Stats.Added != 2 || o.Stats.Removed != 0 {
		t.Errorf("统计错误: %+v", o.Stats)
	}
}

func TestAllDelete(t *testing.T) {
	o := run(t, "a\nb", "")
	if o.Equal || len(o.Rows) != 2 {
		t.Fatalf("全部删除行为异常: equal=%v rows=%d", o.Equal, len(o.Rows))
	}
	if o.Rows[0].Type != "delete" || o.Rows[0].OldNo != 1 {
		t.Errorf("期望 delete 行号 1，实际 %+v", o.Rows[0])
	}
	if o.Stats.Removed != 2 {
		t.Errorf("统计错误: %+v", o.Stats)
	}
}

func TestInlineChange(t *testing.T) {
	// 单行替换："abcXdef" → "abcYdef"：公共前后缀剥离后中间 X/Y 为差异。
	o := run(t, "abcXdef", "abcYdef")
	if o.Equal || len(o.Rows) != 2 {
		t.Fatalf("期望删除+插入两行，实际 %d 行", len(o.Rows))
	}
	del, ins := o.Rows[0], o.Rows[1]
	if del.Type != "delete" || ins.Type != "insert" {
		t.Fatalf("类型错误: %s / %s", del.Type, ins.Type)
	}
	checkSegs(t, del.OldSegments, []seg{{"abc", false}, {"X", true}, {"def", false}})
	checkSegs(t, ins.NewSegments, []seg{{"abc", false}, {"Y", true}, {"def", false}})
}

func TestInlineSubstringInsert(t *testing.T) {
	// "hello world" → "hello brave world"：旧行是新的子串，仅新行中间为差异。
	o := run(t, "hello world", "hello brave world")
	if len(o.Rows) != 2 {
		t.Fatalf("期望 2 行，实际 %d", len(o.Rows))
	}
	del, ins := o.Rows[0], o.Rows[1]
	for _, s := range del.OldSegments {
		if s.Changed {
			t.Errorf("旧行不应有差异高亮片段: %+v", del.OldSegments)
		}
	}
	checkSegs(t, ins.NewSegments, []seg{{"hello ", false}, {"brave ", true}, {"world", false}})
	if o.Stats.Changed != 1 {
		t.Errorf("期望 1 处修改，实际 %d", o.Stats.Changed)
	}
}

func TestMixedBlocks(t *testing.T) {
	// a → a'（修改）、删 c、新增 x：块隔离与顺序。
	o := run(t, "a\nc\nd", "a'\nx\nd")
	// 期望：delete(a)、insert(a')、delete(c)、insert(x)、equal(d)
	types := []string{}
	for _, r := range o.Rows {
		types = append(types, r.Type)
	}
	want := []string{"delete", "insert", "delete", "insert", "equal"}
	if strings.Join(types, ",") != strings.Join(want, ",") {
		t.Fatalf("行类型序列错误: %v（期望 %v）", types, want)
	}
	if o.Stats.Removed != 2 || o.Stats.Added != 2 || o.Stats.Changed != 2 {
		t.Errorf("统计错误: %+v", o.Stats)
	}
}

func TestCRLF(t *testing.T) {
	o := run(t, "a\r\nb\r\n", "a\nb")
	if !o.Equal {
		t.Error("CRLF 与 LF 应视为相同")
	}
}

func TestEmptyVsNonEmpty(t *testing.T) {
	o := run(t, "", "x")
	if o.Equal || len(o.Rows) != 1 || o.Rows[0].Type != "insert" {
		t.Fatalf("空 vs 非空行为异常: %+v", o)
	}
}

func TestBothEmpty(t *testing.T) {
	o := run(t, "", "")
	if !o.Equal || len(o.Rows) != 0 {
		t.Fatalf("空 vs 空应 equal=true 且无行: %+v", o)
	}
}

func TestTooLong(t *testing.T) {
	in, _ := json.Marshal(input{OldText: strings.Repeat("a", maxChars+1)})
	if _, err := (Executor{}).Execute(t.Context(), string(in)); err == nil {
		t.Error("超长字符应报错")
	}
	inLine, _ := json.Marshal(input{OldText: strings.Repeat("x\n", maxLines+1)})
	if _, err := (Executor{}).Execute(t.Context(), string(inLine)); err == nil {
		t.Error("超行数应报错")
	}
}

// checkSegs 校验片段序列。
func checkSegs(t *testing.T, got []seg, want []seg) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("片段数不匹配: 实际 %+v，期望 %+v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("片段 %d 不匹配: 实际 %+v，期望 %+v", i, got[i], want[i])
		}
	}
}