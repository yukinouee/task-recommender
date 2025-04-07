package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitializeDatabase(db *sql.DB) error {
	// テーブルを削除して再作成
	dropTableQuery := `DROP TABLE IF EXISTS tasks;`
	_, err := db.Exec(dropTableQuery)
	if err != nil {
		return err
	}

	// テーブルを新規作成
	createTableQuery := `
    CREATE TABLE tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        done BOOLEAN DEFAULT FALSE,
        priority INT,
        due_date TIMESTAMP,
        estimated_duration INT,
        created_at TIMESTAMP NOT NULL,
        completed_at TIMESTAMP
    );`

	_, err = db.Exec(createTableQuery)
	return err
}

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
