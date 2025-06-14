package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ayush-pandya/Graphql/internal/graphql"
	database "github.com/ayush-pandya/Graphql/internal/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

// createTestTicket creates a test ticket in the database
func createTestTicket(resolver *graphql.Resolver) (*graphql.Ticket, error) {
	log.Println("Attempting to create test ticket...")
	ticket, err := resolver.Mutation().CreateTicket(
		context.Background(),
		"Test Ticket",
		strPtr("Test Description"),
		nil,
		nil,
		"user-123",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %v", err)
	}

	log.Printf("✅ Ticket created successfully! ID: %s, Title: %s", ticket.ID, ticket.Title)
	return ticket, nil
}

// retrieveTicket gets a ticket from the database by ID
func retrieveTicket(resolver *graphql.Resolver, ticketID string) (*graphql.Ticket, error) {
	log.Printf("Attempting to retrieve ticket with ID: %s", ticketID)
	ticket, err := resolver.Query().Ticket(context.Background(), ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ticket: %v", err)
	}

	if ticket != nil {
		log.Printf("✅ Successfully retrieved ticket: %s", ticket.Title)
	} else {
		log.Println("❌ Ticket retrieval returned nil (not found)")
	}

	return ticket, nil
}

func main() {
	// Connect to database (optional - the resolvers work without it for now)
	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection failed (continuing without DB): %v", err)
		db = nil // Set to nil so server can still start
	} else {
		log.Println("✅ Database connection successful")
	}

	// Create resolver
	resolver := graphql.NewResolver(db)
	log.Println("✅ Resolver created")

	// Create GraphQL server
	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	log.Println("✅ Executable schema created")

	// Configure transports and features
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})
	log.Println("✅ GraphQL transports configured")

	// Setup routes
	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Println("🚀 GraphQL server starting on http://localhost:8080")
	log.Println("📊 GraphQL Playground available at http://localhost:8080")
	log.Println("🔍 GraphQL endpoint at http://localhost:8080/query")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
