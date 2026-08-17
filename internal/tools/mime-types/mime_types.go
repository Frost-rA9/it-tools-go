// Package mimetypes 实现 MIME 类型转换器：查询 MIME 类型关联的文件扩展名，及扩展名关联的 MIME 类型。
// 数据源为 mime-db（npm mime-types 同源），以 JSON 嵌入。
package mimetypes

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "mime-types"
	Name        = "MIME 类型转换器"
	Description = "查询 MIME 类型与文件扩展名之间的对应关系"
	Category    = "Web"
	Icon        = "World"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"mime", "types", "extension", "content", "type", "MIME", "扩展名", "类型"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

//go:embed mime-db.json
var mimeDBData []byte

// mimeDBEntry 是 mime-db.json 中单个 MIME 类型的条目（仅使用 extensions 字段）。
type mimeDBEntry struct {
	Extensions []string `json:"extensions"`
}

// mimeToExtensions 保存 MIME 类型 → 扩展名列表（仅含带扩展名的类型）。
var mimeToExtensions map[string][]string

// extensionToMime 保存扩展名 → MIME 类型（取首个声明该扩展名的 MIME 类型，对齐 mime-types）。
var extensionToMime map[string]string

func init() {
	var db map[string]mimeDBEntry
	if err := json.Unmarshal(mimeDBData, &db); err != nil {
		panic("解析 mime-db.json 失败: " + err.Error())
	}

	mimeToExtensions = make(map[string][]string)
	extensionToMime = make(map[string]string)

	// 按键排序保证确定性（首个扩展名归属与 mime-types 一致）。
	keys := make([]string, 0, len(db))
	for k := range db {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, mime := range keys {
		entry := db[mime]
		if len(entry.Extensions) == 0 {
			continue
		}
		mimeToExtensions[mime] = entry.Extensions
		for _, ext := range entry.Extensions {
			if _, exists := extensionToMime[ext]; !exists {
				extensionToMime[ext] = mime
			}
		}
	}
}

// input 是工具的输入结构。
type input struct {
	List        bool   `json:"list"`         // 返回全部 MIME ↔ 扩展名对
	MimeType    string `json:"mime_type"`    // 查询该 MIME 类型的扩展名
	Extension   string `json:"extension"`    // 查询该扩展名的 MIME 类型
}

// mimePair 表示一个 MIME 类型与其扩展名。
type mimePair struct {
	MimeType   string   `json:"mime_type"`
	Extensions []string `json:"extensions"`
}

// output 是工具的输出结构。
type output struct {
	All        []mimePair `json:"all,omitempty"`
	Extensions []string   `json:"extensions,omitempty"`
	MimeTypes  []string   `json:"mime_types,omitempty"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	var out output
	switch {
	case in.List:
		out.All = allPairs()
	case in.MimeType != "":
		out.Extensions = mimeToExtensions[in.MimeType]
	case in.Extension != "":
		ext := normalizeExt(in.Extension)
		if mime, ok := extensionToMime[ext]; ok {
			out.MimeTypes = []string{mime}
		}
	default:
		return "", fmt.Errorf("缺少查询参数（list/mime_type/extension 至少一个）")
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// allPairs 返回全部带扩展名的 MIME 类型及其扩展名（按 MIME 类型排序）。
func allPairs() []mimePair {
	keys := make([]string, 0, len(mimeToExtensions))
	for k := range mimeToExtensions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]mimePair, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, mimePair{MimeType: k, Extensions: mimeToExtensions[k]})
	}
	return pairs
}

// normalizeExt 规范化扩展名：去掉前导点并转为小写。
func normalizeExt(ext string) string {
	return strings.TrimLeft(strings.ToLower(ext), ".")
}