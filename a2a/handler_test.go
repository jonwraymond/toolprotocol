package a2a

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jonwraymond/toolprotocol/task"
	"github.com/jonwraymond/toolprotocol/wire"
)

type fakeAgent struct {
	card   map[string]any
	skills []wire.Tool
}

func (f fakeAgent) AgentCard(ctx context.Context) (any, error) { return f.card, nil }
func (f fakeAgent) ListSkills(ctx context.Context) ([]wire.Tool, error) {
	return f.skills, nil
}
func (f fakeAgent) Invoke(ctx context.Context, skillID string, args map[string]any) (InvokeResult, error) {
	return InvokeResult{
		Content: []wire.Content{{Type: wire.ContentTypeText, Text: "ok"}},
		Meta:    map[string]any{"skillId": skillID},
	}, nil
}

func TestHandler_ServeAgentCardAndSkills(t *testing.T) {
	agent := fakeAgent{
		card: map[string]any{"name": "test-agent"},
		skills: []wire.Tool{
			{Name: "echo", Description: "Echo text"},
		},
	}
	h := NewHandler(agent, task.NewManager())

	cardReq := httptest.NewRequest(http.MethodGet, "/a2a/agent-card", nil)
	cardRec := httptest.NewRecorder()
	h.ServeAgentCard(cardRec, cardReq)
	if cardRec.Code != http.StatusOK {
		t.Fatalf("ServeAgentCard status = %d, want 200", cardRec.Code)
	}
	var card map[string]any
	if err := json.Unmarshal(cardRec.Body.Bytes(), &card); err != nil {
		t.Fatalf("ServeAgentCard JSON error: %v", err)
	}
	if card["name"] != "test-agent" {
		t.Fatalf("ServeAgentCard name = %v, want test-agent", card["name"])
	}

	skillsReq := httptest.NewRequest(http.MethodGet, "/a2a/skills", nil)
	skillsRec := httptest.NewRecorder()
	h.ServeSkills(skillsRec, skillsReq)
	if skillsRec.Code != http.StatusOK {
		t.Fatalf("ServeSkills status = %d, want 200", skillsRec.Code)
	}
	decoded, err := h.Wire.DecodeToolList(context.Background(), skillsRec.Body.Bytes())
	if err != nil {
		t.Fatalf("DecodeToolList error: %v", err)
	}
	if len(decoded) != 1 || decoded[0].Name != "echo" {
		t.Fatalf("skills = %+v, want echo", decoded)
	}
}

func TestHandler_InvokeAndStatus(t *testing.T) {
	agent := fakeAgent{
		card: map[string]any{"name": "test-agent"},
		skills: []wire.Tool{
			{Name: "echo", Description: "Echo text"},
		},
	}
	h := NewHandler(agent, task.NewManager())

	req := &wire.Request{
		ID:        "task-1",
		Method:    "agent/invoke",
		ToolID:    "echo",
		Arguments: map[string]any{"message": "hi"},
	}
	payload, err := h.Wire.EncodeRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("EncodeRequest error: %v", err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "/a2a", bytes.NewReader(payload))
	httpRec := httptest.NewRecorder()
	h.ServeRPC(httpRec, httpReq)
	if httpRec.Code != http.StatusOK {
		t.Fatalf("ServeRPC status = %d, want 200", httpRec.Code)
	}

	resp, err := h.Wire.DecodeResponse(context.Background(), httpRec.Body.Bytes())
	if err != nil {
		t.Fatalf("DecodeResponse error: %v", err)
	}
	taskID, ok := resp.Meta["taskId"].(string)
	if !ok || taskID == "" {
		t.Fatalf("taskId = %v, want non-empty string", resp.Meta["taskId"])
	}

	statusReq := &wire.Request{
		ID:     taskID,
		Method: "agent/status",
	}
	statusPayload, err := h.Wire.EncodeRequest(context.Background(), statusReq)
	if err != nil {
		t.Fatalf("EncodeRequest(status) error: %v", err)
	}
	statusHTTP := httptest.NewRequest(http.MethodPost, "/a2a", bytes.NewReader(statusPayload))
	statusRec := httptest.NewRecorder()
	h.ServeRPC(statusRec, statusHTTP)
	if statusRec.Code != http.StatusOK {
		t.Fatalf("ServeRPC(status) status = %d, want 200", statusRec.Code)
	}
	statusResp, err := h.Wire.DecodeResponse(context.Background(), statusRec.Body.Bytes())
	if err != nil {
		t.Fatalf("DecodeResponse(status) error: %v", err)
	}
	status, ok := statusResp.Meta["status"].(map[string]any)
	if !ok {
		t.Fatalf("status meta missing: %v", statusResp.Meta)
	}
	if status["id"] != taskID {
		t.Fatalf("status id = %v, want %s", status["id"], taskID)
	}
}

func TestHandler_TaskEvents(t *testing.T) {
	ctx := context.Background()
	agent := fakeAgent{card: map[string]any{"name": "test-agent"}}
	tasks := task.NewManager()
	h := NewHandler(agent, tasks)

	taskID := "task-events"
	if _, err := tasks.Create(ctx, taskID); err != nil {
		t.Fatalf("Create task error: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeTaskEvents(w, r, taskID)
	}))
	defer srv.Close()

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = tasks.Update(context.Background(), taskID, 0.5, "running")
	}()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("GET events error: %v", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var dataLine string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read SSE error: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "event: task" {
			for {
				next, err := reader.ReadString('\n')
				if err != nil {
					t.Fatalf("read SSE data error: %v", err)
				}
				next = strings.TrimSpace(next)
				if strings.HasPrefix(next, "data: ") {
					dataLine = strings.TrimPrefix(next, "data: ")
					break
				}
			}
			break
		}
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(dataLine), &payload); err != nil {
		t.Fatalf("decode SSE payload error: %v", err)
	}
	if payload["id"] != taskID && payload["ID"] != taskID {
		t.Fatalf("payload id = %v (ID=%v), want %s", payload["id"], payload["ID"], taskID)
	}
}
