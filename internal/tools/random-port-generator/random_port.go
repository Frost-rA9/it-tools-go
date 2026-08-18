// Package randomport 实现随机端口生成器：按范围与数量生成随机端口号。
package randomport

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sort"

	"it-tools-go/internal/registry"
)

const (
	ID          = "random-port-generator"
	Name        = "随机端口生成器"
	Description = "按范围与数量生成随机端口号，支持排除列表"
	Category    = "开发"
	Icon        = "Plug"
)

var Keywords = []string{"random", "port", "端口", "生成器", "tcp", "udp", "随机", "network"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Count   int   `json:"count"`   // 生成数量，默认 1
	Min     int   `json:"min"`     // 最小值，默认 1024
	Max     int   `json:"max"`     // 最大值，默认 65535
	Exclude []int `json:"exclude"` // 排除端口列表（如 ["80","443"]）
}

// output 是工具的输出结构。
type output struct {
	Ports []int `json:"ports"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回生成的端口列表 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	ports, err := generate(in)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Ports: ports})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// generate 校验参数并从可用端口池中随机取样，结果升序返回。
func generate(in input) ([]int, error) {
	if in.Count == 0 {
		in.Count = 1
	}
	if in.Max == 0 {
		in.Max = 65535
	}
	if in.Min == 0 {
		in.Min = 1024
	}

	if in.Count < 1 {
		return nil, fmt.Errorf("数量必须大于 0（当前 %d）", in.Count)
	}
	if in.Min < 0 || in.Max > 65535 {
		return nil, fmt.Errorf("端口范围必须在 0-65535 之间")
	}
	if in.Min > in.Max {
		return nil, fmt.Errorf("最小值不得大于最大值（%d > %d）", in.Min, in.Max)
	}

	excluded := make(map[int]bool, len(in.Exclude))
	for _, p := range in.Exclude {
		if p >= in.Min && p <= in.Max {
			excluded[p] = true
		}
	}

	available := in.Max - in.Min + 1 - len(excluded)
	if available < in.Count {
		return nil, fmt.Errorf("可用端口不足：范围内共 %d 个，排除 %d 个，需要 %d 个",
			in.Max-in.Min+1, len(excluded), in.Count)
	}

	pool := make([]int, 0, available)
	for p := in.Min; p <= in.Max; p++ {
		if !excluded[p] {
			pool = append(pool, p)
		}
	}

	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	ports := make([]int, in.Count)
	copy(ports, pool[:in.Count])
	sort.Ints(ports)
	return ports, nil
}
