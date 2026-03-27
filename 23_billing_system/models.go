package main

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Accounts struct {
	ID      int  `json:"id"`
	UserID  *int `json:"user_id"`
	Balance int  `json:"balance"`
}
type Card struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Number    string `json:"number"`
}
type Transaction struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Amount    int    `json:"amount"`
	Type      string `json:"type"`
}
