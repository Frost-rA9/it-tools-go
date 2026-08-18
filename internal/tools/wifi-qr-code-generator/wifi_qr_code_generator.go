// Package wifiqrcodegen 实现 WiFi 配置二维码生成器。
package wifiqrcodegen

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/qrutil"
	"it-tools-go/internal/registry"
)

const (
	ID          = "wifi-qr-code-generator"
	Name        = "WiFi 二维码生成器"
	Description = "生成可供手机扫描连接 WiFi 的二维码"
	Category    = "图片和视频"
	Icon        = "Qrcode"
)

var Keywords = []string{"qr", "code", "generator", "wifi", "WPA", "WEP", "SSID", "EAP", "二维码", "无线网络"}

func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

type input struct {
	SSID       string `json:"ssid"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"`
	EAPMethod  string `json:"eap_method"`
	Hidden     bool   `json:"hidden"`
	Anonymous  bool   `json:"anonymous"`
	Identity   string `json:"identity"`
	Phase2     string `json:"phase2"`
	Foreground string `json:"foreground"`
	Background string `json:"background"`
	Level      string `json:"level"`
}

type output struct {
	Payload string `json:"payload"`
	DataURL string `json:"data_url"`
}

type Executor struct{}

func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}
	payload, ok := wifiPayload(in)
	if !ok {
		return `{"payload":"","data_url":""}`, nil
	}
	dataURL, err := qrutil.GeneratePNGDataURL(payload, in.Foreground, in.Background, in.Level)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(output{Payload: payload, DataURL: dataURL})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

func wifiPayload(in input) (string, bool) {
	if in.SSID == "" {
		return "", false
	}
	ssid := escapeString(in.SSID)
	switch in.Encryption {
	case "nopass":
		return fmt.Sprintf("WIFI:S:%s;;", ssid), true
	case "WPA", "WEP":
		if in.Password == "" {
			return "", false
		}
		hidden := ""
		if in.Hidden {
			hidden = "H:true;"
		}
		return fmt.Sprintf("WIFI:S:%s;T:%s;P:%s;%s;", ssid, in.Encryption, escapeString(in.Password), hidden), true
	case "WPA2-EAP":
		if in.Password == "" || in.EAPMethod == "" || (!in.Anonymous && in.Identity == "") {
			return "", false
		}
		if in.EAPMethod == "PEAP" && in.Phase2 == "" {
			return "", false
		}
		identity := "I:" + escapeString(in.Identity)
		if in.Anonymous {
			identity = "A:anon"
		}
		phase2 := ""
		if in.Phase2 != "" && in.Phase2 != "None" {
			phase2 = "PH2:" + in.Phase2 + ";"
		}
		hidden := ""
		if in.Hidden {
			hidden = "H:true;"
		}
		return fmt.Sprintf("WIFI:S:%s;T:WPA2-EAP;P:%s;E:%s;%s%s;%s;", ssid, escapeString(in.Password), in.EAPMethod, phase2, identity, hidden), true
	default:
		return "", false
	}
}

func escapeString(value string) string {
	return strings.NewReplacer(`\`, `\\`, ";", `\;`, ",", `\,`, ":", `\:`, `"`, `\"`).Replace(value)
}
