package transport

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// StreamableHTTPTransport implements Transport for MCP's Streamable HTTP protocol.
//
// This transport is the recommended HTTP transport per MCP spec 2025-03-26,
// replacing the deprecated SSE transport. It provides:
//   - Single endpoint handling POST/GET/DELETE methods
//   - Session management via Mcp-Session-Id header
//   - Bidirectional communication support
//   - Optional stateless mode for simpler deployments
//
// StreamableHTTPTransport is safe for concurrent use.
type StreamableHTTPTransport struct {
	Config StreamableConfig

	mu       sync.Mutex
	listener net.Listener
	server   *http.Server
}

// Name returns "streamable" as the transport identifier.
func (t *StreamableHTTPTransport) Name() string {
	return "streamable"
}

// Info returns runtime information about the transport.
// If the server is running, Info returns the actual bound address.
func (t *StreamableHTTPTransport) Info() Info {
	path := t.Config.Path
	if path == "" {
		path = "/mcp"
	}

	addr := ""
	t.mu.Lock()
	if t.listener != nil {
		addr = t.listener.Addr().String()
	}
	t.mu.Unlock()

	if addr == "" && t.Config.Port != 0 {
		host := t.Config.Host
		if host == "" {
			host = "0.0.0.0"
		}
		addr = fmt.Sprintf("%s:%d", host, t.Config.Port)
	}

	return Info{Name: "streamable", Addr: addr, Path: path}
}

// Serve starts the HTTP server and blocks until ctx is cancelled.
//
// When ctx is cancelled, Serve initiates graceful shutdown with a
// 5-second timeout for in-flight requests to complete.
func (t *StreamableHTTPTransport) Serve(ctx context.Context, server Server) error {
	host := t.Config.Host
	if host == "" {
		host = "0.0.0.0"
	}
	path := t.Config.Path
	if path == "" {
		path = "/mcp"
	}
	addr := fmt.Sprintf("%s:%d", host, t.Config.Port)

	readHeaderTimeout := t.Config.ReadHeaderTimeout
	if readHeaderTimeout == 0 {
		readHeaderTimeout = 10 * time.Second
	}

	// Create handler
	mux := http.NewServeMux()
	if handlerProvider, ok := server.(interface{ Handler() http.Handler }); ok {
		mux.Handle(path, handlerProvider.Handler())
	} else {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.WriteHeader(http.StatusOK)
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Configure TLS if enabled
	if t.Config.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(t.Config.TLS.CertFile, t.Config.TLS.KeyFile)
		if err != nil {
			return fmt.Errorf("load TLS certificate: %w", err)
		}
		httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	t.mu.Lock()
	t.listener = ln
	t.server = httpServer
	t.mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		var serveErr error
		if t.Config.TLS.Enabled {
			serveErr = httpServer.ServeTLS(ln, "", "")
		} else {
			serveErr = httpServer.Serve(ln)
		}
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		_ = t.Close()
		return nil
	case err := <-errCh:
		return err
	}
}

// Close gracefully shuts down the HTTP server with a 5-second timeout.
func (t *StreamableHTTPTransport) Close() error {
	t.mu.Lock()
	srv := t.server
	ln := t.listener
	t.mu.Unlock()

	if srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if ln != nil {
		_ = ln.Close()
	}
	return err
}
