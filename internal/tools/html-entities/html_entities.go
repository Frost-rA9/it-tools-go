// Package htmlentities 实现 HTML 实体转义/反转义工具，行为对齐 lodash 的 escape/unescape。
package htmlentities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "html-entities"
	Name        = "HTML实体转义"
	Description = "在特殊字符与其 HTML 实体形式之间进行转换"
	Category    = "Web"
	Icon        = "Code"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"html", "entities", "escape", "unescape", "special", "characters", "tags", "转义", "反转义", "实体"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text   string `json:"text"`   // 待处理文本
	Action string `json:"action"` // escape | unescape
}

// output 是工具的输出结构。
type output struct {
	Result string `json:"result"`
}

// escape 对齐 lodash escape：仅处理 & < > " ' 五个字符。
var escape = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&#39;",
)

// unescape 对齐 lodash unescape：还原对应的五个实体（单遍替换，不递归）。
var unescape = strings.NewReplacer(
	"&amp;", "&",
	"&lt;", "<",
	"&gt;", ">",
	"&quot;", `"`,
	"&#39;", "'",
)

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	var result string
	switch in.Action {
	case "escape":
		result = escape.Replace(in.Text)
	case "unescape":
		result = unescape.Replace(in.Text)
	default:
		return "", fmt.Errorf("未知操作: %q（仅支持 escape/unescape）", in.Action)
	}

	out, err := json.Marshal(output{Result: result})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}