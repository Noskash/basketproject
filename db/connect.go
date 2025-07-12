package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect_to_database() (*sql.DB, error) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal(err)
	}

	connectStr := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=%s",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSLMODE"),
	)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных")
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS matches(
			id SERIAL PRIMARY KEY,
			game_id VARCHAR(255) NOT NULL,
			url VARCHAR(255) NOT NULL,
			count VARCHAR(255) NOT NULL 
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать таблицу files: %v", err)
	}
	return db, nil
}
