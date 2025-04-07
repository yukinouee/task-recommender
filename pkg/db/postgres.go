package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "require" // 本番環境ではrequireを使用
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		host, user, password, dbname, sslmode)

	return sql.Open("postgres", connStr)
}
