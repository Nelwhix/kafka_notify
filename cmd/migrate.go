package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://nelwhix:admin@localhost:5432/kafka_notify")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id CHAR(26) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = conn.Exec(context.Background(), createUsersTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute query: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migrations ran successfully!")
}
