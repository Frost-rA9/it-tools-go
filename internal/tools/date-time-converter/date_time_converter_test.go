package datetime

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 固定测试时刻：2026-08-14 12:34:56 UTC（周五）。
func fixedTime() time.Time {
	return time.Date(2026, 8, 14, 12, 34, 56, 0, time.UTC)
}

func TestFormatFunctions(t *testing.T) {
	ts := fixedTime()

	checks := []struct {
		name string
		got  string
		want string
	}{
		{"ISO 8601", formatISO8601(ts), "2026-08-14T12:34:56Z"},
		{"ISO 9075", formatISO9075(ts), "2026-08-14 12:34:56"},
		{"RFC 3339", formatRFC3339(ts), "2026-08-14T12:34:56Z"},
		{"RFC 7231", formatRFC7231(ts), "Fri, 14 Aug 2026 12:34:56 GMT"},
		{"UTC", formatUTC(ts), "Fri, 14 Aug 2026 12:34:56 GMT"},
		{"Locale", formatLocale(ts), "Fri Aug 14 12:34:56 UTC 2026"},
		{"Unix", formatUnix(ts), "1786710896"},
		{"Timestamp", formatTimestamp(ts), "1786710896000"},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q，期望 %q", c.name, c.got, c.want)
		}
	}

	// Mongo ObjectID：hex(秒) + 16 个 0，当前纪元下应为 24 位。
	mongo := formatMongo(ts)
	if len(mongo) != 24 {
		t.Errorf("Mongo ObjectID 长度 = %d，期望 24（%q）", len(mongo), mongo)
	}
	if want := strconv.FormatInt(ts.Unix(), 16) + strings.Repeat("0", 16); mongo != want {
		t.Errorf("Mongo ObjectID = %q，期望 %q", mongo, want)
	}

	// Excel：Unix 纪元（1970-01-01）对应序列号 25569。
	if got := formatExcel(time.Unix(0, 0)); got != "25569" {
		t.Errorf("Excel(epoch) = %q，期望 25569", got)
	}
}

func TestParseFunctions(t *testing.T) {
	ts := fixedTime()

	t.Run("ISO8601 往返", func(t *testing.T) {
		got, err := parseISO8601("2026-08-14T12:34:56Z")
		if err != nil || !got.Equal(ts) {
			t.Errorf("parseISO8601 = %v, %v", got, err)
		}
	})
	t.Run("ISO9075 本地解析", func(t *testing.T) {
		want := time.Date(2026, 8, 14, 12, 34, 56, 0, time.Local)
		got, err := parseISO9075("2026-08-14 12:34:56")
		if err != nil || !got.Equal(want) {
			t.Errorf("parseISO9075 = %v, %v", got, err)
		}
	})
	t.Run("RFC3339 往返", func(t *testing.T) {
		got, err := parseRFC3339("2026-08-14T12:34:56Z")
		if err != nil || !got.Equal(ts) {
			t.Errorf("parseRFC3339 = %v, %v", got, err)
		}
	})
	t.Run("RFC7231", func(t *testing.T) {
		got, err := parseRFC7231("Fri, 14 Aug 2026 12:34:56 GMT")
		if err != nil || !got.Equal(ts) {
			t.Errorf("parseRFC7231 = %v, %v", got, err)
		}
	})
	t.Run("Unix", func(t *testing.T) {
		got, err := parseUnix("1786710896")
		if err != nil || got.Unix() != 1786710896 {
			t.Errorf("parseUnix = %v, %v", got, err)
		}
	})
	t.Run("Timestamp", func(t *testing.T) {
		got, err := parseTimestamp("1786710896000")
		if err != nil || got.UnixMilli() != 1786710896000 {
			t.Errorf("parseTimestamp = %v, %v", got, err)
		}
	})
	t.Run("Mongo 往返", func(t *testing.T) {
		got, err := parseMongo(formatMongo(ts))
		if err != nil || got.Unix() != ts.Unix() {
			t.Errorf("parseMongo = %v, %v", got, err)
		}
	})
	t.Run("Excel 纪元", func(t *testing.T) {
		got, err := parseExcel("25569")
		if err != nil || got.Unix() != 0 {
			t.Errorf("parseExcel = %v, %v", got, err)
		}
	})
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"Mongo ObjectID", "62a88d5e0000000000000000", 8},
		{"RFC 7231", "Fri, 14 Aug 2026 12:34:56 GMT", 4},
		{"ISO 9075", "2026-08-14 12:34:56", 2},
		{"RFC 3339", "2026-08-14T12:34:56Z", 3},
		{"RFC 3339 带偏移", "2026-08-14T12:34:56+08:00", 3},
		{"ISO 8601 仅日期", "2026-08-14", 1},
		{"ISO 8601 无偏移时间", "2026-08-14T12:34:56", 1},
		{"Unix 秒", "1786710896", 5},
		{"毫秒", "1786710896000", 6},
		{"Excel 小数", "46248.52425925926", 9},
		{"Excel 负数", "-2209161600", 9},
		{"未识别", "hello world", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectFormat(tt.s); got != tt.want {
				t.Errorf("detectFormat(%q) = %d，期望 %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	t.Run("空输入返回当前时间", func(t *testing.T) {
		out, err := run(t, e, input{Value: "", Format: 6})
		if err != nil {
			t.Fatalf("Execute 返回错误: %v", err)
		}
		if !out.Valid || len(out.Results) != 10 {
			t.Fatalf("Valid=%v, 结果数=%d", out.Valid, len(out.Results))
		}
		for _, r := range out.Results {
			if r.Value == "" {
				t.Errorf("格式 %s 结果为空", r.Name)
			}
		}
	})

	t.Run("Unix 秒自动识别", func(t *testing.T) {
		out, err := run(t, e, input{Value: "1786710896", Format: 6})
		if err != nil {
			t.Fatalf("Execute 返回错误: %v", err)
		}
		if out.Detected != 5 || !out.Valid {
			t.Fatalf("Detected=%d, Valid=%v", out.Detected, out.Valid)
		}
		// RFC 7231 恒为 UTC，与时区无关。
		if got := valueOf(out, "RFC 7231"); got != "Fri, 14 Aug 2026 12:34:56 GMT" {
			t.Errorf("RFC 7231 = %q", got)
		}
		if got := valueOf(out, "Unix timestamp"); got != "1786710896" {
			t.Errorf("Unix timestamp = %q", got)
		}
	})

	t.Run("非法输入", func(t *testing.T) {
		out, err := run(t, e, input{Value: "hello world", Format: 6})
		if err != nil {
			t.Fatalf("Execute 返回错误: %v", err)
		}
		if out.Valid {
			t.Errorf("期望 invalid，但 Valid=true")
		}
		for _, r := range out.Results {
			if r.Value != "" {
				t.Errorf("invalid 输入下格式 %s 应为空，实际 %q", r.Name, r.Value)
			}
		}
	})
}

func run(t *testing.T, e Executor, in input) (output, error) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("序列化输入失败: %v", err)
	}
	outJSON, err := e.Execute(context.Background(), string(b))
	if err != nil {
		return output{}, err
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	return out, nil
}

func valueOf(out output, name string) string {
	for _, r := range out.Results {
		if r.Name == name {
			return r.Value
		}
	}
	return ""
}
