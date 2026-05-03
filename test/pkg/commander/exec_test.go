package commander_test

import (
	"net"
	"testing"

	"go-api/pkg/commander"
)

func TestExecCommand(t *testing.T) {
	out, err := commander.ExecCommand("sh", "-c", "printf hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Fatalf("unexpected output: %q", out)
	}

	_, err = commander.ExecCommand("sh", "-c", "echo boom 1>&2; exit 1")
	if err == nil {
		t.Fatalf("expected command failure")
	}
}

func TestCheckCommandExists(t *testing.T) {
	if err := commander.CheckCommandExists("sh"); err != nil {
		t.Fatalf("expected sh to exist: %v", err)
	}
	if err := commander.CheckCommandExists("__definitely_missing_cmd__"); err == nil {
		t.Fatalf("expected missing command error")
	}
}

func TestCheckPortAvailable(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to allocate test port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	if err := commander.CheckPortAvailable(port); err == nil {
		t.Fatalf("expected occupied port to be unavailable")
	}

	_ = ln.Close()
	if err := commander.CheckPortAvailable(port); err != nil {
		t.Fatalf("expected port to be available after close: %v", err)
	}
}

func TestCheckDiskExists(t *testing.T) {
	if err := commander.CheckDiskExists("/dev/this-disk-does-not-exist"); err == nil {
		t.Fatalf("expected invalid disk path error")
	}
}
