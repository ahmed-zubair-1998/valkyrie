package main

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("should register client to topic", func(t *testing.T) {
		topic := &Topic{
			id:         1,
			clients:    map[*Client]bool{},
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
		go topic.Run()

		client := &Client{
			topic: topic,
			send:  make(chan []byte),
		}
		topic.register <- client

		if len(topic.clients) != 1 {
			t.Errorf("client not registered to a topic")
		}
	})
	t.Run("should broadcast message to all clients", func(t *testing.T) {
		topic := &Topic{
			id:         1,
			clients:    map[*Client]bool{},
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}

		for range 10 {
			client := &Client{
				topic: topic,
				send:  make(chan []byte),
			}
			topic.clients[client] = true
		}
		go topic.Run()

		topic.broadcast <- []byte("test message")

		go func() {
			for client := range topic.clients {
				select {
				case <-client.send:
				case <-time.After(time.Second * 5):
					t.Errorf("message not broadcasted")
				}
			}
		}()
	})
}
