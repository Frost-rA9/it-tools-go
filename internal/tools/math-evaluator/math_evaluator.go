// Package mathevaluator 实现数学表达式求值（自研递归下降解析器，零第三方依赖）。
//
// 支持：四则运算、幂（^，右结合）、一元负号、括号、科学计数法数字、
// 常量 pi/e，以及三角函数/对数等 34 个内置函数（对齐参考项目 it-tools 的 keywords）。
package mathevaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "math-evaluator"
	Name        = "数学表达式求值器"
	Description = "计算数学表达式，支持三角函数、对数、幂等运算"
	Category    = registry.CategoryMath
	Icon        = "Math"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{
	"math", "evaluator", "calculator", "expression", "数学", "表达式", "求值", "计算",
	"abs", "acos", "acosh", "acot", "acoth", "acsc", "acsch", "asec", "asech",
	"asin", "asinh", "atan", "atan2", "atanh", "cos", "cosh", "cot", "coth",
	"csc", "csch", "sec", "sech", "sin", "sinh", "sqrt", "tan", "tanh", "log", "ln",
}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Expression string `json:"expression"` // 待求值的数学表达式
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"` // 求值结果数值字符串（失败时为空）
	Error  string `json:"error"`  // 失败时的中文错误描述（成功时为空）
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	out := output{}
	if strings.TrimSpace(in.Expression) != "" {
		v, err := Eval(in.Expression)
		if err != nil {
			out.Error = err.Error()
		} else {
			out.Result = FormatNumber(v)
		}
	}

	res, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(res), nil
}
