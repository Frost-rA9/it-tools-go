package jsoncsv

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wantCSV  string
		wantRows int
		wantCols int
		wantErr  bool
	}{
		{
			name:     "基本转换",
			in:       `{"json":"[{\"name\":\"Alice\",\"age\":30},{\"name\":\"Bob\",\"age\":25}]"}`,
			wantCSV:  "age,name\n30,Alice\n25,Bob\n",
			wantRows: 2, wantCols: 2,
		},
		{
			name:     "无表头",
			in:       `{"json":"[{\"a\":1},{\"a\":2}]","include_header":false}`,
			wantCSV:  "1\n2\n",
			wantRows: 2, wantCols: 1,
		},
		{
			name:     "逗号转义",
			in:       `{"json":"[{\"city\":\"Bei, Jing\"}]"}`,
			wantCSV:  "city\n\"Bei, Jing\"\n",
			wantRows: 1, wantCols: 1,
		},
		{
			name:     "引号与换行转义",
			in:       `{"json":"[{\"text\":\"say \\\"hi\\\"\\nnext\"}]"}`,
			wantCSV:  "text\n\"say \"\"hi\"\"\nnext\"\n",
			wantRows: 1, wantCols: 1,
		},
		{
			name:     "嵌套对象与数组",
			in:       `{"json":"[{\"id\":1,\"tags\":[\"a\",\"b\"],\"meta\":{\"x\":1}}]"}`,
			wantCSV:  "id,meta,tags\n1,\"{\"\"x\"\":1}\",\"[\"\"a\"\",\"\"b\"\"]\"\n",
			wantRows: 1, wantCols: 3,
		},
		{
			name:     "null与数字精度",
			in:       `{"json":"[{\"a\":null,\"b\":12345678901234567890}]"}`,
			wantCSV:  "a,b\n,12345678901234567890\n",
			wantRows: 1, wantCols: 2,
		},
		{
			name:     "分号分隔",
			in:       `{"json":"[{\"a\":1,\"b\":2}]","delimiter":";"}`,
			wantCSV:  "a;b\n1;2\n",
			wantRows: 1, wantCols: 2,
		},
		{
			name:     "空数组",
			in:       `{"json":"[]"}`,
			wantCSV:  "",
			wantRows: 0, wantCols: 0,
		},
		{name: "非法JSON", in: `{"json":"{bad}"}`, wantErr: true},
		{name: "非对象数组", in: `{"json":"[1,2,3]"}`, wantErr: true},
		{name: "非数组", in: `{"json":"{\"a\":1}"}`, wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			exec := Executor{}
			out, err := exec.Execute(t.Context(), c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际成功: %s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			var o output
			if err := json.Unmarshal([]byte(out), &o); err != nil {
				t.Fatalf("解析输出失败: %v", err)
			}
			if o.CSV != c.wantCSV {
				t.Fatalf("CSV 不符:\ngot  %q\nwant %q", o.CSV, c.wantCSV)
			}
			if o.Rows != c.wantRows || o.Columns != c.wantCols {
				t.Fatalf("行列不符: rows=%d cols=%d, want rows=%d cols=%d", o.Rows, o.Columns, c.wantRows, c.wantCols)
			}
		})
	}
}

func TestCellString(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, ""},
		{"hello", "hello"},
		{json.Number("42"), "42"},
		{json.Number("3.14"), "3.14"},
		{true, "true"},
		{[]any{1, 2}, "[1,2]"},
		{map[string]any{"x": 1}, `{"x":1}`},
	}
	for _, c := range cases {
		if got := cellString(c.in); got != c.want {
			t.Fatalf("cellString(%v) = %q, want %q", c.in, got, c.want)
		}
	}
	// float64 兼容（直接构造场景）
	if got := cellString(float64(3.5)); got != "3.5" {
		t.Fatalf("float64 转换失败: %q", got)
	}
	if !strings.Contains(cellString(3.5), "3.5") {
		t.Fatal("float64 转换异常")
	}
}
