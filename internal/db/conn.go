package db

import (
	"database/sql"

	"log"
	_ "github.com/lib/pq
)

var DB *sql.DB

func Connect() {

	var err error
	// DB, err = sql.Open("postgres", os.Getenv("DB_URL"))

	DB, err = sql.Open("postgres", "postgres://postgres:password@localhost:5433/itami?sslmode=disable")

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	log.Println("Connected to DB")
}
