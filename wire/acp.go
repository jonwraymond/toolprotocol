package wire

import (
	"context"
	"encoding/json"
	"fmt"
)

// ACPVersion is the ACP protocol version.
const ACPVersion = "1.0.0"

// ACPWire implements Wire for IBM's Agent Communication Protocol.
type ACPWire struct{}

// NewACP creates a new ACP wire format handler.
func NewACP() *ACPWire {
	return &ACPWire{}
}

// Name returns "acp".
func (w *ACPWire) Name() string {
	return "acp"
}

// Version returns the ACP spec version.
func (w *ACPWire) Version() string {
	return ACPVersion
}

// acpRequest is the ACP JSON-RPC request format.
type acpRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

// acpResponse is the ACP JSON-RPC response format.
type acpResponse struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Result  map[string]any `json:"result,omitempty"`
	Error   *jsonrpcError  `json:"error,omitempty"`
}

// EncodeRequest encodes a Request to ACP format.
func (w *ACPWire) EncodeRequest(ctx context.Context, req *Request) ([]byte, error) {
	rpc := acpRequest{
		JSONRPC: "2.0",
		ID:      req.ID,
		Method:  req.Method,
		Params: map[string]any{
			"agentId": req.ToolID,
			"input":   req.Arguments,
		},
	}

	if len(req.Meta) > 0 {
		rpc.Params["metadata"] = req.Meta
	}

	return json.Marshal(rpc)
}

// DecodeRequest decodes ACP format to a Request.
func (w *ACPWire) DecodeRequest(ctx context.Context, data []byte) (*Request, error) {
	var rpc acpRequest
	if err := json.Unmarshal(data, &rpc); err != nil {
		return nil, fmt.Errorf("decode acp request: %w", err)
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
		if agentID, ok := rpc.Params["agentId"].(string); ok {
			req.ToolID = agentID
		}
		if input, ok := rpc.Params["input"].(map[string]any); ok {
			req.Arguments = input
		}
		if meta, ok := rpc.Params["metadata"].(map[string]any); ok {
			req.Meta = meta
		}
	}

	return req, nil
}

// EncodeResponse encodes a Response to ACP format.
func (w *ACPWire) EncodeResponse(ctx context.Context, resp *Response) ([]byte, error) {
	rpc := acpResponse{
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
		output := make(map[string]any)
		for i, c := range resp.Content {
			key := fmt.Sprintf("content_%d", i)
			switch c.Type {
			case ContentTypeText:
				output[key] = map[string]any{
					"type": "text",
					"text": c.Text,
				}
			case ContentTypeImage:
				output[key] = map[string]any{
					"type":     "binary",
					"data":     c.Data,
					"mimeType": c.MIMEType,
				}
			case ContentTypeResource:
				output[key] = map[string]any{
					"type": "resource",
					"uri":  c.URI,
				}
			}
		}

		rpc.Result = map[string]any{
			"status": "success",
			"output": output,
		}
		if len(resp.Meta) > 0 {
			rpc.Result["metadata"] = resp.Meta
		}
	}

	return json.Marshal(rpc)
}

// DecodeResponse decodes ACP format to a Response.
func (w *ACPWire) DecodeResponse(ctx context.Context, data []byte) (*Response, error) {
	var rpc acpResponse
	if err := json.Unmarshal(data, &rpc); err != nil {
		return nil, fmt.Errorf("decode acp response: %w", err)
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
		// Extract content from output
		if output, ok := rpc.Result["output"].(map[string]any); ok {
			for _, v := range output {
				if item, ok := v.(map[string]any); ok {
					c := Content{}
					if t, ok := item["type"].(string); ok {
						switch t {
						case "text":
							c.Type = ContentTypeText
							if text, ok := item["text"].(string); ok {
								c.Text = text
							}
						case "binary":
							c.Type = ContentTypeImage
							if mime, ok := item["mimeType"].(string); ok {
								c.MIMEType = mime
							}
						case "resource":
							c.Type = ContentTypeResource
							if uri, ok := item["uri"].(string); ok {
								c.URI = uri
							}
						}
					}
					resp.Content = append(resp.Content, c)
				}
			}
		}
		if meta, ok := rpc.Result["metadata"].(map[string]any); ok {
			resp.Meta = meta
		}
	}

	return resp, nil
}

// acpAgentList is the ACP agents list format.
type acpAgentList struct {
	Agents []acpAgent `json:"agents"`
}

type acpAgent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"inputSchema,omitempty"`
}

// EncodeToolList encodes a tool list to ACP format.
func (w *ACPWire) EncodeToolList(ctx context.Context, tools []Tool) ([]byte, error) {
	list := acpAgentList{
		Agents: make([]acpAgent, 0, len(tools)),
	}

	for _, t := range tools {
		list.Agents = append(list.Agents, acpAgent{
			ID:          t.Name,
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}

	return json.Marshal(list)
}

// DecodeToolList decodes ACP format to a tool list.
func (w *ACPWire) DecodeToolList(ctx context.Context, data []byte) ([]Tool, error) {
	var list acpAgentList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("decode acp tool list: %w", err)
	}

	tools := make([]Tool, 0, len(list.Agents))
	for _, a := range list.Agents {
		name := a.ID
		if name == "" {
			name = a.Name
		}
		tools = append(tools, Tool{
			Name:        name,
			Description: a.Description,
			InputSchema: a.InputSchema,
		})
	}

	return tools, nil
}

// Capabilities returns ACP protocol capabilities.
func (w *ACPWire) Capabilities() *Capabilities {
	return &Capabilities{
		Streaming:     false,
		BatchRequests: true,
		Progress:      false,
		Cancellation:  true,
	}
}
