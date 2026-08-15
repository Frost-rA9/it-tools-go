// Package pwdstrength 实现密码强度分析工具（对齐 it-tools password-strength-analyser）。
package pwdstrength

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "password-strength-analyser"
	Name        = "密码强度分析仪"
	Description = "基于暴力破解时间估算分析密码强度"
	Category    = registry.CategoryCrypto
	Icon        = "ShieldLock"
)

// Keywords 为搜索关键词。
var Keywords = []string{"password", "strength", "analyser", "crack", "time", "estimation", "brute", "force", "attack", "entropy", "密码", "强度"}

// guessesPerSecond 与 it-tools 一致。
const guessesPerSecond = 1e9

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Password string `json:"password"`
}

// output 是工具的输出结构。
type output struct {
	PasswordLength int     `json:"password_length"`
	Entropy        float64 `json:"entropy"`
	CharsetLength  int     `json:"charset_length"`
	CrackDuration  string  `json:"crack_duration"`
	SecondsToCrack float64 `json:"seconds_to_crack"`
	Score          float64 `json:"score"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 分析密码强度并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	charset := charsetLength(in.Password)
	entropy := 0.0
	if in.Password != "" {
		entropy = math.Log2(float64(charset)) * float64(len(in.Password))
	}
	seconds := math.Pow(2, entropy) / guessesPerSecond

	out, err := json.Marshal(output{
		PasswordLength: len(in.Password),
		Entropy:        entropy,
		CharsetLength:  charset,
		CrackDuration:  humanDuration(seconds),
		SecondsToCrack: seconds,
		Score:          math.Min(entropy/128, 1),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

var (
	reLowercase = regexp.MustCompile(`[a-z]`)
	reUppercase = regexp.MustCompile(`[A-Z]`)
	reDigits    = regexp.MustCompile(`[0-9]`)
	reSpecial   = regexp.MustCompile(`[^a-zA-Z0-9]`) // 对齐 JS 的 \W|_（含下划线）
)

// charsetLength 统计密码所覆盖的字符集大小（对齐 it-tools getCharsetLength）。
func charsetLength(password string) int {
	n := 0
	if reLowercase.MatchString(password) {
		n += 26
	}
	if reUppercase.MatchString(password) {
		n += 26
	}
	if reDigits.MatchString(password) {
		n += 10
	}
	if reSpecial.MatchString(password) {
		n += 32
	}
	return n
}

// timeUnit 描述一个时长单位。
type timeUnit struct {
	unit      string
	secondsIn float64
	plural    string
	prettify  bool
}

var timeUnits = []timeUnit{
	{"millenium", 31536000000, "millennia", true},
	{"century", 3153600000, "centuries", false},
	{"decade", 315360000, "decades", false},
	{"year", 31536000, "years", false},
	{"month", 2592000, "months", false},
	{"week", 604800, "weeks", false},
	{"day", 86400, "days", false},
	{"hour", 3600, "hours", false},
	{"minute", 60, "minutes", false},
	{"second", 1, "seconds", false},
}

// humanDuration 将秒数格式化为人类可读时长（对齐 it-tools getHumanFriendlyDuration）。
func humanDuration(seconds float64) string {
	if seconds <= 0.001 {
		return "Instantly"
	}
	if seconds <= 1 {
		return "Less than a second"
	}

	parts := make([]string, 0, 2)
	for _, u := range timeUnits {
		quantity := math.Floor(seconds / u.secondsIn)
		seconds = math.Mod(seconds, u.secondsIn)

		// 注意：NaN <= 0 为 false，与 JS 行为一致（会进入并输出 NaN）。
		if quantity <= 0 {
			continue
		}
		var formatted string
		if u.prettify {
			formatted = prettifyQuantity(quantity)
		} else {
			formatted = identityQuantity(quantity)
		}
		label := u.unit
		if quantity > 1 {
			label = u.plural
		}
		parts = append(parts, formatted+" "+label)
		if len(parts) == 2 {
			break
		}
	}
	return strings.Join(parts, ", ")
}

// prettifyQuantity 对齐 it-tools 的 prettifyExponentialNotation（仅用于 millennia）。
func prettifyQuantity(q float64) string {
	if math.IsInf(q, 1) {
		return "Infinity"
	}
	if math.IsNaN(q) {
		return "NaN"
	}
	// JS number.toString() 对 >= 1e21 使用指数表示。
	if q >= 1e21 {
		exp := int(math.Floor(math.Log10(q)))
		base := q / math.Pow(10, float64(exp))
		var baseStr string
		if base == math.Floor(base) {
			baseStr = strconv.FormatInt(int64(base), 10)
		} else {
			baseStr = strconv.FormatFloat(base, 'f', 2, 64)
		}
		return baseStr + "e+" + strconv.Itoa(exp)
	}
	// < 1e21：JS 输出完整整数并带千分位分组。
	return groupThousands(q)
}

// identityQuantity 对齐 JS 的 _.identity（非 millennia 单位）。
func identityQuantity(q float64) string {
	if math.IsInf(q, 1) {
		return "Infinity"
	}
	if math.IsNaN(q) {
		return "NaN"
	}
	return strconv.FormatFloat(q, 'f', 0, 64)
}

// groupThousands 为整数浮点添加千分位逗号（用最短往返表示对齐 JS toLocaleString）。
func groupThousands(q float64) string {
	s := strconv.FormatFloat(q, 'f', -1, 64)
	var b strings.Builder
	n := len(s)
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
