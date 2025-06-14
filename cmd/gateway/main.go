package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayush-pandya/Graphql/internal/clients"
	"github.com/ayush-pandya/Graphql/internal/graphql"
	database "github.com/ayush-pandya/Graphql/internal/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	log.Println("🚀 Starting GraphQL Gateway Server...")

	// Connect to database (optional)
	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection failed (continuing without DB): %v", err)
		db = nil
	} else {
		log.Println("✅ Database connection successful")
	}

	// Connect to gRPC microservices
	ticketServiceURL := getEnv("TICKET_SERVICE_URL", "localhost:50051")
	log.Printf("🔌 Connecting to Ticket Service at %s", ticketServiceURL)

	ticketClient, err := clients.NewTicketClient(ticketServiceURL)
	if err != nil {
		log.Printf("❌ Failed to connect to ticket service: %v", err)
		log.Println("⚠️  Continuing without ticket service - some features may not work")
		ticketClient = nil
	} else {
		log.Println("✅ Connected to Ticket Service via gRPC")
	}

	// Ensure we close the gRPC connection
	if ticketClient != nil {
		defer func() {
			if err := ticketClient.Close(); err != nil {
				log.Printf("Error closing ticket client: %v", err)
			}
		}()
	}

	// Create GraphQL resolver with gRPC clients
	resolver := graphql.NewResolverWithGRPC(db, ticketClient)
	log.Println("✅ GraphQL Resolver created with gRPC clients")

	// Create GraphQL server
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	log.Println("✅ GraphQL Schema created")

	// Configure server
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})
	log.Println("✅ GraphQL Server configured")

	// Setup HTTP routes
	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL Gateway", "/query"))

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		log.Println("🌐 GraphQL Gateway Server starting on http://localhost:8080")
		log.Println("📊 GraphQL Playground available at http://localhost:8080")
		log.Println("🔍 GraphQL API endpoint at http://localhost:8080/query")
		log.Println("🔄 Gateway communicates with microservices via gRPC")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down GraphQL Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("👋 GraphQL Gateway stopped")
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
