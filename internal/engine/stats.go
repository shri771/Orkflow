package engine

import (
	"sync"
	"time"
)

// ExecutionStats tracks timing and cost for a workflow run
type ExecutionStats struct {
	mu          sync.Mutex
	StartTime   time.Time
	AgentStats  map[string]*AgentStat
	TotalTokens struct {
		Input  int
		Output int
	}
}

// AgentStat tracks per-agent statistics
type AgentStat struct {
	AgentID      string
	Role         string
	Model        string
	StartTime    time.Time
	Duration     time.Duration
	InputTokens  int
	OutputTokens int
	Completed    bool
}

// NewExecutionStats creates a new stats tracker
func NewExecutionStats() *ExecutionStats {
	return &ExecutionStats{
		StartTime:  time.Now(),
		AgentStats: make(map[string]*AgentStat),
	}
}

// StartAgent marks an agent as started
func (s *ExecutionStats) StartAgent(agentID, role, model string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.AgentStats[agentID] = &AgentStat{
		AgentID:   agentID,
		Role:      role,
		Model:     model,
		StartTime: time.Now(),
	}
}

// CompleteAgent marks an agent as completed with token counts
func (s *ExecutionStats) CompleteAgent(agentID string, inputTokens, outputTokens int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if stat, ok := s.AgentStats[agentID]; ok {
		stat.Duration = time.Since(stat.StartTime)
		stat.InputTokens = inputTokens
		stat.OutputTokens = outputTokens
		stat.Completed = true

		s.TotalTokens.Input += inputTokens
		s.TotalTokens.Output += outputTokens
	}
}

// GetElapsedTime returns total elapsed time
func (s *ExecutionStats) GetElapsedTime() time.Duration {
	return time.Since(s.StartTime)
}

// GetCompletedCount returns number of completed agents
func (s *ExecutionStats) GetCompletedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, stat := range s.AgentStats {
		if stat.Completed {
			count++
		}
	}
	return count
}

// EstimateCost calculates estimated cost based on token usage and model
func (s *ExecutionStats) EstimateCost() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Pricing per 1M tokens (input/output)
	pricing := map[string]struct{ Input, Output float64 }{
		"gpt-4o":           {2.50, 10.00},
		"gpt-4o-mini":      {0.15, 0.60},
		"gpt-4-turbo":      {10.00, 30.00},
		"gpt-3.5-turbo":    {0.50, 1.50},
		"gemini-2.0-flash": {0.075, 0.30},
		"gemini-1.5-pro":   {1.25, 5.00},
	}

	var totalCost float64
	for _, stat := range s.AgentStats {
		if p, ok := pricing[stat.Model]; ok {
			inputCost := float64(stat.InputTokens) / 1000000 * p.Input
			outputCost := float64(stat.OutputTokens) / 1000000 * p.Output
			totalCost += inputCost + outputCost
		}
	}

	return totalCost
}
