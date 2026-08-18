package wifiqrcodegen

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWifiPayload(t *testing.T) {
	tests := []struct {
		name string
		in   input
		want string
	}{
		{"nopass", input{SSID: "Cafe", Encryption: "nopass"}, "WIFI:S:Cafe;;"},
		{"WPA", input{SSID: "My;Wifi", Password: `p:a,s\\`, Encryption: "WPA", Hidden: true}, `WIFI:S:My\;Wifi;T:WPA;P:p\:a\,s\\\\;H:true;;`},
		{"WEP", input{SSID: "Home", Password: "secret", Encryption: "WEP"}, "WIFI:S:Home;T:WEP;P:secret;;"},
		{"WPA2-EAP", input{SSID: "Corp", Password: "secret", Encryption: "WPA2-EAP", EAPMethod: "PEAP", Phase2: "MSCHAPV2", Identity: "alice"}, "WIFI:S:Corp;T:WPA2-EAP;P:secret;E:PEAP;PH2:MSCHAPV2;I:alice;;"},
		{"anonymous EAP", input{SSID: "Corp", Password: "secret", Encryption: "WPA2-EAP", EAPMethod: "TLS", Anonymous: true}, "WIFI:S:Corp;T:WPA2-EAP;P:secret;E:TLS;A:anon;;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := wifiPayload(tt.in)
			if !ok || got != tt.want {
				t.Errorf("payload = %q, %v, want %q", got, ok, tt.want)
			}
		})
	}
}

func TestWifiInvalid(t *testing.T) {
	tests := []input{
		{Encryption: "WPA", Password: "secret"},
		{SSID: "Wifi", Encryption: "WPA"},
		{SSID: "Corp", Encryption: "WPA2-EAP", Password: "secret", EAPMethod: "PEAP", Identity: "alice"},
	}
	for _, in := range tests {
		if payload, ok := wifiPayload(in); ok || payload != "" {
			t.Errorf("无效输入应无 payload: %q, %v", payload, ok)
		}
	}
}

func TestExecute(t *testing.T) {
	var e Executor
	out, err := e.Execute(t.Context(), `{"ssid":"Cafe","encryption":"nopass","foreground":"#000000","background":"#ffffff"}`)
	if err != nil {
		t.Fatalf("Execute 意外错误: %v", err)
	}
	var got output
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("解析输出失败: %v", err)
	}
	if got.Payload != "WIFI:S:Cafe;;" || !strings.HasPrefix(got.DataURL, "data:image/png;base64,") {
		t.Errorf("输出不符: %+v", got)
	}
}
