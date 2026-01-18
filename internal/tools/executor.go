package tools

import (
	"fmt"
	"regexp"
	"strings"
)

// ToolCall represents a parsed tool invocation from LLM output
type ToolCall struct {
	Name  string
	Input string
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolName string
	Output   string
	Error    error
}

// ParseToolCalls extracts tool calls from LLM response
// Format: ```tool:<name>\n<input>\n```
func ParseToolCalls(response string) []ToolCall {
	// Match ```tool:<name>\n...\n```
	re := regexp.MustCompile("(?s)```tool:([a-zA-Z_][a-zA-Z0-9_]*)\n(.*?)```")
	matches := re.FindAllStringSubmatch(response, -1)

	var calls []ToolCall
	for _, match := range matches {
		if len(match) >= 3 {
			calls = append(calls, ToolCall{
				Name:  strings.TrimSpace(match[1]),
				Input: strings.TrimSpace(match[2]),
			})
		}
	}
	return calls
}

// ExecuteToolCalls runs all parsed tool calls and returns results
func ExecuteToolCalls(calls []ToolCall) []ToolResult {
	var results []ToolResult

	for _, call := range calls {
		tool, ok := Get(call.Name)
		if !ok {
			results = append(results, ToolResult{
				ToolName: call.Name,
				Error:    fmt.Errorf("unknown tool: %s", call.Name),
			})
			continue
		}

		fmt.Printf("  🔧 Executing tool: %s\n", call.Name)
		output, err := tool.Execute(call.Input)
		results = append(results, ToolResult{
			ToolName: call.Name,
			Output:   output,
			Error:    err,
		})
	}

	return results
}

// FormatToolResults creates a string describing tool results for LLM
func FormatToolResults(results []ToolResult) string {
	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n\n=== Tool Results ===\n")

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("\n[%s]:\n", r.ToolName))
		if r.Error != nil {
			sb.WriteString(fmt.Sprintf("ERROR: %v\n", r.Error))
		} else {
			sb.WriteString(r.Output)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// HasToolCalls checks if a response contains tool calls
func HasToolCalls(response string) bool {
	return strings.Contains(response, "```tool:")
}
