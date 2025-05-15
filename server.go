package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"Graphql/db"
	"Graphql/graphql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

// getDBConfig returns database configuration parameters
func getDBConfig() (string, string, string, string, string, string) {
	host := "localhost"
	port := "5432"
	user := "ayushpandya"
	password := "postgres"
	dbname := "ticketdb"
	sslmode := "disable"

	return host, port, user, password, dbname, sslmode
}

// setupDatabase establishes a connection to the database and initializes it
func setupDatabase(connStr string) (*db.TicketStore, error) {
	// Try to establish DB connection
	store, err := db.NewTicketStore(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket store: %v", err)
	}
	log.Println("✅ Successfully connected to PostgreSQL database")

	// Initialize database schema
	log.Println("Initializing database schema...")
	if err := store.InitDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}
	log.Println("✅ Database schema initialized")

	// Check DB connectivity by pinging
	if err := store.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}
	log.Println("✅ Database ping successful")

	return store, nil
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
	// Get database configuration
	host, port, user, password, dbname, sslmode := getDBConfig()

	log.Printf("Database settings: host=%s, port=%s, user=%s, dbname=%s",
		host, port, user, dbname)

	// Create connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	// Setup database
	store, err := setupDatabase(connStr)
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}

	// Create resolver
	resolver := graphql.NewResolver(store)

	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	// Add these lines to configure transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		log.Printf("Request body: %s", string(body))
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		srv.ServeHTTP(w, r)
	})
	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// Test ticket creation
	ticket, err := createTestTicket(resolver)
	if err != nil {
		log.Fatalf("Test ticket creation failed: %v", err)
	}

	// Test ticket retrieval
	_, err = retrieveTicket(resolver, ticket.ID)
	if err != nil {
		log.Fatalf("Test ticket retrieval failed: %v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tickets" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Simulate a delay
	time.Sleep(2 * time.Second)
	// Return a simple JSON response
	fmt.Fprintf(w, `{"message": "Hello, World!"}`)
	// Return a simple text response
	// fmt.Fprintf(w, `{"message": "Hello, World!"}`)
	fmt.Fprintf(w, "Hello!")
}
