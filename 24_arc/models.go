package main

type Users struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Bank struct {
	ID      int  `json:"id"`
	UserID  *int `json:"user_id"`
	Balance int  `json:"balance"`
}

type Cards struct {
	ID     int    `json:"id"`
	BankID int    `json:"bank_id"`
	Number string `json:"number"`
}
type Logs struct {
	ID     int    `json:"id"`
	BankID int    `json:"bank_id"`
	Amount int    `json:"amount"`
	Info   string `json:"info"`
}
