# gRPC Microservice with PostgreSQL

A complete gRPC microservice that manages tickets with PostgreSQL database backend.

## Architecture

```
Client (gRPC) → Ticket Microservice → PostgreSQL Database
```

## Features

- ✅ **gRPC API** - High-performance binary protocol
- ✅ **PostgreSQL Integration** - Reliable ACID database
- ✅ **CRUD Operations** - Create, Read, Update, Delete tickets
- ✅ **Connection Pooling** - Optimized database connections
- ✅ **Docker Support** - Easy deployment with containers
- ✅ **Environment Configuration** - Flexible configuration
- ✅ **Graceful Shutdown** - Clean resource cleanup

## Quick Start

### 1. Start with Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up --build

# View logs
docker-compose logs -f ticket-service
```

### 2. Manual Setup

**Prerequisites:**
- Go 1.21+
- PostgreSQL 15+
- Protocol Buffers compiler

**Steps:**

```bash
# 1. Install dependencies
go mod tidy

# 2. Start PostgreSQL
docker run -d \
  --name postgres-tickets \
  -e POSTGRES_DB=tickets \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:15-alpine

# 3. Run database migration
psql -h localhost -U postgres -d tickets -f migrations/001_create_tickets_table.sql

# 4. Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=tickets

# 5. Start the microservice
cd cmd/ticket-service-db
go run main.go
```

## API Usage

### Using grpcurl

```bash
# List all tickets
grpcurl -plaintext localhost:50051 ticket.TicketService/ListTickets

# Get specific ticket
grpcurl -plaintext \
  -d '{"id": "550e8400-e29b-41d4-a716-446655440001"}' \
  localhost:50051 ticket.TicketService/GetTicket

# Create new ticket
grpcurl -plaintext \
  -d '{
    "title": "New Bug Report",
    "description": "Something is broken",
    "priority": "TICKET_PRIORITY_HIGH",
    "assignee_id": "user-123",
    "tags": ["bug", "urgent"]
  }' \
  localhost:50051 ticket.TicketService/CreateTicket

# Update ticket
grpcurl -plaintext \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "Updated Title",
    "status": "TICKET_STATUS_IN_PROGRESS"
  }' \
  localhost:50051 ticket.TicketService/UpdateTicket

# Delete ticket
grpcurl -plaintext \
  -d '{"id": "550e8400-e29b-41d4-a716-446655440001"}' \
  localhost:50051 ticket.TicketService/DeleteTicket
```

### Using Go Client

```go
package main

import (
    "context"
    "log"

    ticketpb "github.com/ayush-pandya/Graphql/proto/ticket"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Connect to gRPC server
    conn, err := grpc.NewClient("localhost:50051", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := ticketpb.NewTicketServiceClient(conn)

    // Create ticket
    resp, err := client.CreateTicket(context.Background(), &ticketpb.CreateTicketRequest{
        Title:       "New Ticket",
        Description: "Ticket description",
        Priority:    ticketpb.TicketPriority_TICKET_PRIORITY_HIGH,
        AssigneeId:  "user-123",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created ticket: %s", resp.Ticket.Id)
}
```

## Database Schema

```sql
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',
    assignee_id VARCHAR(100),
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | PostgreSQL username |
| `DB_PASSWORD` | password | PostgreSQL password |
| `DB_NAME` | tickets | Database name |
| `DB_SSLMODE` | disable | SSL mode for connection |
| `GRPC_PORT` | 50051 | gRPC server port |

## Development

```bash
# Generate protobuf code
protoc --go_out=. --go-grpc_out=. proto/ticket/ticket.proto

# Run tests
go test ./...

# Build binary
go build -o ticket-service ./cmd/ticket-service-db

# View database
docker exec -it ticket-postgres psql -U postgres -d tickets
```

## Production Deployment

1. Use proper SSL certificates
2. Set strong passwords
3. Enable SSL mode for database
4. Use environment-specific configurations
5. Set up monitoring and logging
6. Configure proper resource limits

## Troubleshooting

**Connection Issues:**
- Ensure PostgreSQL is running
- Check firewall settings
- Verify environment variables

**Database Issues:**
- Run migrations manually
- Check PostgreSQL logs
- Verify user permissions

**gRPC Issues:**
- Use `grpc_health_probe` for health checks
- Check port availability
- Verify proto definitions 