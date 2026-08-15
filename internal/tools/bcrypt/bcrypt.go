// Package bcrypt 实现 BCrypt 哈希与对比工具（对齐 it-tools bcrypt）。
package bcrypt

import (
	"context"
	"encoding/json"
	"fmt"

	xbcrypt "golang.org/x/crypto/bcrypt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "bcrypt"
	Name        = "BCrypt 加密"
	Description = "使用 BCrypt（基于 Blowfish 密码的密码哈希函数）对文本进行哈希并与哈希对比"
	Category    = registry.CategoryCrypto
	Icon        = "LockSquare"
)

// Keywords 为搜索关键词。
var Keywords = []string{"bcrypt", "hash", "compare", "password", "salt", "round", "storage", "crypto", "加密", "对比"}

// 转换模式。
const (
	ModeHash    = "hash"
	ModeCompare = "compare"
)

// defaultCost 为未指定时的默认 salt 轮数。
const defaultCost = 10

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Mode string `json:"mode"` // hash | compare
	Text string `json:"text"` // 待哈希字符串或待对比字符串
	Cost int    `json:"cost"` // salt 轮数（4..31，默认 10），仅 hash 模式
	Hash string `json:"hash"` // 待对比的 BCrypt 哈希，仅 compare 模式
}

// hashOutput 是 hash 模式的输出结构。
type hashOutput struct {
	Result string `json:"result"`
}

// compareOutput 是 compare 模式的输出结构。
type compareOutput struct {
	Match bool `json:"match"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 执行 BCrypt 哈希或对比并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	switch in.Mode {
	case ModeHash:
		cost := in.Cost
		if cost == 0 {
			cost = defaultCost
		}
		if cost < xbcrypt.MinCost || cost > xbcrypt.MaxCost {
			return "", fmt.Errorf("salt 轮数必须在 %d..%d 之间（当前 %d）", xbcrypt.MinCost, xbcrypt.MaxCost, cost)
		}
		hash, err := xbcrypt.GenerateFromPassword([]byte(in.Text), cost)
		if err != nil {
			return "", fmt.Errorf("生成 BCrypt 哈希失败: %w", err)
		}
		out, err := json.Marshal(hashOutput{Result: string(hash)})
		if err != nil {
			return "", fmt.Errorf("序列化输出失败: %w", err)
		}
		return string(out), nil
	case ModeCompare:
		match := xbcrypt.CompareHashAndPassword([]byte(in.Hash), []byte(in.Text)) == nil
		out, err := json.Marshal(compareOutput{Match: match})
		if err != nil {
			return "", fmt.Errorf("序列化输出失败: %w", err)
		}
		return string(out), nil
	default:
		return "", fmt.Errorf("未知模式: %q（仅支持 hash/compare）", in.Mode)
	}
}
