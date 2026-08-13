package app

import (
	"context"
	"fmt"
)

// App 是绑定到前端的应用结构体（类比传统 web 应用的 controller）。
type App struct {
	ctx context.Context
}

// NewApp 创建 App 实例。
func NewApp() *App {
	return &App{}
}

// Startup 在应用启动时调用，保存 context 以便调用 runtime 方法。
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet 返回示例问候语（骨架阶段的示例绑定）。
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
