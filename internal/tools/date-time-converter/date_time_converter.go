// Package datetime 实现在多种日期时间格式之间进行转换。
package datetime

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "date-time-converter"
	Name        = "日期时间转换器"
	Description = "在不同日期时间格式之间进行转换（时间戳、ISO 8601、RFC 系列、Excel 等）"
	Category    = "转换器"
	Icon        = "Calendar"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"date", "time", "日期", "时间", "时间戳", "timestamp", "unix", "iso8601", "excel"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Value  string `json:"value"`  // 待解析的日期时间字符串
	Format int    `json:"format"` // 用于解析的格式索引（自动识别失败时的回退）
}

// result 是单个格式的转换结果。
type result struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// output 是工具的输出结构。
type output struct {
	Results  []result `json:"results"`  // 全部格式的转换结果
	Detected int      `json:"detected"` // 自动识别到的格式索引（-1 表示未识别）
	Valid    bool     `json:"valid"`    // 输入是否有效
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 解析输入并返回全部格式的转换结果。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	value := strings.TrimSpace(in.Value)
	var (
		t        time.Time
		detected = -1
		valid    = true
	)

	if value == "" {
		t = time.Now()
	} else {
		formatIndex := in.Format
		if d := detectFormat(value); d >= 0 {
			formatIndex = d
			detected = d
		}
		if formatIndex < 0 || formatIndex >= len(formats) {
			return "", fmt.Errorf("未知格式索引: %d", formatIndex)
		}
		parsed, err := formats[formatIndex].toDate(value)
		if err != nil {
			valid = false
		} else {
			t = parsed
		}
	}

	results := make([]result, 0, len(formats))
	for _, f := range formats {
		v := ""
		if valid {
			v = f.fromDate(t)
		}
		results = append(results, result{Name: f.name, Value: v})
	}

	out, err := json.Marshal(output{Results: results, Detected: detected, Valid: valid})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// formatDef 描述一种日期时间格式及其解析/格式化逻辑。
type formatDef struct {
	name     string
	fromDate func(time.Time) string
	toDate   func(string) (time.Time, error)
	matcher  func(string) bool
}

// formats 为全部支持的格式（顺序即索引，前端下拉框与此对应）。
var formats = []formatDef{
	{"Locale string", formatLocale, parseLocale, func(string) bool { return false }},
	{"ISO 8601", formatISO8601, parseISO8601, isISO8601},
	{"ISO 9075", formatISO9075, parseISO9075, isISO9075},
	{"RFC 3339", formatRFC3339, parseRFC3339, isRFC3339},
	{"RFC 7231", formatRFC7231, parseRFC7231, isRFC7231},
	{"Unix timestamp", formatUnix, parseUnix, isUnixTimestamp},
	{"Timestamp", formatTimestamp, parseTimestamp, isTimestamp},
	{"UTC format", formatUTC, parseRFC7231, isRFC7231},
	{"Mongo ObjectID", formatMongo, parseMongo, isMongoObjectID},
	{"Excel date/time", formatExcel, parseExcel, isExcelFormat},
}

// 各格式的「格式化」实现（time.Time → string）。
func formatLocale(t time.Time) string    { return t.Format(time.UnixDate) }
func formatISO8601(t time.Time) string   { return t.Format("2006-01-02T15:04:05Z07:00") }
func formatISO9075(t time.Time) string   { return t.Format("2006-01-02 15:04:05") }
func formatRFC3339(t time.Time) string   { return t.Format(time.RFC3339) }
func formatRFC7231(t time.Time) string   { return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") }
func formatUnix(t time.Time) string      { return strconv.FormatInt(t.Unix(), 10) }
func formatTimestamp(t time.Time) string { return strconv.FormatInt(t.UnixMilli(), 10) }
func formatUTC(t time.Time) string       { return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") }
func formatMongo(t time.Time) string     { return strconv.FormatInt(t.Unix(), 16) + strings.Repeat("0", 16) }
func formatExcel(t time.Time) string {
	return strconv.FormatFloat(float64(t.UnixMilli())/86400000+25569, 'f', -1, 64)
}

// 各格式的「解析」实现（string → time.Time）。
func parseLocale(s string) (time.Time, error) {
	for _, layout := range []string{time.UnixDate, time.ANSIC, time.RubyDate} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析为本地日期: %q", s)
}

func parseISO8601(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	for _, layout := range []string{"2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析为 ISO 8601: %q", s)
}

func parseISO9075(s string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析为 ISO 9075: %q", s)
}

func parseRFC3339(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析为 RFC 3339: %q", s)
}

func parseRFC7231(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC1123, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析为 RFC 7231: %q", s)
	}
	return t, nil
}

func parseUnix(s string) (time.Time, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析为 Unix 时间戳: %q", s)
	}
	return time.Unix(n, 0), nil
}

func parseTimestamp(s string) (time.Time, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析为毫秒时间戳: %q", s)
	}
	return time.UnixMilli(n), nil
}

func parseMongo(s string) (time.Time, error) {
	if !isMongoObjectID(s) {
		return time.Time{}, fmt.Errorf("无法解析为 Mongo ObjectID: %q", s)
	}
	sec, err := strconv.ParseInt(s[:8], 16, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析为 Mongo ObjectID: %q", s)
	}
	return time.Unix(sec, 0), nil
}

func parseExcel(s string) (time.Time, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析为 Excel 日期: %q", s)
	}
	totalSec := (n - 25569) * 86400
	sec := int64(totalSec)
	nsec := int64((totalSec - float64(sec)) * 1e9)
	return time.Unix(sec, nsec), nil
}

// 格式识别正则（顺序见 detectFormat）。
var (
	reMongoObjectID = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
	reRFC7231       = regexp.MustCompile(`^[A-Za-z]{3},\s\d{2}\s[A-Za-z]{3}\s\d{4}\s\d{2}:\d{2}:\d{2}\sGMT$`)
	reISO9075       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
	reRFC3339       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`)
	reISO8601       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}(:\d{2})?)?$`)
	reUnixTimestamp = regexp.MustCompile(`^\d{1,10}$`)
	reTimestamp     = regexp.MustCompile(`^\d{1,13}$`)
	reExcelFormat   = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
)

func isMongoObjectID(s string) bool { return reMongoObjectID.MatchString(s) }
func isRFC7231(s string) bool       { return reRFC7231.MatchString(s) }
func isISO9075(s string) bool       { return reISO9075.MatchString(s) }
func isRFC3339(s string) bool       { return reRFC3339.MatchString(s) }
func isISO8601(s string) bool       { return reISO8601.MatchString(s) }
func isUnixTimestamp(s string) bool { return reUnixTimestamp.MatchString(s) }
func isTimestamp(s string) bool     { return reTimestamp.MatchString(s) }
func isExcelFormat(s string) bool   { return reExcelFormat.MatchString(s) }

// detectFormat 按优先级识别输入字符串所属的格式索引，未识别返回 -1。
func detectFormat(s string) int {
	switch {
	case isMongoObjectID(s):
		return 8
	case isRFC7231(s):
		return 4
	case isISO9075(s):
		return 2
	case isRFC3339(s):
		return 3
	case isISO8601(s):
		return 1
	case isUnixTimestamp(s):
		return 5
	case isTimestamp(s):
		return 6
	case isExcelFormat(s):
		return 9
	default:
		return -1
	}
}
