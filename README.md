# IT Tools Go

<p align="center">
  <b>跨平台桌面 IT 工具集</b> —— 常用开发、运维与 IT 工具合集，单文件离线可用。<br>
  <a href="https://github.com/Frost-rA9/it-tools-go/releases">下载即用</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/Frost-rA9/it-tools-go" alt="Release">
  <img src="https://img.shields.io/badge/Go-1.26-blue" alt="Go">
  <img src="https://img.shields.io/badge/License-GPLv3-blue.svg" alt="License: GPLv3">
</p>

## 特性

- 工具业务逻辑**全部由 Go 实现**（`internal/`），前端仅负责渲染与调用
- 单一可执行文件，跨平台（Windows / macOS / Linux），**无需联网与外部依赖**
- it-tools 风格界面：亮/暗主题、Command Palette、侧边栏分类导航、首页工具网格
- 表驱动单元测试，输出对齐 it-tools / crypto-js / node-forge 等参考实现
- 完整工具清单见 [SPEC.md](SPEC.md) §6（目前 10 个分类）

## 安装与使用

从 [GitHub Releases](https://github.com/Frost-rA9/it-tools-go/releases) 下载对应平台的二进制
（Windows `.exe` / macOS `.dmg` / Linux `AppImage`），免安装直接运行。

## 技术栈

| 组件 | 技术 |
|---|---|
| 桌面壳 | [Wails v2](https://wails.io/)（单二进制、原生渲染引擎） |
| 后端 | Go 1.26 |
| 前端 | Vue 3 + TypeScript + Vite + [Naive UI](https://www.naiveui.com/) |

## 开发

环境要求：Go、gcc（Windows 下为 mingw）、Node.js 与 Wails CLI。

```sh
wails dev                   # 热重载开发
wails build                 # 构建桌面二进制
go test ./...               # Go 单元测试
go generate ./internal/app  # 新增工具后刷新注册（tools_gen.go）
npm run build               # 前端类型检查 + 构建（vue-tsc + vite）
```

## 新增工具

1. 新建 Go 包 `internal/tools/<toolId>/`，导出 `Tool()`（元数据）与 `Executor`（JSON string 协议）
2. 新建前端组件 `frontend/src/views/tools/<toolId>.vue`（文件名必须等于 toolId）
3. 执行 `go generate ./internal/app` 刷新注册

详细约定见 [SPEC.md](SPEC.md)。

## 发布

推送 `v*` 标签自动触发 [GitHub Actions](.github/workflows/release.yml)：
三平台交叉编译（Windows amd64 / macOS Intel+Silicon / Linux amd64）并创建 GitHub Release，
版本号由 git tag 决定（如 `v0.5.0`）。

## 致谢

本项目灵感来源于 [IT-Tools](https://it-tools.tech/)，
其仓库地址为 [github.com/CorentinTh/it-tools](https://github.com/CorentinTh/it-tools)。
感谢原作者的杰出工作。

## License

本项目基于 [GNU GPLv3](LICENSE)。