package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	AccountNo int64  `json:"accountnumber"`
	Balance   int64  `json:"balance"`
}

func NewAccount(fn, ln string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: fn,
		LastName:  ln,
		AccountNo: int64(rand.Intn(1000000)),
	}
}
