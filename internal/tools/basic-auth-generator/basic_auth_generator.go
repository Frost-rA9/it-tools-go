// Package basicauth 实现基本身份验证生成器：根据用户名与密码生成 Authorization: Basic 请求头。
package basicauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/registry"
)

// 工具元数据。
const (
	ID          = "basic-auth-generator"
	Name        = "基本身份验证生成器"
	Description = "根据用户名与密码生成 Basic 身份验证请求头"
	Category    = "Web"
	Icon        = "Lock"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"basic", "auth", "generator", "username", "password", "base64", "authentication", "header", "authorization", "身份验证", "认证", "请求头"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// output 是工具的输出结构。
type output struct {
	Header string `json:"header"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(in.Username + ":" + in.Password))
	out := output{Header: "Authorization: Basic " + credentials}

	outJSON, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(outJSON), nil
}