package jsonmin

import (
	"encoding/json"
	"testing"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		want     string
		wantOrig int
		wantMin  int
		wantErr  bool
	}{
		{
			name:     "基本压缩",
			in:       `{"json":"{\n  \"a\": 1,\n  \"b\": [1, 2]\n}"}`,
			want:     `{"a":1,"b":[1,2]}`,
			wantOrig: 27, wantMin: 17,
		},
		{
			name:     "已压缩",
			in:       `{"json":"{\"a\":1}"}`,
			want:     `{"a":1}`,
			wantOrig: 7, wantMin: 7,
		},
		{name: "空对象", in: `{"json":"{ }"}`, want: `{}`},
		{name: "非法JSON", in: `{"json":"{bad}"}`, wantErr: true},
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
			if o.Minified != c.want {
				t.Fatalf("压缩结果不符: got %q, want %q", o.Minified, c.want)
			}
			if c.wantOrig > 0 && o.OriginalSize != c.wantOrig {
				t.Fatalf("原始大小不符: got %d, want %d", o.OriginalSize, c.wantOrig)
			}
			if c.wantMin > 0 && o.MinifiedSize != c.wantMin {
				t.Fatalf("压缩大小不符: got %d, want %d", o.MinifiedSize, c.wantMin)
			}
			if o.Saved != o.OriginalSize-o.MinifiedSize {
				t.Fatalf("节省大小不符: %+v", o)
			}
		})
	}
}

func TestSavedPercent(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"json":"{\n  \"a\": 1\n}"}`)
	if err != nil {
		t.Fatal(err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	if o.SavedPercent <= 0 || o.SavedPercent >= 100 {
		t.Fatalf("节省百分比异常: %v", o.SavedPercent)
	}
}
