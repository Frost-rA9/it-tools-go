// Package svgplaceholder 实现 SVG 占位符生成器。
package svgplaceholder

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"

	"it-tools-go/internal/registry"
)

const (
	ID          = "svg-placeholder-generator"
	Name        = "SVG 占位符生成器"
	Description = "生成自定义尺寸、颜色与文本的 SVG 占位图"
	Category    = "图片和视频"
	Icon        = "Photo"
)

var Keywords = []string{"svg", "placeholder", "generator", "image", "size", "mockup", "占位图", "图片"}

func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

type input struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FontSize  int    `json:"font_size"`
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	ExactSize bool   `json:"exact_size"`
	Text      string `json:"text"`
}

type output struct {
	SVG    string `json:"svg"`
	Base64 string `json:"base64"`
}

type Executor struct{}

func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	if in.Width < 1 || in.Height < 1 || in.FontSize < 1 {
		return "", fmt.Errorf("宽度、高度和字体大小必须大于 0")
	}
	if in.Background == "" {
		in.Background = "#cccccc"
	}
	if in.Foreground == "" {
		in.Foreground = "#333333"
	}
	text := in.Text
	if text == "" {
		text = fmt.Sprintf("%dx%d", in.Width, in.Height)
	}
	text = html.EscapeString(text)

	size := ""
	if in.ExactSize {
		size = fmt.Sprintf(` width="%d" height="%d"`, in.Width, in.Height)
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d"%s>
  <rect width="%d" height="%d" fill="%s"></rect>
  <text x="50%%" y="50%%" dominant-baseline="middle" text-anchor="middle" font-family="monospace" font-size="%dpx" fill="%s">%s</text>
</svg>`, in.Width, in.Height, size, in.Width, in.Height, in.Background, in.FontSize, in.Foreground, text)

	out, err := json.Marshal(output{
		SVG:    svg,
		Base64: "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg)),
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}
