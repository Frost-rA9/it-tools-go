// Package metatag 实现开放式图形元生成器：动态表单生成 Open Graph / Twitter 等 HTML meta 标签。
// 逻辑对齐 it-tools 的 meta-tag-generator 与 @it-tools/oggen。
package metatag

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "meta-tag-generator"
	Name        = "开放式图形元生成器"
	Description = "通过表单生成 Open Graph 与 Twitter 等社交分享 meta 标签"
	Category    = "Web"
	Icon        = "Tags"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"meta", "tag", "generator", "social", "title", "description", "image", "share", "website", "open", "graph", "og", "元", "标签", "分享"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Action   string         `json:"action"` // schemas | meta
	Metadata map[string]any `json:"metadata"`
}

// output 是工具的输出结构。
type output struct {
	Schemas *schemasResult `json:"schemas,omitempty"`
	Meta    string         `json:"meta,omitempty"`
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
	switch in.Action {
	case "schemas":
		s := schemas()
		out = output{Schemas: &s}
	case "meta":
		meta, err := generateMetaFromMetadata(in.Metadata)
		if err != nil {
			return "", err
		}
		out = output{Meta: meta}
	default:
		return "", fmt.Errorf("未知操作: %q（仅支持 schemas/meta）", in.Action)
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}

// generateMetaFromMetadata 将完整 metadata 拆分为 og 与 twitter 两部分并生成 meta 标签。
func generateMetaFromMetadata(metadata map[string]any) (string, error) {
	if metadata == nil {
		return "", fmt.Errorf("缺少 metadata")
	}

	ogPart := map[string]any{}
	twitterPart := map[string]any{}
	for k, v := range metadata {
		if strings.HasPrefix(k, "twitter:") {
			twitterPart[strings.TrimPrefix(k, "twitter:")] = v
		} else {
			ogPart[k] = v
		}
	}

	return generateMeta(ogPart, twitterPart, true), nil
}

// metaFlat 表示扁平化后的单个键值对。
type metaFlat struct {
	Key   string
	Value string
}

// generateMeta 对齐 oggen 的 generateMeta。
func generateMeta(ogMetadata, twitterMetadata map[string]any, generateTwitterCompatible bool) string {
	ogFlat := flattenMetadata(ogMetadata, "og")
	twitterFlat := flattenMetadata(twitterMetadata, "twitter")

	twitterFinal := append([]metaFlat{}, twitterFlat...)
	if generateTwitterCompatible {
		twitterFinal = append(twitterFinal, pickTwitterCompatible(ogFlat, twitterFlat)...)
	}

	var groups []string
	if len(ogFlat) > 0 {
		groups = append(groups, "<!-- og meta -->\n"+buildMetaStrings(ogFlat, "property"))
	}
	if len(twitterFinal) > 0 {
		groups = append(groups, "<!-- twitter meta -->\n"+buildMetaStrings(twitterFinal, "name"))
	}
	return strings.Join(groups, "\n\n")
}

// flattenMetadata 递归扁平化元数据（对齐 oggen 的 flattenMetadata）。
func flattenMetadata(metadata map[string]any, basePrefix string) []metaFlat {
	var acc []metaFlat

	var walk func(node any, prefix string)
	walk = func(node any, prefix string) {
		switch v := node.(type) {
		case nil:
			return
		case string:
			if v == "" {
				return
			}
			acc = append(acc, metaFlat{Key: prefix, Value: v})
		case bool:
			acc = append(acc, metaFlat{Key: prefix, Value: strconv.FormatBool(v)})
		case float64:
			acc = append(acc, metaFlat{Key: prefix, Value: formatNumber(v)})
		case []any:
			for _, item := range v {
				walk(item, prefix)
			}
		case map[string]any:
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				prefixed := k
				if prefix != "" {
					prefixed = prefix + ":" + toSnakeCase(k)
				}
				walk(v[k], prefixed)
			}
		default:
			acc = append(acc, metaFlat{Key: prefix, Value: fmt.Sprintf("%v", v)})
		}
	}

	walk(metadata, basePrefix)
	return acc
}

// formatNumber 将 JSON 数值格式化为字符串（整数不带小数点）。
func formatNumber(f float64) string {
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// toSnakeCase 将对冒号分隔的各段转成 snake_case（对齐 oggen 的 toSnakeCase）。
func toSnakeCase(s string) string {
	parts := strings.Split(s, ":")
	for i, p := range parts {
		parts[i] = toSnakeCaseStrict(p)
	}
	return strings.Join(parts, ":")
}

// toSnakeCaseStrict 将驼峰命名的单段转成 snake_case。
// 注意：Go 的 regexp（RE2）不支持 lookahead，故手写逐字符拆分。
func toSnakeCaseStrict(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				prev := runes[i-1]
				next := rune(0)
				if i+1 < len(runes) {
					next = runes[i+1]
				}
				prevLower := prev >= 'a' && prev <= 'z'
				prevDigit := prev >= '0' && prev <= '9'
				nextLower := next >= 'a' && next <= 'z'
				if prevLower || prevDigit || (nextLower && prev >= 'A' && prev <= 'Z') {
					sb.WriteByte('_')
				}
			}
			sb.WriteRune(r + ('a' - 'A'))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// twitterCompatibility 为 og → twitter 的兼容映射。
var twitterCompatibility = map[string]string{
	"og:description": "twitter:description",
	"og:title":       "twitter:title",
	"og:image":       "twitter:image",
	"og:image:url":   "twitter:image",
	"og:image:alt":   "twitter:image:alt",
}

// pickTwitterCompatible 补齐 twitter meta 中缺失的兼容字段（对齐 oggen）。
func pickTwitterCompatible(existing, twitter []metaFlat) []metaFlat {
	seen := make(map[string]bool, len(twitter))
	for _, tm := range twitter {
		seen[tm.Key] = true
	}
	var out []metaFlat
	for _, em := range existing {
		if tkey, ok := twitterCompatibility[em.Key]; ok && !seen[tkey] {
			out = append(out, metaFlat{Key: tkey, Value: em.Value})
		}
	}
	return out
}

// buildMetaStrings 生成 meta 标签字符串列表。
func buildMetaStrings(flats []metaFlat, attr string) string {
	var sb strings.Builder
	for i, f := range flats {
		if i > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, `<meta %s="%s" value="%s" />`, attr, f.Key, f.Value)
	}
	return sb.String()
}