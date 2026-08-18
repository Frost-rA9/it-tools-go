// Package sqlfmt 的关键字集合：格式化时按关键字识别子句与语法元素。
package sqlfmt

// keywords 是 SQL 关键字集合（大写）。格式化为子句起始关键字时换行，其余行内。
var keywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true, "AND": true, "OR": true,
	"NOT": true, "IN": true, "LIKE": true, "IS": true, "NULL": true,
	"GROUP": true, "BY": true, "ORDER": true, "ASC": true, "DESC": true,
	"HAVING": true, "LIMIT": true, "OFFSET": true, "DISTINCT": true,
	"JOIN": true, "LEFT": true, "RIGHT": true, "FULL": true, "INNER": true,
	"OUTER": true, "CROSS": true, "ON": true, "UNION": true, "ALL": true,
	"INSERT": true, "INTO": true, "VALUES": true, "UPDATE": true, "SET": true,
	"DELETE": true, "CREATE": true, "TABLE": true, "ALTER": true, "DROP": true,
	"AS": true, "CASE": true, "WHEN": true, "THEN": true, "ELSE": true, "END": true,
	"EXISTS": true, "WITH": true, "PRIMARY": true, "KEY": true, "FOREIGN": true,
	"REFERENCES": true, "INDEX": true, "TRUNCATE": true, "GRANT": true, "REVOKE": true,
}

// clauseKeywords 是子句起始关键字：位于行首（换行触发）。
var clauseKeywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true, "GROUP": true, "ORDER": true,
	"HAVING": true, "LIMIT": true, "OFFSET": true, "UNION": true,
	"JOIN": true, "LEFT": true, "RIGHT": true, "FULL": true, "INNER": true,
	"OUTER": true, "CROSS": true, "ON": true,
	"INSERT": true, "VALUES": true, "UPDATE": true, "SET": true,
	"DELETE": true, "CREATE": true, "TABLE": true, "ALTER": true, "DROP": true,
	"WITH": true, "TRUNCATE": true, "GRANT": true, "REVOKE": true,
}
