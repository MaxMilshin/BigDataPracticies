package main

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"log"
	"sync"
)

var Snapshot string = "{}"
var Journal []Transaction

type Transaction struct {
	Source  string
	Id      uint64
	Payload string
}

type TransactionManager struct {
	transactions chan Transaction
	mutex        sync.Mutex
	vClock       map[string]uint64
}

func getTransactionManager() TransactionManager {
	return TransactionManager{transactions: make(chan Transaction), vClock: make(map[string]uint64)}
}

func (transactionManager *TransactionManager) startTransactionManager() {
	for {
		transactionManager.ApplyTransaction()
	}
}

func (transactionManager *TransactionManager) ApplyTransaction() {
	transaction := <-transactionManager.transactions
	transactionManager.mutex.Lock()
	defer transactionManager.mutex.Unlock()
	log.Printf("Got transaction: {Source: %s, Id: %v, Payload: %s}\n", transaction.Source, transaction.Id, transaction.Payload)

	transactionAlreadyApplied := transactionManager.vClock[transaction.Source] >= transaction.Id
	
	if transactionAlreadyApplied {
		log.Printf("Transaction had been applied already\n")
		return
	}

	transactionManager.vClock[transaction.Source] = transaction.Id
	Journal = append(Journal, transaction)
	
	patch, error := jsonpatch.DecodePatch([]byte(transaction.Payload))
	if error != nil {
		log.Printf("Transaction application had been failed: %s", error)
		return
	}

	newsnap, error := patch.Apply([]byte(Snapshot))
	if error != nil {
		log.Printf("Transaction application had been failed: %s", error)
		return
	}
	Snapshot = string(newsnap)
	log.Printf("Transaction had been applied successfully\n")
}
