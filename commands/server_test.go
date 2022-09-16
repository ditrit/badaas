package commands

// This files holds the tests for the commands/server.go file.

import (
	"net/http"
	"testing"
	"time"
)

func Test_addrFromConf(t *testing.T) {
	expected := "192.168.236.222:25100"
	addr := addrFromConf("192.168.236.222", 25100)
	if addr != expected {
		t.Errorf("expected %s, got %s", expected, addr)
	}
}
func Test_createServer(t *testing.T) {
	handl := http.NewServeMux()
	timeout := time.Duration(time.Second)
	srv := createServer(
		handl,
		"localhost:8000",
		timeout, timeout,
	)
	if srv == nil {
		t.Error("createServer should not return a nil value")
	}
}