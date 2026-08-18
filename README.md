# IT Tools Go

跨平台桌面 IT 工具集，使用 [Wails v2](https://wails.io/)（Go 后端）+ Vue 3 + Naive UI 实现，参考 [it-tools.tech](https://it-tools.tech/)。工具业务逻辑全部由 Go 实现（`internal/`），前端仅负责渲染与调用。

## 特性

- 工具业务逻辑全部在 Go 侧（`internal/`），前端仅负责渲染与调用
- 单一可执行文件，跨平台（Windows / macOS / Linux），离线可用
- 亮/暗主题、Command Palette、工具按分类组织
- Go 工具均有单元测试，输出对齐 it-tools / crypto-js / node-forge 等参考实现

## 已实现工具（55）

**加密（10）**：Token 生成器 · Hash 文本 · 加密/解密 · BCrypt · UUID 生成器 · ULID 生成器 · BIP39 密码生成器 · HMAC 生成器 · RSA 密钥对生成器 · 密码强度分析仪

**转换器（17）**：Base64 · 罗马数字 · 大小写转换 · 日期时间转换 · 整数进制 · 文本↔ASCII 二进制 · 文本↔Unicode · 列表转换 · Markdown→HTML · JSON↔YAML/TOML/XML · TOML↔YAML

**Web（15）**：URL 编码/解码 · HTML实体转义 · URL 分析器 · JWT 解析器 · HTTP 状态码 · JSON 差异比较 · 设备信息 · 用户代理分析器 · 基本身份验证生成器 · OTP 代码生成器 · 开放式图形元生成器 · MIME 类型转换器 · Slug 化字符串 · Outlook 安全链接解码器 · 按键码信息

**图片和视频（3）**：二维码生成器 · WiFi 二维码生成器 · SVG 占位符生成器

**开发（10）**：Git 备忘录 · 随机端口生成器 · Crontab 表达式生成器 · Chmod 计算器 · JSON 格式化 · JSON 压缩 · JSON 转 CSV · SQL 格式化 · XML 格式化 · YAML 格式化

## 技术栈

| 组件 | 技术 |
|---|---|
| 桌面壳 | [Wails v2](https://wails.io/)（原生渲染引擎，单二进制） |
| 后端 | Go 1.26 |
| 前端 | Vue 3 + TypeScript + Vite + [Naive UI](https://www.naiveui.com/) |

## 开发

环境要求：Go、gcc（Windows 下为 mingw）、Node.js 与 Wails CLI。

```sh
wails dev                    # 热重载开发
wails build                  # 构建桌面二进制
go test ./...                # Go 测试
go generate ./internal/app   # 新增工具后刷新注册（tools_gen.go）
npm run build                # 前端类型检查 + 构建
```

## 新增工具

1. 新建 Go 包 `internal/tools/<name>/`，导出 `Tool()`（元数据）与 `Executor`（`Execute(ctx, input) (output, error)`，JSON string 协议）
2. 新建前端组件 `frontend/src/views/tools/<toolId>.vue`（文件名必须等于 toolId）
3. 执行 `go generate ./internal/app` 重新生成注册

详见 [SPEC.md](SPEC.md)。

## 发布

推送 `v*` 标签自动触发 [GitHub Actions](.github/workflows/release.yml)：三平台交叉编译（Windows amd64 / macOS Intel+Silicon / Linux amd64），创建 GitHub Release 并附带各平台产物。

## 致谢

本项目灵感来源于 [IT-Tools](https://it-tools.tech/)，其仓库地址为 [github.com/CorentinTh/it-tools](https://github.com/CorentinTh/it-tools)。感谢原作者的杰出工作。

## License

[GNU GPLv3](LICENSE)
