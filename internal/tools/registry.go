package tools

import (
	"fmt"
	"sync"
)

// Tool interface for all executable tools
type Tool interface {
	Name() string
	Description() string
	Execute(input string) (string, error)
}

// Registry holds all available tools
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// Global registry
var globalRegistry = &Registry{
	tools: make(map[string]Tool),
}

// Register adds a tool to the global registry
func Register(tool Tool) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.tools[tool.Name()] = tool
}

// Get retrieves a tool by name
func Get(name string) (Tool, bool) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	tool, ok := globalRegistry.tools[name]
	return tool, ok
}

// GetAll returns all registered tools
func GetAll() []Tool {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	result := make([]Tool, 0, len(globalRegistry.tools))
	for _, tool := range globalRegistry.tools {
		result = append(result, tool)
	}
	return result
}

// GetByNames returns tools matching the given names
func GetByNames(names []string) ([]Tool, error) {
	result := make([]Tool, 0, len(names))
	for _, name := range names {
		tool, ok := Get(name)
		if !ok {
			return nil, fmt.Errorf("tool not found: %s", name)
		}
		result = append(result, tool)
	}
	return result, nil
}

// GetByPrefix returns tools that start with the given prefix
func GetByPrefix(prefix string) []Tool {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	result := make([]Tool, 0)
	for name, tool := range globalRegistry.tools {
		// Simple prefix match, e.g. "filesystem." matches "filesystem.list"
		if len(name) > len(prefix) && name[:len(prefix)] == prefix {
			result = append(result, tool)
		}
	}
	return result
}

// ListNames returns all tool names
func ListNames() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	names := make([]string, 0, len(globalRegistry.tools))
	for name := range globalRegistry.tools {
		names = append(names, name)
	}
	return names
}

// FormatToolsForPrompt creates a description of available tools for the LLM
func FormatToolsForPrompt(tools []Tool) string {
	if len(tools) == 0 {
		return ""
	}

	result := "You have access to the following tools:\n\n"
	for _, tool := range tools {
		result += fmt.Sprintf("- **%s**: %s\n", tool.Name(), tool.Description())
	}
	result += "\nTo use a tool, write your response in this format:\n"
	result += "```tool:<tool_name>\n<input for the tool>\n```\n"
	result += "\nThe tool output will be provided to you for further processing.\n"
	return result
}
