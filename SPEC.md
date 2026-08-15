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
- 工具包自持完整元数据：每个工具导出 `Tool()`（返回 `registry.Tool`）与 `Executor`；集中注册函数 `registerTools`（`internal/app/app.go`）用 `reg.Register(xxx.Tool(), xxx.Executor{})` 一行注册一个工具。
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
├── assets/
│   └── logo.svg                 # 品牌 logo 设计源（唯一源，脚本派生各尺寸资产）
├── internal/
│   ├── app/                     # App 结构体 + 生命周期 + 绑定（package app）
│   ├── registry/                # Tool 接口 + 注册表实现
│   └── tools/                   # 各工具实现（每工具一个包）
│       └── base64-string-converter/  # Base64 工具 + 单测（目录短杠，文件下划线）
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── index.html
│   ├── scripts/
│   │   └── gen-brand.mjs        # 品牌资产生成（SVG→PNG，拷贝 favicon/侧边栏 logo）
│   ├── public/
│   │   └── favicon.svg          # 生成物（源为 assets/logo.svg）
│   └── src/
│       ├── main.ts
│       ├── App.vue              # n-config-provider（暗色主题）
│       ├── theme.ts             # Naive UI 主题覆盖（it-tools 风格）
│       ├── components/          # 通用组件（ToolCard、ToolMenu、ToolTextarea）
│       ├── router/              # vue-router 配置
│       ├── layouts/             # BaseLayout（侧边栏）+ ToolLayout（工具页）
│       ├── views/               # HomeView + ToolView + tools/（按 toolId 命名）
│       ├── stores/              # Pinia（tools store）
│       ├── composables/         # useToolComponent（glob 动态加载）
│       ├── assets/              # 主题资源（hero-gradient.svg、logo.svg 等）
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

> 注：Go 模块代理为官方源 `GOPROXY=https://proxy.golang.org,direct`，访问需本地代理（见全局 AGENTS.md）。

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

- 应用形态：Wails v2 桌面应用（Go 后端 + Vue 3 + TypeScript + Naive UI 前端），Go 模块 `it-tools-go`（go 1.25.0，wails v2.14.0）。
- 后端目录：`internal/`（`app` / `registry` / `tools`），App 绑定路径 `wailsjs/go/app/App`。
- 工具注册机制：`registry`（Tool 接口 + 注册表）+ JSON string 协议 + `ListTools`/`RunTool` 绑定；工具包自持元数据（导出 `Tool()` 与 `Executor`），`registerTools` 一行注册一个；前端 `import.meta.glob` 按 toolId 动态加载组件。
- 已实现 15 个「转换器」工具：Base64、罗马数字、大小写、日期时间、整数基（radix）、文本↔ASCII 二进制、文本↔Unicode、YAML 转 JSON、YAML 转 TOML、列表（listconv）、Markdown 转 HTML（goldmark）、TOML 转 JSON、TOML 转 YAML、XML 转 JSON、JSON 转 XML（mxj）。
- 加密分类首个工具：Token 生成器（token-generator，包 `tokengen`），用 `crypto/rand` 均匀采样生成（大写/小写/数字/符号可开关，长度 1–512，支持自定义字符集），前端 2×2 开关 + 长度滑块 + 复制/重新生成。
- 加密分类：哈希文本（hash-text，包 `hashtext`）——8 算法（MD5/SHA1/SHA224/SHA256/SHA384/SHA512/SHA3/RIPEMD160，SHA3 为 legacy Keccak 对齐 crypto-js）× 4 编码（Hex/Base64/Base64url/Bin）；加密/解密（encryption，包 `encryption`）——AES/TripleDES/Rabbit/RC4，passphrase 走 EVP_BytesToKey（MD5、1 迭代）+ OpenSSL 格式（`Salted__`+盐），流密码无填充，Rabbit 为 crypto-js 移植；测试用 crypto-js 生成向量保证字节级一致（testdata/crypto-js-vectors.json）。
- 等宽字体：`@fontsource/cascadia-code`（latin 子集）打包，`ToolTextarea` 新增 `monospace` prop（默认 false），加密工具密文用 Cascadia Code 展示。
- 通用组件 `ToolTextarea`：封装"可选 label + 可拉伸多行输入框"（`defineModel`，`resizable` 默认开启）；已全量迁移所有工具的 textarea 输入框（12 文件 29 处），case-converter 结果行的单行只读输入与 list-converter 的配置输入框（非 textarea）不使用该组件。
- 工具元数据含 `Icon` 字段（tabler 图标名）；侧边栏与首页卡片显示图标；首页为扁平网格（响应式 1/2/3/4 列、等宽等高）。
- 前端 it-tools 风格：Naive UI 亮/暗主题、侧边栏分类菜单、顶栏 + Command Palette 搜索。
- 依赖：前端 `@vueuse/core`；后端 `gopkg.in/yaml.v3` + `github.com/pelletier/go-toml/v2` + `github.com/yuin/goldmark` + `github.com/clbanning/mxj/v2`。
- XML↔JSON 转换对齐 it-tools 的 `_attributes`/`_text` 约定。
- 命名规范：目录短杠（toolId）、文件下划线、包名去 `-converter` 后缀拼接。
- 品牌标识：`assets/logo.svg` 为唯一设计源（深色圆角底 + 青色开口扳手剪影，主色对齐主题 `#00ADD8`）；`frontend/scripts/gen-brand.mjs`（@resvg/resvg-js 渲染多尺寸 PNG，PIL 打 ICO）派生 `build/appicon.png`、`build/windows/icon.ico`、`frontend/public/favicon.svg`、`frontend/src/assets/logo.svg`；favicon 已接入。
- 已验证可构建：`go build/vet/test ./...`、`npm run build`、`wails build` 均通过。

---

## 11. 变更记录

| 日期 | 变更 | 说明 |
|---|---|---|
| 2026-08-14 | 初版 + M1 | 技术选型、架构、目录结构；Wails 骨架，后端采用 `internal/` 目录 |
| 2026-08-14 | M2 注册机制 | Tool 接口 / registry / JSON string 协议 / `ListTools`+`RunTool` 绑定 + Base64 工具 |
| 2026-08-14 | 前端 it-tools 风格 | 亮/暗主题、顶栏 + Command Palette、分类中文化、工具图标、首页扁平网格 |
| 2026-08-14 | 转换器工具扩展 | 分批实现 8 个转换器（罗马/大小写/日期时间/整数基/文本二进制/文本Unicode/YAML→JSON/YAML→TOML），引入 yaml.v3 + go-toml/v2 |
| 2026-08-14 | 命名规范 | 目录短杠、文件下划线、包名去 `-converter` 后缀拼接 |
| 2026-08-14 | 注册重构 | 工具包自持完整元数据（导出 `Tool()` + `Executor`），`registerTools` 精简为一行注册 |
| 2026-08-14 | 新增 2 个转换器 | 列表转换（listconv）、Markdown 转 HTML（引入 goldmark） |
| 2026-08-14 | 新增 4 个转换器 | TOML→JSON、TOML→YAML、XML→JSON、JSON→XML（引入 mxj，XML↔JSON 对齐 `_attributes`/`_text` 约定），转换器分类基本完成 |
| 2026-08-14 | 品牌 logo | `assets/logo.svg` 设计源 + `gen-brand.mjs` 生成链路；接入 appicon/ico/favicon/侧边栏 hero 品牌位 |
| 2026-08-15 | Token 生成器 | 加密分类首个工具：`crypto/rand` 均匀采样，字符集开关 + 长度 1–512，支持自定义字符集；前端复刻 it-tools 布局 |
| 2026-08-15 | ToolTextarea 组件 | 封装"可选 label + 可拉伸多行输入框"（`resizable` 默认开）；先接入 token-generator，其余工具待观察后迁移 |
| 2026-08-15 | ToolTextarea 全量迁移 | 12 个工具 29 处 textarea 输入框迁移至 ToolTextarea；清理无用的 `NInput`/`useThemeVars`/`.field` 样式 |
| 2026-08-15 | Hash 文本工具 | 8 算法 × 4 编码；SHA3 用 legacy Keccak 对齐 crypto-js；引入 `x/crypto` 直接依赖（sha3/ripemd160） |
| 2026-08-15 | 加密/解密工具 | AES/TripleDES/Rabbit/RC4 + EVP KDF + OpenSSL 格式；Rabbit 移植 crypto-js；测试用 crypto-js 向量校验 |
| 2026-08-15 | 等宽字体 | `@fontsource/cascadia-code`（latin）；ToolTextarea 新增 `monospace` prop |
