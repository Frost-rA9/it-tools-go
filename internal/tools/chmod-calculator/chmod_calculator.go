// Package chmodcalc 实现 Chmod 计算器：八进制与符号权限双向转换，展示特殊位。
package chmodcalc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/registry"
)

const (
	ID          = "chmod-calculator"
	Name        = "Chmod 计算器"
	Description = "八进制与符号模式权限互转，展示 owner/group/others 与特殊位"
	Category    = "开发"
	Icon        = "LockAccess"
)

var Keywords = []string{"chmod", "permission", "权限", "rwx", "755", "suid", "sgid", "sticky", "文件权限"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Mode string `json:"mode"` // 如 "755"、"4755"、"rwxr-xr-x"、"-rwxr-xr-x"
}

// output 是工具的输出结构。
type output struct {
	Octal       string `json:"octal"`        // 八进制（含特殊位时为 4 位）
	Symbolic    string `json:"symbolic"`     // 符号位串（含 s/t 特殊位标记）
	Owner       string `json:"owner"`        // 属主权限（3 位）
	Group       string `json:"group"`        // 属组权限（3 位）
	Others      string `json:"others"`       // 其他用户权限（3 位）
	Special     int    `json:"special"`      // 特殊位数值（0-7）
	SpecialText string `json:"special_text"` // 特殊位中文说明
	HasSUID     bool   `json:"has_suid"`     // 设置用户 ID
	HasSGID     bool   `json:"has_sgid"`     // 设置组 ID
	HasSticky   bool   `json:"has_sticky"`   // 粘滞位
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回解析结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	mode, err := parseMode(in.Mode)
	if err != nil {
		return "", err
	}

	out := buildOutput(mode)
	raw, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(raw), nil
}

// parseMode 解析八进制（3-4 位）或符号（9/10 字符）权限表示。
func parseMode(s string) (uint32, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("权限为空")
	}

	// 全数字 → 八进制
	if isOctal(s) {
		if len(s) < 3 || len(s) > 4 {
			return 0, fmt.Errorf("八进制模式长度必须为 3 或 4 位（当前 %q）", s)
		}
		var v uint32
		for _, c := range s {
			v = v*8 + uint32(c-'0')
		}
		return v, nil
	}

	// 符号模式：可选一位文件类型前缀 + 9 位权限
	body := s
	if len(s) > 9 && (s[0] == '-' || s[0] == 'd' || s[0] == 'l' || s[0] == 'b' || s[0] == 'c') {
		body = s[1:]
	}
	if len(body) != 9 {
		return 0, fmt.Errorf("符号模式长度必须为 9 位（可选类型前缀，当前 %q）", s)
	}

	const rwx = "rwx"
	var mode uint32
	// 特殊位标志：owner/group/others 的 x 位分别承载 s/S/t/T
	groups := []struct {
		start   int
		special uint32
	}{{0, 0o4000}, {3, 0o2000}, {6, 0o1000}}

	for gi, g := range groups {
		for i := 0; i < 3; i++ {
			ch := body[g.start+i]
			base := uint32(1) << (8 - uint(gi*3+i)) // owner rwx=0o400/200/100, group=0o040/020/010, others=0o004/002/001
			switch ch {
			case rwx[i], '-':
				if ch == rwx[i] {
					mode |= base
				}
			case 's', 'S':
				if i != 2 || gi == 2 {
					return 0, fmt.Errorf("'%c' 只能出现在属主或属组的执行位上", ch)
				}
				mode |= g.special
				if ch == 's' {
					mode |= base
				}
			case 't', 'T':
				if gi != 2 || i != 2 {
					return 0, fmt.Errorf("'%c' 只能出现在其他用户的执行位", ch)
				}
				mode |= 0o1000
				if ch == 't' {
					mode |= base
				}
			default:
				return 0, fmt.Errorf("非法字符 %q（取 %d，期望 %c/-/s/S/t/T）", ch, g.start+i, rwx[i])
			}
		}
	}
	return mode, nil
}

func isOctal(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '7' {
			return false
		}
	}
	return true
}

// buildOutput 由 mode 构造全部展示字段。
func buildOutput(mode uint32) output {
	sym := symbolicString(mode)
	return output{
		Octal:       octalString(mode),
		Symbolic:    sym,
		Owner:       sym[0:3],
		Group:       sym[3:6],
		Others:      sym[6:9],
		Special:     int((mode >> 9) & 0o7),
		SpecialText: specialText(mode),
		HasSUID:     mode&0o4000 != 0,
		HasSGID:     mode&0o2000 != 0,
		HasSticky:   mode&0o1000 != 0,
	}
}

// octalString 输出八进制：无特殊位 3 位，有特殊位 4 位。
func octalString(mode uint32) string {
	if mode&0o7000 != 0 {
		return fmt.Sprintf("%04o", mode)
	}
	return fmt.Sprintf("%03o", mode&0o777)
}

// symbolicString 生成 9 位符号串，特殊位在对应执行位显示 s/S/t/T。
func symbolicString(mode uint32) string {
	b := []byte("---------")
	perm := []byte{'r', 'w', 'x'}
	for i := 0; i < 9; i++ {
		if mode&(1<<(8-uint(i))) != 0 {
			b[i] = perm[i%3]
		}
	}
	if mode&0o4000 != 0 {
		b[2] = specialChar(b[2], 's', 'S')
	}
	if mode&0o2000 != 0 {
		b[5] = specialChar(b[5], 's', 'S')
	}
	if mode&0o1000 != 0 {
		b[8] = specialChar(b[8], 't', 'T')
	}
	return string(b)
}

// specialChar 将执行位替换为特殊位字符：有执行权限时小写，否则大写。
// 属主/属组执行位用 s/S，其他用户执行位用 t/T。
func specialChar(x byte, lower, upper byte) byte {
	switch x {
	case 'x':
		return lower
	case '-':
		return upper
	}
	return x
}

// specialText 生成特殊位中文说明。
func specialText(mode uint32) string {
	var parts []string
	if mode&0o4000 != 0 {
		parts = append(parts, "SUID（设置用户 ID，程序以属主身份运行）")
	}
	if mode&0o2000 != 0 {
		parts = append(parts, "SGID（设置组 ID，目录内新文件继承属组）")
	}
	if mode&0o1000 != 0 {
		parts = append(parts, "Sticky 粘滞位（仅属主/属组可删除目录内文件）")
	}
	if len(parts) == 0 {
		return "无特殊位"
	}
	return strings.Join(parts, "；")
}
