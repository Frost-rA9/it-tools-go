package sqlfmt

import (
	"fmt"
	"strings"
)

// tokenType 是 SQL token 的类型。
type tokenType int

const (
	tokIdent tokenType = iota
	tokKeyword
	tokString
	tokQuotedIdent
	tokComment
	tokNumber
	tokOperator
	tokLParen
	tokRParen
	tokComma
	tokSemicolon
	tokDot
)

// token 是词法单元：保留原始文本与大写形式。
type token struct {
	typ  tokenType
	text string
	// upper 是关键字的规范大写形式（用于子句匹配与输出）。
	upper string
}

// tokenize 将 SQL 文本切分为 token 序列。
// 支持的语法：字符串（'...'，” 转义，N'/E'/x' 前缀）、双引号标识符（"" 转义）、
// -- 行注释、/* */ 块注释、数字、标识符、关键字、运算符、括号、逗号、分号。
func tokenize(sql string) ([]token, error) {
	var toks []token
	i := 0
	n := len(sql)
	for i < n {
		c := sql[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c == '\'':
			t, next, err := scanQuoted(sql, i, '\'')
			if err != nil {
				return nil, err
			}
			toks = append(toks, token{typ: tokString, text: t})
			i = next
		case c == '"':
			t, next, err := scanQuoted(sql, i, '"')
			if err != nil {
				return nil, err
			}
			toks = append(toks, token{typ: tokQuotedIdent, text: t})
			i = next
		case c == '-' && i+1 < n && sql[i+1] == '-':
			end := i + 2
			for end < n && sql[end] != '\n' {
				end++
			}
			toks = append(toks, token{typ: tokComment, text: sql[i:end]})
			i = end
		case c == '/' && i+1 < n && sql[i+1] == '*':
			end := strings.Index(sql[i+2:], "*/")
			if end < 0 {
				return nil, fmt.Errorf("块注释未闭合")
			}
			end += i + 2 + 2
			toks = append(toks, token{typ: tokComment, text: sql[i:end]})
			i = end
		case c == '.':
			toks = append(toks, token{typ: tokDot, text: "."})
			i++
		case c == '(':
			toks = append(toks, token{typ: tokLParen, text: "("})
			i++
		case c == ')':
			toks = append(toks, token{typ: tokRParen, text: ")"})
			i++
		case c == ',':
			toks = append(toks, token{typ: tokComma, text: ","})
			i++
		case c == ';':
			toks = append(toks, token{typ: tokSemicolon, text: ";"})
			i++
		case isDigit(c) || (c == '.' && i+1 < n && isDigit(sql[i+1])):
			start := i
			for i < n && (isDigit(sql[i]) || sql[i] == '.') {
				i++
			}
			toks = append(toks, token{typ: tokNumber, text: sql[start:i]})
		case isIdentStart(c):
			// 标识符 / 关键字；支持 N'...'、E'...'、x'...' 字符串前缀
			if i+1 < n && sql[i+1] == '\'' && (c == 'N' || c == 'E' || c == 'n' || c == 'e' || c == 'x' || c == 'X') {
				t, next, err := scanQuoted(sql, i+1, '\'')
				if err != nil {
					return nil, err
				}
				toks = append(toks, token{typ: tokString, text: sql[i:i+1] + t})
				i = next
				continue
			}
			start := i
			for i < n && isIdentPart(sql[i]) {
				i++
			}
			word := sql[start:i]
			upper := strings.ToUpper(word)
			if keywords[upper] {
				toks = append(toks, token{typ: tokKeyword, text: word, upper: upper})
			} else {
				toks = append(toks, token{typ: tokIdent, text: word})
			}
		default:
			// 运算符或单字符
			op := operator(sql, i)
			toks = append(toks, token{typ: tokOperator, text: op})
			i += len(op)
		}
	}
	return toks, nil
}

// scanQuoted 扫描引号包裹的内容（quote 为 ' 或 "），支持双写转义，返回含引号的完整文本。
func scanQuoted(sql string, start int, quote byte) (string, int, error) {
	i := start + 1
	n := len(sql)
	for i < n {
		if sql[i] == quote {
			if i+1 < n && sql[i+1] == quote {
				i += 2
				continue
			}
			return sql[start : i+1], i + 1, nil
		}
		i++
	}
	return "", 0, fmt.Errorf("引号未闭合")
}

// operator 识别多字符运算符。
func operator(sql string, i int) string {
	two := ""
	if i+1 < len(sql) {
		two = sql[i : i+2]
	}
	switch two {
	case "<=", ">=", "<>", "!=", "||":
		return two
	}
	return sql[i : i+1]
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdentStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}
