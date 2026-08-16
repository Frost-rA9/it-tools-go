# SPEC — it-tools-go 工程规格书

> 跨平台 IT 工具集桌面应用（参考 [it-tools.tech](https://it-tools.tech/)），Go 后端 + Wails 桌面壳。
> 本文档为项目唯一权威规格，后续实现以本文为准。

---

## 1. 项目概述

- **形态**：跨平台（Windows / macOS / Linux）桌面 IT 工具集，单二进制分发、离线可用。
- **定位**：工具业务逻辑全部在 Go 侧（`internal/`），前端仅负责渲染与交互。
- **参考**：[IT Tools](https://it-tools.tech/)（Vue 3 纯前端 SPA）；本项目以 Go 为核心实现其工具能力。
- **非目标**：云托管 / i18n / 一次对齐全部 40+ 工具（分阶段扩展）/ 多用户与数据持久化。

## 2. 技术选型

| 维度 | 选型 | 版本 | 说明 |
|---|---|---|---|
| 桌面壳 | Wails v2 | 2.14.0 | 单二进制、原生渲染引擎，不内嵌浏览器 |
| 后端 | Go | 1.26.x | 工具逻辑全部在 Go 侧 |
| 前端 | Vue 3 + TypeScript + Vite | 3.x / 7.x | 与参考项目同栈 |
| UI | Naive UI + 自定义主题 | 2.x | 暗色模式、表单密集场景 |
| 路由 | vue-router | 4.x | 前端 SPA 路由 |
| 包管理 | pnpm / Go Modules | pnpm 11.x | pnpm 可选，npm 亦可 |

### 关键架构决策

工具逻辑放在 Go 侧：每个工具在 Go 中实现核心计算（`Execute`），前端经 Wails Binding 调用。
好处：体现 Go 技术栈、便于单测、逻辑与 UI 解耦、可复用于 CLI/服务端。

## 3. 架构与数据流

Wails 双层架构：桌面窗口（原生渲染引擎）承载 `frontend/`（Vue 3），经 Wails Binding（IPC）
调用 Go 后端（`internal/`：Tool 接口 + registry + `internal/tools/*`）。

```
前端工具表单
   │ JSON.stringify(input)
   ▼
RunTool(id, input) ──Binding──▶ registry 查找工具
                                   │
                                   ▼
                              tool.Execute(input)
                                   │
                                   ▼
前端 ◀──Binding── JSON.stringify(output)
```

## 4. 工具注册机制（核心设计）

### Tool 接口

每个工具实现统一接口（定义于 `internal/registry`）：

```go
type Tool struct {
    ID          string   `json:"id"`          // 唯一标识，如 "base64-string-converter"
    Name        string   `json:"name"`        // 展示名称
    Description string   `json:"description"` // 一句话描述
    Category    string   `json:"category"`    // 所属分类（见 §6）
    Keywords    []string `json:"keywords"`    // 搜索关键词
    Icon        string   `json:"icon"`        // 图标名（对应前端 @vicons/tabler）
}

type Executor interface {
    Execute(ctx context.Context, input string) (output string, err error)
}
```

> 传输协议：`Execute` 的 input/output 均为 **JSON 字符串**。前端用原生
> `JSON.stringify`/`JSON.parse`，Go 用 `encoding/json`，前后端零序列化依赖。

### registry

- 聚合工具，提供 `List()` / `Get()` / `Execute()`。
- `internal/app/tools_gen.go` 由代码生成器 `internal/toolsgen` 扫描 `internal/tools/`
  自动生成（按目录名排序，**勿手改**）；新增工具后执行 `go generate ./internal/app` 刷新。
- 单一 Wails Binding `RunTool(id, input)` 暴露执行入口。

### 前端动态渲染

- 启动时 `ListTools()` 获取工具元数据，按分类渲染侧边栏与首页。
- 路由 `/tool/:id` 用 Vite `import.meta.glob('./tools/*.vue')` 按 toolId（文件名）匹配组件。

## 5. 项目目录结构

```
it-tools-go/
├── go.mod / go.sum / main.go        # Go 模块 + Wails 入口（embed frontend/dist）
├── wails.json                       # Wails 构建配置
├── assets/
│   └── logo.svg                     # 品牌 logo 唯一设计源（脚本派生各尺寸资产）
├── internal/
│   ├── app/                         # App 结构体 + 生命周期 + 绑定 + tools_gen.go（勿手改）
│   ├── registry/                    # Tool 接口 + 注册表实现
│   ├── toolsgen/                    # 工具注册代码生成器（go generate ./internal/app）
│   └── tools/                       # 各工具实现（每工具一个包，含单测）
├── frontend/
│   ├── scripts/gen-brand.mjs        # 品牌资产生成（SVG→PNG，拷贝 favicon/logo）
│   ├── public/favicon.svg           # 生成物（源为 assets/logo.svg）
│   └── src/
│       ├── components/              # 通用组件（ToolCard、ToolMenu、ToolTextarea、ToolCodeBlock）
│       ├── layouts/                 # BaseLayout（侧边栏）+ ToolLayout（工具页）
│       ├── views/                   # HomeView + ToolView + tools/（按 toolId 命名）
│       ├── router/ stores/ composables/
│       ├── assets/                  # 主题资源（logo.svg、Cascadia Code 字体等）
│       └── wailsjs/                 # Wails 自动生成绑定（勿手改）
├── build/
│   ├── bin/                         # wails build 输出（it-tools-go.exe）
│   └── windows/                     # 图标、Windows 清单、NSIS 安装器
└── SPEC.md
```

> `main.go` 位于根目录的 `package main`（Wails 强制要求），业务代码全部置于 `internal/`。

## 6. 工具分类

对齐参考项目的分类体系（中文展示，常量定义于 `internal/registry/registry.go`）：

| 分类 | 示例工具 |
|---|---|
| 加密 | UUID、哈希、Token、BIP39、HMAC、RSA 密钥对、密码强度 |
| 转换器 | Base64、罗马数字、大小写、日期时间、进制、文本↔二进制/Unicode、列表、Markdown→HTML、TOML/XML/YAML 互转 |
| Web / 图片 / 开发 / 网络 / 数学 / 测量 / 文本 / 数据 | 后续扩展按此归类 |

## 7. 开发与构建

### 环境要求

Go、gcc（mingw）、Node.js、Wails CLI；本机安装细节见全局 AGENTS.md。

### 常用命令

```bash
wails dev                       # 热重载开发
wails build                     # 构建当前平台二进制
go test ./...                   # 后端单测
go generate ./internal/app      # 新增工具后刷新注册（tools_gen.go）
npm run build                   # 前端类型检查 + 构建（vue-tsc + vite）
```

### 发布

推送 `v*` 标签自动触发 [GitHub Actions](.github/workflows/release.yml)：
三平台交叉编译（Windows amd64 / macOS Intel+Silicon / Linux amd64），
创建 GitHub Release 并附带各平台产物。版本号由 git tag 决定（如 `v0.2.0`）。

## 8. 测试与质量

- Go：表驱动单测（各工具 + registry）；前端：Vitest + `vue-tsc` 类型检查。
- 提交遵循 Conventional Commits；Go 用 gofmt/vet，前端 ESLint + Prettier。

## 9. 里程碑

| 阶段 | 内容 | 验收标准 |
|---|---|---|
| M0 环境准备 | 安装 Go/gcc/Wails/pnpm | ✅ 已完成 |
| M1 骨架搭建 | `wails init` + `internal/` 目录结构 | ✅ 已完成 |
| M2 注册机制 | Tool 接口 + registry + RunTool Binding + 前端动态渲染 | ✅ 已完成 |
| M3 工具扩展 | 高频工具实现（当前 25 个） | ✅ 各工具均有单测 |
| M4 打包发布 | 三平台交叉编译 + GitHub Release 流程 | ✅ 已完成（v0.1.0/v0.2.0） |

## 10. 当前状态

- **版本**：v0.2.0「El Shaddoll Wendigo」（git tag `v0.2.0` 驱动 CI 发布）。
- 注册机制：`registry` + JSON string 协议 + `ListTools`/`RunTool` 绑定；注册由 `internal/toolsgen`
  扫描 `internal/tools/` 生成 `internal/app/tools_gen.go`；前端 `import.meta.glob` 按 toolId 动态加载。
- 已实现工具（26 个）：「转换器」15（Base64、罗马数字、大小写、日期时间、整数基、文本↔ASCII 二进制、
  文本↔Unicode、列表、Markdown→HTML、TOML↔JSON/YAML、XML↔JSON、YAML→JSON/TOML）；「加密」10
  （Token、Hash 文本、加密/解密、BCrypt、UUID、ULID、BIP39、HMAC、RSA 密钥对、密码强度分析）；「Web」1（URL 编码/解码）。
- 前端 it-tools 风格：亮/暗主题、侧边栏分类菜单、Command Palette、首页网格；通用组件
  `ToolTextarea` / `ToolCodeBlock`；等宽字体 Cascadia Code 随包分发。
- 品牌标识：`assets/logo.svg` 唯一源 → `build/appicon.png`、`build/windows/icon.ico`、favicon。
- 命名规范：目录短杠（toolId）、文件下划线、包名去 `-converter` 后缀拼接。
- 已验证：`go build/vet/test ./...`、`npm run build`、`wails build` 均通过。

## 11. 变更记录

| 日期 | 变更 | 说明 |
|---|---|---|
| 2026-08-14 | 初版 + M1 | 技术选型、架构、目录结构；Wails 骨架 |
| 2026-08-14 | M2 注册机制 | Tool 接口 / registry / JSON string 协议 / `ListTools`+`RunTool` 绑定 + Base64 工具 |
| 2026-08-14 | 前端 it-tools 风格 | 亮/暗主题、顶栏 + Command Palette、分类中文化、工具图标、首页网格 |
| 2026-08-14 | 转换器工具扩展 | 罗马/大小写/日期时间/整数基/文本二进制/文本Unicode/YAML→JSON/YAML→TOML 等 8 个 |
| 2026-08-14 | 注册重构 | 工具包自持元数据（`Tool()` + `Executor`），集中一行注册 |
| 2026-08-14 | 转换器补齐 | 列表、Markdown→HTML、TOML↔JSON/YAML、XML↔JSON（引入 goldmark、mxj） |
| 2026-08-14 | 品牌 logo | `assets/logo.svg` 设计源 + `gen-brand.mjs` 生成链路，接入 appicon/ico/favicon |
| 2026-08-15 | 加密工具 | Token、Hash（SHA3 legacy Keccak）、加密/解密（AES/TripleDES/Rabbit/RC4）、BCrypt、UUID、ULID |
| 2026-08-15 | 通用组件 | `ToolTextarea`（可拉伸+monospace）、`ToolCodeBlock`（只读等宽块）；Cascadia Code 随包 |
| 2026-08-15 | 加密工具扩展 | BIP39（10 语言）、HMAC（8 哈希 × 4 编码）、RSA 密钥对（256–16384 bits）、密码强度分析 |
| 2026-08-15 | 注册生成器 | `internal/toolsgen` 扫描 `internal/tools/` 自动生成 `tools_gen.go`；`app.go` 精简 |
| 2026-08-15 | 品牌 logo 重绘 | 齿轮 + 扳手造型（青色齿轮 `#00ADD8` + 白色扳手），同步 favicon/appicon/ico |
| 2026-08-15 | 文档精简 | SPEC.md 精简为要点式规格；README.md 重写并补充工具清单与发布说明 |
| 2026-08-15 | URL 编码/解码 | Web 分类首个工具；编码/解码对齐 JS `encodeURIComponent`/`decodeURIComponent`（含 UTF-8 校验），向量单测 |
