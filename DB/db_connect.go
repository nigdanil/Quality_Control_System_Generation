package DB

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // SQLite драйвер
)

// OpenDatabase открывает соединение с SQLite-базой.
func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping к БД: %v", err)
	}

	return db, nil
}
