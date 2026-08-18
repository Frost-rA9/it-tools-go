package svgplaceholder

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"width":600,"height":350,"font_size":26,"background":"#cccccc","foreground":"#333333","exact_size":true,"text":""}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if !strings.Contains(got.SVG, `viewBox="0 0 600 350" width="600" height="350"`) || !strings.Contains(got.SVG, ">600x350</text>") {
		t.Errorf("SVG 默认内容不符: %s", got.SVG)
	}
	if !strings.HasPrefix(got.Base64, "data:image/svg+xml;base64,") {
		t.Fatalf("Base64 Data URL 前缀不符")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(got.Base64, "data:image/svg+xml;base64,"))
	if err != nil || string(decoded) != got.SVG {
		t.Errorf("Base64 内容与 SVG 不一致: %v", err)
	}
}

func TestCustomTextEscaped(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"width":100,"height":50,"font_size":12,"background":"#fff","foreground":"#000","exact_size":false,"text":"<Hello & goodbye>"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	json.Unmarshal([]byte(out), &got)
	if !strings.Contains(got.SVG, "&lt;Hello &amp; goodbye&gt;") || strings.Contains(got.SVG, "<Hello & goodbye>") {
		t.Errorf("自定义文本未正确转义: %s", got.SVG)
	}
}

func TestInvalidSize(t *testing.T) {
	var e Executor
	if _, err := e.Execute(t.Context(), `{"width":0,"height":100,"font_size":12}`); err == nil {
		t.Error("无效尺寸应报错")
	}
}
