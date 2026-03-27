package main

import (
	"fmt"
)

type BMW struct{}
type Auto interface {
	StepOnGas()
	StepOnBracke()
}

func (b *BMW) StepOnGas() {
	fmt.Println("Я БМВ! Жми на газ!")
}

type Zhiga struct{}

func (z *Zhiga) StepOnGas() {
	fmt.Println("Я Жига! Пробую не развалиться!")
}
func (z *Zhiga) StepOnBracke() {
	fmt.Println("Я Жига! Чуть тормоза не отвалились!")
}

type Mazda struct{}

func (m *Mazda) StepOnGas() {

	fmt.Println("Я Мазда! Жми на газ!")
}
func (m *Mazda) StepOnBracke() {
	fmt.Println("Я Мазда! Жму на тормоз!")
}
func (m *Mazda) BipBip() {
	fmt.Println("Бип-бип!")
}

func ride(a Auto) {
	fmt.Println("Я водитель!")
	fmt.Println("Я сажусь в свою машину!")
	fmt.Println("Жму на газ!")

}

func main() {
	fmt.Println("Тренировочный полигон для Ошибок (Error Handling) готов!")
	// Место для примеров из видео
	//bmw := &BMW{}
	//zhiga := &Zhiga{}
	mazda := &Mazda{}
	//ride(bmw)
	//ride(zhiga)
	ride(mazda)

}
