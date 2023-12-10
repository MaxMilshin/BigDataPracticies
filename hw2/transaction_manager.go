package main

import (
	"sync"
	"time"
)

var lastTransaction string
var journal []string
var snapshot string

const timeInterval = 60 * time.Second

type TransactionManager struct {
	transactions chan string
	mutex        sync.Mutex
	timer        *time.Ticker
}

func getTransactionManager() TransactionManager {
	return TransactionManager{transactions: make(chan string), timer: time.NewTicker(timeInterval)}
}

func (manager *TransactionManager) startTransactionManager() {
	for {
		select {
			case transaction := <-manager.transactions:
				manager.mutex.Lock()
				journal = append(journal, transaction)
				lastTransaction = transaction
				manager.mutex.Unlock()
			case <-manager.timer.C:
				manager.mutex.Lock()
				snapshot = lastTransaction
				journal = nil
				manager.mutex.Unlock()
		}
	}
}

func (manager *TransactionManager) pushTransaction(transaction string) {
	manager.transactions <- transaction
}
