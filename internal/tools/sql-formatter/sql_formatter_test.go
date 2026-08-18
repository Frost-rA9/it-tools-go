package sqlfmt

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "基本SELECT",
			sql:  "select id,name from users where age>18 and city='beijing' order by id desc;",
			want: "SELECT id,\n  name\nFROM users\nWHERE age > 18\n  AND city = 'beijing'\nORDER BY id DESC;\n",
		},
		{
			name: "JOIN",
			sql:  "select u.id,o.amount from users u join orders o on u.id=o.user_id where u.active=true;",
			want: "SELECT u.id,\n  o.amount\nFROM users u\nJOIN orders o\n  ON u.id = o.user_id\nWHERE u.active = true;\n",
		},
		{
			name: "INSERT VALUES",
			sql:  "insert into t(a,b) values(1,'x'),(2,'y');",
			want: "INSERT INTO t(a, b)\nVALUES(1, 'x'),\n(2, 'y');\n",
		},
		{
			name: "子查询",
			sql:  "select * from t where id in (select id from u where x=1);",
			want: "SELECT *\nFROM t\nWHERE id IN (\n  SELECT id\n  FROM u\n  WHERE x = 1);\n",
		},
		{
			name: "行注释",
			sql:  "select a -- 注释\nfrom t;",
			want: "SELECT a -- 注释\nFROM t;\n",
		},
		{
			name: "块注释",
			sql:  "select a /* 注释 */ from t;",
			want: "SELECT a /* 注释 */\nFROM t;\n",
		},
		{
			name: "字符串含关键字逗号",
			sql:  "select 'select, from where' as s;",
			want: "SELECT 'select, from where' AS s;\n",
		},
		{
			name: "函数与IN列表",
			sql:  "select count(*),max(a) from t where b in (1,2,3);",
			want: "SELECT count(*),\n  max(a)\nFROM t\nWHERE b IN (1, 2, 3);\n",
		},
		{
			name: "多语句",
			sql:  "select 1; select 2;",
			want: "SELECT 1;\nSELECT 2;\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Format(c.sql, true)
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("格式化不符:\ngot:\n%s\nwant:\n%s", got, c.want)
			}
		})
	}
}

func TestFormatErrors(t *testing.T) {
	bad := []string{"", "   ", "select 'unterminated", "select 1; /* unclosed", "select (1"}
	for _, sql := range bad {
		t.Run(sql, func(t *testing.T) {
			if _, err := Format(sql, true); err == nil {
				t.Fatalf("期望报错: %q", sql)
			}
		})
	}
}

func TestFormatLowerKeywords(t *testing.T) {
	got, err := Format("select a from t;", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "select a") || !strings.Contains(got, "from t") {
		t.Fatalf("关键字原样保留失败: %s", got)
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"sql":"select 1"}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if o.Formatted != "SELECT 1\n" || o.LineCount != 1 {
		t.Fatalf("输出不符: %+v", o)
	}
	if _, err := exec.Execute(t.Context(), `{"sql":"bad '"`); err == nil {
		t.Fatal("非法输入应报错")
	}
}
