// Package app 承载绑定到前端的 App 结构体，负责组装工具注册表并暴露调用入口。
// 注册函数由 internal/toolsgen 生成（tools_gen.go），按 internal/tools 目录名排序。
// 新增工具后执行：go generate ./internal/app
//go:generate go run ../../internal/toolsgen
package app

import (
	"context"
	"fmt"

	"it-tools-go/internal/registry"
)

// App 是绑定到前端的应用结构体（类比传统 web 应用的 controller）。
type App struct {
	ctx      context.Context
	registry *registry.Registry
}

// NewApp 创建 App 实例并注册全部工具。
func NewApp() *App {
	reg := registry.New()
	registerAllTools(reg)
	return &App{registry: reg}
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
