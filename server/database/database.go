package database

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4/pgxpool"
    "os"
)

// InitDB initializes the database connection pool.
func InitDB() (*pgxpool.Pool, error) {
  dbUser := os.Getenv("DB_USER")
  dbPassword := os.Getenv("DB_PASSWORD")
  dbName := os.Getenv("DB_NAME")
  dbHost := os.Getenv("DB_HOST")
  dbPort := os.Getenv("DB_PORT")

  connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

  dbpool, err := pgxpool.Connect(context.Background(), connStr)
  if err != nil {
    return nil, fmt.Errorf("unable to connect to database: %v", err)
  }

  return dbpool, nil
}
