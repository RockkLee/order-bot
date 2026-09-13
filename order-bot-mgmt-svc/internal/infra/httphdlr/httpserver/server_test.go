package httpserver

import "testing"

func TestNewHTTPServerUsesProvidedAddr(t *testing.T) {
	app := NewServerContainer(8080, nil, nil)

	srv := NewHTTPServer(app, "debug", "127.0.0.1:8080")

	if srv == nil {
		t.Fatal("NewHTTPServer() returned nil")
	}
	if srv.Addr != "127.0.0.1:8080" {
		t.Fatalf("srv.Addr = %q, want %q", srv.Addr, "127.0.0.1:8080")
	}
	if srv.Handler == nil {
		t.Fatal("srv.Handler is nil")
	}
}
