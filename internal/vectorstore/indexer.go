package vectorstore

import (
	"fmt"

	"Orkflow/internal/memory"
)

// IndexSession indexes all messages from a session into the vector store
func IndexSession(store *ChromemStore, session *memory.Session) error {
	for i, msg := range session.Messages {
		// Create a unique document ID
		docID := fmt.Sprintf("%s_%d", session.ID, i)

		// Store metadata
		metadata := map[string]string{
			"session_id": session.ID,
			"workflow":   session.Workflow,
			"agent_id":   msg.AgentID,
			"role":       msg.Role,
			"timestamp":  msg.Timestamp.Format("2006-01-02 15:04:05"),
		}

		// Add document to vector store
		if err := store.AddDocument(docID, msg.Content, metadata); err != nil {
			return fmt.Errorf("failed to index message %d: %w", i, err)
		}
	}
	return nil
}
