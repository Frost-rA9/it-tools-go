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
	"it-tools-go/internal/tools/list-converter"
	"it-tools-go/internal/tools/markdown-to-html"
	"it-tools-go/internal/tools/roman-numeral-converter"
	"it-tools-go/internal/tools/text-to-binary"
	"it-tools-go/internal/tools/text-to-unicode"
	"it-tools-go/internal/tools/yaml-to-json-converter"
	"it-tools-go/internal/tools/yaml-to-toml"
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
	reg.Register(base64string.Tool(), base64string.Executor{})
	reg.Register(romannumeral.Tool(), romannumeral.Executor{})
	reg.Register(caseconv.Tool(), caseconv.Executor{})
	reg.Register(datetime.Tool(), datetime.Executor{})
	reg.Register(radix.Tool(), radix.Executor{})
	reg.Register(textbinary.Tool(), textbinary.Executor{})
	reg.Register(textunicode.Tool(), textunicode.Executor{})
	reg.Register(yamljson.Tool(), yamljson.Executor{})
	reg.Register(yamltoml.Tool(), yamltoml.Executor{})
	reg.Register(listconv.Tool(), listconv.Executor{})
	reg.Register(markdownhtml.Tool(), markdownhtml.Executor{})
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
