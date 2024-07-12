package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
)

type MockWebsocketClient struct {
	msg chan string
}

func (ws *MockWebsocketClient) WriteMessage(messageType int, data []byte) error {
	var jsonData BroadcastEventData
	json.Unmarshal(data, &jsonData)

	if jsonData.TopicId == -1 {
		return errors.New("error while writing message")
	}
	ws.msg <- jsonData.Message
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

func TestSubscribeToTopics(t *testing.T) {
	data := &SubscribeTopicData{TopicId: 1}
	jsonData, _ := json.Marshal(data)

	tests := []struct {
		name             string
		message          string
		subscribedTopics int
	}{
		{"should return if websocket connection is closed by FE server", "close", 0},
		{"should return if error is received from FE server", "error", 0},
		{"should return if invalid message is received from FE server", "invalid message", 0},
		{"should subscribe to topic if valid message is received from FE server", string(jsonData), 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			channel := make(chan string)
			server := &FrontendServer{
				topics:    make(map[int]bool),
				broadcast: make(chan *BroadcastEventData),
				conn:      &MockWebsocketClient{channel},
			}
			go func() {
				channel <- test.message
				channel <- "close"
			}()

			server.SubscribeToTopics()
			got := len(server.topics)

			if got != test.subscribedTopics {
				t.Errorf("got %q, wanted %q", got, test.subscribedTopics)
			}
		})
	}
}

func TestSendRequest(t *testing.T) {
	t.Run("should return error if sending message to FE server fails", func(t *testing.T) {
		data := &BroadcastEventData{
			TopicId: -1,
			Message: "Hello",
		}
		server := &FrontendServer{
			topics:    make(map[int]bool),
			broadcast: make(chan *BroadcastEventData),
			conn:      &MockWebsocketClient{},
		}

		got := server.SendRequest(data)
		if got == nil {
			t.Errorf("should have returned error")
		}
	})
	t.Run("should not return error if message is succesfully sent to FE server", func(t *testing.T) {
		data := &BroadcastEventData{
			TopicId: 1,
			Message: "Hello",
		}
		channel := make(chan string)
		server := &FrontendServer{
			topics:    make(map[int]bool),
			broadcast: make(chan *BroadcastEventData),
			conn:      &MockWebsocketClient{channel},
		}

		go func() {
			<-channel
		}()

		got := server.SendRequest(data)
		if got != nil {
			t.Errorf("should not have returned error")
		}
	})
}

func TestBroadcastEvents(t *testing.T) {
	t.Run("should send message to FE server on receiving broadcast", func(t *testing.T) {
		channel := make(chan string)
		server := &FrontendServer{
			topics:    make(map[int]bool),
			broadcast: make(chan *BroadcastEventData),
			conn:      &MockWebsocketClient{channel},
		}
		want := "Hello"
		data := &BroadcastEventData{
			TopicId: 1,
			Message: want,
		}

		go server.BroadcastEvents()
		server.broadcast <- data
		close(server.broadcast)

		got := <-channel

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}
