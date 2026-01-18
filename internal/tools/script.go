package tools

import (
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

// ScriptTool executes Tengo scripts (Go-like syntax)
type ScriptTool struct{}

func init() {
	Register(&ScriptTool{})
}

func (s *ScriptTool) Name() string {
	return "script"
}

func (s *ScriptTool) Description() string {
	return "Execute scripts using Tengo (Go-like syntax). Supports variables, loops, functions, math, and string operations. Set 'output' variable to return a value."
}

func (s *ScriptTool) Execute(input string) (string, error) {
	// Wrap script to capture output variable
	wrappedScript := fmt.Sprintf(`
output := ""
__run := func() {
%s
}
__run()
`, input)

	// Create script with standard library
	script := tengo.NewScript([]byte(wrappedScript))

	// Add standard library modules
	script.SetImports(stdlib.GetModuleMap(
		"fmt",
		"math",
		"text",
		"times",
		"rand",
		"json",
	))

	// Run the script
	compiled, err := script.Run()
	if err != nil {
		return "", fmt.Errorf("script error: %w", err)
	}

	// Get the output variable
	output := compiled.Get("output")
	if output.IsUndefined() || output.String() == "" {
		return "Script executed successfully", nil
	}

	return output.String(), nil
}

// Example usage in prompts:
// ```tool:script
// a := 10
// b := 20
// output = a + b
// ```
