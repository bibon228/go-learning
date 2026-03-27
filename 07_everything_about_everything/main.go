package main

import "fmt"

func main() { /*
		score := 15     //перемменная int
		if score > 15 { // больше 15?
			fmt.Println("Мега хорош")

		} else if score > 10 { //иначе если не больше 15 , может больше 10?
			fmt.Println("Просто хорош")

		} else if score > 5 { // ну раз мееньше 15 и меньше 10 тогда ты лоутаб
			fmt.Println("Лоутаб")

		} else { // пиздец ты нищенка у тебя меньше 15 , меньше 10 и меньше 5, ты нуб!
			fmt.Println("Нуб")
		}
	*/
	/*
		YaKrasavcheg := true  // булевая переменная

		if YaKrasavcheg {
			fmt.Println("Да, ты Krasavcheg")    // true
		} else {
			fmt.Println("Нет ты не Krasavcheg!!!")   //false
		}
	*/

	number := 15
	ravno5 := number == 5   // равно 5, если да TRUE, если нет FALSE
	bolshe12 := number > 12 // больше 12, если да TRUE, если нет FALSE
	if ravno5 {
		fmt.Println("number равно пяти")
	}
	if bolshe12 {
		fmt.Println("number больше 12")
	}
}
