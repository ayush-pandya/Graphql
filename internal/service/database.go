package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

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

func Connect() (*sql.DB, error) {

	// Get database configuration
	host, port, user, password, dbname, sslmode := getDBConfig()

	log.Printf("Database settings: host=%s, port=%s, user=%s, dbname=%s",
		host, port, user, dbname)

	// Create connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	log.Println("Connection string", connStr)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	return DB, nil
}
