# IT Tools Go

跨平台桌面 IT 工具集，使用 [Wails v2](https://wails.io/)（Go 后端）+ Vue 3 + Naive UI 实现，参考 [it-tools.tech](https://it-tools.tech/)。

## 特性

- 工具业务逻辑全部由 Go 实现（`internal/`），前端仅负责渲染与调用
- 单一可执行文件，跨平台（Windows / macOS / Linux）

## 开发

环境要求：Go、gcc（Windows 下为 mingw）、Node.js 与 Wails CLI。

```sh
wails dev      # 热重载开发
wails build    # 构建桌面二进制
go test ./...  # Go 测试
```

## 新增工具

1. 新建 Go 包 `internal/tools/<name>/`
2. 新建前端组件 `frontend/src/views/tools/<toolId>.vue`（文件名必须等于 toolId）
3. 在 `internal/app/app.go` 的 `registerTools` 中注册

详见 [SPEC.md](SPEC.md)。

## 致谢

本项目灵感来源于 [IT-Tools](https://it-tools.tech/)，其仓库地址为 [github.com/CorentinTh/it-tools](https://github.com/CorentinTh/it-tools)。感谢原作者的杰出工作。

## License

[GNU GPLv3](LICENSE)
