package chmodcalc

import (
	"encoding/json"
	"testing"
)

func TestParseModeOctal(t *testing.T) {
	cases := []struct {
		in      string
		want    uint32
		wantErr bool
	}{
		{in: "755", want: 0o755},
		{in: "000", want: 0},
		{in: "777", want: 0o777},
		{in: "4755", want: 0o4755},
		{in: "2755", want: 0o2755},
		{in: "1777", want: 0o1777},
		{in: "6755", want: 0o6755},
		{in: "644", want: 0o644},
		{in: "888", wantErr: true},
		{in: "75", wantErr: true},
		{in: "77777", wantErr: true},
		{in: "", wantErr: true},
		{in: "abc", wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, err := parseMode(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际成功: %o", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("解析不符: got %o, want %o", got, c.want)
			}
		})
	}
}

func TestParseModeSymbolic(t *testing.T) {
	cases := []struct {
		in   string
		want uint32
	}{
		{in: "rwxr-xr-x", want: 0o755},
		{in: "-rwxr-xr-x", want: 0o755},
		{in: "drwxr-xr-x", want: 0o755},
		{in: "rw-------", want: 0o600},
		{in: "rwxrwxrwx", want: 0o777},
		{in: "rwsr-xr-x", want: 0o4755},
		{in: "rwsr-sr-x", want: 0o6755},
		{in: "rwxr-xr-t", want: 0o1755},
		{in: "rwS------", want: 0o4600}, // SUID（无执行权限）+ owner rw
		{in: "r-xr-x---", want: 0o550},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, err := parseMode(c.in)
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("解析不符: got %o, want %o", got, c.want)
			}
		})
	}
}

func TestParseModeSymbolicErrors(t *testing.T) {
	bad := []string{"rwx", "rrrrrrrrr", "rwsr-t---", "rwxrwxrws", "rwxrwxrwxrwx", "-rwxr-xr-"}
	for _, s := range bad {
		t.Run(s, func(t *testing.T) {
			if _, err := parseMode(s); err == nil {
				t.Fatalf("期望报错: %q", s)
			}
		})
	}
	// "rwsr-t---": 组 x 位出现 't' 非法
	// "rwxrwxrws": 其他用户 x 位出现 's' 非法
	// "-rwxr-xr-": 长度不足 9 位
}

func TestRoundTrip(t *testing.T) {
	for _, m := range []uint32{0, 0o755, 0o644, 0o4755, 0o6755, 0o1777, 0o4000, 0o2000, 0o1000} {
		s := symbolicString(m)
		back, err := parseMode(s)
		if err != nil {
			t.Fatalf("%o → %q 解析失败: %v", m, s, err)
		}
		if back != m {
			t.Fatalf("往返不符: %o → %q → %o", m, s, back)
		}
	}
}

func TestBuildOutput(t *testing.T) {
	o := buildOutput(0o4755)
	if o.Octal != "4755" || o.Symbolic != "rwsr-xr-x" {
		t.Fatalf("输出不符: %+v", o)
	}
	if o.Owner != "rws" || o.Group != "r-x" || o.Others != "r-x" {
		t.Fatalf("分组不符: %+v", o)
	}
	if !o.HasSUID || o.HasSGID || o.HasSticky || o.Special != 4 {
		t.Fatalf("特殊位不符: %+v", o)
	}
	if o.SpecialText == "" {
		t.Fatal("特殊位说明为空")
	}

	o2 := buildOutput(0o755)
	if o2.Octal != "755" || o2.SpecialText != "无特殊位" {
		t.Fatalf("普通权限输出不符: %+v", o2)
	}

	o3 := buildOutput(0o7000)
	if o3.Symbolic != "--S--S--T" {
		t.Fatalf("仅特殊位符号串不符: %q", o3.Symbolic)
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"mode":"755"}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if o.Octal != "755" || o.Symbolic != "rwxr-xr-x" {
		t.Fatalf("输出不符: %+v", o)
	}

	if _, err := exec.Execute(t.Context(), `{"mode":"999"}`); err == nil {
		t.Fatal("非法输入应报错")
	}
}
