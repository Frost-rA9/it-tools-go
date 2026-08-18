// Package sqlfmt 实现 SQL 格式化工具：关键字大写、子句换行、括号缩进。
package sqlfmt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/registry"
)

const (
	ID          = "sql-formatter"
	Name        = "SQL 格式化"
	Description = "美化 SQL：关键字大写、子句换行、括号缩进"
	Category    = "开发"
	Icon        = "Database"
)

var Keywords = []string{"sql", "format", "pretty", "美化", "格式化", "缩进", "select", "database"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	SQL           string `json:"sql"`            // 待格式化 SQL
	UpperKeywords *bool  `json:"upper_keywords"` // 关键字是否大写（默认 true）
}

// output 是工具的输出结构。
type output struct {
	Formatted string `json:"formatted"`
	LineCount int    `json:"line_count"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回格式化结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	upper := true
	if in.UpperKeywords != nil {
		upper = *in.UpperKeywords
	}

	formatted, err := Format(in.SQL, upper)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{
		Formatted: formatted,
		LineCount: strings.Count(strings.TrimRight(formatted, "\n"), "\n") + 1,
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Format 对 SQL 文本执行格式化。upper 控制关键字是否大写。
func Format(sql string, upper bool) (string, error) {
	if strings.TrimSpace(sql) == "" {
		return "", fmt.Errorf("SQL 为空")
	}
	toks, err := tokenize(sql)
	if err != nil {
		return "", fmt.Errorf("SQL 解析失败: %w", err)
	}
	return render(toks, upper)
}

// render 按 token 流生成格式化输出，并校验括号闭合。
func render(toks []token, upper bool) (string, error) {
	var sb strings.Builder
	depth := 0
	lineStart := true
	needSpace := false
	selectList := false // SELECT 与 FROM 之间（字段列表逗号缩进 2）
	inWhere := false    // WHERE 子句内：AND/OR 换行缩进 2
	inJoin := false     // JOIN 子句内：ON 换行缩进 2
	var last *token

	write := func(s string) {
		sb.WriteString(s)
	}
	newline := func(level int) {
		sb.WriteByte('\n')
		sb.WriteString(strings.Repeat("  ", level))
		lineStart = true
		needSpace = false
	}

	for i := range toks {
		t := toks[i]
		switch t.typ {
		case tokKeyword:
			if (t.upper == "AND" || t.upper == "OR") && inWhere && !lineStart {
				newline(1)
			} else if t.upper == "ON" && inJoin && !lineStart {
				newline(1)
			} else if clauseKeywords[t.upper] && !lineStart {
				newline(depth)
			}
			if needSpace && !lineStart {
				write(" ")
			}
			if upper {
				write(t.upper)
			} else {
				write(t.text)
			}
			lineStart = false
			needSpace = true
			switch t.upper {
			case "SELECT":
				selectList = true
				inWhere = false
				inJoin = false
			case "WHERE":
				selectList = false
				inWhere = true
				inJoin = false
			case "GROUP", "ORDER", "HAVING", "LIMIT", "OFFSET", "UNION", "VALUES", "SET", "UPDATE", "DELETE":
				selectList = false
				inWhere = false
				inJoin = false
			case "JOIN", "LEFT", "RIGHT", "FULL", "INNER", "OUTER", "CROSS":
				selectList = false
				inWhere = false
				inJoin = true
			case "FROM", "ON":
				selectList = false
				inWhere = false
				inJoin = false
			}
		case tokIdent, tokNumber:
			if needSpace && !lineStart {
				write(" ")
			}
			write(t.text)
			lineStart = false
			needSpace = true
		case tokString, tokQuotedIdent:
			if needSpace && !lineStart {
				write(" ")
			}
			write(t.text)
			lineStart = false
			needSpace = true
		case tokComment:
			if needSpace && !lineStart {
				write(" ")
			}
			write(t.text)
			lineStart = false
			needSpace = true
			// 行注释后强制换行
			if strings.HasPrefix(t.text, "--") {
				newline(depth)
			}
		case tokOperator:
			// 通配符 *（行首或紧跟左括号）不加空格；其余运算符两侧空格
			if t.text == "*" && (lineStart || (last != nil && last.typ == tokLParen)) {
				write("*")
			} else {
				if !lineStart {
					write(" ")
				}
				write(t.text)
			}
			lineStart = false
			needSpace = true
		case tokDot:
			// 限定符（u.id）：两侧不加空格
			write(".")
			lineStart = false
			needSpace = false
		case tokLParen:
			// 函数调用/表名后紧跟；IN/LIKE/EXISTS/ON 等语法关键字后加空格
			space := false
			if last != nil && (last.upper == "IN" || last.upper == "LIKE" || last.upper == "EXISTS" ||
				last.upper == "ON" || last.upper == "NOT" || last.typ == tokOperator) {
				space = true
			}
			if space && !lineStart {
				write(" ")
			}
			write("(")
			lineStart = false
			needSpace = false
			depth++
		case tokRParen:
			if depth > 0 {
				depth--
			}
			write(")")
			lineStart = false
			needSpace = false
		case tokComma:
			write(",")
			if depth == 0 {
				// 顶层逗号 → 换行：SELECT 列表缩进 2，其余（VALUES 元组间）缩进当前深度
				if selectList {
					newline(1)
				} else {
					newline(depth)
				}
			} else {
				needSpace = true
				lineStart = false
			}
		case tokSemicolon:
			write(";")
			newline(0)
		}
		last = &toks[i]
	}

	if depth != 0 {
		return "", fmt.Errorf("括号未闭合")
	}
	return strings.TrimRight(sb.String(), " \t\n") + "\n", nil
}
