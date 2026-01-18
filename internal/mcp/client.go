package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerConfig defines an MCP server configuration from YAML
type ServerConfig struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
	Env     []string `yaml:"env,omitempty"`
}

// Client manages connections to MCP servers
type Client struct {
	mu      sync.RWMutex
	servers map[string]*mcpServer
	client  *mcp.Client
	ctx     context.Context
	cancel  context.CancelFunc
}

type mcpServer struct {
	config  ServerConfig
	session *mcp.ClientSession
	tools   []*mcp.Tool
}

// NewClient creates a new MCP client manager
func NewClient() *Client {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a single MCP client instance
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "orka",
		Version: "1.0.0",
	}, nil)

	return &Client{
		servers: make(map[string]*mcpServer),
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Connect starts an MCP server and connects to it
func (c *Client) Connect(name string, config ServerConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create command transport
	cmd := exec.Command(config.Command, config.Args...)
	if len(config.Env) > 0 {
		cmd.Env = append(cmd.Env, config.Env...)
	}

	transport := &mcp.CommandTransport{Command: cmd}

	// Connect to the server
	session, err := c.client.Connect(c.ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}

	// List available tools
	toolsResult, err := session.ListTools(c.ctx, nil)
	if err != nil {
		session.Close()
		return fmt.Errorf("failed to list tools: %w", err)
	}

	server := &mcpServer{
		config:  config,
		session: session,
		tools:   toolsResult.Tools,
	}

	c.servers[name] = server
	fmt.Printf("🔌 Connected to MCP server '%s' with %d tools\n", name, len(server.tools))

	return nil
}

// GetTools returns tools from a specific server
func (c *Client) GetTools(serverName string) ([]*mcp.Tool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	server, ok := c.servers[serverName]
	if !ok {
		return nil, fmt.Errorf("MCP server not found: %s", serverName)
	}

	return server.tools, nil
}

// GetAllTools returns tools from all connected servers
func (c *Client) GetAllTools() map[string][]*mcp.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string][]*mcp.Tool)
	for name, server := range c.servers {
		result[name] = server.tools
	}
	return result
}

// CallTool executes a tool on an MCP server
func (c *Client) CallTool(serverName, toolName string, args map[string]interface{}) (string, error) {
	c.mu.RLock()
	server, ok := c.servers[serverName]
	c.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("MCP server not found: %s", serverName)
	}

	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	}

	result, err := server.session.CallTool(c.ctx, params)
	if err != nil {
		return "", fmt.Errorf("tool call failed: %w", err)
	}

	// Extract text content from result
	var output string
	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			output += textContent.Text
		}
	}

	return output, nil
}

// Close shuts down all MCP server connections
func (c *Client) Close() error {
	c.cancel()

	c.mu.Lock()
	defer c.mu.Unlock()

	for name, server := range c.servers {
		if server.session != nil {
			server.session.Close()
		}
		delete(c.servers, name)
	}

	return nil
}

// ListServerNames returns all connected server names
func (c *Client) ListServerNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.servers))
	for name := range c.servers {
		names = append(names, name)
	}
	return names
}
