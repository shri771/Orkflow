package tools

import (
	"testing"
)

func TestCalcTool(t *testing.T) {
	calc := &CalcTool{}

	tests := []struct {
		input    string
		expected string
	}{
		{"2 + 2", "4"},
		{"10 * 5", "50"},
		{"100 / 4", "25"},
		{"2 ^ 10", "1024"},
		{"10 > 5", "true"},
		{"len([1,2,3])", "3"},
	}

	for _, tt := range tests {
		result, err := calc.Execute(tt.input)
		if err != nil {
			t.Errorf("calc.Execute(%q) error: %v", tt.input, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("calc.Execute(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestScriptTool(t *testing.T) {
	script := &ScriptTool{}

	// Simple addition
	result, err := script.Execute(`
		a := 10
		b := 20
		output = a + b
	`)
	if err != nil {
		t.Errorf("script error: %v", err)
	}
	if result != "30" {
		t.Errorf("expected '30', got %q", result)
	}

	// Loop
	result, err = script.Execute(`
		sum := 0
		for i := 1; i <= 5; i++ {
			sum = sum + i
		}
		output = sum
	`)
	if err != nil {
		t.Errorf("script error: %v", err)
	}
	if result != "15" {
		t.Errorf("expected '15', got %q", result)
	}
}

func TestFileTool(t *testing.T) {
	file := &FileTool{}

	// Test exists
	result, err := file.Execute("exists:/tmp")
	if err != nil {
		t.Errorf("file.exists error: %v", err)
	}
	if result != "true (directory)" {
		t.Errorf("expected 'true (directory)', got %q", result)
	}

	// Test list
	result, err = file.Execute("list:/tmp")
	if err != nil {
		t.Errorf("file.list error: %v", err)
	}
	// Just check it doesn't error
	t.Logf("list /tmp: %s", result[:min(100, len(result))])
}

func TestRegistry(t *testing.T) {
	// Tools should be auto-registered via init()
	names := ListNames()
	if len(names) < 3 {
		t.Errorf("expected at least 3 tools, got %d: %v", len(names), names)
	}

	// Check specific tools exist
	for _, name := range []string{"script", "calc", "file"} {
		if _, ok := Get(name); !ok {
			t.Errorf("tool %q should be registered", name)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
