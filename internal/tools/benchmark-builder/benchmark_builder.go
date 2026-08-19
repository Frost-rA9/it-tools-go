// Package benchmarkbuilder 实现基准测试建构器：对多组测量值（suite）计算
// 均值/方差，与最佳结果对比，并生成 Markdown 表格或列表便于分享。
//
// 统计口径与参考项目 it-tools（benchmark-builder.models.ts）保持一致：
// 均值为算术平均，方差为平方差均值（总体方差），比率/差值按 3 位小数四舍五入。
package benchmarkbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "benchmark-builder"
	Name        = "基准测试建构器"
	Description = "对比多组测量数据，计算均值与方差并生成可分享的结果"
	Category    = registry.CategoryMeasurement
	Icon        = "Gauge"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"benchmark", "builder", "execution", "duration", "mean", "variance", "suite", "基准", "测试", "平均", "方差", "性能"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// suite 是一组命名测量值。
type suite struct {
	Title string    `json:"title"`
	Data  []float64 `json:"data"`
}

// input 是工具的输入结构。
type input struct {
	Unit   string  `json:"unit"`   // 结果单位（如 ms），可空
	Suites []suite `json:"suites"` // 待对比的测量组
}

// result 是单个 suite 的处理结果。
type result struct {
	Position        int     `json:"position"`        // 排名（按均值升序，从 1 开始）
	Title           string  `json:"title"`           // 套件名
	Size            int     `json:"size"`            // 有效样本数
	Mean            float64 `json:"mean"`            // 均值（原始值）
	Variance        float64 `json:"variance"`        // 总体方差（原始值）
	MeanDisplay     string  `json:"meanDisplay"`     // 展示用均值（含单位与对比）
	VarianceDisplay string  `json:"varianceDisplay"` // 展示用方差（含单位）
}

// output 是工具的输出结构。
type output struct {
	Results    []result `json:"results"`    // 按均值升序排列的结果
	Markdown   string   `json:"markdown"`   // Markdown 表格（供复制）
	BulletList string   `json:"bulletList"` // 缩进列表（供复制）
}

// 结果展示列与生成文本用的英文表头。
var headerNames = map[string]string{
	"position": "Position",
	"title":    "Suite",
	"size":     "Samples",
	"mean":     "Mean",
	"variance": "Variance",
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	unit := strings.TrimSpace(in.Unit)
	results := BuildResults(in.Suites, unit)

	out, err := json.Marshal(output{
		Results:    results,
		Markdown:   ArrayToMarkdownTable(results, unit),
		BulletList: ArrayToBulletList(results),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// BuildResults 计算各 suite 的均值/方差并排序、对比，返回展示结果。
func BuildResults(suites []suite, unit string) []result {
	type stats struct {
		idx      int
		title    string
		size     int
		mean     float64
		variance float64
	}
	computed := make([]stats, 0, len(suites))
	for _, s := range suites {
		data := s.Data
		meanV := ComputeAverage(data)
		computed = append(computed, stats{
			idx:      len(computed),
			title:    s.Title,
			size:     len(data),
			mean:     meanV,
			variance: ComputeVariance(data, meanV),
		})
	}
	sort.SliceStable(computed, func(i, j int) bool { return computed[i].mean < computed[j].mean })

	results := make([]result, 0, len(computed))
	for i, c := range computed {
		bestMean := computed[0].mean
		meanTxt := formatFloat(round3(c.mean)) + unit
		if i != 0 && bestMean != c.mean {
			delta := round3(c.mean - bestMean)
			ratio := "∞"
			if bestMean != 0 {
				ratio = formatFloat(round3(c.mean / bestMean))
			}
			meanTxt += fmt.Sprintf(" (+%s%s ; x%s)", formatFloat(delta), unit, ratio)
		}
		varianceTxt := formatFloat(round3(c.variance)) + unit
		if unit != "" {
			varianceTxt += "²"
		}
		results = append(results, result{
			Position:        i + 1,
			Title:           c.title,
			Size:            c.size,
			Mean:            c.mean,
			Variance:        c.variance,
			MeanDisplay:     meanTxt,
			VarianceDisplay: varianceTxt,
		})
	}
	return results
}

// ComputeAverage 计算算术平均，空数组返回 0。
func ComputeAverage(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// ComputeVariance 计算总体方差（平方差的均值），空数组返回 0。
func ComputeVariance(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sqSum := 0.0
	for _, v := range data {
		d := v - mean
		sqSum += d * d
	}
	return sqSum / float64(len(data))
}

// round3 保留 3 位小数（四舍五入）。
func round3(v float64) float64 { return math.Round(v*1000) / 1000 }

// formatFloat 将数值转为最短十进制表示（无尾随零，与 JS toString 一致）。
func formatFloat(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

// ArrayToMarkdownTable 生成 Markdown 表格。
//
// 列顺序与参考实现（object key 顺序）一致：Position | Suite | Mean | Variance | Samples。
func ArrayToMarkdownTable(results []result, unit string) string {
	if len(results) == 0 {
		return ""
	}
	order := []string{"position", "title", "mean", "variance", "size"}

	header := make([]string, len(order))
	for i, k := range order {
		header[i] = headerNames[k]
	}
	lines := []string{
		"| " + strings.Join(header, " | ") + " |",
	}
	sep := make([]string, len(order))
	for i := range sep {
		sep[i] = "---"
	}
	lines = append(lines, "| "+strings.Join(sep, " | ")+" |")
	for _, r := range results {
		row := []string{
			strconv.Itoa(r.Position),
			r.Title,
			r.MeanDisplay,
			r.VarianceDisplay,
			strconv.Itoa(r.Size),
		}
		lines = append(lines, "| "+strings.Join(row, " | ")+" |")
	}
	return strings.Join(lines, "\n")
}

// ArrayToBulletList 生成缩进列表，格式对齐参考实现：
//
//   - title
//   - Position: 1
//   - Mean: ...
//   - Variance: ...
//   - Samples: 2
func ArrayToBulletList(results []result) string {
	if len(results) == 0 {
		return ""
	}
	var b strings.Builder
	for i, r := range results {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, " - %s", r.Title)
		items := []struct {
			key   string
			value string
		}{
			{"position", strconv.Itoa(r.Position)},
			{"mean", r.MeanDisplay},
			{"variance", r.VarianceDisplay},
			{"size", strconv.Itoa(r.Size)},
		}
		for _, it := range items {
			fmt.Fprintf(&b, "\n    - %s: %s", headerNames[it.key], it.value)
		}
	}
	return b.String()
}
