package main

import "fmt"

type Blacklist map[string]bool

func IsAllowed(name string, list Blacklist) bool {
	_, ok := list[name]
	if ok {
		return false
	}
	return true

}

func main() {
	myList := Blacklist{
		"Gimli": true,
		"Orc":   true,
	}
	yesOrNo := IsAllowed("Gimli", myList)
	fmt.Println("Можно ли войти Gimli? :", yesOrNo)
	yesOrNoCoder := IsAllowed("LiveCoder", myList)
	fmt.Println("Можно ли войти LiveCoder? :", yesOrNoCoder)
}
