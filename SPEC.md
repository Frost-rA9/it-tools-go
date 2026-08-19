# SPEC — it-tools-go 工程规格书

> 跨平台 IT 工具集桌面应用（参考 [it-tools.tech](https://it-tools.tech/)），Go 后端 + Wails 桌面壳。
> 本文档为项目唯一权威规格，后续实现以本文为准。

## 1. 项目概述

- **形态**：跨平台（Windows / macOS / Linux）桌面 IT 工具集，单二进制分发、离线可用。
- **定位**：工具业务逻辑全部在 Go 侧（`internal/`），前端仅负责渲染与交互。
- **参考**：[IT Tools](https://it-tools.tech/)（Vue 3 纯前端 SPA），本项目以 Go 为核心实现其工具能力。
- **非目标**：云托管 / i18n / 一次对齐全部 40+ 工具 / 多用户与数据持久化。

## 2. 技术选型

| 维度 | 选型 | 版本 |
|---|---|---|
| 桌面壳 | Wails v2 | 2.14.0 |
| 后端 | Go | 1.26.x |
| 前端 | Vue 3 + TypeScript + Vite | 3.x / 7.x |
| UI | Naive UI + 自定义主题 | 2.x |
| 路由 | vue-router | 4.x |
| 包管理 | pnpm / Go Modules | pnpm 11.x |

关键决策：工具逻辑放 Go（`Execute`），前端经 Wails Binding 调用——体现 Go 技术栈、便于单测、逻辑与 UI 解耦。

## 3. 架构与数据流

Wails 双层架构：桌面窗口承载 `frontend/`（Vue 3），经 Wails Binding（IPC）调用 Go 后端
（Tool 接口 + registry + `internal/tools/*`）。

```
前端工具表单 --RunTool(id, input)--> registry 查找 --> tool.Execute(input) --> 返回 output
```

传输均为 JSON string。

## 4. 工具注册机制

每个工具实现统一接口（`internal/registry`）：

```go
type Tool struct {
    ID, Name, Description, Category, Icon string
    Keywords []string
}

type Executor interface {
    Execute(ctx context.Context, input string) (output string, err error)
}
```

- 传输协议：input/output 均为 **JSON 字符串**（前端 `JSON.stringify/parse`，Go `encoding/json`）。
- `internal/app/tools_gen.go` 由代码生成器 `internal/toolsgen` 扫描 `internal/tools/` 生成
  （按目录名排序，**勿手改**）；新增工具后执行 `go generate ./internal/app` 刷新。
- 单一 Binding `RunTool(id, input)` 暴露执行入口；前端 `ListTools()` 拉元数据、
  路由 `/tool/:id` 用 `import.meta.glob('./tools/*.vue')` 按 toolId（文件名）匹配组件。

## 5. 目录结构

```
it-tools-go/
├── main.go / go.mod / go.sum    # Wails 入口 + Go 模块
├── wails.json / assets/logo.svg # 构建配置 / 品牌 logo 唯一源
├── internal/
│   ├── app/                     # 绑定 + tools_gen.go（生成，勿改）
│   ├── registry/                # Tool 接口 + 注册表
│   ├── toolsgen/                # 注册代码生成器
│   └── tools/                   # 各工具实现（每包含单测）
├── frontend/
│   ├── scripts/gen-brand.mjs    # 品牌资产生成（SVG→PNG/favicon/ico）
│   └── src/
│       ├── components/ layouts/ views/ router/ stores/ composables/ assets/
│       └── wailsjs/             # Wails 生成绑定（勿改）
├── build/                       # 图标 / Windows 资源 / NSIS 安装器
└── SPEC.md
```

## 6. 开发与发布

环境：Go、gcc（mingw）、Node.js、Wails CLI（本机安装细节见全局 AGENTS.md）。

```bash
wails dev                       # 热重载开发
wails build                     # 构建当前平台二进制
go test ./...                   # 后端单测
go generate ./internal/app      # 新增工具后刷新注册（tools_gen.go）
npm run build                   # 前端类型检查 + 构建（vue-tsc + vite）
```

发布：推送 `v*` 标签自动触发 [GitHub Actions](.github/workflows/release.yml) 三平台交叉编译
（Windows amd64 / macOS Intel+Silicon / Linux amd64）并创建 Release；版本号由 git tag 决定（如 `v0.3.0`）。

## 7. 测试与质量

- Go：表驱动单测（各工具 + registry）；前端：`vue-tsc` 类型检查。
- 提交遵循 Conventional Commits；Go 用 gofmt/vet。

