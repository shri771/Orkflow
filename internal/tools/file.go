package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileTool provides file system operations
type FileTool struct{}

func init() {
	Register(&FileTool{})
}

func (f *FileTool) Name() string {
	return "file"
}

func (f *FileTool) Description() string {
	return "File operations. Commands: 'read:<path>' to read file, 'write:<path>:<content>' to write, 'list:<dir>' to list directory, 'exists:<path>' to check existence."
}

func (f *FileTool) Execute(input string) (string, error) {
	input = strings.TrimSpace(input)

	// Parse command
	parts := strings.SplitN(input, ":", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid format. Use 'read:<path>', 'write:<path>:<content>', 'list:<dir>', or 'exists:<path>'")
	}

	cmd := strings.ToLower(parts[0])
	arg := parts[1]

	switch cmd {
	case "read":
		return f.readFile(arg)
	case "write":
		writeParts := strings.SplitN(arg, ":", 2)
		if len(writeParts) < 2 {
			return "", fmt.Errorf("write requires path and content: 'write:<path>:<content>'")
		}
		return f.writeFile(writeParts[0], writeParts[1])
	case "list":
		return f.listDir(arg)
	case "exists":
		return f.exists(arg)
	default:
		return "", fmt.Errorf("unknown command: %s. Use read, write, list, or exists", cmd)
	}
}

func (f *FileTool) readFile(path string) (string, error) {
	path = filepath.Clean(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

func (f *FileTool) writeFile(path, content string) (string, error) {
	path = filepath.Clean(path)

	// Create parent directories if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return fmt.Sprintf("Successfully wrote %d bytes to %s", len(content), path), nil
}

func (f *FileTool) listDir(path string) (string, error) {
	path = filepath.Clean(path)
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to list directory: %w", err)
	}

	var result strings.Builder
	for _, entry := range entries {
		if entry.IsDir() {
			result.WriteString(fmt.Sprintf("[DIR]  %s\n", entry.Name()))
		} else {
			info, _ := entry.Info()
			size := int64(0)
			if info != nil {
				size = info.Size()
			}
			result.WriteString(fmt.Sprintf("[FILE] %s (%d bytes)\n", entry.Name(), size))
		}
	}
	return result.String(), nil
}

func (f *FileTool) exists(path string) (string, error) {
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "false", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to check path: %w", err)
	}
	if info.IsDir() {
		return "true (directory)", nil
	}
	return fmt.Sprintf("true (file, %d bytes)", info.Size()), nil
}
