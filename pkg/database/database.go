package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("Lỗi mở kết nối : %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB Không phản hồi: %w", err)
	}
	return db, nil
}
