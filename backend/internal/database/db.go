package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost:5432/grub_exchange?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
