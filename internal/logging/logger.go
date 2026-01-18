package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger handles file-based execution logging
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	filePath string
	enabled  bool
}

// NewLogger creates a new file logger
func NewLogger(sessionID string, logDir string) (*Logger, error) {
	if logDir == "" {
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, ".orka", "logs")
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.log", timestamp, sessionID)
	filePath := filepath.Join(logDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	logger := &Logger{
		file:     file,
		filePath: filePath,
		enabled:  true,
	}

	// Write header
	logger.writeHeader(sessionID)

	return logger, nil
}

// writeHeader writes the log file header
func (l *Logger) writeHeader(sessionID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	header := fmt.Sprintf(`╔══════════════════════════════════════════════════════════════╗
║                    ORKFLOW EXECUTION LOG                     ║
╠══════════════════════════════════════════════════════════════╣
║  Session: %-50s ║
║  Started: %-50s ║
╚══════════════════════════════════════════════════════════════╝

`, sessionID, time.Now().Format("2006-01-02 15:04:05"))

	l.file.WriteString(header)
}

// Log writes a message to the log file
func (l *Logger) Log(format string, args ...interface{}) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] %s\n", timestamp, msg)
	l.file.WriteString(line)
}

// LogAgent logs agent-specific events
func (l *Logger) LogAgent(agentID, event, details string) {
	l.Log("[%s] %s: %s", agentID, event, details)
}

// LogSection writes a section header
func (l *Logger) LogSection(title string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	section := fmt.Sprintf("\n═══════════════════════════════════════════════════════════════\n")
	section += fmt.Sprintf("  %s\n", title)
	section += fmt.Sprintf("═══════════════════════════════════════════════════════════════\n\n")
	l.file.WriteString(section)
}

// LogAgentOutput logs the full output from an agent
func (l *Logger) LogAgentOutput(agentID, role, output string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	section := fmt.Sprintf("\n┌─────────────────────────────────────────────────────────────┐\n")
	section += fmt.Sprintf("│ Agent: %-54s │\n", agentID)
	section += fmt.Sprintf("│ Role: %-55s │\n", role)
	section += fmt.Sprintf("└─────────────────────────────────────────────────────────────┘\n")
	section += output + "\n"
	l.file.WriteString(section)
}

// LogError logs an error
func (l *Logger) LogError(err error) {
	l.Log("ERROR: %v", err)
}

// LogToolCall logs a tool execution
func (l *Logger) LogToolCall(toolName, input, output string) {
	l.Log("TOOL [%s] Input: %s", toolName, truncate(input, 100))
	l.Log("TOOL [%s] Output: %s", toolName, truncate(output, 200))
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Write footer
	footer := fmt.Sprintf("\n╔══════════════════════════════════════════════════════════════╗\n")
	footer += fmt.Sprintf("║  Completed: %-48s ║\n", time.Now().Format("2006-01-02 15:04:05"))
	footer += fmt.Sprintf("╚══════════════════════════════════════════════════════════════╝\n")
	l.file.WriteString(footer)

	return l.file.Close()
}

// GetFilePath returns the log file path
func (l *Logger) GetFilePath() string {
	return l.filePath
}

// truncate shortens a string for logging
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// NullLogger is a no-op logger when logging is disabled
type NullLogger struct{}

func (n *NullLogger) Log(format string, args ...interface{})      {}
func (n *NullLogger) LogAgent(agentID, event, details string)     {}
func (n *NullLogger) LogSection(title string)                     {}
func (n *NullLogger) LogAgentOutput(agentID, role, output string) {}
func (n *NullLogger) LogError(err error)                          {}
func (n *NullLogger) LogToolCall(toolName, input, output string)  {}
func (n *NullLogger) Close() error                                { return nil }
func (n *NullLogger) GetFilePath() string                         { return "" }
