package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Не удалось подключиться к БД", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД", err)
	}

	fmt.Println("Успешно подключено к PostgreSQL!")

}
