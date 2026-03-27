package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleBuyItem(w http.ResponseWriter, r *http.Request) {
	// 1. Декодируем BuyRequest из JSON Body
	// 2. Ищем героя (проверяем хватает ли золота)
	// 3. Ищем предмет (проверяем существует ли он в таблице items)
	// 4. Снимаем золото у героя (UPDATE heroes)
	// 5. Записываем покупку в инвентарь (INSERT INTO inventory)
	var input BuyRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal("Не удалось прочитать героя из БД", err)
		return
	}
	var hero Hero
	var item Item
	err = DB.QueryRow("SELECT * FROM heroes WHERE name = $1", input.HeroName).Scan(&hero.ID, &hero.Name, &hero.Class, &hero.Level, &hero.Gold, &hero.Dmg, &hero.HP)
	if err != nil {
		log.Fatal("Не удалось найти героя", err)
		return
	}
	err = DB.QueryRow("SELECT * FROM items WHERE name = $1", input.ItemName).Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		log.Fatal("Не удалось найти предмет", err)
		return
	}
	if hero.Gold < item.Price {
		fmt.Fprint(w, "У героя недостаточно золота")
		return
	}
	_, err = DB.Exec("UPDATE heroes SET gold = gold - $1 WHERE id = $2", item.Price, hero.ID)
	if err != nil {
		log.Fatal("Не удалось купить предмет", err)
		return
	}
	_, err = DB.Exec("INSERT INTO inventory (hero_id, item_id) VALUES ($1, $2)", hero.ID, item.ID)
	if err != nil {
		log.Fatal("Не удалось добавить предмет в инвентарь", err)
		return
	}
	fmt.Fprintf(w, "Предмет %s успешно куплен героем %s!", item.Name, hero.Name)

}

func handleInventory(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем имя героя из ссылки (например: /inventory?hero=Арагорн)
	heroName := r.URL.Query().Get("hero")
	if heroName == "" {
		fmt.Fprint(w, "Укажите имя героя (параметр ?hero=)")
		return
	}
	// 2. Самая главная магия — запрос с двумя склейками (JOIN)
	// Сделай так:
	// "Дай мне id, имя и цену..."
	// "...ИСХОДЯ ИЗ таблицы предметов (items)..."
	// "...ПРИКЛЕИВ к ней таблицу инвентаря там, где ID предмета = вещи в инвентаре..."
	// "...И ПРИКЛЕИВ к этому таблицу героев там, где ID владельца = ID героя..."
	// "...НО ТОЛЬКО ДЛЯ ТОГО героя, чье имя мы передали в $1"
	query := `
		SELECT items.id, items.name, items.price
		FROM items
		JOIN inventory ON inventory.item_id = items.id
		JOIN heroes ON inventory.hero_id = heroes.id
		WHERE heroes.name = $1
	`

	rows, err := DB.Query(query, heroName)
	if err != nil {
		log.Fatal("Ошибка при запросе инвентаря:", err)
		return
	}
	// Обязательно закрываем rows, чтобы не было утечки памяти
	defer rows.Close()
	// 3. Создаем пустой "ящик" (слайс) для наших предметов
	var list []Item
	// 4. Начинаем крутить цикл по всем строчкам, которые вернула база
	for rows.Next() {
		var item Item

		// Забираем три колонки из базы и кладем в нашу структуру
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			log.Fatal("Ошибка чтения строки из БД:", err)
		}

		// Добавляем заполненный предмет в наш "ящик"
		list = append(list, item)
	}
	// 5. Упаковываем весь наш ящик со структурами в JSON и отдаем пользователю!
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
