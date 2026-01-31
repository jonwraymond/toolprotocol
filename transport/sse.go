package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// SSETransport implements Transport using Server-Sent Events over HTTP.
//
// This is the legacy HTTP transport. For new implementations, prefer
// StreamableHTTPTransport per MCP spec 2025-03-26.
//
// SSETransport is safe for concurrent use.
type SSETransport struct {
	Config SSEConfig

	mu       sync.Mutex
	listener net.Listener
	server   *http.Server
}

// Name returns "sse" as the transport identifier.
func (t *SSETransport) Name() string {
	return "sse"
}

// Info returns runtime information about the transport.
// If the server is running, Info returns the actual bound address.
func (t *SSETransport) Info() Info {
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

	return Info{Name: "sse", Addr: addr, Path: path}
}

// Serve starts the SSE HTTP server and blocks until ctx is cancelled.
func (t *SSETransport) Serve(ctx context.Context, server Server) error {
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

	// Create handler - for standalone transport, use a simple handler
	mux := http.NewServeMux()
	if handlerProvider, ok := server.(interface{ Handler() http.Handler }); ok {
		mux.Handle(path, handlerProvider.Handler())
	} else {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
		})
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
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
		err := httpServer.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
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
func (t *SSETransport) Close() error {
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
