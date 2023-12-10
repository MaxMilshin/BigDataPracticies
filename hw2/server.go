package main

import (
	"io"
	"net/http"
)

func replaceHandler(transactionManager *TransactionManager) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, error := io.ReadAll(request.Body)
		if error != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactionManager.pushTransaction(string(body))
		writer.WriteHeader(http.StatusOK)
	}
}

func getHandler(transactionManager *TransactionManager) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, error := writer.Write([]byte(snapshot))
		if error != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func main() {
	transactionManager := getTransactionManager()
	go transactionManager.startTransactionManager()
	
	http.HandleFunc("/replace", replaceHandler(&transactionManager))
	http.HandleFunc("/get", getHandler(&transactionManager))

	println("Server is run up..")

	if http.ListenAndServe(":8080", nil) != nil {
        println("Error occured while serving the address")
    }
}
