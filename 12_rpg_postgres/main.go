package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Этот импорт с нижним подчеркиванием нужен, чтобы просто загрузить драйвер Postgres
)

func main() {
	// Твоя строка подключения (замени на свои данные!)
	// connStr := "user=пользователь password=пароль dbname=имя_базы sslmode=disable"

	// db, err := sql.Open("postgres", connStr)
	// if err != nil { log.Fatal(err) }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil { log.Fatal("Не удалось подключиться к БД:", err) }

	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	fmt.Println("Успешно подключено к PostgreSQL!")
	_, err = db.Exec("INSERT INTO heroes (name, class) VALUES ('Орк', 'Танк')")
	if err != nil {
		log.Fatal("Не удалось добавить героя:", err)
	}

	fmt.Println("Герой добавлен!")
	rows, err := db.Query("SELECT id, name, class, level, gold, dmg FROM heroes")
	if err != nil {
		log.Fatal("Не удалось извлечь героя из БД", err)
	}
	defer db.Close()
	for rows.Next() {
		var id, level, gold, dmg int
		var name, class string
		err := rows.Scan(&id, &name, &class, &level, &gold, &dmg)
		if err != nil {
			log.Fatal("Не удалось прочитать героя из БД", err)
		}
		fmt.Printf("ID: %d, Имя: %s, Класс: %s, Уровень: %d, Золото: %d, Урон: %d\n", id, name, class, level, gold, dmg)
	}
}
