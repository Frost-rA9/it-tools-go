// Package dockerruncompose 实现将 docker run 命令转换为 docker-compose YAML。
package dockerruncompose

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"it-tools-go/internal/registry"
	"gopkg.in/yaml.v3"
)

// 工具元数据。
const (
	ID          = "docker-run-to-docker-compose-converter"
	Name        = "Docker Run 转 Compose"
	Description = "将 docker run 命令转换为 docker-compose YAML"
	Category    = "开发"
	Icon        = "BrandDocker"
)

// Keywords 为搜索关键词（Go 常量不能是 slice，故用 var）。
var Keywords = []string{"docker", "run", "compose", "yaml", "yml", "convert", "转换", "container"}

// Tool 返回工具的完整注册元数据。
func Tool() registry.Tool {
	return registry.Tool{ID: ID, Name: Name, Description: Description, Category: Category, Keywords: Keywords, Icon: Icon}
}

// input 是工具的输入结构。
type input struct {
	Text string `json:"text"`
}

// output 是工具的输出结构。
type output struct {
	Compose  string   `json:"compose"`
	Warnings []string `json:"warnings"`
}

// Executor 实现 registry.Executor 接口。
type Executor struct{}

// Execute 处理输入并返回结果 JSON。
func (Executor) Execute(_ context.Context, inputJSON string) (string, error) {
	var in input
	if err := json.Unmarshal([]byte(inputJSON), &in); err != nil {
		return "", fmt.Errorf("解析输入失败: %w", err)
	}

	compose, warnings, err := Convert(in.Text)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(output{Compose: compose, Warnings: warnings})
	if err != nil {
		return "", fmt.Errorf("序列化输出失败: %w", err)
	}
	return string(out), nil
}

// Convert 将 docker run 命令转换为 docker-compose YAML。
func Convert(cmd string) (string, []string, error) {
	if strings.TrimSpace(cmd) == "" {
		return "", nil, nil
	}

	tokens, err := tokenize(cmd)
	if err != nil {
		return "", nil, err
	}

	svc, warnings, err := parse(tokens)
	if err != nil {
		return "", nil, err
	}
	if svc.Image == "" && len(svc.Command) == 0 {
		return "", nil, fmt.Errorf("未找到 Docker 镜像或命令")
	}
	if svc.Name == "" {
		svc.Name = sanitizeName(svc.Image)
	}
	if svc.Name == "" {
		return "", nil, fmt.Errorf("无法从镜像推导服务名，请使用 --name 指定")
	}

	cf := composeFile{
		Services: map[string]composeService{svc.Name: svc.toCompose()},
	}
	if svc.Network != "" {
		cf.Networks = map[string]interface{}{svc.Network: nil}
	}

	b, err := yaml.Marshal(cf)
	if err != nil {
		return "", nil, fmt.Errorf("序列化 YAML 失败: %w", err)
	}
	return strings.TrimSpace(string(b)), warnings, nil
}

// tokenize 按 shell 风格分词：支持单/双引号、反斜杠转义、连续空白折叠。
func tokenize(s string) ([]string, error) {
	var tokens []string
	var cur strings.Builder
	inSingle, inDouble := false, false
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			} else {
				cur.WriteByte(c)
			}
		case inDouble:
			if c == '"' {
				inDouble = false
			} else if c == '\\' && i+1 < len(s) && (s[i+1] == '"' || s[i+1] == '\\') {
				cur.WriteByte(s[i+1])
				i++
			} else {
				cur.WriteByte(c)
			}
		default:
			switch c {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '\\':
				if i+1 < len(s) {
					cur.WriteByte(s[i+1])
					i++
				}
			case ' ', '\t', '\n', '\r':
				if cur.Len() > 0 {
					tokens = append(tokens, cur.String())
					cur.Reset()
				}
			default:
				cur.WriteByte(c)
			}
		}
		i++
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("未闭合的引号")
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens, nil
}

// service 是解析过程中的中间表示。
type service struct {
	Name       string
	Image      string
	Ports      []string
	Volumes    []string
	Env        []string
	Restart    string
	Network    string
	EnvFile    []string
	ExtraHosts []string
	Privileged bool
	Entrypoint string
	Workdir    string
	User       string
	DNS        []string
	Command    []string
	LogDriver  string
	LogOpts    map[string]string
}

// toCompose 转换为 compose YAML 结构。
func (s service) toCompose() composeService {
	svc := composeService{
		Image:       s.Image,
		Ports:       s.Ports,
		Volumes:     s.Volumes,
		Environment: s.Env,
		Restart:     s.Restart,
		Networks:    optList(s.Network),
		EnvFile:     s.EnvFile,
		ExtraHosts:  s.ExtraHosts,
		Privileged:  s.Privileged,
		Entrypoint:  s.Entrypoint,
		WorkingDir:  s.Workdir,
		User:        s.User,
		DNS:         s.DNS,
		Command:     s.Command,
	}
	if s.LogDriver != "" || len(s.LogOpts) > 0 {
		svc.Logging = &composeLogging{Driver: s.LogDriver, Options: s.LogOpts}
	}
	return svc
}

func optList(v string) []string {
	if v == "" {
		return nil
	}
	return []string{v}
}

// optDef 描述一个选项：needsValue 表示是否消费下一个 token；apply 为 nil 表示忽略并提示。
type optDef struct {
	needsValue bool
	apply      func(s *service, val string, warnings *[]string)
}

// 长选项表（--xxx）。
var longOpts = map[string]optDef{
	"--publish":    {true, addList("ports")},
	"--env":        {true, addEnv},
	"--volume":     {true, addVolume},
	"--name":       {true, setField("Name")},
	"--restart":    {true, setField("Restart")},
	"--network":    {true, setField("Network")},
	"--env-file":   {true, addList("EnvFile")},
	"--add-host":   {true, addList("ExtraHosts")},
	"--entrypoint": {true, setField("Entrypoint")},
	"--workdir":    {true, setField("Workdir")},
	"--user":       {true, setField("User")},
	"--dns":        {true, addList("DNS")},
	"--log-driver": {true, setField("LogDriver")},
	"--log-opt":    {true, addLogOpt},
	"--label":      {true, nil},
	"--pull":       {true, nil},
	"--mount":      {true, nil},
	"--cpus":       {true, nil},
	"--gpus":       {true, nil},
	"--memory":     {true, nil},
	"--shm-size":   {true, nil},
	"--ulimit":     {true, nil},
	"--cap-add":    {true, nil},
	"--cap-drop":   {true, nil},
	"--security-opt": {true, nil},
	"--device":     {true, nil},
	"--tmpfs":      {true, nil},
	"--expose":     {true, nil},
	"--link":       {true, nil},
	"--health-cmd": {true, nil},
	"--health-interval": {true, nil},
	"--health-retries":  {true, nil},
	"--health-timeout":  {true, nil},
	"--health-start-period": {true, nil},
	"--stop-timeout": {true, nil},
	"--stop-signal": {true, nil},
	"--init-path":  {true, nil},
	"--privileged": {false, func(s *service, _ string, _ *[]string) { s.Privileged = true }},
}

// 常见无副作用布尔选项：忽略且不提示。
var silentBoolOpts = map[string]bool{
	"-d": true, "--detach": true, "-i": true, "--interactive": true,
	"-t": true, "--tty": true, "--rm": true, "--init": true,
	"--read-only": true, "--sig-proxy": true, "--no-healthcheck": true,
	"-q": true, "--quiet": true,
}

// 短选项表（-x）。
var shortOpts = map[string]optDef{
	"-p": {true, addPort},
	"-e": {true, addEnv},
	"-v": {true, addVolume},
	"-w": {true, setField("Workdir")},
	"-u": {true, setField("User")},
}

// parse 解析 token 序列。
func parse(tokens []string) (service, []string, error) {
	var svc service
	var warnings []string

	i := 0
	// 跳过前导 docker / run 子命令。
	if i < len(tokens) && tokens[i] == "docker" {
		i++
	}
	if i < len(tokens) && tokens[i] == "run" {
		i++
	}

	for i < len(tokens) {
		tok := tokens[i]

		// 镜像已确定：后续 token 全部为 command/arg（含选项形式，对齐 docker 语义）。
		if svc.Image != "" {
			svc.Command = append(svc.Command, tokens[i:]...)
			break
		}

		if tok == "--" {
			svc.Command = append(svc.Command, tokens[i+1:]...)
			break
		}

		if strings.HasPrefix(tok, "--") {
			name, val, hasVal := tok, "", false
			if eq := strings.Index(tok, "="); eq >= 0 {
				name, val, hasVal = tok[:eq], tok[eq+1:], true
			}
			opt, known := longOpts[name]
			if !known {
				// 布尔或未知长选项
				if silentBoolOpts[name] {
					i++
					continue
				}
				warnings = append(warnings, fmt.Sprintf("选项 %s 未转换", name))
				i++
				continue
			}
			if opt.needsValue {
				if !hasVal {
					if i+1 >= len(tokens) {
						return svc, warnings, fmt.Errorf("选项 %s 缺少参数", name)
					}
					i++
					val = tokens[i]
				}
				if opt.apply != nil {
					opt.apply(&svc, val, &warnings)
				} else {
					warnings = append(warnings, fmt.Sprintf("选项 %s 未转换", name))
				}
				i++
				continue
			}
			if opt.apply != nil {
				opt.apply(&svc, "", &warnings)
			}
			i++
			continue
		}

		if strings.HasPrefix(tok, "-") && len(tok) > 1 && tok != "-" {
			// 短选项，支持组合（如 -it）。
			short := tok[1:]
			for j := 0; j < len(short); j++ {
				opt := "-" + string(short[j])
				if o, ok := shortOpts[opt]; ok && o.needsValue {
					val := short[j+1:]
					if val == "" {
						if i+1 >= len(tokens) {
							return svc, warnings, fmt.Errorf("选项 %s 缺少参数", opt)
						}
						i++
						val = tokens[i]
					}
					o.apply(&svc, val, &warnings)
					break
				}
				if silentBoolOpts[opt] {
					continue
				}
				warnings = append(warnings, fmt.Sprintf("选项 %s 未转换", opt))
			}
			i++
			continue
		}

		// 非选项：第一个为镜像，其余为 command。
		if svc.Image == "" {
			svc.Image = tok
		} else {
			svc.Command = append(svc.Command, tok)
		}
		i++
	}

	return svc, warnings, nil
}

// sanitizeName 从镜像名推导服务名：去 tag、去注册表前缀、仅保留合法字符。
func sanitizeName(image string) string {
	name := image
	// 去 tag。
	if idx := strings.LastIndex(name, ":"); idx > 0 && !strings.Contains(name[idx+1:], "/") {
		name = name[:idx]
	}
	// 去注册表前缀（最后一段）。
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	var sb strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '.', r == '-':
			sb.WriteRune(r)
		default:
			sb.WriteRune('-')
		}
	}
	return sb.String()
}

// addPort 处理 -p/--publish。
func addPort(s *service, val string, _ *[]string) {
	s.Ports = append(s.Ports, val)
}

// addEnv 处理 -e/--env：KEY=VALUE 保序；无等号的 KEY 记 KEY=。
func addEnv(s *service, val string, _ *[]string) {
	if !strings.Contains(val, "=") {
		val += "="
	}
	s.Env = append(s.Env, val)
}

// addVolume 处理 -v/--volume。
func addVolume(s *service, val string, _ *[]string) {
	s.Volumes = append(s.Volumes, val)
}

// addList 生成向指定字段追加的 handler。
func addList(field string) func(s *service, val string, _ *[]string) {
	return func(s *service, val string, _ *[]string) {
		switch field {
		case "EnvFile":
			s.EnvFile = append(s.EnvFile, val)
		case "ExtraHosts":
			s.ExtraHosts = append(s.ExtraHosts, val)
		case "DNS":
			s.DNS = append(s.DNS, val)
		case "ports":
			s.Ports = append(s.Ports, val)
		}
	}
}

// setField 生成设置指定字符串字段的 handler。
func setField(field string) func(s *service, val string, _ *[]string) {
	return func(s *service, val string, _ *[]string) {
		switch field {
		case "Name":
			s.Name = val
		case "Restart":
			s.Restart = val
		case "Network":
			s.Network = val
		case "Entrypoint":
			s.Entrypoint = val
		case "Workdir":
			s.Workdir = val
		case "User":
			s.User = val
		case "LogDriver":
			s.LogDriver = val
		}
	}
}

// addLogOpt 处理 --log-opt KEY=VALUE，累加到日志选项映射。
func addLogOpt(s *service, val string, _ *[]string) {
	if s.LogOpts == nil {
		s.LogOpts = make(map[string]string)
	}
	k, v, ok := strings.Cut(val, "=")
	if !ok {
		k, v = val, ""
	}
	s.LogOpts[k] = v
}

// composeLogging 是 compose 的日志配置。
type composeLogging struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

// composeService 是 compose 中单个服务的 YAML 结构。
type composeService struct {
	Image       string                 `yaml:"image,omitempty"`
	Ports       []string               `yaml:"ports,omitempty"`
	Volumes     []string               `yaml:"volumes,omitempty"`
	Environment []string               `yaml:"environment,omitempty"`
	Restart     string                 `yaml:"restart,omitempty"`
	Networks    []string               `yaml:"networks,omitempty"`
	EnvFile     []string               `yaml:"env_file,omitempty"`
	ExtraHosts  []string               `yaml:"extra_hosts,omitempty"`
	Privileged  bool                   `yaml:"privileged,omitempty"`
	Entrypoint  string                 `yaml:"entrypoint,omitempty"`
	WorkingDir  string                 `yaml:"working_dir,omitempty"`
	User        string                 `yaml:"user,omitempty"`
	DNS         []string               `yaml:"dns,omitempty"`
	Command     []string               `yaml:"command,omitempty"`
	Logging     *composeLogging        `yaml:"logging,omitempty"`
}

// composeFile 是 compose 顶层结构。
type composeFile struct {
	Services map[string]composeService `yaml:"services"`
	Networks map[string]interface{}    `yaml:"networks,omitempty"`
}
