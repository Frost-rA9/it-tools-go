package mathevaluator

import (
	"math"
	"strings"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-12
}

func TestEvalBasic(t *testing.T) {
	cases := []struct {
		expr string
		want float64
	}{
		{"2+3*4", 14},
		{"(1+2)*3", 9},
		{"10-4-3", 3},
		{"7/2", 3.5},
		{"2^3^2", 512}, // 右结合：2^(3^2)
		{"-2^2", -4},   // 一元负号低于幂：-(2^2)
		{"2^-2", 0.25}, // 负指数
		{"2*-3", -6},   // 一元负号参与乘法
		{"-(-3)", 3},
		{"1.5e3", 1500},
		{"1e-2", 0.01},
		{".5", 0.5},
		{"2.", 2},
		{"pi", math.Pi},
		{"2*pi", 2 * math.Pi},
		{"e", math.E},
		{"8/4/2", 1}, // 左结合
		{"2+2*2^2", 10},
	}
	for _, c := range cases {
		got, err := Eval(c.expr)
		if err != nil {
			t.Errorf("Eval(%q) 报错: %v", c.expr, err)
			continue
		}
		if !almostEqual(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.expr, got, c.want)
		}
	}
}

func TestEvalFunctions(t *testing.T) {
	cases := []struct {
		expr string
		want float64
	}{
		{"sqrt(16)", 4},
		{"abs(-5)", 5},
		{"sin(pi/2)", 1},
		{"cos(0)", 1},
		{"tan(pi/4)", 1},
		{"log(100)", 2}, // 以 10 为底
		{"ln(e)", 1},
		{"log2(8)", 3},
		{"exp(0)", 1},
		{"pow(3,2)", 9},
		{"atan2(1,1)", math.Pi / 4},
		{"min(3,7)", 3},
		{"max(3,7)", 7},
		{"mod(7,3)", 1},
		{"floor(2.7)", 2},
		{"ceil(2.1)", 3},
		{"round(2.5)", 3},
		{"sec(0)", 1},
		{"csc(pi/2)", 1},
		{"cot(pi/4)", 1},
		{"asinh(0)", 0},
		{"acosh(1)", 0},
		{"acot(1)", math.Pi / 4},
		{"coth(1)", 1.3130352854993312},
		{"cosh(0)", 1},
		{"sinh(0)", 0},
		{"tanh(0)", 0},
		{"2*sqrt(6)", 2 * math.Sqrt(6)},
		{"sqrt(2)^2", 2},
	}
	for _, c := range cases {
		got, err := Eval(c.expr)
		if err != nil {
			t.Errorf("Eval(%q) 报错: %v", c.expr, err)
			continue
		}
		if !almostEqual(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.expr, got, c.want)
		}
	}
}

func TestEvalErrors(t *testing.T) {
	cases := []struct {
		expr string
		want string // 期望错误信息的子串
	}{
		{"1+", "意外的表达式结尾"},
		{"*3", "意外的记号"},
		{"(1+2", "缺少 )"},
		{"1/0", "除以零"},
		{"sqrt(-1)", "定义域"},
		{"log(0)", "定义域"},
		{"unknownfn(2)", "未知的函数"},
		{"foo", "未知的标识符"},
		{"1 2", "意外的记号"},
		{"sin()", "需要 1 个参数"},
		{"max(1)", "需要 2 个参数"},
		{"2+", "意外的表达式结尾"},
	}
	for _, c := range cases {
		_, err := Eval(c.expr)
		if err == nil {
			t.Errorf("Eval(%q) 应报错，却成功返回", c.expr)
			continue
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("Eval(%q) 错误 = %q, 期望包含 %q", c.expr, err.Error(), c.want)
		}
	}
}

func TestFormatNumber(t *testing.T) {
	cases := []struct {
		v    float64
		want string
	}{
		{14, "14"},
		{3.5, "3.5"},
		{1, "1"},
		{0.1, "0.1"},
		{-2, "-2"},
		{512, "512"},
	}
	for _, c := range cases {
		if got := FormatNumber(c.v); got != c.want {
			t.Errorf("FormatNumber(%v) = %q, want %q", c.v, got, c.want)
		}
	}
}
