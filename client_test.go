package main

import (
	"testing"
)

type MockWebsocketConnection struct {
	msg chan string
}

func (ws *MockWebsocketConnection) WriteMessage(messageType int, data []byte) error {
	ws.msg <- string(data)
	return nil
}

func TestListenToEvents(t *testing.T) {
	topic := &Topic{
		id:         1,
		clients:    map[*Client]bool{},
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	channel := make(chan string)
	client := &Client{
		topic: topic,
		conn:  &MockWebsocketConnection{channel},
		send:  make(chan []byte),
	}
	want := "test message"

	go client.ListenToEvents()

	client.send <- []byte(want)
	close(client.send)

	if got := <-channel; got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
