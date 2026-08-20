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

- 开发：`wails dev -tags webkit2_41`（热重载）；构建：`wails build -tags webkit2_41`
- Go：`go build ./...` / `go test ./...`
- 前端：`npm run build`（`vue-tsc --noEmit && vite build`）；`npm run dev`（在 `frontend/` 运行）
- 重新生成工具注册：`go generate ./internal/app`（新增工具后）
- 重新生成 Wails 绑定：`wails generate module`（后端绑定方法变更后）

## 环境要求（本机为 WSL2 Ubuntu 26.04 LTS，工具链已装齐）

- Go（`/usr/local/go/bin`）、Node（nvm、npm）、Wails CLI、Linux gcc 均已装。
- Linux 构建依赖（apt）：`build-essential libgtk-3-dev libwebkit2gtk-4.1-dev \
  libayatana-appindicator3-dev librsvg2-dev pkg-config`
  （本机源只有 WebKitGTK **4.1**；CI 的 ubuntu-22.04 用 `libwebkit2gtk-4.0-dev`，不带 tag）
- 本机编译/开发命令需加 `-tags webkit2_41`（选用 WebKitGTK 4.1），见上文「命令」节。
- 显示由 WSLg 提供（`DISPLAY=:0`、`WAYLAND_DISPLAY=wayland-0`），`wails dev` 可直接弹出窗口。
- 本机 PATH 前缀与（如需要）网络代理配置见全局 AGENTS.md（`~/.pi/agent/AGENTS.md`）。

## 指导原则

- 工具逻辑属于 Go；前端只负责渲染和调用绑定；`frontend/wailsjs/` 由 Wails 生成，切勿手动编辑。
- 新增工具或改动功能后，保持 `SPEC.md` 与实现一致（工具清单、变更历史等易变内容不再维护，
  以代码扫描（toolsgen 注册）与 git 历史为准）。
- 完成改动后先构建（`wails build` 出当前平台二进制，WSL2 下为 Linux ELF）并停下，等用户查看
  效果并确认后再执行 `git commit`/`git push`。
- 提交信息遵循 Conventional Commits（见全局 AGENTS.md）。
