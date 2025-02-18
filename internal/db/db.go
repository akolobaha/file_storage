package db

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

var DB *sqlx.DB

func Init() error {

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s  sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	)

	DB, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Println("Error opening database connection:", err)
		return err
	}
	err = DB.Ping()
	if err != nil {
		log.Println("Error pinging database connection:", err)
		return err
	}

	return nil
}
