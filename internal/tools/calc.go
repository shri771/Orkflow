package tools

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// CalcTool evaluates mathematical expressions
type CalcTool struct{}

func init() {
	Register(&CalcTool{})
}

func (c *CalcTool) Name() string {
	return "calc"
}

func (c *CalcTool) Description() string {
	return "Evaluate mathematical expressions. Supports +, -, *, /, %, ^, comparisons, and functions like abs(), max(), min(), len()."
}

func (c *CalcTool) Execute(input string) (string, error) {
	program, err := expr.Compile(input)
	if err != nil {
		return "", fmt.Errorf("expression error: %w", err)
	}

	result, err := expr.Run(program, nil)
	if err != nil {
		return "", fmt.Errorf("evaluation error: %w", err)
	}

	return fmt.Sprintf("%v", result), nil
}
