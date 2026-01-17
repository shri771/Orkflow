/*
Copyright © 2026 Orkflow Authors
*/
package cli

import (
	"fmt"
	"os"

	"Orkflow/internal/vectorstore"

	"github.com/spf13/cobra"
)

var searchLimit int

var sessionsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search sessions semantically",
	Long: `Search through past session messages using semantic similarity.

Requires Ollama running locally with an embedding model (e.g., nomic-embed-text).

Examples:
  orka sessions search "API design patterns"
  orka sessions search "database optimization" --limit 5`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		// Create vector store with Ollama embeddings
		store, err := vectorstore.NewChromemStoreWithOllama("nomic-embed-text")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not connect to vector store: %v\n", err)
			fmt.Println("\n💡 Make sure Ollama is running with an embedding model:")
			fmt.Println("   ollama pull nomic-embed-text")
			os.Exit(1)
		}
		defer store.Close()

		fmt.Printf("🔍 Searching for: \"%s\"\n\n", query)

		results, err := store.Search(query, searchLimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Println("No matching sessions found.")
			fmt.Println("Run some workflows first to build up session history.")
			return
		}

		for i, r := range results {
			fmt.Printf("─── Result %d (%.1f%% match) ───\n", i+1, r.Score*100)
			if sessionID, ok := r.Metadata["session_id"]; ok {
				fmt.Printf("Session: %s\n", sessionID)
			}
			if agentID, ok := r.Metadata["agent_id"]; ok {
				fmt.Printf("Agent: %s\n", agentID)
			}

			// Show truncated content
			content := r.Content
			if len(content) > 300 {
				content = content[:300] + "..."
			}
			fmt.Printf("\n%s\n\n", content)
		}
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsSearchCmd)
	sessionsSearchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 3, "Number of results to return")
}
