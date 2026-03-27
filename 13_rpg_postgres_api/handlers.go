package main

// Здесь живут обработчики (handleCreate, handleMarket, handleListHeroes)

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var input CreateHeroRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal("Не удалось прочитать героя из БД", err)
		return
	}
	_, err = DB.Exec("INSERT INTO heroes (name, class) VALUES ($1, $2)", input.Name, input.Class)
	if err != nil {
		log.Fatal("Не удалось добавить героя", err)
		return
	}

	fmt.Fprintf(w, "Герой %s успешно создан в БД!", input.Name)

}
func handleListHeroes(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM heroes")
	if err != nil {
		log.Fatal("Не удалось получить список героев", err)
		return
	}
	defer rows.Close()
	list := []Hero{}
	for rows.Next() {
		var hero Hero
		err := rows.Scan(&hero.ID, &hero.Name, &hero.Class, &hero.Level, &hero.Gold, &hero.Dmg, &hero.HP)
		if err != nil {
			log.Fatal("Не удалось прочитать героя", err)
			return
		}
		list = append(list, hero)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
func handleLevelUp(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле Name не может быть пустым")
		return
	}

	DB.QueryRow("SELECT * FROM heroes WHERE name = $1", name)
	_, err := DB.Exec("UPDATE heroes SET level= level+1, gold= gold+500 WHERE name = $1", name)
	if err != nil {
		log.Fatal("Не удалось повысить героя", err)
		return
	}

	fmt.Fprintf(w, "Герой %s успешно повышен в БД!", name)

}

func handleMarket(w http.ResponseWriter, r *http.Request) {
	var correntGold int

	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле Name не может быть пустым")
		return
	}

	DB.QueryRow("SELECT gold FROM heroes WHERE name = $1", name).Scan(&correntGold)
	if correntGold < 100 {
		fmt.Fprint(w, "У героя недостаточно золота")
		return
	}
	_, err := DB.Exec("UPDATE heroes SET gold= gold-100, dmg= dmg+2 WHERE name = $1", name)
	if err != nil {
		log.Fatal("Не удалось купить предмет", err)
		return
	}

	fmt.Fprintf(w, "Герой %s успешно купил предмет в БД!", name)

}
func handleFight(w http.ResponseWriter, r *http.Request) {
	var input FightRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal("Не удалось прочитать героя из БД", err)
		return
	}
	var dmg1, hp1 int
	var dmg2, hp2 int

	// 1. ОДИН РАЗ ДОСТАЕМ ВСЕ СТАТЫ ДО БОЯ
	// У QueryRow можно сканировать сразу несколько колонок через запятую!
	err = DB.QueryRow("SELECT dmg, hp FROM heroes WHERE name = $1", input.Hero1).Scan(&dmg1, &hp1)
	if err != nil { /* обработка */
	}

	err = DB.QueryRow("SELECT dmg, hp FROM heroes WHERE name = $1", input.Hero2).Scan(&dmg2, &hp2)
	if err != nil { /* обработка */
	}
	// 2. БОЙ ПРОИСХОДИТ ТОЛЬКО В ПАМЯТИ GO
	for hp1 > 0 && hp2 > 0 {
		hp1 -= dmg2
		hp2 -= dmg1
		fmt.Fprintf(w, "%s ударил %s. У %s осталось %d HP\n", input.Hero1, input.Hero2, input.Hero2, hp2)
		fmt.Fprintf(w, "%s ударил %s. У %s осталось %d HP\n", input.Hero2, input.Hero1, input.Hero1, hp1)
	}
	// 3. ПРОВЕРЯЕМ КТО ВЫЖИЛ ПОСЛЕ ЦИКЛА
	if hp1 > 0 {
		fmt.Fprintf(w, "\n🏆 Победил %s!\n", input.Hero1)
	} else if hp2 > 0 {
		fmt.Fprintf(w, "\n🏆 Победил %s!\n", input.Hero2)
	} else {
		fmt.Fprintf(w, "\n💀 Оба погибли!\n")
	}
	// 4. ОДИН РАЗ СОХРАНЯЕМ ФИНАЛЬНОЕ ЗДОРОВЬЕ В БАЗУ!
	DB.Exec("UPDATE heroes SET hp = $1 WHERE name = $2", hp1, input.Hero1)
	DB.Exec("UPDATE heroes SET hp = $1 WHERE name = $2", hp2, input.Hero2)
}
