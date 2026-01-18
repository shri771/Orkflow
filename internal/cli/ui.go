package cli

import "fmt"

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// Agent emojis - auto-assigned based on index
var agentEmojis = []string{"🤖", "🎨", "🔧", "📊", "🧠", "✨", "🚀", "💡", "📝", "🎯", "⚙️", "🔍"}

// GetAgentEmoji returns a unique emoji for each agent based on index
func GetAgentEmoji(index int) string {
	return agentEmojis[index%len(agentEmojis)]
}

// ColorText wraps text with ANSI color codes
func ColorText(text, color string) string {
	return color + text + ColorReset
}

// ProgressBar generates a text-based progress bar
func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))
	empty := width - filled

	bar := "["
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}
	bar += "]"

	return fmt.Sprintf("%s %d/%d (%.0f%%)", bar, current, total, percent*100)
}

// FormatDuration formats seconds into human readable string
func FormatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	}
	mins := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%dm %ds", mins, secs)
}

// ModelPricing stores cost per 1M tokens (input/output)
var ModelPricing = map[string]struct{ Input, Output float64 }{
	"gpt-4o":           {2.50, 10.00},
	"gpt-4o-mini":      {0.15, 0.60},
	"gpt-4-turbo":      {10.00, 30.00},
	"gpt-3.5-turbo":    {0.50, 1.50},
	"gemini-2.0-flash": {0.075, 0.30},
	"gemini-1.5-pro":   {1.25, 5.00},
}

// EstimateCost calculates cost based on token counts
func EstimateCost(model string, inputTokens, outputTokens int) float64 {
	pricing, ok := ModelPricing[model]
	if !ok {
		return 0
	}
	inputCost := float64(inputTokens) / 1000000 * pricing.Input
	outputCost := float64(outputTokens) / 1000000 * pricing.Output
	return inputCost + outputCost
}
