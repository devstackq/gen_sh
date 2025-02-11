package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDB() (*sql.DB, error) {
	connStr := "postgres://user:password@localhost/dbname?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
		return nil, err
	}

	return db, nil
}
