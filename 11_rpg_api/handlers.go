package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// handleListHeroes возвращает всех героев в JSON

// handleCreateHero принимает параметры из JSON (Тела запроса) и создает героя 1 уровня со 100 золота
func handleCreateHero(w http.ResponseWriter, r *http.Request) {
	// Подсказка: Как читать JSON из тела запроса (r.Body):
	//
	// type CreateRequest struct {
	//	 Name  string `json:"name"`
	//	 Class string `json:"class"`
	// }
	//
	// var input CreateRequest
	// err := json.NewDecoder(r.Body).Decode(&input)
	// if err != nil { ... }

	// Твой код здесь
	type CreateRequest struct {
		Name  string `json:"name"`
		Class string `json:"class"`
		HP    int    `json:"hp"`
		MP    int    `json:"mp"`
		Gold  int    `json:"gold"`
		Level int    `json:"level"`
		Dmg   int    `json:"dmg"`
	}
	var input CreateRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		fmt.Fprint(w, "Ошибка при создании героя")
		return
	}
	if _, ok := heroes[input.Name]; ok {
		fmt.Fprint(w, "Герой с таким именем уже существует")
		return
	}
	heroes[input.Name] = &Hero{
		Name:  input.Name,
		Class: input.Class,
		HP:    100,
		MP:    100,
		Gold:  100,
		Level: 1,
		Dmg:   10,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heroes[input.Name])

	fmt.Fprintf(w, "Герой %s создан\n", input.Name)
	saveHeroes()

	// Не забыть вызвать saveHeroes() после добавления!
}

// handleLevelUp повышает уровень на 1 и добавляет 500 золота
func handleLevelUp(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле Name не может быть пустым")
		return
	}
	if _, ok := heroes[name]; !ok {
		fmt.Fprint(w, "Герой с таким именем не найден")
		return
	}
	heroes[name].Level++
	heroes[name].Gold += 500
	saveHeroes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heroes[name])
	fmt.Fprintf(w, "Герой %s повышен до уровня %d\n", name, heroes[name].Level)

	// Твой код здесь
}

// handleTransferGold переводит золото от одного героя другому
func handleTransferGold(w http.ResponseWriter, r *http.Request) {
	loadHeroes()
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Fprint(w, "Количество золота должно быть числом")
		return
	}
	if from == "" || to == "" || count == 0 {
		fmt.Fprint(w, "Поля from, to и count не могут быть пустыми")
		return
	}
	if _, ok := heroes[from]; !ok {
		fmt.Fprint(w, "Герой с таким именем не найден")
		return
	}
	if _, ok := heroes[to]; !ok {
		fmt.Fprint(w, "Герой с таким именем не найден")
		return
	}
	heroes[from].Gold -= 100
	heroes[to].Gold += 100
	saveHeroes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heroes[from])
	json.NewEncoder(w).Encode(heroes[to])
	fmt.Fprintf(w, "Герой %s перевел %d золота герою %s\n", from, count, to)

	// Твой код здесь
}
func handleFight(w http.ResponseWriter, r *http.Request) {
	loadHeroes()
	// 1. Получаем ИМЕНА из ссылки
	name1 := r.URL.Query().Get("hero1")
	name2 := r.URL.Query().Get("hero2")

	if name1 == "" || name2 == "" {
		fmt.Fprint(w, "Поля hero1 и hero2 не могут быть пустыми")
		return
	}
	// 2. Достаем ГЕРОЕВ из мапы
	h1, ok1 := heroes[name1]
	h2, ok2 := heroes[name2]
	if !ok1 || !ok2 {
		fmt.Fprint(w, "Герой с таким именем не найден")
		return
	}
	h1.HP = 100
	h2.HP = 100

	// 3. БОЙ
	for h1.HP > 0 && h2.HP > 0 {
		// Пусть пока они бьют друг друга по очереди на 10 урона
		// Независимо от класса (чтобы сервер не завис)
		h2.HP -= h1.Dmg
		h1.HP -= h2.Dmg

		fmt.Fprintf(w, "%s бьет %s! У %s осталось %d HP\n", h1.Name, h2.Name, h2.Name, h2.HP)
	}
	// 4. КТО ПОБЕДИЛ? (Если HP ушло в минус, например -5, то проверяем кто выжил)
	if h1.HP > 0 {
		fmt.Fprintf(w, "\n🏆 Победил %s!\n", h1.Name)
	} else if h2.HP > 0 {
		fmt.Fprintf(w, "\n🏆 Победил %s!\n", h2.Name)
	} else {
		fmt.Fprintf(w, "\n💀 Оба погибли!\n")
	}
	saveHeroes()

}
func handleMarket(w http.ResponseWriter, r *http.Request) {
	loadHeroes()
	name := r.URL.Query().Get("name")
	if _, ok := heroes[name]; !ok {
		fmt.Fprint(w, "Герой с таким именем не найден")
		return
	}
	if heroes[name].Gold < 100 {
		fmt.Fprint(w, "У героя недостаточно золота")
		return
	}
	heroes[name].Gold -= 100
	heroes[name].Dmg += 2
	saveHeroes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heroes[name])
	fmt.Fprintf(w, "Герой %s купил меч и стал сильнее\n", name)

}
func handleListHeroes(w http.ResponseWriter, r *http.Request) {
	loadHeroes()
	list := []Hero{}
	for _, hero := range heroes {
		list = append(list, *hero)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}
