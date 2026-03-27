package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	// ВНИМАНИЕ: Если ты назовешь свою базу в PgAdmin по-другому, не забудь исправить тут "arcade_db"
	connStr := "postgres://postgres:postgres@localhost/Company?sslmode=disable"

	DB, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/billing?sslmode=disable") // Я специально оставил ошибку, чтобы ты потом исправил на нужную базу! (arcade_db)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе", err)
	}

	fmt.Println("Успешно подключено к PostgreSQL!")
}
