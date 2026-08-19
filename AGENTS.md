# AGENTS.md

## 项目介绍

it-tools-go 是跨平台桌面 IT 工具集，参考 [it-tools.tech](https://it-tools.tech/)，使用
Wails v2（Go 后端）+ Vue 3 + Naive UI 实现。工具业务逻辑全部在 Go 侧（`internal/`），
前端只负责渲染与调用绑定。设计规格见 [SPEC.md](SPEC.md)。

## 工具开发约定

- 新增一个工具 = Go 包（`internal/tools/<toolId>/`）+ 前端组件
  （`frontend/src/views/tools/<toolId>.vue`，文件名必须等于 toolId，`import.meta.glob` 据此匹配）。
- 命名：目录用短杠（完整 toolId）、`.go` 文件用下划线、`package` 名用去短杠拼接的小写单字
  （如 `base64string`、`romannumeral`，因 Go 包名不能含短杠/下划线，且 `case` 为关键字）。
- 注册：`internal/toolsgen` 扫描 `internal/tools/` 自动生成 `internal/app/tools_gen.go`
  （按目录名排序，勿手改）；新增工具后执行 `go generate ./internal/app` 刷新。
  每个工具包须导出 `Tool()`（返回 `registry.Tool` 元数据）与 `Executor`。
- 传输协议为 JSON string（前端 `JSON.stringify/parse`，Go `encoding/json`）。

## 命令

- 开发：`wails dev`（热重载）；构建：`wails build`
- Go：`go build ./...` / `go test ./...`
- 前端：`npm run build`（`vue-tsc --noEmit && vite build`）；`npm run dev`（在 `frontend/` 运行）
- 重新生成工具注册：`go generate ./internal/app`（新增工具后）
- 重新生成 Wails 绑定：`wails generate module`（后端绑定方法变更后）

## 环境要求

需要 Go、gcc（mingw）、Node.js 及 Wails CLI。本机 PATH 前缀与网络代理配置见全局 AGENTS.md
（`~/.config/opencode/AGENTS.md`）。

## 指导原则

- 工具逻辑属于 Go；前端只负责渲染和调用绑定；`frontend/wailsjs/` 由 Wails 生成，切勿手动编辑。
- 新增工具或改动功能后，须同步更新 `SPEC.md`（§10 当前状态保持最新快照、§11 变更记录追加或精简）。
  §11 变更记录**只记录**：功能添加与实现、bug 修复、功能性重构、工程结构重构；
  样式/文案/预置值/非功能性重构等细节性改动一般不记录（确有必要才追加）。
- 完成改动后先构建（含 `wails build` 出 exe）并停下，等用户查看效果并确认后再执行
  `git commit`/`git push`。
- 提交信息遵循 Conventional Commits（见全局 AGENTS.md）。
