package mathevaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// tokenKind 是词法单元的类型。
type tokenKind int

const (
	tokNumber tokenKind = iota
	tokIdent
	tokPlus
	tokMinus
	tokStar
	tokSlash
	tokCaret
	tokLParen
	tokRParen
	tokComma
	tokEOF
)

// token 是一个词法单元。
type token struct {
	kind  tokenKind
	text  string
	value float64
}

// isFunction 标记某函数是否接受两个参数。
var arity = map[string]int{}

func init() {
	// 单参数函数。
	unaryNames := []string{
		"abs", "acos", "acosh", "acot", "acoth", "acsc", "acsch",
		"asec", "asech", "asin", "asinh", "atan", "atanh",
		"ceil", "cos", "cosh", "cot", "coth", "csc", "csch",
		"exp", "floor", "ln", "log", "log2", "round",
		"sec", "sech", "sin", "sinh", "sqrt", "tan", "tanh",
	}
	// 双参数函数。
	binaryNames := []string{"atan2", "max", "min", "mod", "pow"}
	for _, n := range unaryNames {
		arity[n] = 1
	}
	for _, n := range binaryNames {
		arity[n] = 2
	}
}

// Eval 求值数学表达式，返回结果与错误（错误为中文描述，便于前端展示）。
func Eval(expr string) (float64, error) {
	toks, err := lex(expr)
	if err != nil {
		return 0, err
	}
	p := &parser{tokens: toks}
	v, err := p.parseExpression()
	if err != nil {
		return 0, err
	}
	if p.peek().kind != tokEOF {
		return 0, fmt.Errorf("意外的记号 %q", p.peek().text)
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("结果无定义（可能包含除以零或非法的数学运算）")
	}
	return v, nil
}

// lex 将表达式切分为词法单元。
func lex(expr string) ([]token, error) {
	var toks []token
	i := 0
	for i < len(expr) {
		c := expr[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c == '+':
			toks = append(toks, token{kind: tokPlus, text: "+"})
			i++
		case c == '-':
			toks = append(toks, token{kind: tokMinus, text: "-"})
			i++
		case c == '*':
			toks = append(toks, token{kind: tokStar, text: "*"})
			i++
		case c == '/':
			toks = append(toks, token{kind: tokSlash, text: "/"})
			i++
		case c == '^':
			toks = append(toks, token{kind: tokCaret, text: "^"})
			i++
		case c == '(':
			toks = append(toks, token{kind: tokLParen, text: "("})
			i++
		case c == ')':
			toks = append(toks, token{kind: tokRParen, text: ")"})
			i++
		case c == ',':
			toks = append(toks, token{kind: tokComma, text: ","})
			i++
		case c == '.' || (c >= '0' && c <= '9'):
			start := i
			for i < len(expr) && (expr[i] == '.' || (expr[i] >= '0' && expr[i] <= '9')) {
				i++
			}
			// 科学计数法指数部分：e/E [±] digits。
			if i < len(expr) && (expr[i] == 'e' || expr[i] == 'E') {
				j := i + 1
				if j < len(expr) && (expr[j] == '+' || expr[j] == '-') {
					j++
				}
				if j < len(expr) && expr[j] >= '0' && expr[j] <= '9' {
					i = j
					for i < len(expr) && expr[i] >= '0' && expr[i] <= '9' {
						i++
					}
				}
			}
			s := expr[start:i]
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf("无效的数字 %q", s)
			}
			toks = append(toks, token{kind: tokNumber, text: s, value: v})
		case isIdentStart(c):
			start := i
			for i < len(expr) && isIdentPart(expr[i]) {
				i++
			}
			toks = append(toks, token{kind: tokIdent, text: expr[start:i]})
		default:
			return nil, fmt.Errorf("无法识别的字符 %q", string(c))
		}
	}
	toks = append(toks, token{kind: tokEOF})
	return toks, nil
}

func isIdentStart(c byte) bool { return c == '_' || unicode.IsLetter(rune(c)) }
func isIdentPart(c byte) bool  { return c == '_' || unicode.IsLetter(rune(c)) || (c >= '0' && c <= '9') }

// parser 是递归下降语法解析器。
//
// 语法（优先级从低到高）：
//
//	expression → term (('+' | '-') term)*
//	term       → unary (('*' | '/') unary)*
//	unary      → '-' unary | power
//	power      → primary ('^' unary)?          // 幂右结合，且优先于一元负号（-2^2 = -(2^2)）
//	primary    → number | '(' expression ')' | fn '(' args ')' | 常量
type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token { return p.tokens[p.pos] }
func (p *parser) advance() token {
	t := p.tokens[p.pos]
	if t.kind != tokEOF {
		p.pos++
	}
	return t
}

// parseExpression 处理加减（最低优先级）。
func (p *parser) parseExpression() (float64, error) {
	v, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		switch p.peek().kind {
		case tokPlus:
			p.advance()
			r, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			v += r
		case tokMinus:
			p.advance()
			r, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			v -= r
		default:
			return v, nil
		}
	}
}

// parseTerm 处理乘除。
func (p *parser) parseTerm() (float64, error) {
	v, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	for {
		switch p.peek().kind {
		case tokStar:
			p.advance()
			r, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			v *= r
		case tokSlash:
			p.advance()
			r, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			if r == 0 {
				return 0, fmt.Errorf("除以零")
			}
			v /= r
		default:
			return v, nil
		}
	}
}

// parseUnary 处理一元负号。
func (p *parser) parseUnary() (float64, error) {
	if p.peek().kind == tokMinus {
		p.advance()
		v, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return -v, nil
	}
	return p.parsePower()
}

// parsePower 处理幂（右结合）。
func (p *parser) parsePower() (float64, error) {
	base, err := p.parsePrimary()
	if err != nil {
		return 0, err
	}
	if p.peek().kind == tokCaret {
		p.advance()
		exp, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return math.Pow(base, exp), nil
	}
	return base, nil
}

// parsePrimary 处理数字、括号表达式、函数调用与常量。
func (p *parser) parsePrimary() (float64, error) {
	t := p.advance()
	switch t.kind {
	case tokNumber:
		return t.value, nil
	case tokLParen:
		v, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		if p.peek().kind != tokRParen {
			return 0, errExpect(")")
		}
		p.advance()
		return v, nil
	case tokIdent:
		if p.peek().kind == tokLParen {
			return p.parseCall(t.text)
		}
		switch t.text {
		case "pi":
			return math.Pi, nil
		case "e":
			return math.E, nil
		default:
			return 0, fmt.Errorf("未知的标识符 %q", t.text)
		}
	default:
		return 0, errUnexpected(strings.TrimSpace(t.text))
	}
}

// parseCall 解析函数调用并求值。
func (p *parser) parseCall(name string) (float64, error) {
	n, ok := arity[name]
	if !ok {
		return 0, fmt.Errorf("未知的函数 %q", name)
	}
	if p.peek().kind != tokLParen {
		return 0, errExpect("(")
	}
	p.advance()

	if p.peek().kind == tokRParen {
		return 0, fmt.Errorf("函数 %s 需要 %d 个参数，实际 0 个", name, n)
	}

	args := make([]float64, 0, n)
	for {
		a, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		args = append(args, a)
		if p.peek().kind == tokComma {
			p.advance()
			continue
		}
		break
	}
	if p.peek().kind != tokRParen {
		return 0, errExpect(")")
	}
	p.advance()

	if len(args) != n {
		return 0, fmt.Errorf("函数 %s 需要 %d 个参数，实际 %d 个", name, n, len(args))
	}
	v, err := callFunc(name, n, args)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("函数 %s 的输入超出定义域", name)
	}
	return v, nil
}

// callFunc 按函数名与元数分发到实现。
func callFunc(name string, n int, args []float64) (float64, error) {
	if n == 2 {
		switch name {
		case "atan2":
			return math.Atan2(args[0], args[1]), nil
		case "max":
			return math.Max(args[0], args[1]), nil
		case "min":
			return math.Min(args[0], args[1]), nil
		case "mod":
			return math.Mod(args[0], args[1]), nil
		case "pow":
			return math.Pow(args[0], args[1]), nil
		}
	}
	x := args[0]
	switch name {
	case "abs":
		return math.Abs(x), nil
	case "acos":
		return math.Acos(x), nil
	case "acosh":
		return math.Acosh(x), nil
	case "acot":
		return math.Pi/2 - math.Atan(x), nil
	case "acoth":
		return 0.5 * math.Log((x+1)/(x-1)), nil
	case "acsc":
		return math.Asin(1 / x), nil
	case "acsch":
		return math.Asinh(1 / x), nil
	case "asec":
		return math.Acos(1 / x), nil
	case "asech":
		return math.Acosh(1 / x), nil
	case "asin":
		return math.Asin(x), nil
	case "asinh":
		return math.Asinh(x), nil
	case "atan":
		return math.Atan(x), nil
	case "atanh":
		return math.Atanh(x), nil
	case "ceil":
		return math.Ceil(x), nil
	case "cos":
		return math.Cos(x), nil
	case "cosh":
		return math.Cosh(x), nil
	case "cot":
		return 1 / math.Tan(x), nil
	case "coth":
		return math.Cosh(x) / math.Sinh(x), nil
	case "csc":
		return 1 / math.Sin(x), nil
	case "csch":
		return 1 / math.Sinh(x), nil
	case "exp":
		return math.Exp(x), nil
	case "floor":
		return math.Floor(x), nil
	case "ln":
		return math.Log(x), nil
	case "log":
		return math.Log10(x), nil
	case "log2":
		return math.Log2(x), nil
	case "round":
		return math.Round(x), nil
	case "sec":
		return 1 / math.Cos(x), nil
	case "sech":
		return 1 / math.Cosh(x), nil
	case "sin":
		return math.Sin(x), nil
	case "sinh":
		return math.Sinh(x), nil
	case "sqrt":
		return math.Sqrt(x), nil
	case "tan":
		return math.Tan(x), nil
	case "tanh":
		return math.Tanh(x), nil
	}
	return 0, fmt.Errorf("未知的函数 %q", name)
}

// FormatNumber 将结果序列化为数值字符串（Go 最短表示，消除浮点尾差）。
func FormatNumber(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func errExpect(what string) error {
	return fmt.Errorf("缺少 %s", what)
}

func errUnexpected(tok string) error {
	if tok == "" {
		return fmt.Errorf("意外的表达式结尾")
	}
	return fmt.Errorf("意外的记号 %q", tok)
}
