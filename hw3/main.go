package main

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	ws "nhooyr.io/websocket"
	"os"

	"log"
)

var port string
var peers []string
const Source = "127.0.0.1"

var IdCounter uint64 = 1

var content embed.FS

func main() {
	port = os.Args[1]
	peers = os.Args[2:]

	transactionManager := getTransactionManager()
	go transactionManager.startTransactionManager()
	
	go websocketClient(&transactionManager);
	
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.FS(content))))
	http.HandleFunc("/vclock", handleVClock(&transactionManager))
	http.HandleFunc("/post", handlePost(&transactionManager))
	http.HandleFunc("/get", handleGet())
	http.HandleFunc("/ws", handleWS())
	
	if http.ListenAndServe(":" + port, nil) != nil {
        println("Error occured while serving the address")
    }
}

func handleVClock(transactionManager *TransactionManager) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bytes, error := json.Marshal(transactionManager.vClock)
		if error != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, error = writer.Write(bytes)
		if error != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func handlePost(transactionManager *TransactionManager) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, error := io.ReadAll(request.Body)
		if error != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var id = IdCounter
		IdCounter++
		var transaction = Transaction{Source: Source, Id: id, Payload: string(body)}
		transactionManager.transactions <- transaction
	}
}

func handleGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, error := writer.Write([]byte(Snapshot))
		if error != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func handleWS() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		c, error := ws.Accept(writer, request, &ws.AcceptOptions{
			InsecureSkipVerify: true,
			OriginPatterns:     []string{"*"},
		})
		if error != nil {
			panic(error)
		}
		ctx := request.Context()
		handlePeerDownstream(c, ctx, request.Host)
		c.Close(ws.StatusNormalClosure, "Connection closed")
	}
}
