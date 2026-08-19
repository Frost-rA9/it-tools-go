package benchmarkbuilder

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"it-tools-go/internal/registry"
)

func TestComputeAverage(t *testing.T) {
	if v := ComputeAverage(nil); v != 0 {
		t.Errorf("空数组均值 = %v, want 0", v)
	}
	if v := ComputeAverage([]float64{5, 10}); v != 7.5 {
		t.Errorf("均值 = %v, want 7.5", v)
	}
	if v := ComputeAverage([]float64{2, 4, 6}); v != 4 {
		t.Errorf("均值 = %v, want 4", v)
	}
}

func TestComputeVariance(t *testing.T) {
	if v := ComputeVariance(nil, 0); v != 0 {
		t.Errorf("空数组方差 = %v, want 0", v)
	}
	// [5,10] 均值 7.5，平方差 (2.5²+2.5²)/2 = 6.25
	if v := ComputeVariance([]float64{5, 10}, 7.5); v != 6.25 {
		t.Errorf("方差 = %v, want 6.25", v)
	}
	// [8,12,10] 均值 10，平方差 (4+4+0)/3 = 8/3
	if v := ComputeVariance([]float64{8, 12, 10}, 10); v != 8.0/3 {
		t.Errorf("方差 = %v, want 8/3", v)
	}
}

func TestBuildResultsSortAndCompare(t *testing.T) {
	// 参考示例：Suite1 [5,10] 均值 7.5 / Suite2 [8,12] 均值 10
	res := BuildResults([]suite{
		{Title: "Suite 1", Data: []float64{5, 10}},
		{Title: "Suite 2", Data: []float64{8, 12}},
	}, "ms")
	if len(res) != 2 {
		t.Fatalf("结果数 = %d, want 2", len(res))
	}
	// 排序：Suite1(7.5) 第一
	if res[0].Title != "Suite 1" || res[0].Position != 1 {
		t.Errorf("第一名 = %+v", res[0])
	}
	if res[1].Position != 2 {
		t.Errorf("第二名 position = %d, want 2", res[1].Position)
	}
	// 第一名不含对比（无 (+delta）
	if strings.Contains(res[0].MeanDisplay, "(+") {
		t.Errorf("第一名 meanDisplay 不应含对比: %q", res[0].MeanDisplay)
	}
	// 第二名含对比：delta=2.5, ratio=10/7.5=1.333
	second := res[1]
	if second.MeanDisplay != "10ms (+2.5ms ; x1.333)" {
		t.Errorf("第二名 meanDisplay = %q", second.MeanDisplay)
	}
	if second.VarianceDisplay != "4ms²" {
		t.Errorf("第二名 varianceDisplay = %q", second.VarianceDisplay)
	}
}

func TestBuildResultsNoUnit(t *testing.T) {
	res := BuildResults([]suite{
		{Title: "A", Data: []float64{10}},
		{Title: "B", Data: []float64{20}},
	}, "")
	if !strings.Contains(res[1].MeanDisplay, "(+10 ; x2)") {
		t.Errorf("无单位时对比格式 = %q", res[1].MeanDisplay)
	}
	// 无单位时不带 ²
	if res[1].VarianceDisplay != "0" {
		t.Errorf("无单位方差 = %q, want 0", res[1].VarianceDisplay)
	}
}

func TestBuildResultsBestMeanZero(t *testing.T) {
	// 最佳均值为 0 时，比率显示 ∞
	res := BuildResults([]suite{
		{Title: "Zero", Data: []float64{0}},
		{Title: "Five", Data: []float64{5}},
	}, "s")
	if res[0].Title != "Zero" {
		t.Fatalf("第一名 = %+v", res[0])
	}
	if !strings.Contains(res[1].MeanDisplay, " x∞") {
		t.Errorf("均值为 0 时比率应为 ∞: %q", res[1].MeanDisplay)
	}
}

func TestBuildResultsEmpty(t *testing.T) {
	res := BuildResults(nil, "ms")
	if len(res) != 0 {
		t.Errorf("空输入结果数 = %d, want 0", len(res))
	}
}

func TestMarkdown(t *testing.T) {
	res := BuildResults([]suite{
		{Title: "Suite 1", Data: []float64{5, 10}},
		{Title: "Suite 2", Data: []float64{8, 12}},
	}, "ms")
	md := ArrayToMarkdownTable(res, "ms")
	wantLines := []string{
		"| Position | Suite | Mean | Variance | Samples |",
		"| --- | --- | --- | --- | --- |",
		"| 1 | Suite 1 | 7.5ms | 6.25ms² | 2 |",
		"| 2 | Suite 2 | 10ms (+2.5ms ; x1.333) | 4ms² | 2 |",
	}
	got := strings.Split(md, "\n")
	if len(got) != len(wantLines) {
		t.Fatalf("markdown 行数 = %d, want %d\n%q", len(got), len(wantLines), md)
	}
	for i := range wantLines {
		if got[i] != wantLines[i] {
			t.Errorf("第 %d 行 = %q, want %q", i+1, got[i], wantLines[i])
		}
	}
}

func TestMarkdownEmpty(t *testing.T) {
	if md := ArrayToMarkdownTable(nil, "ms"); md != "" {
		t.Errorf("空结果 markdown = %q, want empty", md)
	}
}

func TestBulletList(t *testing.T) {
	res := BuildResults([]suite{
		{Title: "Suite 2", Data: []float64{8, 12}},
		{Title: "Suite 1", Data: []float64{5, 10}},
	}, "ms")
	bl := ArrayToBulletList(res)
	want := " - Suite 1\n    - Position: 1\n    - Mean: 7.5ms\n    - Variance: 6.25ms²\n    - Samples: 2\n - Suite 2\n    - Position: 2\n    - Mean: 10ms (+2.5ms ; x1.333)\n    - Variance: 4ms²\n    - Samples: 2"
	if bl != want {
		t.Errorf("bulletList =\n%q\nwant\n%q", bl, want)
	}
}

func TestExecute(t *testing.T) {
	outStr, err := (Executor{}).Execute(context.Background(),
		`{"unit":"ms","suites":[{"title":"Suite 1","data":[5,10]},{"title":"Suite 2","data":[8,12]}]}`)
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outStr), &out); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(out.Results) != 2 || out.Markdown == "" || len(out.BulletList) == 0 {
		t.Errorf("输出不完整: %+v", out)
	}
}

func TestToolMeta(t *testing.T) {
	meta := Tool()
	if meta.ID != ID || meta.Category != registry.CategoryMeasurement || meta.Icon != "Gauge" {
		t.Errorf("元数据不符: %+v", meta)
	}
}
