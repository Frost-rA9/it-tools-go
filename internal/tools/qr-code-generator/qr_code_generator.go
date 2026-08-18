// Package qrcodegen 实现通用二维码生成器。
package qrcodegen

import (
	"context"
	"encoding/json"
	"fmt"

	"it-tools-go/internal/qrutil"
	"it-tools-go/internal/registry"
)

const (
	ID          = "qr-code-generator"
	Name        = "二维码生成器"
	Description = "将文本或链接生成可下载的二维码"
	Category    = "图片和视频"
	Icon        = "Qrcode"
)

var Keywords = []string{"qr", "code", "generator", "square", "color", "link", "low", "medium", "quartile", "high", "transparent", "二维码"}

func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

type input struct {
	Text       string `json:"text"`
	Foreground string `json:"foreground"`
	Background string `json:"background"`
	Level      string `json:"level"`
}

type output struct {
	DataURL string `json:"data_url"`
}

type Executor struct{}

func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if in.Text == "" {
		return `{"data_url":""}`, nil
	}

	dataURL, err := qrutil.GeneratePNGDataURL(in.Text, in.Foreground, in.Background, in.Level)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(output{DataURL: dataURL})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
