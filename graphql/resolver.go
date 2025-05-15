package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
)

// Database interface defines methods our resolvers need
type Database interface {
	CreateTicket(ctx context.Context, ticket *Ticket) error
	GetTicket(ctx context.Context, id string) (*Ticket, error)
	// Add other database methods as needed
}

// Resolver holds dependencies that will be used across resolvers
type Resolver struct {
	db Database
}

// NewResolver creates a new Resolver with the required dependencies
func NewResolver(db Database) *Resolver {
	return &Resolver{
		db: db,
	}
}
