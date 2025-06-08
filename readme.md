# Ticket Management System (GraphQL + Go + PostgreSQL)

A backend microservice for managing tickets, built with GraphQL, Go, and PostgreSQL, and containerized using Docker.

---

## Table of Contents
- [Technologies Used](#technologies-used)
- [Setup Instructions](#setup-instructions)
- [Application Structure](#application-structure)
- [Resolvers and PostgreSQL Integration](#resolvers-and-postgresql-integration)
- [Docker Integration](#docker-integration)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

---

## Technologies Used
- **Go**: Backend language
- **GraphQL**: API query language (using gqlgen)
- **PostgreSQL**: Database
- **Docker**: Containerization
- **gqlgen**: Go GraphQL server library
- **pq**: PostgreSQL driver for Go

---

## Setup Instructions

### 1. Clone the Repository
```bash
git clone <your-repo-url>
cd <your-project>
```

### 2. Set Up PostgreSQL with Docker
Create a `docker-compose.yml` file (if not present):

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: ayushpandya
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ticketdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Start PostgreSQL:
```bash
docker-compose up -d postgres
```

### 3. Initialize the Go Application
Install dependencies:
```bash
go get github.com/99designs/gqlgen github.com/lib/pq
```
Generate GraphQL boilerplate code (if using gqlgen):
```bash
go run github.com/99designs/gqlgen generate
```

### 4. Run the Application
```bash
go run main.go
```

---

## Application Structure
```
.
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── database/
│   └── database.go      # PostgreSQL connection logic
├── graph/
│   ├── generated/       # Auto-generated GraphQL code
│   ├── model/           # GraphQL schema models
│   └── resolvers.go     # Resolver implementations
└── schema.graphqls      # GraphQL schema definitions
```

---

## Resolvers and PostgreSQL Integration

### What Are Resolvers?
Resolvers are functions that fetch or modify data for GraphQL queries/mutations. They interact with PostgreSQL to retrieve or update ticket data.

#### Example Resolver
```go
// graph/resolvers.go
func (r *Resolver) Ticket(ctx context.Context, id string) (*model.Ticket, error) {
    var ticket model.Ticket
    err := r.DB.QueryRow(`
        SELECT id, title, description 
        FROM tickets 
        WHERE id = $1
    `, id).Scan(&ticket.ID, &ticket.Title, &ticket.Description)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch ticket: %v", err)
    }
    return &ticket, nil
}
```

#### Connecting to PostgreSQL
Example connection string:
```go
connStr := "user=ayushpandya password=postgres dbname=ticketdb host=localhost sslmode=disable"
```
Initialize the database in `database/database.go`:
```go
func Connect() error {
    DB, err = sql.Open("postgres", connStr)
    // Handle error and ping the database
}
```

---

## Docker Integration

- **PostgreSQL in Docker**: Runs in an isolated container with persistent storage.
- **Networking**: Use `host: postgres` (not `localhost`) if your Go app is also in Docker.

Example `docker-compose.yml` snippet for the Go app:
```yaml
your_go_app:
  build: .
  ports:
    - "8080:8080"
  depends_on:
    - postgres
```

---

## Troubleshooting

- **PostgreSQL Connection Refused**:
  - Ensure Docker is running and the PostgreSQL container is up (`docker ps`).
  - Verify credentials in the connection string match `docker-compose.yml`.
- **Missing Dependencies**:
  - Run `go mod tidy` to install missing packages.
- **Docker Port Conflicts**:
  - Stop other PostgreSQL instances: `sudo service postgresql stop` (if running locally).

---

## Next Steps
- Add authentication (JWT or OAuth)
- Write unit/integration tests for resolvers
- Extend the schema for user management or comments
- Deploy using Kubernetes or a cloud provider (AWS, GCP, etc.)
