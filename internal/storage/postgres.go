package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/shoksin/quotes-service/configs"
	"log"
)

func NewPostgresConnection(databaseConfig configs.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", databaseConfig.Host, databaseConfig.Port, databaseConfig.User, databaseConfig.Password, databaseConfig.Name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Printf("successfully connected to database %s", databaseConfig.Name)
	return db, nil
}
