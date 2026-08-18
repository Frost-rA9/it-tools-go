package dockerruncompose

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name       string
		cmd        string
		wantSub    []string // 期望包含的子串
		wantWarn   []string // 期望包含的警告子串
		wantErr    bool
	}{
		{
			name: "完整示例",
			cmd:  "docker run -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro --restart always --log-opt max-size=1g nginx",
			wantSub: []string{
				"services:", "nginx:", "image: nginx",
				"ports:", "- 80:80",
				"volumes:", "- /var/run/docker.sock:/tmp/docker.sock:ro",
				"restart: always",
				"logging:", "options:", "max-size: 1g",
			},
		},
		{
			name: "日志驱动与选项",
			cmd:  "docker run --log-driver json-file --log-opt max-size=1g --log-opt max-file=3 nginx",
			wantSub: []string{
				"logging:", "driver: json-file", "options:", "max-size: 1g", `max-file: "3"`,
			},
		},
		{
			name: "指定服务名",
			cmd:  "docker run --name myapp -p 8080:80 nginx:latest",
			wantSub: []string{
				"myapp:", "image: nginx:latest", "- 8080:80",
			},
		},
		{
			name:    "等号形式",
			cmd:     "docker run --name=web --restart=on-failure -e FOO=bar alpine",
			wantSub: []string{"web:", "restart: on-failure", "FOO=bar"},
		},
		{
			name: "环境变量与引号",
			cmd:  `docker run -e "FOO=bar baz" -e EMPTY -e K=V redis`,
			wantSub: []string{
				"FOO=bar baz", "EMPTY=", "K=V",
			},
		},
		{
			name: "网络与主机映射",
			cmd:  "docker run --network mynet --add-host db:127.0.0.1 --privileged postgres",
			wantSub: []string{
				"networks:", "- mynet", "extra_hosts:", "- db:127.0.0.1", "privileged: true",
			},
		},
		{
			name: "命令参数",
			cmd:  "docker run --entrypoint /bin/sh alpine -c 'echo hi'",
			wantSub: []string{
				"entrypoint: /bin/sh", "command:", "- -c", "- echo hi",
			},
		},
		{
			name: "无 docker 前缀",
			cmd:  "run -d -p 3000:3000 node:18",
			wantSub: []string{"node:", "- 3000:3000"},
		},
		{
			name: "组合短选项",
			cmd:  "docker run -it --rm -p 3000:3000 node",
			wantSub: []string{"node:", "- 3000:3000"},
		},
		{
			name:    "未知选项警告",
			cmd:     "docker run --cpus 2 --gpus all python:3.11",
			wantSub: []string{"python:", "image: python:3.11"},
			wantWarn: []string{"--cpus", "--gpus"},
		},
		{
			name: "空输入",
			cmd:  "",
			wantSub: []string{""},
		},
		{
			name:    "缺少参数",
			cmd:     "docker run -p",
			wantErr: true,
		},
		{
			name:    "未闭合引号",
			cmd:     "docker run -e 'FOO=bar nginx",
			wantErr: true,
		},
		{
			name:    "无镜像无命令",
			cmd:     "docker run --rm",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, warnings, err := Convert(tt.cmd)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望出错，但成功（%q / %v）", got, warnings)
				}
				return
			}
			if err != nil {
				t.Fatalf("返回错误: %v", err)
			}
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("输出缺少 %q：\n%s", sub, got)
				}
			}
			gotWarn := strings.Join(warnings, "\n")
			for _, w := range tt.wantWarn {
				if !strings.Contains(gotWarn, w) {
					t.Errorf("警告缺少 %q：%v", w, warnings)
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	e := Executor{}

	in, _ := json.Marshal(input{Text: "docker run -p 80:80 nginx"})
	outJSON, err := e.Execute(context.Background(), string(in))
	if err != nil {
		t.Fatalf("Execute 返回错误: %v", err)
	}
	var out output
	if err := json.Unmarshal([]byte(outJSON), &out); err != nil {
		t.Fatalf("反序列化输出失败: %v", err)
	}
	if !strings.Contains(out.Compose, "image: nginx") || !strings.Contains(out.Compose, "- 80:80") {
		t.Errorf("结果 = %q", out.Compose)
	}
}
