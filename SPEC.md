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

## 6. 工具分类

对齐参考项目分类（常量见 `internal/registry/registry.go`）。当前已实现五类：

| 分类 | 数量 | 工具 |
|---|---|---|
| 加密 | 10 | Token、Hash、加密/解密、BCrypt、UUID、ULID、BIP39、HMAC、RSA 密钥对、密码强度 |
| 转换器 | 17 | Base64、罗马数字、大小写、日期时间、进制、文本↔二进制/Unicode、列表、Markdown→HTML、JSON/TOML/XML/YAML 互转 |
| Web | 15 | URL 编码/解码、HTML 实体、URL 分析器、JWT、HTTP 状态码、JSON 差异、设备信息、UA 分析、Basic Auth、OTP、OG 元生成、MIME、Slug、SafeLink 解码、按键码 |
| 图片和视频 | 3 | 二维码生成器、WiFi 二维码生成器、SVG 占位符生成器 |
| 开发 | 13 | Git 备忘录、随机端口生成器、Crontab 表达式生成器、Chmod 计算器、JSON 格式化、JSON 压缩、JSON 转 CSV、SQL/XML/YAML 格式化、Docker Run 转 Compose、Regex 测试器、正则表达式速查表 |

其余分类（网络/数学/测量/文本）后续扩展按此归类。

## 7. 开发与发布

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

## 8. 测试与质量

- Go：表驱动单测（各工具 + registry）；前端：`vue-tsc` 类型检查。
- 提交遵循 Conventional Commits；Go 用 gofmt/vet。

## 9. 里程碑

| 阶段 | 内容 | 状态 |
|---|---|---|
| M0 环境准备 | 安装 Go/gcc/Wails/pnpm | ✅ |
| M1 骨架搭建 | `wails init` + `internal/` 目录结构 | ✅ |
| M2 注册机制 | Tool 接口 + registry + RunTool Binding + 前端动态渲染 | ✅ |
| M3 工具扩展 | 高频工具实现（当前 43 个） | ✅ 各工具均有单测 |
| M4 打包发布 | 三平台交叉编译 + GitHub Release 流程 | ✅（v0.1.0/v0.2.0） |

## 10. 当前状态

- **版本**：最新发布 v0.4.0「El Shaddoll Grysta」（git tag `v0.4.0`）；main 与发布 tag 同步。
- 注册机制：`registry` + JSON string 协议 + `ListTools`/`RunTool` 绑定；注册由 `internal/toolsgen`
  扫描生成 `internal/app/tools_gen.go`；前端 `import.meta.glob` 按 toolId 动态加载。
- 已实现工具（58 个）：转换器 17 + 加密 10 + Web 15 + 图片和视频 3 + 开发 13（清单见 §6）。前三类已与 it-tools 对齐；开发分类已完结。
- 前端 it-tools 风格：亮/暗主题、侧边栏分类菜单、Command Palette、首页网格；通用组件
  `ToolTextarea` / `ToolCodeBlock`；等宽字体 Cascadia Code 随包分发。
- 品牌标识：`assets/logo.svg` 唯一源 → `build/appicon.png`、`build/windows/icon.ico`、favicon。
- 命名规范：目录短杠（toolId）、文件下划线、包名去 `-converter` 后缀拼接。
- 已验证（2026-08-18 复核）：`go build/vet/test ./...`、`vue-tsc --noEmit`、`npm run build` 通过；`wails build` 在 v0.3.0 发布周期通过。

## 11. 变更记录

| 日期 | 变更 | 说明 |
|---|---|---|
| 2026-08-14 | 初版搭建 | M1 骨架 + M2 注册机制（Tool 接口/registry/JSON 协议/RunTool）+ 前端 it-tools 风格 + Base64 首个工具 |
| 2026-08-14 | 转换器 15 个 | 罗马/大小写/日期时间/进制/文本↔二进制·Unicode/列表/Markdown→HTML/TOML·XML·YAML 互转（goldmark、mxj、yaml.v3、go-toml） |
| 2026-08-14 | 品牌与注册重构 | `assets/logo.svg` 设计源 + `gen-brand.mjs` 链路；`internal/toolsgen` 代码生成器；`app.go` 精简 |
| 2026-08-15 | 加密 10 个 | Token/Hash/加密解密/BCrypt/UUID/ULID/BIP39/HMAC/RSA 密钥对/密码强度 |
| 2026-08-15 | 通用组件 | `ToolTextarea`（可拉伸+monospace）、`ToolCodeBlock`（只读等宽块）；Cascadia Code 随包 |
| 2026-08-15 | Web 15 个 | URL/HTML 实体/URL 分析/JWT/HTTP 状态码/JSON 差异/设备/UA/Basic Auth/OTP/OG 元/MIME/Slug/SafeLink/按键码 —— Web 分类与 it-tools 对齐完成 |
| 2026-08-15 | 文档精简 | SPEC.md / AGENTS.md 精简为要点式；README.md 重写并补充工具清单与发布说明 |
| 2026-08-18 | 图片和视频 3 个 | 二维码生成器、WiFi 二维码生成器、SVG 占位符生成器（go-qrcode v2）；同步注册、图标与文档 |
| 2026-08-18 | Git 备忘录 | 开发分类首个工具；纯静态 Git 命令速查页，用于验证长文本展示布局 |
| 2026-08-18 | 开发 3 个 | 随机端口生成器、Crontab 表达式生成器、Chmod 计算器（math/rand/v2、标准库）；工具总数 44→47，开发分类 1→4 |
| 2026-08-18 | 开发 JSON 3 个 | JSON 格式化（缩进/排序）、JSON 压缩（大小对比）、JSON 转 CSV（分隔符/表头，encoding/csv 转义）；toolsgen 改为包级契约检查；工具总数 47→50；三工具双卡布局并封装公共 .tool-card 样式 |
| 2026-08-18 | 格式化 3 个 | SQL 格式化（自研 tokenizer：子句换行/括号缩进/关键字大写）、XML 格式化（encoding/xml Token 流重建）、YAML 格式化（yaml.v3 Node 保留注释）；工具总数 50→53，开发分类 7→10 |
| 2026-08-18 | crontab 增强 | 新增表达式直接解析（Go 侧自研解析器：*/步长/范围/列表/@简写/秒字段/名称别名）+ it-tools 风格帮助表格与描述区 |
| 2026-08-18 | JSON 互转 2 个 | JSON 转 YAML、JSON 转 TOML（补齐 it-tools Converter 互转矩阵）；工具总数 53→55，转换器 15→17 |
| 2026-08-18 | 开发收尾 3 个 | Docker Run 转 Compose（自研 shell 分词+参数解析）、Regex 测试器（Go regexp 捕获组/命名组）、正则表达式速查表（静态页）；工具总数 55→58，开发分类 10→13 完结 |
| 2026-08-18 | 修复首页布局 | 全局 .tool-card 样式（双卡布局用）误伤主页网格（min-width 400px 撑破列数）；主页卡片类名改为 tool-card-item，样式只作用于工具页卡片 |
| 2026-08-18 | 统一首页卡片高度 | 描述区 line-clamp 最多两行导致短描述卡片矮一截；固定 line-height 1.5 + min-height 两行，所有工具卡片等高 |
| 2026-08-18 | 首页卡片等高加固 | 名称区显式 line-height 24px、图标区固定 40×40 + flex-shrink:0、描述区改为固定 height 42px；彻底消除中英文行高/盒子尺寸差异导致的卡片高度不齐 |
| 2026-08-19 | 统一首页网格行高 | HomeView 网格增加 `grid-auto-rows: 180px`，统一不同行的工具卡片高度 |
| 2026-08-18 | 发布 v0.4.0 | 代号 El Shaddoll Grysta；58 个工具（加密 10 / 转换器 17 / Web 15 / 图片和视频 3 / 开发 13），开发分类完结 |
