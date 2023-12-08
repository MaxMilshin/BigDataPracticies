package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func clearFile(filename string) error {
	file, error := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
	if os.IsNotExist(error) {
		return nil
	} else if error != nil {
		return error
	}
	defer file.Close()

	_, error = file.WriteString("")
	if error != nil {
		return error
	}

	return nil
}

func TestHandlers(t *testing.T) {
	error := clearFile(storageFile)
	if error != nil {
		panic(error)
	}
	testItems := []struct {
		description  string
		method       string
		route        string
		handler      func(http.ResponseWriter, *http.Request)
		inputData    string
		outputData   string
		statusCode   int
	}{
		{"GET before any POST requests", http.MethodGet, "/get", getHandler, "", "", 200},
		{"POST request", http.MethodPost, "/replace", replaceHandler, "sample data for testing purpose", "", 200},
		{"GET after POST request", http.MethodGet, "/get", getHandler, "", "sample data for testing purpose", 200},
		{"one more POST request", http.MethodPost, "/replace", replaceHandler, "complicated sample data for testing purpose", "", 200},
		{"one more GET after last POST request", http.MethodGet, "/get", getHandler, "", "complicated sample data for testing purpose", 200},
	}
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