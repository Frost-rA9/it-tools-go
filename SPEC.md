# SPEC — it-tools-go 工程规格书

> 跨平台 IT 工具集桌面应用（参考 [it-tools.tech](https://it-tools.tech/)），使用 Go 语言实现。
> 本文档为项目唯一权威规格，后续实现须以本文为准。

---

## 1. 项目概述

### 1.1 目标

实现一个**跨平台**（Windows / macOS / Linux）的桌面 IT 工具集应用，提供开发者与运维人员常用的便捷小工具（UUID 生成、哈希计算、Base64 编解码、URL 编解码、JSON 格式化、时间戳转换、颜色转换、二维码生成等），以**单二进制**形式分发，无需联网即可本地运行。

### 1.2 参考项目

- 名称：IT Tools（https://it-tools.tech/）
- 原项目技术栈：Vue 3 + TypeScript（纯前端 SPA）
- 本项目的定位：**以 Go 为后端核心**，工具业务逻辑全部落在 Go 侧，前端仅负责展示与交互。

### 1.3 非目标（Non-Goals）

首期明确不做：

- 云部署 / 在线托管（本项目为本地桌面应用）。
- 多语言国际化（i18n），后续可扩展。
- 一次性对齐原版全部 40+ 工具（分阶段扩展）。
- 多用户、鉴权、数据持久化等（工具为无状态纯计算）。

---

## 2. 技术选型

| 维度 | 选型 | 版本 | 说明 |
|---|---|---|---|
| 应用形态 | 桌面 GUI（Wails） | **v2（稳定版，当前 2.14.0）** | 单二进制、原生渲染引擎（WebView2/WebKitGTK），不内嵌浏览器 |
| 后端语言 | Go | **1.26.x（当前 1.26.5）** | 工具业务逻辑全部在 Go 侧 |
| 前端框架 | Vue 3 + TypeScript + Vite | Vue 3.x / Vite 7.x | 与参考项目同栈，交互还原成本最低 |
| UI 库 | Naive UI + 自定义主题 | Naive UI 2.x | 暗色模式/国际化开箱即用，适合工具类密集表单 |
| 路由 | vue-router | 4.x | 前端 SPA 路由 |
| 包管理 | 前端 pnpm / 后端 Go Modules | pnpm 11.x | pnpm 可选，npm 亦可（见 §7.1） |
| 工具逻辑归属 | **全部在 Go 后端** | — | 通过 Wails Binding 暴露给前端（自动生成 TS 类型） |

### 2.1 关键架构决策

**工具逻辑放在 Go 侧**（而非前端 TS）。这是本项目与参考项目最大的差异，也是「用 Go 实现」的核心：

- 每个工具在 Go 中实现核心计算逻辑（`Execute` 方法）。
- 前端通过 Wails Binding 调用 Go 方法，Wails 自动生成 TypeScript 类型定义。
- 好处：真正体现 Go 技术栈、便于单测、逻辑与 UI 解耦、可复用于 CLI/服务端。

---

## 3. 架构设计

### 3.1 整体架构（Wails 双层架构）

```
┌─────────────────────────────────────────────────┐
│                  桌面窗口（原生渲染引擎）            │
│  ┌───────────────────────────────────────────┐  │
│  │        frontend/  (Vue 3 + Naive UI)      │  │
│  │  - 布局、路由、工具表单、结果展示            │  │
│  └──────────────┬────────────────────────────┘  │
│                 │  Wails Binding（IPC）          │
│  ┌──────────────▼────────────────────────────┐  │
│  │           Go 后端 (internal/)              │  │
│  │  - Tool 接口 + 注册表 registry             │  │
│  │  - 各工具实现 internal/tools/*             │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

### 3.2 数据流

```
前端工具表单
   │ 用户输入（JSON 参数）
   ▼
RunTool(id, input)  ──Wails Binding──▶  Go 后端 registry 查找工具
                                          │
                                          ▼
                                     tool.Execute(input)
                                          │
                                          ▼
                                       输出结果
                                          │
  前端展示结果 ◀──Wails Binding─── 返回 output
```

---

## 4. 工具注册机制（核心设计）

### 4.1 `Tool` 接口

每个工具实现统一接口（定义于 `internal/registry`）：

```go
type Tool struct {
    ID          string   `json:"id"`          // 唯一标识，如 "base64-string-converter"
    Name        string   `json:"name"`        // 展示名称
    Description string   `json:"description"` // 一句话描述
    Category    string   `json:"category"`    // 所属分类（见 §6）
    Keywords    []string `json:"keywords"`    // 搜索关键词
    Icon        string   `json:"icon"`        // 图标名（对应前端 @vicons/tabler 组件，如 "FileDigit"）
}

type Executor interface {
    Execute(ctx context.Context, input string) (output string, err error)
}
```

> 传输协议：`Execute` 的 input/output 均为 **JSON 字符串**。前端用原生
> `JSON.stringify`/`JSON.parse`，Go 用标准库 `encoding/json`，前后端零序列化依赖，
> 且 Wails 生成 `string → string` 的类型安全绑定。每个工具自行约定 input/output 的 JSON 结构。

### 4.2 注册表（registry）

- `registry` 聚合所有工具，提供 `List()`（返回全部工具元数据，供前端渲染侧边栏）与 `Get(id)`（按 ID 查找）、`Execute(id, input)`。
- 新工具通过集中注册函数 `registerTools`（`internal/app/app.go`）调用 `registry.Register(tool, executor)` 加入。
- 单一 Wails Binding `RunTool(id, input)` 暴露执行入口，前端无需为每个工具单独绑定。

### 4.3 前端动态渲染

- 前端启动时调用 `ListTools()` 获取工具元数据，按分类渲染侧边栏与首页。
- 路由 `/tool/:id` 通过 Vite `import.meta.glob('./tools/*.vue')` 自动按 toolId 匹配组件
  （文件名即 toolId，如 `base64-string-converter.vue`），表单提交后调用 `RunTool`。

---

## 5. 项目目录结构

```
it-tools-go/
├── go.mod / go.sum              # Go 模块
├── main.go                      # Wails 入口（wails.Run + embed frontend/dist）
├── wails.json                   # Wails 构建配置
├── internal/
│   ├── app/                     # App 结构体 + 生命周期 + 绑定（package app）
│   ├── registry/                # Tool 接口 + 注册表实现
│   └── tools/                   # 各工具实现（每工具一个包）
│       └── base64/              # Base64 工具 + 单测
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── index.html
│   └── src/
│       ├── main.ts
│       ├── App.vue              # n-config-provider（暗色主题）
│       ├── theme.ts             # Naive UI 主题覆盖（it-tools 风格）
│       ├── components/          # 通用组件（ToolCard、ToolMenu）
│       ├── router/              # vue-router 配置
│       ├── layouts/             # BaseLayout（侧边栏）+ ToolLayout（工具页）
│       ├── views/               # HomeView + ToolView + tools/（按 toolId 命名）
│       ├── stores/              # Pinia（tools store）
│       ├── composables/         # useToolComponent（glob 动态加载）
│       └── wailsjs/             # Wails 自动生成的绑定代码（勿手改）
├── build/                       # 构建产物（图标、打包配置）
│   └── bin/                     # wails build 输出（it-tools-go.exe）
└── SPEC.md
```

> 注：`main.go` 位于根目录的 `package main`（Wails 官方强制要求），业务代码全部置于
> `internal/`（Go 社区惯例的私有包位置）。App 绑定路径为 `wailsjs/go/app/App`。

---

## 6. 工具分类

对齐参考项目的分类体系（中文展示，后续工具扩展据此归类）：

| 分类 | 示例工具 |
|---|---|
| 加密 | UUID 生成、哈希计算、Token 生成 |
| 转换器 | Base64 编解码、时间戳转换、颜色转换、大小写转换、日期转换 |
| Web | URL 编解码、HTML 实体、User-Agent 解析 |
| 图片和视频 | 二维码生成、SVG 占位图 |
| 开发 | JSON 格式化、Cron 解析、正则测试 |
| 网络 | IPv4 子网计算、MAC 地址查询 |
| 数学 | 进制转换、ETA 计算 |
| 测量 | 计时器、温度/单位换算 |
| 文本 | Slugify、词频统计、字符统计 |
| 数据 | YAML↔JSON、XML 格式化、Phone 解析 |

> 分类名当前仅实现中文（不引入 i18n 框架），常量定义于 `internal/registry/registry.go`。
> 若后续有开发者贡献多语言 PR，再评估引入 i18n。

---

## 7. 开发与构建

### 7.1 环境要求（M0，已在本机完成）

| 组件 | 版本 | 安装方式 |
|---|---|---|
| Go | 1.26.5 | `scoop install go` |
| MinGW (gcc) | 16.1.0 | `scoop install mingw`（Wails Windows 构建必需，用于 CGO/资源编译） |
| Node.js | 24.19.0 | 已装（scoop） |
| pnpm | 11.21.0 | `scoop install pnpm`（可选，npm 亦可） |
| Wails CLI | 2.14.0 | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| WebView2 Runtime | — | Win10/11 通常内置，Wails 依赖 |

> 注：Go 模块代理已配置为 `GOPROXY=https://goproxy.cn,direct`（本机需此代理才能拉取模块）。

### 7.2 常用命令

```bash
# 开发（热重载）
wails dev

# 构建当前平台二进制
wails build

# 交叉编译（示例）
wails build -platform windows/amd64

# 后端单测
go test ./...

# 前端类型检查 / 构建 / 测试 / lint
pnpm typecheck
pnpm build
pnpm test
pnpm lint
```

---

## 8. 测试与质量

- **Go 后端**：表驱动单元测试（`internal/tools/*` 及 `internal/registry`），覆盖各工具核心逻辑与边界情况。
- **前端**：Vitest 组件测试；`vue-tsc` 类型检查。
- **提交规范**：遵循 Conventional Commits（见全局 AGENTS.md），如 `feat(registry): ...`、`fix(uuid): ...`。
- **代码风格**：Go 用 `gofmt`/`go vet`；前端用 ESLint + Prettier。

---

## 9. 里程碑（Roadmap）

| 阶段 | 内容 | 验收标准 |
|---|---|---|
| **M0 环境准备** | 安装 Go/gcc/Wails/pnpm | ✅ 已完成（见 §7.1） |
| **M1 骨架搭建** | `wails init` 生成项目，配置 `internal/` 目录结构 | ✅ 已完成（`wails build` 产出 exe） |
| **M2 注册机制** | 实现 `Tool` 接口 + `registry` + `RunTool` Binding + 前端动态渲染 + Base64 示例工具 | ✅ 已完成（Base64 工具可执行并返回结果） |
| **M3 首批工具** | 实现 15-20 个高频工具（UUID、哈希、Base64、URL、JSON、时间戳、颜色、二维码等） | 各工具均有单测并通过 |
| **M4 打包发布** | 三平台交叉编译、图标/签名、发布产物 | 产出 Windows/macOS/Linux 单二进制 |

---

## 10. 当前状态

- 已搭建 Wails v2 桌面应用（Vue 3 + TypeScript 前端，Go 后端）。
- 存在 `go.mod`（模块 `it-tools-go`，go 1.25.0，wails v2.14.0）。
- 后端采用 `internal/` 目录结构，`App` 结构体在 `internal/app`（绑定路径 `wailsjs/go/app/App`）。
- 已实现工具注册机制（`registry` + JSON string 协议 + `ListTools`/`RunTool` 绑定）与首个工具 Base64。
- Base64 工具为双卡片布局（文本↔Base64 双向），支持 URL-safe 编解码与复制，输入实时调用 Go 后端转换。
- 新增 3 个转换器工具：罗马数字（1–3999 双向转换 + 合法性校验）、大小写（14 种格式：驼峰/蛇形/常量/mocking 等）、日期时间（10 种格式：时间戳/ISO 8601/ISO 9075/RFC 3339/RFC 7231/UTC/Mongo ObjectID/Excel 等，含格式自动识别与下拉回退）。
- 新增 3 个转换器工具：整数基转换（2~64 进制任意精度，`math/big`）、文本↔ASCII 二进制、文本↔Unicode 转义。
- 前端已还原 it-tools 风格：Naive UI 主题（默认亮色 + 可切换暗色）+ 侧边栏分类菜单（中文分类）+ 顶栏（侧边栏/主页/GitHub/主题切换按钮 + Command Palette 搜索）+ 首页卡片网格 + `import.meta.glob` 动态加载工具组件。
- 工具元数据新增 `Icon` 字段（Go 端存 `@vicons/tabler` 图标名，前端 `src/tools/icons.ts` 映射为组件）；侧边栏与首页卡片均显示工具图标；首页改为扁平网格（不分分类），卡片等宽等高（响应式 1/2/3/4 列）。
- 主色采用 Go 官方蓝 `#00ADD8`；侧边栏顶部为 Go 蓝渐变色块。
- 引入 `@vueuse/core`（`useDark`/`useToggle`/`useStorage`/`useMediaQuery`/`useMagicKeys`）实现主题切换、侧边栏折叠持久化、搜索快捷键。
- it-tools 参考仓库已 clone 到 `D:\Code\it-tools`（非本项目目录，仅作设计参考）。
- 已验证可构建：`go build ./...`、`go vet ./...`、`go test ./...`、`npm run build`、`wails build` 均通过。

---

## 11. 变更记录

| 日期 | 变更 | 说明 |
|---|---|---|
| 2026-08-14 | 初版 | 确定技术选型、架构、目录结构、里程碑 |
| 2026-08-14 | M1 骨架 + 目录重构 | 搭建 Wails 骨架；后端采用 `internal/` 目录（app/registry/tools），App 结构体迁入 `internal/app`，绑定路径改为 `wailsjs/go/app/App` |
| 2026-08-14 | M2 注册机制 | 完成 Tool 接口/registry/JSON string 协议/ListTools+RunTool 绑定；Base64 工具；前端还原 it-tools 暗色主题 + 侧边栏 + 卡片网格 + glob 动态加载 |
| 2026-08-14 | 亮色主题 + 顶栏 + 品牌色 | 默认亮色 + 可切换暗色；新增顶栏（侧边栏/主页/GitHub/主题切换 + Command Palette 搜索）；主色改为 Go 蓝 `#00ADD8`，侧边栏顶部加渐变色块 |
| 2026-08-14 | 分类中文化 + Base64 双卡片 | 分类常量改为中文（含「图片和视频」）；Base64 归类「转换器」并改造为双卡片（URL-safe + 复制 + 实时调用 Go）；修复侧边栏分类折叠 bug |
| 2026-08-14 | 新增 3 个转换器 | 罗马数字（1–3999 双向+校验）、大小写（14 种格式）、日期时间（10 种格式+自动识别） |
| 2026-08-14 | 工具图标 + 首页扁平化 | Tool 元数据新增 Icon 字段（tabler 图标名）；侧边栏与首页卡片加图标；首页去分类标题、卡片等宽等高（响应式列） |
| 2026-08-14 | 新增 3 个转换器 | 整数基转换（2~64 进制）、文本↔ASCII 二进制、文本↔Unicode |
