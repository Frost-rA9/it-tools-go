// Package ulidgen 实现 ULID 生成工具（对齐 it-tools ulid-generator）。
package ulidgen

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "ulid-generator"
	Name        = "ULID 生成器"
	Description = "生成随机通用字典序可排序唯一标识符（ULID，128 位）"
	Category    = registry.CategoryCrypto
	Icon        = "SortDescendingNumbers"
)

// Keywords 为搜索关键词。
var Keywords = []string{"ulid", "generator", "random", "id", "identifier", "unique", "token", "string", "标识"}

// 数量上限。
const countMax = 100

// crockford 为 Crockford Base32 编码表。
const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// 输出格式。
const (
	FormatRaw  = "raw"
	FormatJSON = "json"
)

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Count  int    `json:"count"`  // 生成数量（1..100）
	Format string `json:"format"` // raw | json
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 生成 ULID 列表并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if in.Count < 1 || in.Count > countMax {
		return "", fmt.Errorf("数量必须在 1..%d 之间（当前 %d）", countMax, in.Count)
	}
	if in.Format != FormatRaw && in.Format != FormatJSON {
		return "", fmt.Errorf("未知格式: %q（仅支持 raw/json）", in.Format)
	}

	ids := make([]string, 0, in.Count)
	for i := 0; i < in.Count; i++ {
		id, err := newULID()
		if err != nil {
			return "", err
		}
		ids = append(ids, id)
	}

	var result string
	if in.Format == FormatJSON {
		b, err := json.MarshalIndent(ids, "", "  ")
		if err != nil {
			return "", fmt.Errorf("序列化输出失败: %w", err)
		}
		result = string(b)
	} else {
		result = strings.Join(ids, "\n")
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// newULID 生成一个新的 ULID：48 位毫秒时间戳（大端）+ 80 位随机数，Crockford Base32 编码为 26 字符。
func newULID() (string, error) {
	ms := uint64(time.Now().UnixMilli())
	var random [10]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("生成随机数失败: %w", err)
	}
	return encodeULID(ms, random), nil
}

// encodeULID 将毫秒时间戳与 80 位随机数编码为 26 字符 ULID。
// 时间戳部分为 48 位数值的 Crockford Base32（10 字符，高位补零），
// 随机部分为 80 位 → 16 字符，与 oklog/ulid、npm ulid 等参考实现一致。
func encodeULID(ms uint64, random [10]byte) string {
	var sb strings.Builder
	sb.Grow(26)

	// 时间戳：48 位数值 base32，10 字符。
	var tsChars [10]byte
	ts := ms
	for i := 9; i >= 0; i-- {
		tsChars[i] = crockford[ts%32]
		ts /= 32
	}
	sb.Write(tsChars[:])

	// 随机数：80 位按 5 位一组编码为 16 字符。
	var acc uint64
	accBits := 0
	byteIdx := 0
	for i := 0; i < 16; i++ {
		for accBits < 5 {
			acc = (acc << 8) | uint64(random[byteIdx])
			accBits += 8
			byteIdx++
		}
		accBits -= 5
		sb.WriteByte(crockford[(acc>>accBits)&0x1f])
	}
	return sb.String()
}
