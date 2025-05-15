package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"Graphql/graphql"

	_ "github.com/lib/pq"
)

// Add this method to TicketStore
func (s *TicketStore) Ping() error {
	return s.db.Ping()
}

type TicketStore struct {
	db *sql.DB
}

// InitDatabase creates the necessary tables if they don't exist
func (s *TicketStore) InitDatabase() error {
	// Create tickets table
	query := `
    CREATE TABLE IF NOT EXISTS tickets (
        id VARCHAR(36) PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        status VARCHAR(20) NOT NULL,
        priority VARCHAR(20) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        assignee_id VARCHAR(36),
        reporter_id VARCHAR(36) NOT NULL,
        tags JSONB
    );`

	_, err := s.db.Exec(query)
	return err
}

func NewTicketStore(connStr string) (*TicketStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &TicketStore{db: db}, nil
}

func (s *TicketStore) CreateTicket(ctx context.Context, ticket *graphql.Ticket) error {
	query := `
		INSERT INTO tickets (
			id, title, description, status, priority,
			created_at, updated_at, assignee_id, reporter_id, tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// Convert tags to JSON for storage
	tagsJSON, err := json.Marshal(ticket.Tags)
	if err != nil {
		return err
	}

	// Parse timestamps
	createdAt, _ := time.Parse(time.RFC3339, ticket.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, ticket.UpdatedAt)

	var assigneeID *string
	if ticket.Assignee != nil {
		assigneeID = &ticket.Assignee.ID
	}

	_, err = s.db.ExecContext(ctx, query,
		ticket.ID,
		ticket.Title,
		ticket.Description,
		ticket.Status,
		ticket.Priority,
		createdAt,
		updatedAt,
		assigneeID,
		ticket.Reporter.ID,
		tagsJSON,
	)
	return err
}

func (s *TicketStore) GetTicket(ctx context.Context, id string) (*graphql.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, created_at, updated_at, assignee_id, reporter_id, tags
		FROM tickets
		WHERE id = $1
	`

	var ticket graphql.Ticket
	var description sql.NullString
	var assigneeID sql.NullString
	var reporterID string
	var createdAt, updatedAt time.Time
	var tagsJSON []byte

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID,
		&ticket.Title,
		&description,
		&ticket.Status,
		&ticket.Priority,
		&createdAt,
		&updatedAt,
		&assigneeID,
		&reporterID,
		&tagsJSON,
	)

	if err != nil {
		return nil, err
	}

	// Set optional description if provided
	if description.Valid {
		ticket.Description = &description.String
	}

	// Convert timestamps to string format
	ticket.CreatedAt = createdAt.Format(time.RFC3339)
	ticket.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Set reporter
	ticket.Reporter = &graphql.User{ID: reporterID}

	// Set assignee if exists
	if assigneeID.Valid {
		ticket.Assignee = &graphql.User{ID: assigneeID.String}
	}

	// Parse tags JSON if exists
	if tagsJSON != nil {
		var tagStrings []string
		if err := json.Unmarshal(tagsJSON, &tagStrings); err != nil {
			return nil, err
		}

		tags := make([]*string, len(tagStrings))
		for i, tag := range tagStrings {
			tagCopy := tag
			tags[i] = &tagCopy
		}
		ticket.Tags = tags
	}

	return &ticket, nil
}
