package wire

import (
	"context"
	"encoding/json"
	"fmt"
)

// A2AVersion is the A2A protocol version.
const A2AVersion = "0.2.1"

// A2AWire implements Wire for Google's Agent-to-Agent protocol.
type A2AWire struct{}

// NewA2A creates a new A2A wire format handler.
func NewA2A() *A2AWire {
	return &A2AWire{}
}

// Name returns "a2a".
func (w *A2AWire) Name() string {
	return "a2a"
}

// Version returns the A2A spec version.
func (w *A2AWire) Version() string {
	return A2AVersion
}

// a2aRequest is the A2A JSON-RPC request format.
type a2aRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

// a2aResponse is the A2A JSON-RPC response format.
type a2aResponse struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Result  map[string]any `json:"result,omitempty"`
	Error   *jsonrpcError  `json:"error,omitempty"`
}

// EncodeRequest encodes a Request to A2A format.
func (w *A2AWire) EncodeRequest(ctx context.Context, req *Request) ([]byte, error) {
	rpc := a2aRequest{
		JSONRPC: "2.0",
		ID:      req.ID,
		Method:  req.Method,
		Params: map[string]any{
			"id": req.ID,
			"message": map[string]any{
				"role": "user",
				"parts": []map[string]any{
					{
						"kind": "text",
						"text": fmt.Sprintf("invoke %s with %v", req.ToolID, req.Arguments),
					},
				},
			},
		},
	}

	// Store tool info in params
	if req.ToolID != "" {
		rpc.Params["skillId"] = req.ToolID
	}
	if len(req.Arguments) > 0 {
		rpc.Params["arguments"] = req.Arguments
	}
	if len(req.Meta) > 0 {
		rpc.Params["_meta"] = req.Meta
	}

	return json.Marshal(rpc)
}

// DecodeRequest decodes A2A format to a Request.
func (w *A2AWire) DecodeRequest(ctx context.Context, data []byte) (*Request, error) {
	var rpc a2aRequest
	if err := json.Unmarshal(data, &rpc); err != nil {
		return nil, fmt.Errorf("decode a2a request: %w", err)
	}

	req := &Request{
		Method: rpc.Method,
	}

	// Handle ID
	switch id := rpc.ID.(type) {
	case string:
		req.ID = id
	case float64:
		req.ID = fmt.Sprintf("%.0f", id)
	}

	if rpc.Params != nil {
		// Try to get task ID from params
		if taskID, ok := rpc.Params["id"].(string); ok {
			req.ID = taskID
		}
		if skillID, ok := rpc.Params["skillId"].(string); ok {
			req.ToolID = skillID
		}
		if args, ok := rpc.Params["arguments"].(map[string]any); ok {
			req.Arguments = args
		}
		if meta, ok := rpc.Params["_meta"].(map[string]any); ok {
			req.Meta = meta
		}
	}

	return req, nil
}

// EncodeResponse encodes a Response to A2A format.
func (w *A2AWire) EncodeResponse(ctx context.Context, resp *Response) ([]byte, error) {
	rpc := a2aResponse{
		JSONRPC: "2.0",
		ID:      resp.ID,
	}

	if resp.IsError && resp.Error != nil {
		rpc.Error = &jsonrpcError{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
			Data:    resp.Error.Data,
		}
	} else {
		// A2A uses status and artifacts
		artifacts := make([]map[string]any, 0)
		for _, c := range resp.Content {
			part := map[string]any{}
			switch c.Type {
			case ContentTypeText:
				part["kind"] = "text"
				part["text"] = c.Text
			case ContentTypeImage:
				part["kind"] = "data"
				part["data"] = c.Data
				part["mimeType"] = c.MIMEType
			case ContentTypeResource:
				part["kind"] = "file"
				part["uri"] = c.URI
				if c.MIMEType != "" {
					part["mimeType"] = c.MIMEType
				}
			}
			artifacts = append(artifacts, map[string]any{"parts": []map[string]any{part}})
		}

		status := map[string]any{
			"state": "completed",
		}
		if resp.Meta != nil {
			if metaStatus, ok := resp.Meta["status"].(map[string]any); ok {
				if state, ok := metaStatus["state"].(string); ok {
					status["state"] = state
				}
				for k, v := range metaStatus {
					status[k] = v
				}
			} else if state, ok := resp.Meta["state"].(string); ok {
				status["state"] = state
			}
		}

		rpc.Result = map[string]any{
			"status":    status,
			"artifacts": artifacts,
		}
		if resp.Meta != nil {
			if taskID, ok := resp.Meta["taskId"]; ok {
				rpc.Result["taskId"] = taskID
			}
		}
		if len(resp.Meta) > 0 {
			rpc.Result["_meta"] = resp.Meta
		}
	}

	return json.Marshal(rpc)
}

// DecodeResponse decodes A2A format to a Response.
func (w *A2AWire) DecodeResponse(ctx context.Context, data []byte) (*Response, error) {
	var rpc a2aResponse
	if err := json.Unmarshal(data, &rpc); err != nil {
		return nil, fmt.Errorf("decode a2a response: %w", err)
	}

	resp := &Response{}

	// Handle ID
	switch id := rpc.ID.(type) {
	case string:
		resp.ID = id
	case float64:
		resp.ID = fmt.Sprintf("%.0f", id)
	}

	if rpc.Error != nil {
		resp.IsError = true
		resp.Error = &Error{
			Code:    rpc.Error.Code,
			Message: rpc.Error.Message,
			Data:    rpc.Error.Data,
		}
	} else if rpc.Result != nil {
		if resp.Meta == nil {
			resp.Meta = map[string]any{}
		}
		if status, ok := rpc.Result["status"].(map[string]any); ok {
			resp.Meta["status"] = status
		}
		if taskID, ok := rpc.Result["taskId"]; ok {
			resp.Meta["taskId"] = taskID
		}
		// Extract content from artifacts
		if artifacts, ok := rpc.Result["artifacts"].([]any); ok {
			for _, artifact := range artifacts {
				if artMap, ok := artifact.(map[string]any); ok {
					if parts, ok := artMap["parts"].([]any); ok {
						for _, part := range parts {
							if partMap, ok := part.(map[string]any); ok {
								c := Content{}
								if kind, ok := partMap["kind"].(string); ok {
									switch kind {
									case "text":
										c.Type = ContentTypeText
										if text, ok := partMap["text"].(string); ok {
											c.Text = text
										}
									case "data":
										c.Type = ContentTypeImage
										if mime, ok := partMap["mimeType"].(string); ok {
											c.MIMEType = mime
										}
									case "file":
										c.Type = ContentTypeResource
										if uri, ok := partMap["uri"].(string); ok {
											c.URI = uri
										}
									}
								}
								resp.Content = append(resp.Content, c)
							}
						}
					}
				}
			}
		}
		if meta, ok := rpc.Result["_meta"].(map[string]any); ok {
			resp.Meta = meta
		}
	}

	return resp, nil
}

// a2aSkillList is the A2A skills list format.
type a2aSkillList struct {
	Skills []a2aSkill `json:"skills"`
}

type a2aSkill struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"inputSchema,omitempty"`
}

// EncodeToolList encodes a tool list to A2A format.
func (w *A2AWire) EncodeToolList(ctx context.Context, tools []Tool) ([]byte, error) {
	list := a2aSkillList{
		Skills: make([]a2aSkill, 0, len(tools)),
	}

	for _, t := range tools {
		list.Skills = append(list.Skills, a2aSkill{
			ID:          t.Name,
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}

	return json.Marshal(list)
}

// DecodeToolList decodes A2A format to a tool list.
func (w *A2AWire) DecodeToolList(ctx context.Context, data []byte) ([]Tool, error) {
	var list a2aSkillList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("decode a2a tool list: %w", err)
	}

	tools := make([]Tool, 0, len(list.Skills))
	for _, s := range list.Skills {
		name := s.ID
		if name == "" {
			name = s.Name
		}
		tools = append(tools, Tool{
			Name:        name,
			Description: s.Description,
			InputSchema: s.InputSchema,
		})
	}

	return tools, nil
}

// Capabilities returns A2A protocol capabilities.
func (w *A2AWire) Capabilities() *Capabilities {
	return &Capabilities{
		Streaming:     true,
		BatchRequests: false,
		Progress:      true,
		Cancellation:  true,
	}
}
