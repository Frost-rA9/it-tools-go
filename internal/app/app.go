// Package app 承载绑定到前端的 App 结构体，负责组装工具注册表并暴露调用入口。
package app

import (
	"context"
	"fmt"

	"it-tools-go/internal/registry"
	"it-tools-go/internal/tools/base64-string-converter"
	"it-tools-go/internal/tools/case-converter"
	"it-tools-go/internal/tools/date-time-converter"
	"it-tools-go/internal/tools/integer-base-converter"
	"it-tools-go/internal/tools/roman-numeral-converter"
	"it-tools-go/internal/tools/text-to-binary"
	"it-tools-go/internal/tools/text-to-unicode"
	"it-tools-go/internal/tools/yaml-to-json-converter"
)

// App 是绑定到前端的应用结构体（类比传统 web 应用的 controller）。
type App struct {
	ctx      context.Context
	registry *registry.Registry
}

// NewApp 创建 App 实例并注册全部工具。
func NewApp() *App {
	reg := registry.New()
	registerTools(reg)
	return &App{registry: reg}
}

// registerTools 集中注册所有工具。
func registerTools(reg *registry.Registry) {
	reg.Register(registry.Tool{
		ID:          base64string.ID,
		Name:        base64string.Name,
		Description: "在普通文本与其 Base64 编码形式之间进行转换",
		Category:    base64string.Category,
		Keywords:    []string{"base64", "编码", "解码", "encode", "decode"},
		Icon:        "FileDigit",
	}, base64string.Executor{})

	reg.Register(registry.Tool{
		ID:          romannumeral.ID,
		Name:        romannumeral.Name,
		Description: "在阿拉伯数字与罗马数字之间进行转换",
		Category:    romannumeral.Category,
		Keywords:    []string{"roman", "罗马", "数字", "numeral", "阿拉伯"},
		Icon:        "LetterX",
	}, romannumeral.Executor{})

	reg.Register(registry.Tool{
		ID:          caseconv.ID,
		Name:        caseconv.Name,
		Description: "在多种大小写格式之间转换字符串（驼峰、蛇形、常量等）",
		Category:    caseconv.Category,
		Keywords:    []string{"case", "大小写", "驼峰", "camel", "snake", "upper", "lower"},
		Icon:        "LetterCaseToggle",
	}, caseconv.Executor{})

	reg.Register(registry.Tool{
		ID:          datetime.ID,
		Name:        datetime.Name,
		Description: "在不同日期时间格式之间进行转换（时间戳、ISO 8601、RFC 系列、Excel 等）",
		Category:    datetime.Category,
		Keywords:    []string{"date", "time", "日期", "时间", "时间戳", "timestamp", "unix", "iso8601", "excel"},
		Icon:        "Calendar",
	}, datetime.Executor{})

	reg.Register(registry.Tool{
		ID:          radix.ID,
		Name:        radix.Name,
		Description: "在不同进制之间转换整数（二进制、八进制、十进制、十六进制、Base64 等）",
		Category:    radix.Category,
		Keywords:    []string{"integer", "base", "进制", "转换", "binary", "decimal", "hex", "octal", "radix"},
		Icon:        "ArrowsLeftRight",
	}, radix.Executor{})

	reg.Register(registry.Tool{
		ID:          textbinary.ID,
		Name:        textbinary.Name,
		Description: "在文本与 ASCII 二进制表示之间进行转换",
		Category:    textbinary.Category,
		Keywords:    []string{"text", "binary", "文本", "二进制", "ascii", "encode", "decode"},
		Icon:        "Binary",
	}, textbinary.Executor{})

	reg.Register(registry.Tool{
		ID:          textunicode.ID,
		Name:        textunicode.Name,
		Description: "在文本与 Unicode 转义表示之间进行转换",
		Category:    textunicode.Category,
		Keywords:    []string{"text", "unicode", "文本", "转义", "编码"},
		Icon:        "TextWrap",
	}, textunicode.Executor{})

	reg.Register(registry.Tool{
		ID:          yamljson.ID,
		Name:        yamljson.Name,
		Description: "将 YAML 转换为 JSON",
		Category:    yamljson.Category,
		Keywords:    []string{"yaml", "json", "转换", "convert"},
		Icon:        "AlignJustified",
	}, yamljson.Executor{})
}

// Startup 在应用启动时调用，保存 context 以便调用 runtime 方法。
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// ListTools 返回全部已注册工具的元数据，供前端渲染侧边栏与首页。
func (a *App) ListTools() []registry.Tool {
	return a.registry.List()
}

// RunTool 按 ID 执行工具，input/output 均为 JSON 字符串。
func (a *App) RunTool(id, input string) (string, error) {
	output, ok, err := a.registry.Execute(a.ctx, id, input)
	if !ok {
		return "", fmt.Errorf("未找到工具: %q", id)
	}
	return output, err
}
