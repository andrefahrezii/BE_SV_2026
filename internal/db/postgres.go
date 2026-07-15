package db

import (
	"database/sql"
	"fmt"

	"github.com/you/sharing-vision-backend-v2/internal/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	c := config.Conf
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name, c.DB.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	config.Log.Info("database connected")
	return db, nil
}
