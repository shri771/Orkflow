package vectorstore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/philippgille/chromem-go"
)

const (
	CollectionName = "orka_sessions"
	VectorDBPath   = ".orka/vectordb"
)

// VectorStore interface for vector storage backends
type VectorStore interface {
	AddDocument(id string, content string, metadata map[string]string) error
	Search(query string, limit int) ([]SearchResult, error)
	DeleteDocument(id string) error
	Close() error
}

// SearchResult represents a search result from vector store
type SearchResult struct {
	ID       string
	Content  string
	Score    float32
	Metadata map[string]string
}

// ChromemStore implements VectorStore using chromem-go (embedded)
type ChromemStore struct {
	db         *chromem.DB
	collection *chromem.Collection
	ctx        context.Context
}

// NewChromemStoreWithOllama creates a store with Ollama embeddings
func NewChromemStoreWithOllama(model string) (*ChromemStore, error) {
	ef := chromem.NewEmbeddingFuncOllama(model, "")
	return newChromemStore(ef)
}

// NewChromemStoreWithOpenAI creates a store with OpenAI-compatible embeddings
func NewChromemStoreWithOpenAI(apiKey string) (*ChromemStore, error) {
	ef := chromem.NewEmbeddingFuncOpenAI(apiKey, chromem.EmbeddingModelOpenAI3Small)
	return newChromemStore(ef)
}

// NewChromemStoreWithMistral creates a store with Mistral embeddings
func NewChromemStoreWithMistral(apiKey string) (*ChromemStore, error) {
	ef := chromem.NewEmbeddingFuncMistral(apiKey)
	return newChromemStore(ef)
}

func newChromemStore(ef chromem.EmbeddingFunc) (*ChromemStore, error) {
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, VectorDBPath)

	// Ensure directory exists
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create vectordb directory: %w", err)
	}

	ctx := context.Background()

	// Create persistent DB
	db, err := chromem.NewPersistentDB(dbPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create chromem db: %w", err)
	}

	// Get or create collection with embedding function
	collection, err := db.GetOrCreateCollection(CollectionName, nil, ef)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create collection: %w", err)
	}

	return &ChromemStore{
		db:         db,
		collection: collection,
		ctx:        ctx,
	}, nil
}

// AddDocument adds a document to the vector store
func (c *ChromemStore) AddDocument(id string, content string, metadata map[string]string) error {
	return c.collection.AddDocument(c.ctx, chromem.Document{
		ID:       id,
		Content:  content,
		Metadata: metadata,
	})
}

// AddDocuments adds multiple documents at once
func (c *ChromemStore) AddDocuments(docs []chromem.Document) error {
	return c.collection.AddDocuments(c.ctx, docs, 4) // Use 4 concurrent goroutines
}

// Search finds similar documents
func (c *ChromemStore) Search(query string, limit int) ([]SearchResult, error) {
	results, err := c.collection.Query(c.ctx, query, limit, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var searchResults []SearchResult
	for _, r := range results {
		searchResults = append(searchResults, SearchResult{
			ID:       r.ID,
			Content:  r.Content,
			Score:    r.Similarity,
			Metadata: r.Metadata,
		})
	}

	return searchResults, nil
}

// DeleteDocument removes a document from the store
func (c *ChromemStore) DeleteDocument(id string) error {
	return c.collection.Delete(c.ctx, nil, nil, id)
}

// Close is a no-op for chromem-go (persistent storage handles cleanup)
func (c *ChromemStore) Close() error {
	return nil
}

// GetCollection returns the underlying collection for advanced operations
func (c *ChromemStore) GetCollection() *chromem.Collection {
	return c.collection
}
