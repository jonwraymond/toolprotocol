package transport

import (
	"testing"
	"time"
)

func TestHTTPConfig_Defaults(t *testing.T) {
	cfg := HTTPConfig{}

	// Verify zero values
	if cfg.Host != "" {
		t.Errorf("default Host = %q, want empty", cfg.Host)
	}
	if cfg.Port != 0 {
		t.Errorf("default Port = %d, want 0", cfg.Port)
	}
	if cfg.Path != "" {
		t.Errorf("default Path = %q, want empty", cfg.Path)
	}
	if cfg.ReadHeaderTimeout != 0 {
		t.Errorf("default ReadHeaderTimeout = %v, want 0", cfg.ReadHeaderTimeout)
	}
}

func TestHTTPConfig_CustomValues(t *testing.T) {
	cfg := HTTPConfig{
		Host:              "127.0.0.1",
		Port:              8080,
		Path:              "/api/mcp",
		ReadHeaderTimeout: 30 * time.Second,
	}

	if cfg.Host != "127.0.0.1" {
		t.Errorf("Host = %q, want %q", cfg.Host, "127.0.0.1")
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want %d", cfg.Port, 8080)
	}
	if cfg.Path != "/api/mcp" {
		t.Errorf("Path = %q, want %q", cfg.Path, "/api/mcp")
	}
	if cfg.ReadHeaderTimeout != 30*time.Second {
		t.Errorf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 30*time.Second)
	}
}

func TestTLSConfig_Defaults(t *testing.T) {
	cfg := TLSConfig{}

	if cfg.Enabled {
		t.Error("default Enabled = true, want false")
	}
	if cfg.CertFile != "" {
		t.Errorf("default CertFile = %q, want empty", cfg.CertFile)
	}
	if cfg.KeyFile != "" {
		t.Errorf("default KeyFile = %q, want empty", cfg.KeyFile)
	}
}

func TestTLSConfig_Enabled(t *testing.T) {
	cfg := TLSConfig{
		Enabled:  true,
		CertFile: "/path/to/cert.pem",
		KeyFile:  "/path/to/key.pem",
	}

	if !cfg.Enabled {
		t.Error("Enabled = false, want true")
	}
	if cfg.CertFile != "/path/to/cert.pem" {
		t.Errorf("CertFile = %q, want %q", cfg.CertFile, "/path/to/cert.pem")
	}
	if cfg.KeyFile != "/path/to/key.pem" {
		t.Errorf("KeyFile = %q, want %q", cfg.KeyFile, "/path/to/key.pem")
	}
}

func TestStreamableConfig_Defaults(t *testing.T) {
	cfg := StreamableConfig{}

	if cfg.Stateless {
		t.Error("default Stateless = true, want false")
	}
	if cfg.JSONResponse {
		t.Error("default JSONResponse = true, want false")
	}
	if cfg.SessionTimeout != 0 {
		t.Errorf("default SessionTimeout = %v, want 0", cfg.SessionTimeout)
	}
}

func TestStreamableConfig_AllOptions(t *testing.T) {
	cfg := StreamableConfig{
		HTTPConfig: HTTPConfig{
			Host: "localhost",
			Port: 9000,
			Path: "/mcp",
		},
		TLS: TLSConfig{
			Enabled:  true,
			CertFile: "/cert.pem",
			KeyFile:  "/key.pem",
		},
		Stateless:      true,
		JSONResponse:   true,
		SessionTimeout: 5 * time.Minute,
	}

	if cfg.Host != "localhost" {
		t.Errorf("Host = %q, want %q", cfg.Host, "localhost")
	}
	if cfg.Port != 9000 {
		t.Errorf("Port = %d, want %d", cfg.Port, 9000)
	}
	if !cfg.TLS.Enabled {
		t.Error("TLS.Enabled = false, want true")
	}
	if !cfg.Stateless {
		t.Error("Stateless = false, want true")
	}
	if !cfg.JSONResponse {
		t.Error("JSONResponse = false, want true")
	}
	if cfg.SessionTimeout != 5*time.Minute {
		t.Errorf("SessionTimeout = %v, want %v", cfg.SessionTimeout, 5*time.Minute)
	}
}

func TestSSEConfig_Defaults(t *testing.T) {
	cfg := SSEConfig{}

	if cfg.Host != "" {
		t.Errorf("default Host = %q, want empty", cfg.Host)
	}
	if cfg.Port != 0 {
		t.Errorf("default Port = %d, want 0", cfg.Port)
	}
}

func TestSSEConfig_CustomValues(t *testing.T) {
	cfg := SSEConfig{
		HTTPConfig: HTTPConfig{
			Host:              "0.0.0.0",
			Port:              3000,
			Path:              "/events",
			ReadHeaderTimeout: 15 * time.Second,
		},
	}

	if cfg.Host != "0.0.0.0" {
		t.Errorf("Host = %q, want %q", cfg.Host, "0.0.0.0")
	}
	if cfg.Port != 3000 {
		t.Errorf("Port = %d, want %d", cfg.Port, 3000)
	}
	if cfg.Path != "/events" {
		t.Errorf("Path = %q, want %q", cfg.Path, "/events")
	}
}
