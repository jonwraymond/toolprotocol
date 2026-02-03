package a2a

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jonwraymond/toolprotocol/task"
	"github.com/jonwraymond/toolprotocol/wire"
)

// Logger is the interface for logging.
//
// Contract:
// - Concurrency: implementations must be safe for concurrent use.
// - Errors: logging must be best-effort and must not panic.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Agent defines the callbacks needed to serve A2A requests.
type Agent interface {
	// AgentCard returns the A2A agent card payload.
	AgentCard(ctx context.Context) (any, error)

	// ListSkills returns the skills the agent supports.
	ListSkills(ctx context.Context) ([]wire.Tool, error)

	// Invoke runs a skill with the given arguments.
	Invoke(ctx context.Context, skillID string, args map[string]any) (InvokeResult, error)
}

// InvokeResult captures invocation output.
type InvokeResult struct {
	Content []wire.Content
	Meta    map[string]any
}

// Handler serves A2A JSON-RPC and REST endpoints.
type Handler struct {
	Agent  Agent
	Tasks  task.Manager
	Wire   *wire.A2AWire
	Logger Logger
}

// NewHandler creates a new A2A handler.
func NewHandler(agent Agent, tasks task.Manager) *Handler {
	if tasks == nil {
		tasks = task.NewManager()
	}
	return &Handler{
		Agent: agent,
		Tasks: tasks,
		Wire:  wire.NewA2A(),
	}
}

// ServeRPC handles A2A JSON-RPC requests.
func (h *Handler) ServeRPC(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := readBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	req, err := h.Wire.DecodeRequest(ctx, data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var resp *wire.Response
	switch strings.ToLower(req.Method) {
	case "agent/invoke":
		resp = h.handleInvoke(ctx, req)
	case "agent/status", "task/status":
		resp = h.handleStatus(ctx, req)
	default:
		resp = errorResponse(req.ID, fmt.Errorf("unsupported method %q", req.Method))
	}

	out, err := h.Wire.EncodeResponse(ctx, resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, out)
}

// ServeAgentCard returns the A2A agent card document.
func (h *Handler) ServeAgentCard(w http.ResponseWriter, r *http.Request) {
	if h.Agent == nil {
		writeError(w, http.StatusNotImplemented, errors.New("agent card not supported"))
		return
	}
	card, err := h.Agent.AgentCard(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSONValue(w, card)
}

// ServeSkills returns the list of skills in A2A format.
func (h *Handler) ServeSkills(w http.ResponseWriter, r *http.Request) {
	if h.Agent == nil {
		writeError(w, http.StatusNotImplemented, errors.New("skills not supported"))
		return
	}
	tools, err := h.Agent.ListSkills(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	data, err := h.Wire.EncodeToolList(r.Context(), tools)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, data)
}

// ServeTask returns task metadata as JSON.
func (h *Handler) ServeTask(w http.ResponseWriter, r *http.Request, taskID string) {
	if taskID == "" {
		writeError(w, http.StatusBadRequest, errors.New("task id required"))
		return
	}
	t, err := h.Tasks.Get(r.Context(), taskID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSONValue(w, t)
}

// ServeTaskList returns the list of tasks.
func (h *Handler) ServeTaskList(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Tasks.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSONValue(w, tasks)
}

// ServeTaskEvents streams task updates using SSE.
func (h *Handler) ServeTaskEvents(w http.ResponseWriter, r *http.Request, taskID string) {
	if taskID == "" {
		writeError(w, http.StatusBadRequest, errors.New("task id required"))
		return
	}
	ch, err := h.Tasks.Subscribe(r.Context(), taskID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, errors.New("streaming not supported"))
		return
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, "event: heartbeat\ndata: {}\n\n")
			flusher.Flush()
		case taskUpdate, ok := <-ch:
			if !ok {
				return
			}
			payload, _ := json.Marshal(taskUpdate)
			fmt.Fprintf(w, "event: task\nid: %s\ndata: %s\n\n", taskUpdate.ID, payload)
			flusher.Flush()
		}
	}
}

func (h *Handler) handleInvoke(ctx context.Context, req *wire.Request) *wire.Response {
	if h.Agent == nil {
		return errorResponse(req.ID, errors.New("agent not configured"))
	}
	taskID := req.ID
	if taskID == "" {
		taskID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	}
	if _, err := h.Tasks.Create(ctx, taskID); err != nil && !errors.Is(err, task.ErrTaskExists) {
		return errorResponse(req.ID, err)
	}

	go func() {
		_ = h.Tasks.Update(context.Background(), taskID, 0.1, "running")
		result, err := h.Agent.Invoke(context.Background(), req.ToolID, req.Arguments)
		if err != nil {
			_ = h.Tasks.Fail(context.Background(), taskID, err)
			return
		}
		_ = h.Tasks.Complete(context.Background(), taskID, result)
	}()

	return &wire.Response{
		ID: req.ID,
		Meta: map[string]any{
			"status": map[string]any{
				"state": "running",
				"id":    taskID,
			},
			"taskId": taskID,
		},
	}
}

func (h *Handler) handleStatus(ctx context.Context, req *wire.Request) *wire.Response {
	taskID := req.ID
	if taskID == "" {
		if v, ok := req.Arguments["id"].(string); ok {
			taskID = v
		}
	}
	if taskID == "" {
		return errorResponse(req.ID, errors.New("task id required"))
	}
	t, err := h.Tasks.Get(ctx, taskID)
	if err != nil {
		return errorResponse(req.ID, err)
	}
	return &wire.Response{
		ID: req.ID,
		Meta: map[string]any{
			"status": map[string]any{
				"state": string(t.State),
				"id":    t.ID,
			},
		},
	}
}

func errorResponse(id string, err error) *wire.Response {
	return &wire.Response{
		ID:      id,
		IsError: true,
		Error: &wire.Error{
			Code:    -32000,
			Message: err.Error(),
		},
	}
}

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func writeJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func writeJSONValue(w http.ResponseWriter, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	payload := map[string]any{"error": err.Error()}
	data, _ := json.Marshal(payload)
	_, _ = w.Write(data)
}
