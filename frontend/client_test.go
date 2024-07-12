package main

import (
	"errors"
	"testing"

	"github.com/gorilla/websocket"
)

type MockWebsocketClient struct {
	msg chan string
}

func (ws *MockWebsocketClient) WriteMessage(messageType int, data []byte) error {
	if string(data) == `{"topic_id":-1}` {
		ws.msg <- "error"
		return errors.New("error while writing message")
	}
	ws.msg <- string(data)
	return nil
}

func (ws *MockWebsocketClient) ReadMessage() (messageType int, p []byte, err error) {
	msg := <-ws.msg
	switch msg {
	case "error":
		return websocket.TextMessage, []byte(msg), errors.New("error while reading")
	case "close":
		return websocket.TextMessage, []byte(msg), &websocket.CloseError{Code: websocket.CloseNormalClosure, Text: "Closing"}
	default:
		return websocket.TextMessage, []byte(msg), nil
	}
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
		conn:  &MockWebsocketClient{channel},
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
