package main

import "testing"

func TestRun(t *testing.T) {
	t.Run("should register client to topic", func(t *testing.T) {
		topic := &Topic{
			id:         1,
			clients:    map[*Client]bool{},
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			stop:       make(chan bool),
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
			stop:       make(chan bool),
		}
		go topic.Run()

		client1 := &Client{
			topic: topic,
			send:  make(chan []byte),
		}
		topic.register <- client1

		topic.broadcast <- []byte("test message")

		go func() {
			for {
				select {
				case <-client1.send:
				default:
					t.Errorf("message not broadcasted")
				}
			}
		}()
	})
}
