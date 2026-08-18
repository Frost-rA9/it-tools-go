// Package xmlfmt 实现 XML 格式化工具：基于 encoding/xml Token 流重建缩进。
package xmlfmt

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"it-tools-go/internal/registry"
)

const (
	ID          = "xml-formatter"
	Name        = "XML 格式化"
	Description = "美化 XML：按层级缩进，保留注释与文本"
	Category    = "开发"
	Icon        = "Code"
)

var Keywords = []string{"xml", "format", "pretty", "美化", "格式化", "缩进", "indent"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	XML    string `json:"xml"`    // 待格式化 XML
	Indent string `json:"indent"` // "2" | "4" | "\t"（默认 "2"）
}

// output 是工具的输出结构。
type output struct {
	Formatted string `json:"formatted"`
	LineCount int    `json:"line_count"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回格式化结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	indent := in.Indent
	if indent == "" {
		indent = "2"
	}
	switch indent {
	case "2":
		indent = "  "
	case "4":
		indent = "    "
	case "\t":
	default:
		return "", fmt.Errorf("缩进仅支持 2 / 4 / Tab（当前 %q）", indent)
	}

	formatted, err := Format(in.XML, indent)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{
		Formatted: formatted,
		LineCount: strings.Count(strings.TrimRight(formatted, "\n"), "\n") + 1,
	})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Format 基于 encoding/xml 的 Token 流重建带缩进的 XML。
func Format(xmlText, indent string) (string, error) {
	dec := xml.NewDecoder(strings.NewReader(xmlText))
	var sb bytes.Buffer
	depth := 0
	lastWasStart := false // 上一个 token 是 StartElement（可能输出空元素 <a/>）
	lastWasText := false  // 上一个输出是非空文本（EndElement 需同行）

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("XML 解析失败: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if sb.Len() > 0 {
				sb.WriteByte('\n')
				writeIndent(&sb, depth, indent)
			}
			writeStartElement(&sb, t)
			depth++
			lastWasStart = true
			lastWasText = false
		case xml.EndElement:
			depth--
			if depth < 0 {
				depth = 0
			}
			switch {
			case lastWasStart:
				// 空元素：<a> → <a/>
				s := sb.String()
				if strings.HasSuffix(s, ">") {
					sb.Reset()
					sb.WriteString(strings.TrimSuffix(s, ">") + "/>")
				} else {
					sb.WriteString("</" + t.Name.Local + ">")
				}
			case lastWasText:
				// 文本节点同行闭合：<a>text</a>
				sb.WriteString("</" + t.Name.Local + ">")
			default:
				sb.WriteByte('\n')
				writeIndent(&sb, depth, indent)
				sb.WriteString("</" + t.Name.Local + ">")
			}
			lastWasStart = false
			lastWasText = false
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				if !lastWasStart && !lastWasText {
					sb.WriteByte('\n')
					writeIndent(&sb, depth, indent)
				}
				sb.WriteString(text)
				lastWasStart = false
				lastWasText = true
			}
		case xml.Comment:
			sb.WriteByte('\n')
			writeIndent(&sb, depth, indent)
			sb.WriteString("<!--")
			sb.WriteString(string(t))
			sb.WriteString("-->")
			lastWasStart = false
			lastWasText = false
		case xml.ProcInst:
			sb.WriteByte('\n')
			writeIndent(&sb, depth, indent)
			sb.WriteString("<?")
			sb.WriteString(t.Target)
			sb.WriteByte(' ')
			sb.WriteString(string(t.Inst))
			sb.WriteString("?>")
			lastWasStart = false
			lastWasText = false
		case xml.Directive:
			sb.WriteByte('\n')
			writeIndent(&sb, depth, indent)
			sb.WriteString("<!")
			sb.WriteString(string(t))
			sb.WriteString(">")
			lastWasStart = false
			lastWasText = false
		}
	}

	return strings.TrimRight(sb.String(), " \t\n") + "\n", nil
}

func writeIndent(sb *bytes.Buffer, depth int, indent string) {
	sb.WriteString(strings.Repeat(indent, depth))
}

func writeStartElement(sb *bytes.Buffer, t xml.StartElement) {
	sb.WriteByte('<')
	sb.WriteString(t.Name.Local)
	for _, attr := range t.Attr {
		sb.WriteByte(' ')
		if attr.Name.Space != "" {
			sb.WriteString(attr.Name.Space)
			sb.WriteByte(':')
		}
		sb.WriteString(attr.Name.Local)
		sb.WriteString(`="`)
		sb.WriteString(escapeAttr(attr.Value))
		sb.WriteByte('"')
	}
	sb.WriteByte('>')
}

// escapeAttr 转义属性值中的特殊字符（xml.EscapeText 行为）。
func escapeAttr(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}
