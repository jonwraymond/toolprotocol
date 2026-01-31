package wire

import (
	"context"
	"encoding/json"
	"fmt"
)

// MCPVersion is the MCP specification version.
const MCPVersion = "2025-11-25"

// MCPWire implements Wire for the Model Context Protocol.
type MCPWire struct{}

// NewMCP creates a new MCP wire format handler.
func NewMCP() *MCPWire {
	return &MCPWire{}
}

// Name returns "mcp".
func (w *MCPWire) Name() string {
	return "mcp"
}

// Version returns the MCP spec version.
func (w *MCPWire) Version() string {
	return MCPVersion
}

// jsonrpcRequest is the JSON-RPC 2.0 request format.
type jsonrpcRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

// jsonrpcResponse is the JSON-RPC 2.0 response format.
type jsonrpcResponse struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Result  map[string]any `json:"result,omitempty"`
	Error   *jsonrpcError  `json:"error,omitempty"`
}

// jsonrpcError is the JSON-RPC 2.0 error format.
type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// EncodeRequest encodes a Request to MCP JSON-RPC format.
func (w *MCPWire) EncodeRequest(ctx context.Context, req *Request) ([]byte, error) {
	rpc := jsonrpcRequest{
		JSONRPC: "2.0",
		ID:      req.ID,
		Method:  req.Method,
		Params: map[string]any{
			"name":      req.ToolID,
			"arguments": req.Arguments,
		},
	}

	if len(req.Meta) > 0 {
		rpc.Params["_meta"] = req.Meta
	}

	return json.Marshal(rpc)
}

// DecodeRequest decodes MCP JSON-RPC format to a Request.
func (w *MCPWire) DecodeRequest(ctx context.Context, data []byte) (*Request, error) {
	var rpc jsonrpcRequest
	if err := json.Unmarshal(data, &rpc); err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}

	req := &Request{
		Method: rpc.Method,
	}

	// Handle ID which can be string or number
	switch id := rpc.ID.(type) {
	case string:
		req.ID = id
	case float64:
		req.ID = fmt.Sprintf("%.0f", id)
	}

	if rpc.Params != nil {
		if name, ok := rpc.Params["name"].(string); ok {
			req.ToolID = name
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

// EncodeResponse encodes a Response to MCP JSON-RPC format.
func (w *MCPWire) EncodeResponse(ctx context.Context, resp *Response) ([]byte, error) {
	rpc := jsonrpcResponse{
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
		content := make([]map[string]any, 0, len(resp.Content))
		for _, c := range resp.Content {
			item := map[string]any{"type": string(c.Type)}
			switch c.Type {
			case ContentTypeText:
				item["text"] = c.Text
			case ContentTypeImage:
				item["data"] = c.Data
				item["mimeType"] = c.MIMEType
			case ContentTypeResource:
				item["uri"] = c.URI
				if c.MIMEType != "" {
					item["mimeType"] = c.MIMEType
				}
			}
			content = append(content, item)
		}
		rpc.Result = map[string]any{"content": content}
		if len(resp.Meta) > 0 {
			rpc.Result["_meta"] = resp.Meta
		}
	}

	return json.Marshal(rpc)
}

// DecodeResponse decodes MCP JSON-RPC format to a Response.
func (w *MCPWire) DecodeResponse(ctx context.Context, data []byte) (*Response, error) {
	var rpc jsonrpcResponse
	if err := json.Unmarshal(data, &rpc); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
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
		if contentArr, ok := rpc.Result["content"].([]any); ok {
			for _, item := range contentArr {
				if cm, ok := item.(map[string]any); ok {
					c := Content{}
					if t, ok := cm["type"].(string); ok {
						c.Type = ContentType(t)
					}
					if text, ok := cm["text"].(string); ok {
						c.Text = text
					}
					if uri, ok := cm["uri"].(string); ok {
						c.URI = uri
					}
					if mime, ok := cm["mimeType"].(string); ok {
						c.MIMEType = mime
					}
					resp.Content = append(resp.Content, c)
				}
			}
		}
		if meta, ok := rpc.Result["_meta"].(map[string]any); ok {
			resp.Meta = meta
		}
	}

	return resp, nil
}

// mcpToolList is the MCP tools/list response format.
type mcpToolList struct {
	Tools []mcpTool `json:"tools"`
}

type mcpTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"inputSchema,omitempty"`
}

// EncodeToolList encodes a tool list to MCP format.
func (w *MCPWire) EncodeToolList(ctx context.Context, tools []Tool) ([]byte, error) {
	list := mcpToolList{
		Tools: make([]mcpTool, 0, len(tools)),
	}

	for _, t := range tools {
		list.Tools = append(list.Tools, mcpTool(t))
	}

	return json.Marshal(list)
}

// DecodeToolList decodes MCP format to a tool list.
func (w *MCPWire) DecodeToolList(ctx context.Context, data []byte) ([]Tool, error) {
	var list mcpToolList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("decode tool list: %w", err)
	}

	tools := make([]Tool, 0, len(list.Tools))
	for _, t := range list.Tools {
		tools = append(tools, Tool(t))
	}

	return tools, nil
}

// Capabilities returns MCP protocol capabilities.
func (w *MCPWire) Capabilities() *Capabilities {
	return &Capabilities{
		Streaming:     true,
		BatchRequests: false,
		Progress:      true,
		Cancellation:  true,
	}
}
