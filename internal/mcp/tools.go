package mcp

import (
	"fmt"

	"Orkflow/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPTool wraps an MCP tool to implement the Tool interface
type MCPTool struct {
	ServerName string
	ToolDef    *mcp.Tool
	Client     *Client
}

func (t *MCPTool) Name() string {
	return fmt.Sprintf("%s.%s", t.ServerName, t.ToolDef.Name)
}

func (t *MCPTool) Description() string {
	return t.ToolDef.Description
}

func (t *MCPTool) Execute(input string) (string, error) {
	// Parse input as simple key=value or just pass as single arg
	args := map[string]interface{}{
		"input": input,
	}

	return t.Client.CallTool(t.ServerName, t.ToolDef.Name, args)
}

// RegisterMCPTools registers all tools from an MCP server with the tool registry
func RegisterMCPTools(client *Client, serverName string) error {
	mcpTools, err := client.GetTools(serverName)
	if err != nil {
		return err
	}

	for _, toolDef := range mcpTools {
		tool := &MCPTool{
			ServerName: serverName,
			ToolDef:    toolDef,
			Client:     client,
		}
		tools.Register(tool)
		fmt.Printf("  📦 Registered MCP tool: %s\n", tool.Name())
	}

	return nil
}

// FormatMCPToolsForPrompt creates descriptions of MCP tools for the LLM
func FormatMCPToolsForPrompt(client *Client) string {
	allTools := client.GetAllTools()
	if len(allTools) == 0 {
		return ""
	}

	result := "You have access to the following MCP tools:\n\n"
	for serverName, serverTools := range allTools {
		for _, tool := range serverTools {
			result += fmt.Sprintf("- **%s.%s**: %s\n", serverName, tool.Name, tool.Description)
		}
	}
	result += "\nTo use an MCP tool, write your response in this format:\n"
	result += "```tool:<server>.<toolname>\n<input>\n```\n"
	return result
}
