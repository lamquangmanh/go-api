package utils_test

import (
	"net"
	"testing"

	"go-api/pkg/utils"
)

func TestIsLANInterface(t *testing.T) {
	cases := map[string]bool{
		"eth0":                  true,
		"en0":                   true,
		"Ethernet0":             true,
		"Local Area Connection": true,
		"lo":                    false,
	}
	for name, expected := range cases {
		if got := utils.IsLANInterface(name); got != expected {
			t.Fatalf("IsLANInterface(%q)=%v, expected %v", name, got, expected)
		}
	}
}

func TestIsPrivateIP(t *testing.T) {
	if !utils.IsPrivateIP(net.ParseIP("10.1.2.3")) {
		t.Fatalf("expected 10.x.x.x to be private")
	}
	if !utils.IsPrivateIP(net.ParseIP("172.16.1.1")) {
		t.Fatalf("expected 172.16.x.x to be private")
	}
	if !utils.IsPrivateIP(net.ParseIP("192.168.1.1")) {
		t.Fatalf("expected 192.168.x.x to be private")
	}
	if utils.IsPrivateIP(net.ParseIP("8.8.8.8")) {
		t.Fatalf("expected 8.8.8.8 to be public")
	}
}

func TestGetLocalIPv4AndGetAllIPs(t *testing.T) {
	ips := utils.GetAllIPs()
	if ips == nil {
		t.Fatalf("expected non-nil slice")
	}

	ip, err := utils.GetLocalIPv4()
	if err == nil && ip == "" {
		t.Fatalf("expected non-empty ip when no error")
	}
}
