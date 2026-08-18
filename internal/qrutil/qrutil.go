// Package qrutil 提供项目内 QR 码 PNG 渲染辅助函数。
package qrutil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// GeneratePNGDataURL 将文本编码为 PNG QR 码 Data URL。
func GeneratePNGDataURL(text, foreground, background, level string) (string, error) {
	levelOption, err := errorCorrectionLevel(level)
	if err != nil {
		return "", err
	}

	qrc, err := qrcode.NewWith(text, levelOption)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %w", err)
	}

	options := []standard.ImageOption{
		standard.WithQRWidth(8),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithFgColorRGBHex(normalizeColor(foreground, "#000000")),
	}
	if isTransparent(background) {
		options = append(options, standard.WithBgTransparent())
	} else {
		options = append(options, standard.WithBgColorRGBHex(normalizeColor(background, "#FFFFFF")))
	}

	var buf bytes.Buffer
	writer := standard.NewWithWriter(&writeCloser{Writer: &buf}, options...)
	if err := qrc.Save(writer); err != nil {
		return "", fmt.Errorf("渲染二维码失败: %w", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func errorCorrectionLevel(level string) (qrcode.EncodeOption, error) {
	switch strings.ToLower(level) {
	case "", "medium":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium), nil
	case "low":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow), nil
	case "quartile":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart), nil
	case "high":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest), nil
	default:
		return nil, fmt.Errorf("未知二维码纠错级别: %q", level)
	}
}

func normalizeColor(color, fallback string) string {
	color = strings.TrimSpace(color)
	if color == "" {
		return fallback
	}
	if strings.HasPrefix(color, "#") {
		color = color[1:]
	}
	if len(color) == 8 {
		color = color[:6]
	}
	if len(color) != 3 && len(color) != 6 {
		return fallback
	}
	if len(color) == 3 {
		color = strings.Map(func(r rune) rune { return r }, color)
	}
	return "#" + color
}

func isTransparent(color string) bool {
	color = strings.TrimPrefix(strings.TrimSpace(color), "#")
	return len(color) == 8 && strings.EqualFold(color[6:], "00")
}

type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }
