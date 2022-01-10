package main

import (
	"fmt"
	"time"
)

type (
	Account struct {
		bal Money
		ch chan<- func()
	}
)

func NewAccount() *Account {
	ch := make(chan func())
	go func() {
		for f := range ch { f() }
	}()
	return &Account{bal: 100*Dollars, ch: ch}
}

// Deposit adds money to an account.
func (ac *Account) Deposit(amt Money) {
	ac.ch <- func() {
		current := ac.bal
		time.Sleep(1*time.Millisecond)
		ac.bal = current + amt
	}
}

// Adds transfers money to/from an account.
func (ac *Account) Add(amt Money, callback func(error)) {
	ac.ch <- func() {
		if ac.bal + amt < 0 {
			callback(fmt.Errorf("insuff. funds %v for w/d %v",
				ac.bal, amt))
			return
		}
		ac.bal += amt
		callback(nil)   // successful transfer
	}
}

// Balance provides funds available.
func (ac *Account) Balance(callback func(Money)) {
	ac.ch <- func() {
		callback(ac.bal)
	}
}

func main() {
	ac := NewAccount()
	ac.Add(1000 * Dollars, func(error) {} )
	ac.Add(200 * Dollars, func(error) {} )
	ac.Add(-1e6 * Dollars, func(err error) {
		if err != nil { fmt.Println(err) }
	})
	ac.Balance(func(bal Money) {
		fmt.Printf("Balance: $%v\n", bal/100)
	})
	time.Sleep(100*time.Millisecond)
}

func (ac *Account) TransferTo(to *Account, amt Money,
	callback func(error)) {
	ac.ch <- func() {
		if amt > ac.bal {
			callback(fmt.Errorf("Insuff. funds %v for tfr %v",
				ac.bal, amt))
			return
		}
		ac.bal -= amt
		to.Add(amt, callback)
	}
}

func (ac *Account) TransferTo(to *Account, amt Money,
	callback func(error)) {
	ac.ch <- func() {
		if amt > ac.bal {
			callback(fmt.Errorf("Insuff. funds %v for tfr %v",
				ac.bal, amt))
			return
		} else if amt < 0 && -amt > to.bal {
			callback(fmt.Errorf("Insuff. funds %v for tfr %v",
				to.bal, -amt))
			return
		}
		ac.bal -= amt
		to.Add(amt, callback)
	}
}

func (ac *Account) TransferTo(to *Account, amt Money,
	callback func(error)) {
	ac.ch <- func() {
		if amt < 0 {
			to.TransferTo(ac, -amt, callback)
			return
		}
		if amt > ac.bal {
			callback(fmt.Errorf("Insuff. funds %v for tfr %v",
				ac.bal, amt))
			return
		}
		ac.bal -= amt
		to.Add(amt, callback)
	}
}
