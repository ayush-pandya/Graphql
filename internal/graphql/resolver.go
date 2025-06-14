package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"database/sql"

	"github.com/ayush-pandya/Graphql/internal/clients"
)

// Database interface defines methods our resolvers need
type Database interface {
	CreateTicket(ctx context.Context, ticket *Ticket) error
	GetTicket(ctx context.Context, id string) (*Ticket, error)
	// Add other database methods as needed
}

// Resolver holds dependencies for GraphQL resolvers
type Resolver struct {
	db           *sql.DB
	ticketClient *clients.TicketClient
}

// NewResolver creates a new GraphQL resolver
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{
		db: db,
	}
}

// NewResolverWithGRPC creates a new GraphQL resolver with gRPC clients
func NewResolverWithGRPC(db *sql.DB, ticketClient *clients.TicketClient) *Resolver {
	return &Resolver{
		db:           db,
		ticketClient: ticketClient,
	}
}
