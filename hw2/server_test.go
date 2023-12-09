package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	transactionManager := getTransactionManager()

	replaceHandlerFunc := replaceHandler(&transactionManager)
	getHandlerFunc := getHandler(&transactionManager)

	testItems := []struct {
		description  string
		method       string
		route        string
		handler      func(http.ResponseWriter, *http.Request)
		inputData    string
		outputData   string
		statusCode   int
	}{
		{"GET before any POST requests", http.MethodGet, "/get", getHandlerFunc, "", "", 200},
		{"POST request", http.MethodPost, "/replace", replaceHandlerFunc, "sample data for testing purpose", "", 200},
		{"GET after POST request", http.MethodGet, "/get", getHandlerFunc, "", "sample data for testing purpose", 200},
		{"one more POST request", http.MethodPost, "/replace", replaceHandlerFunc, "complicated sample data for testing purpose", "", 200},
		{"one more GET after last POST request", http.MethodGet, "/get", getHandlerFunc, "", "complicated sample data for testing purpose", 200},
	}
	go transactionManager.startTransactionManager()
	for _, testItem := range testItems {
		t.Run(testItem.description, func(t *testing.T) {
			request := httptest.NewRequest(testItem.method, testItem.route, bytes.NewReader([]byte(testItem.inputData)))
			responseRecorder := httptest.NewRecorder()

			testItem.handler(responseRecorder, request)

			if responseRecorder.Code != testItem.statusCode {
				t.Errorf("Response code is incorrect: expected %d, but got %d", testItem.statusCode, responseRecorder.Code)
			}

			body, _ := io.ReadAll(responseRecorder.Body)
			stringBody := string(body)

			if stringBody != testItem.outputData {
				t.Errorf("Result is incorrect: expected %s, but got %s", testItem.outputData, stringBody)
			}
		})
	}
}