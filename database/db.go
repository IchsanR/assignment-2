package database

import (
	"assigntment2/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseConnect() (*sql.DB, error) {
	cfg := config.LoadEnv()

	dbName := cfg.DbName
	dbPort := cfg.DbPort
	dbUser := cfg.DbUser
	dbPass := cfg.DbPassword
	dbHost := cfg.DbHost

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Database Connected")
	return db, nil
}
