package memory

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type Store interface {
	SaveMemory(ctx context.Context, userID string, content string, embedding []float32) error
	FindSimilar(ctx context.Context, userID string, embedding []float32, limit int) ([]MemoryRecord, error)
	Close() error
}

type MemoryRecord struct {
	ID        int
	UserID    string
	Content   string
	Embedding []float32
	CreatedAt time.Time
}

type pgStore struct {
	db *sql.DB
}

func NewPGStore(connStr string) (Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("memory: failed to open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("memory: failed to ping db: %w", err)
	}

	store := &pgStore{db: db}
	if err := store.initSchema(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *pgStore) initSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS vector;")
	if err != nil {
		return fmt.Errorf("memory: failed to create vector extension: %w", err)
	}

	// OpenAI text-embedding-3-small uses 1536 dimensions
	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS memories (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			embedding vector(1536),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("memory: failed to create memories table: %w", err)
	}

	return nil
}

func (s *pgStore) SaveMemory(ctx context.Context, userID string, content string, embedding []float32) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO memories (user_id, content, embedding) VALUES ($1, $2, $3)",
		userID, content, pgvector.NewVector(embedding),
	)
	if err != nil {
		return fmt.Errorf("memory: failed to save: %w", err)
	}
	slog.Debug("saved memory to pgvector")
	return nil
}

func (s *pgStore) FindSimilar(ctx context.Context, userID string, embedding []float32, limit int) ([]MemoryRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, content, created_at
		FROM memories
		WHERE user_id = $1
		ORDER BY embedding <=> $2
		LIMIT $3
	`, userID, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("memory: failed to query similar: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []MemoryRecord
	for rows.Next() {
		var rec MemoryRecord
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.Content, &rec.CreatedAt); err != nil {
			return nil, fmt.Errorf("memory: row scan error: %w", err)
		}
		results = append(results, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("memory: row iteration error: %w", err)
	}
	return results, nil
}

func (s *pgStore) Close() error {
	return s.db.Close()
}
