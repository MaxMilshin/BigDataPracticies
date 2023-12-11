package main

import (
	"context"
	"fmt"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"

	"log"
)

func handlePeerDownstream(c *ws.Conn, ctx context.Context, peer string) {
	lastProcessed := 0
	for {
		for ; lastProcessed < len(Journal); lastProcessed++ {
			transaction := Journal[lastProcessed]
			error := wsjson.Write(ctx, c, transaction)
			if error != nil {
				break
			}
		}
		time.Sleep(time.Second)
	}
}

func websocketClient(manager *TransactionManager) {
	for _, peer := range peers {
		go handlePeerUpstream(manager, peer)
	}
}

func handlePeerUpstream(manager *TransactionManager, peer string) {
	for {
		ctx := context.Background()
		url := fmt.Sprintf("ws://%s:%s/ws", Source, peer)
		c, _, error := ws.Dial(ctx, url, nil)
		if error != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		acceptReplication(manager, c, ctx, peer)
		time.Sleep(15 * time.Second)
	}
}

func acceptReplication(manager *TransactionManager, c *ws.Conn, ctx context.Context, peer string) {
	for {
		var transaction Transaction
		error := wsjson.Read(ctx, c, &transaction)
		log.Printf("Some transaction happened on peer: %s\n", peer)
		if error != nil {
			return
		}
		manager.transactions <- transaction
	}
}
