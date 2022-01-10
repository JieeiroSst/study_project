package main

import (
	"fmt"
	"sync"
	"time"
)

type (
	Money int64 // cents
	Account struct {
		bal Money
		mutex sync.Mutex
	}
)


const Dollars Money = 100  // 100 cents to the dollar

// NewAccount creates a new account with bonus $100.
func NewAccount() *Account {
	return &Account{bal: 100 * Dollars}
}

// Deposit adds money to an account.
func (ac *Account) Deposit(amt Money) {
	ac.mutex.Lock()
	defer ac.mutext.Unlock()

	current := ac.bal
	time.Sleep(1*time.Millisecond)
	ac.bal = current + amt
}

// Balance returns funds available.
func (ac *Account) Balance() Money {
	ac.mutex.Lock()
	defer ac.mutext.Unlock()

	return ac.bal
}

func main() {
	ac := NewAccount()
	go ac.Deposit(1000 * Dollars)
	go ac.Deposit(200 * Dollars)
	time.Sleep(100*time.Millisecond)
	fmt.Printf("Balance: $%2.2f\n", ac.Balance()/100.0)
}