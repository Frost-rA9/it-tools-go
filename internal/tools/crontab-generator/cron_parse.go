// Package crontab 的表达式解析器：将标准 cron 表达式解析为中文描述。
package crontab

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// atAlias 支持 @ 简写表达式。
var atAlias = map[string]string{
	"@yearly":   "每年 1 月 1 日 0 时 0 分执行",
	"@annually": "每年 1 月 1 日 0 时 0 分执行",
	"@monthly":  "每月 1 日 0 时 0 分执行",
	"@weekly":   "每周日 0 时 0 分执行",
	"@daily":    "每天 0 时 0 分执行",
	"@midnight": "每天 0 时 0 分执行",
	"@hourly":   "每小时第 0 分钟执行",
	"@reboot":   "系统启动时执行",
}

// monthNames / weekNames 支持 cron 名称别名。
var monthNames = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

var weekNames = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
}

// fieldDef 描述单个 cron 字段的边界与中文表述。
type fieldDef struct {
	name    string         // 中文名（用于错误提示）
	min     int            // 允许最小值
	max     int            // 允许最大值
	names   map[string]int // 名称别名（月/周），nil 表示不支持名称
	unit    string         // 步长单位，如 "分钟" / "小时"
	single  string         // 单值表述模板，如 "每月 %d 日"
	rangeT  string         // 范围表述模板，如 "每月 %d 至 %d 日"
	listSep string         // 列表分隔字符（"、"）
	weekday bool           // 是否星期字段（显示为 周X）
}

var fieldDefs = []fieldDef{
	{name: "分钟", min: 0, max: 59, unit: "分钟", single: "第 %d 分钟", rangeT: "第 %d 至 %d 分钟"},
	{name: "小时", min: 0, max: 23, unit: "小时", single: "第 %d 点", rangeT: "%d 至 %d 点"},
	{name: "日", min: 1, max: 31, unit: "天", single: "每月 %d 日", rangeT: "每月 %d 至 %d 日"},
	{name: "月", min: 1, max: 12, unit: "个月", names: monthNames, single: "%d 月", rangeT: "%d 至 %d 月"},
	{name: "星期", min: 0, max: 7, unit: "周", names: weekNames, single: "周%s", rangeT: "周%s至周%s", weekday: true},
}

// parsedField 是单个字段 token 的解析结果。
type parsedField struct {
	every  bool  // 通配 *
	step   int   // 步长（0 表示无）
	from   int   // 范围起点或单值
	to     int   // 范围终点（单值时等于 from）
	ranged bool  // 是否为范围
	values []int // 列表值（升序去重）
}

// parseCron 将 cron 表达式解析为中文描述。
func parseCron(expr string) (string, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return "", fmt.Errorf("表达式为空")
	}
	if desc, ok := atAlias[strings.ToLower(expr)]; ok {
		return desc, nil
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 && len(fields) != 6 {
		return "", fmt.Errorf("表达式应包含 5 或 6 个字段（当前 %d 个）", len(fields))
	}

	// 可选秒字段（6 字段时第一个是秒）
	var secDesc string
	start := 0
	if len(fields) == 6 {
		sp, err := parseFieldToken(fields[0], fieldDefs[0])
		if err != nil {
			return "", fmt.Errorf("秒字段: %w", err)
		}
		secDesc = describeFixed(sp, 0, 59, "秒")
		start = 1
	}

	parts := make([]fieldText, 0, 5)
	for i, def := range fieldDefs {
		p, err := parseFieldToken(fields[start+i], def)
		if err != nil {
			return "", fmt.Errorf("%s字段: %w", def.name, err)
		}
		parts = append(parts, describe(p, def))
	}

	// 特判：分/时均为单值时合并为 "N 时 M 分"
	if !parts[0].empty && !parts[1].empty && parts[0].singleVal && parts[1].singleVal {
		parts[0].text = fmt.Sprintf("%d 时 %d 分", parts[1].val, parts[0].val)
		parts[1].text = ""
	} else if parts[0].singleVal && parts[1].empty {
		// 分钟单值、小时通配 → 每小时第 N 分钟
		parts[0].text = fmt.Sprintf("每小时第 %d 分钟", parts[0].val)
	}

	desc := combineDesc(parts)
	if secDesc != "" {
		desc = secDesc + "，" + desc
	}
	return desc, nil
}

// fieldText 是描述片段的中间表示。
type fieldText struct {
	text      string
	empty     bool // 通配 *，跳过
	singleVal bool // 单值（用于分/时合并特判）
	val       int
}

func combineDesc(parts []fieldText) string {
	texts := make([]string, 0, 5)
	for _, p := range parts {
		if !p.empty && p.text != "" {
			texts = append(texts, p.text)
		}
	}
	if len(texts) == 0 {
		return "每分钟"
	}
	return strings.Join(texts, "，")
}

// parseFieldToken 解析单个字段 token（*、*/N、A-B、A-B/N、A,B,C、数字、名称）。
func parseFieldToken(tok string, def fieldDef) (parsedField, error) {
	tok = strings.ToLower(strings.TrimSpace(tok))
	if tok == "" {
		return parsedField{}, fmt.Errorf("字段为空")
	}
	if tok == "*" {
		return parsedField{every: true}, nil
	}

	// */N 步长
	if strings.HasPrefix(tok, "*/") {
		n, err := strconv.Atoi(tok[2:])
		if err != nil || n < 1 {
			return parsedField{}, fmt.Errorf("非法步长 %q", tok[2:])
		}
		return parsedField{step: n}, nil
	}

	// 列表
	if strings.Contains(tok, ",") {
		parts := strings.Split(tok, ",")
		vals := make([]int, 0, len(parts))
		seen := map[int]bool{}
		for _, part := range parts {
			v, err := resolveValue(part, def)
			if err != nil {
				return parsedField{}, err
			}
			if !seen[v] {
				seen[v] = true
				vals = append(vals, v)
			}
		}
		sort.Ints(vals)
		return parsedField{values: vals}, nil
	}

	// 范围（可带步长 A-B/N）
	if strings.Contains(tok, "-") {
		rangePart := tok
		step := 0
		if idx := strings.Index(tok, "/"); idx >= 0 {
			rangePart = tok[:idx]
			n, err := strconv.Atoi(tok[idx+1:])
			if err != nil || n < 1 {
				return parsedField{}, fmt.Errorf("非法步长 %q", tok[idx+1:])
			}
			step = n
		}
		ends := strings.SplitN(rangePart, "-", 2)
		from, err := resolveValue(ends[0], def)
		if err != nil {
			return parsedField{}, err
		}
		to, err := resolveValue(ends[1], def)
		if err != nil {
			return parsedField{}, err
		}
		if from > to {
			return parsedField{}, fmt.Errorf("范围起点 %d 大于终点 %d", from, to)
		}
		return parsedField{from: from, to: to, ranged: true, step: step}, nil
	}

	// 单值（数字或名称）
	v, err := resolveValue(tok, def)
	if err != nil {
		return parsedField{}, err
	}
	return parsedField{from: v, to: v, values: []int{v}}, nil
}

// resolveValue 解析单个值（数字或名称），并校验范围。
func resolveValue(s string, def fieldDef) (int, error) {
	v := -1
	if n, err := strconv.Atoi(s); err == nil {
		v = n
	} else if def.names != nil {
		if n, ok := def.names[s]; ok {
			v = n
		}
	}
	if v < 0 {
		return 0, fmt.Errorf("非法取值 %q", s)
	}
	if v < def.min || v > def.max {
		return 0, fmt.Errorf("取值 %d 超出 %d-%d", v, def.min, def.max)
	}
	return v, nil
}

// describe 将解析结果渲染为中文片段。
func describe(p parsedField, def fieldDef) fieldText {
	switch {
	case p.every:
		return fieldText{empty: true}
	case p.step > 0 && !p.ranged:
		// */N：从字段最小值开始的步长
		return fieldText{text: fmt.Sprintf("每 %d %s", p.step, def.unit)}
	case p.ranged:
		if p.step > 0 {
			return fieldText{text: fmt.Sprintf(def.rangeT+"（每 %d 个）", p.from, p.to, p.step)}
		}
		return fieldText{text: formatRange(p, def)}
	case len(p.values) > 0:
		if len(p.values) == 1 {
			t := formatSingle(p.values[0], def)
			return fieldText{text: t, singleVal: true, val: p.values[0]}
		}
		return fieldText{text: formatList(p.values, def)}
	}
	return fieldText{}
}

func formatSingle(v int, def fieldDef) string {
	if def.weekday {
		return fmt.Sprintf(def.single, weekdayName(v))
	}
	return fmt.Sprintf(def.single, v)
}

func formatRange(p parsedField, def fieldDef) string {
	if def.weekday {
		return fmt.Sprintf(def.rangeT, weekdayName(p.from), weekdayName(p.to))
	}
	return fmt.Sprintf(def.rangeT, p.from, p.to)
}

func formatList(vals []int, def fieldDef) string {
	strs := make([]string, len(vals))
	for i, v := range vals {
		if def.weekday {
			strs[i] = weekdayName(v)
		} else {
			strs[i] = fmt.Sprintf("%d", v)
		}
	}
	if def.weekday {
		return "周" + strings.Join(strs, "、")
	}
	joined := strings.Join(strs, "、")
	switch def.name {
	case "分钟":
		return "第 " + joined + " 分钟"
	case "小时":
		return joined + " 点"
	case "日":
		return "每月 " + joined + " 日"
	case "月":
		return joined + " 月"
	}
	return joined
}

func weekdayName(v int) string {
	names := []string{"日", "一", "二", "三", "四", "五", "六"}
	if v == 7 {
		v = 0 // cron 中 7 也代表周日
	}
	return names[v]
}

// describeFixed 描述秒字段（6 字段表达式的首段）。
func describeFixed(p parsedField, min, max int, unit string) string {
	switch {
	case p.every:
		return ""
	case p.step > 0:
		return fmt.Sprintf("每 %d %s", p.step, unit)
	case p.ranged:
		return fmt.Sprintf("第 %d 至 %d %s", p.from, p.to, unit)
	case len(p.values) > 0:
		if len(p.values) == 1 {
			return fmt.Sprintf("第 %d %s", p.values[0], unit)
		}
		strs := make([]string, len(p.values))
		for i, v := range p.values {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return fmt.Sprintf("第 %s %s", strings.Join(strs, "、"), unit)
	}
	return ""
}
