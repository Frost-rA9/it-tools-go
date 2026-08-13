# AGENTS.md

## 当前状态

- 已搭建 Wails v2 桌面应用骨架（Vue 3 + TypeScript 前端，Go 后端）。
- 存在 `go.mod`（模块 `it-tools-go`，go 1.25.0，wails v2.14.0）。
- 后端采用 `internal/` 目录结构，`App` 结构体在 `internal/app`（绑定路径 `wailsjs/go/app/App`）。
- 已验证可构建：`go build ./...`、`go vet ./...`、`npm run build`、`wails build` 均通过。

## 项目结构

```
it-tools-go/
├── go.mod / go.sum              # Go 模块
├── main.go                      # Wails 入口（wails.Run + embed frontend/dist，package main）
├── wails.json                   # Wails 构建配置（frontend:install/build 命令）
├── internal/
│   ├── app/                     # App 结构体 + 生命周期 + 绑定
│   ├── registry/                # Tool 接口 + 注册表
│   └── tools/                   # 各工具实现（每工具一个包，M2 填充）
├── build/                       # 构建资源（图标、Windows 清单、NSIS 安装器）
│   └── bin/                     # wails build 产物（it-tools-go.exe）
├── frontend/
│   ├── package.json             # npm 脚本：dev / build / preview
│   ├── vite.config.ts
│   ├── tsconfig.json / tsconfig.node.json
│   └── src/
│       ├── main.ts / App.vue / style.css
│       ├── components/HelloWorld.vue
│       ├── assets/              # 字体、图片
│       └── wailsjs/             # Wails 自动生成的绑定（切勿手动编辑）
└── SPEC.md                      # 工程规格书（权威）
```

注：`frontend/src` 下的 `router/`、`layouts/`、`views/`、`composables/` 尚未创建 —— 属于 M2 里程碑。
`internal/tools/` 目前仅占位，具体工具待 M2 加入。

## 命令

- 开发（热重载）：`wails dev`
- 构建桌面二进制：`wails build`
- Go 构建：`go build ./...`
- Go 测试：`go test ./...`
- 前端类型检查 + 构建：`npm run build`（执行 `vue-tsc --noEmit && vite build`）
- 前端开发服务器：`npm run dev`（在 `frontend/` 目录运行）

## 环境要求

- 需要 Go、gcc（mingw）、Node.js 及 Wails CLI。本机具体的 PATH 前缀与网络代理配置
  见全局 AGENTS.md（`~/.config/opencode/AGENTS.md`）。

## 指导原则

- 工具逻辑属于 Go（遵循 SPEC.md §2.1）；前端只负责渲染和调用绑定。
- `frontend/wailsjs/` 由 Wails 生成 —— 切勿手动编辑。
- 提交信息遵循 Conventional Commits（见全局 AGENTS.md）。
