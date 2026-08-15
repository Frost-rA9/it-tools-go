# AGENTS.md

## 项目介绍

it-tools-go 是一个跨平台桌面 IT 工具集，参考 [it-tools.tech](https://it-tools.tech/)，使用
Wails v2（Go 后端）+ Vue 3 + Naive UI 实现。工具业务逻辑全部在 Go 侧（`internal/`），
前端只负责渲染与调用绑定。设计规格见 SPEC.md。

## 项目结构

```
it-tools-go/
├── go.mod / go.sum              # Go 模块
├── main.go                      # Wails 入口（wails.Run + embed frontend/dist，package main）
├── wails.json                   # Wails 构建配置（frontend:install/build 命令）
├── internal/
│   ├── app/                     # App 结构体 + 生命周期 + 绑定（ListTools/RunTool）
│   ├── registry/                # Tool 接口 + 注册表 + 分类常量
│   └── tools/                   # 各工具实现（每工具一个包）
│       └── base64-string-converter/  # Base64 工具 + 单测（目录用短杠 toolId）
├── build/                       # 构建资源（图标、Windows 清单、NSIS 安装器）
│   └── bin/                     # wails build 产物（it-tools-go.exe）
├── frontend/
│   ├── package.json             # npm 脚本：dev / build / preview
│   ├── vite.config.ts
│   ├── tsconfig.json / tsconfig.node.json
│   └── src/
│       ├── main.ts / App.vue / style.css
│       ├── theme.ts             # Naive UI 主题覆盖（亮/暗）
│       ├── components/          # Toolbar、ToolMenu、ToolCard、CommandPalette
│       ├── router/              # vue-router 配置
│       ├── layouts/             # BaseLayout（侧边栏+顶栏）、ToolLayout（工具页）
│       ├── views/               # HomeView、ToolView、tools/（按 toolId 命名）
│       ├── stores/              # Pinia（tools、ui）
│       ├── composables/         # useToolComponent（glob 动态加载）
│       ├── assets/              # 字体、图片
│       └── wailsjs/             # Wails 自动生成的绑定（切勿手动编辑）
└── SPEC.md                      # 工程规格书（权威）
```

## 工具开发约定

- 新增一个工具 = Go 包（`internal/tools/<name>/`）+ 前端组件（`frontend/src/views/tools/<toolId>.vue`）。
- 前端组件文件名必须等于 toolId（如 `base64-string-converter.vue`），`import.meta.glob` 据此自动匹配。
- 命名规范：目录名用短杠 `-` 连接的完整 toolId（如 `base64-string-converter/`）；`.go` 文件名用下划线 `_`（如 `base64_string_converter.go`）；`package` 名用去短杠拼接、去掉冗余 `-converter` 后缀的小写单字（如 `base64string`、`romannumeral`、`caseconv`、`radix`，因 Go 包名不能含短杠/下划线，且 `case` 为关键字）。
- 工具注册：由代码生成器 `internal/toolsgen` 扫描 `internal/tools/` 自动生成
  `internal/app/tools_gen.go`（按目录名排序，勿手改）。新增工具后执行
  `go generate ./internal/app` 刷新。每个工具包须导出 `Tool()`（返回 `registry.Tool`
  元数据）与 `Executor`。
- 传输协议为 JSON string（前端 `JSON.stringify/parse`，Go `encoding/json`）。

## 命令

- 开发（热重载）：`wails dev`
- 构建桌面二进制：`wails build`
- Go 构建：`go build ./...`
- Go 测试：`go test ./...`
- 前端类型检查 + 构建：`npm run build`（执行 `vue-tsc --noEmit && vite build`）
- 前端开发服务器：`npm run dev`（在 `frontend/` 目录运行）
- 重新生成 Wails 绑定：`wails generate module`（后端绑定方法变更后需执行）
- 重新生成工具注册：`go generate ./internal/app`（新增工具后执行）

## 环境要求

- 需要 Go、gcc（mingw）、Node.js 及 Wails CLI。本机具体的 PATH 前缀与网络代理配置
  见全局 AGENTS.md（`~/.config/opencode/AGENTS.md`）。

## 指导原则

- 工具逻辑属于 Go（遵循 SPEC.md §2.1）；前端只负责渲染和调用绑定。
- `frontend/wailsjs/` 由 Wails 生成 —— 切勿手动编辑。
- 新增工具或改动功能后，须同步更新 `SPEC.md`（§10 当前状态保持为最新快照、不堆叠历史；§11 变更记录追加或精简）。
- 完成改动后先构建（含 `wails build` 出 exe）并停下，等用户查看效果并确认后再执行 `git commit`/`git push`。
- 提交信息遵循 Conventional Commits（见全局 AGENTS.md）。
