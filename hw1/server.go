package main

import (
    "io"
    "net/http"
    "os"
)

const storageFile = "storage.txt"

func replaceHandler(writer http.ResponseWriter, request *http.Request) {
    buffer, error := io.ReadAll(request.Body)
    if error != nil {
        writer.WriteHeader(http.StatusBadRequest)
        return
    }
    error = os.WriteFile(storageFile, buffer, 0644)
    if error != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        return
    }
    writer.WriteHeader(http.StatusOK)
}

func getHandler(writer http.ResponseWriter, _ *http.Request) {
    contents, error := os.ReadFile(storageFile)
    if os.IsNotExist(error) {
        writer.WriteHeader(http.StatusOK)
        return
    } else if error != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        return
    }
    writer.Write(contents)
    writer.WriteHeader(http.StatusOK)
}

func main() {
    http.HandleFunc("/get", getHandler)
    http.HandleFunc("/replace", replaceHandler)

    println("Server is run up..")

    if http.ListenAndServe(":8080", nil) != nil {
        println("Error occured while serving the address")
    }
}