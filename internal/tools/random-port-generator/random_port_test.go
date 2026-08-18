package randomport

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		name    string
		in      input
		wantCnt int
		wantErr bool
	}{
		{name: "默认值", in: input{}, wantCnt: 1, wantErr: false},
		{name: "多端口", in: input{Count: 10, Min: 1024, Max: 65535}, wantCnt: 10, wantErr: false},
		{name: "边界全范围", in: input{Count: 5, Min: 0, Max: 65535}, wantCnt: 5, wantErr: false},
		{name: "排除端口", in: input{Count: 3, Min: 1, Max: 10, Exclude: []int{1, 2, 3, 4, 5, 6, 7}}, wantCnt: 3, wantErr: false},
		{name: "count为0取默认1", in: input{Count: 0, Min: 1, Max: 10}, wantCnt: 1, wantErr: false},
		{name: "数量为负报错", in: input{Count: -1}, wantErr: true},
		{name: "min越界报错", in: input{Min: -1}, wantErr: true},
		{name: "max越界报错", in: input{Max: 70000}, wantErr: true},
		{name: "min大于max报错", in: input{Min: 5000, Max: 1000}, wantErr: true},
		{name: "可用数不足报错", in: input{Count: 10, Min: 1, Max: 5}, wantErr: true},
		{name: "排除全部报错", in: input{Count: 1, Min: 1, Max: 3, Exclude: []int{1, 2, 3}}, wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := generate(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际成功: %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if len(got) != c.wantCnt {
				t.Fatalf("数量不符: got %d, want %d", len(got), c.wantCnt)
			}
			// 结果必在当前范围且不包含排除项
			for _, p := range got {
				if c.in.Min == 0 && c.in.Max == 0 {
					if p < 1024 || p > 65535 {
						t.Fatalf("端口 %d 超出默认范围 [1024, 65535]", p)
					}
				} else if p < c.in.Min || p > c.in.Max {
					t.Fatalf("端口 %d 超出范围 [%d, %d]", p, c.in.Min, c.in.Max)
				}
				for _, e := range c.in.Exclude {
					if p == e {
						t.Fatalf("端口 %d 不应出现在排除列表中", p)
					}
				}
			}
			// 结果去重且升序
			for i := 1; i < len(got); i++ {
				if got[i] <= got[i-1] {
					t.Fatalf("结果未去重/未排序: %v", got)
				}
			}
		})
	}
}

func TestGenerateUniquePorts(t *testing.T) {
	// 范围极小且全量取时，结果应为该范围内全部端口（验证去重取样）。
	got, err := generate(input{Count: 3, Min: 100, Max: 102})
	if err != nil {
		t.Fatalf("意外错误: %v", err)
	}
	want := []int{100, 101, 102}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("结果不符: got %v, want %v", got, want)
	}
}

func TestExecuteJSON(t *testing.T) {
	exec := Executor{}
	out, err := exec.Execute(t.Context(), `{"count":2,"min":1000,"max":1005,"exclude":[1001]}`)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	var o output
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if len(o.Ports) != 2 {
		t.Fatalf("端口数量不符: got %v", o.Ports)
	}
	for _, p := range o.Ports {
		if p < 1000 || p > 1005 || p == 1001 {
			t.Fatalf("端口 %d 不合法", p)
		}
	}

	// 非法 JSON 输入应报错
	if _, err := exec.Execute(t.Context(), `not-json`); err == nil {
		t.Fatal("非法 JSON 应报错")
	}
}
