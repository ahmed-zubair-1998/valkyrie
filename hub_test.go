package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestSubscribeToTopic(t *testing.T) {
	t.Run("should fail if invalid topic id", func(t *testing.T) {
		hub := NewHub()
		request, _ := http.NewRequest(http.MethodGet, "/topics/subscribe?topicId=abc", nil)
		response := httptest.NewRecorder()

		hub.SubscribeToTopic(response, request)

		got := response.Body.String()
		want := "Invalid topic id provided"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should subscribe client to topic", func(t *testing.T) {
		hub := NewHub()
		server := httptest.NewServer(http.HandlerFunc(hub.SubscribeToTopic))
		defer server.Close()

		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/topics/subscribe?topicId=1"
		client, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("could not connect to WebSocket server: %v", err)
		}
		defer client.Close()

		want := "Successfully subscribed to topic id: 1"
		_, got, err := client.ReadMessage()
		if err != nil {
			t.Fatalf("could not read message from Websocket server: %v", err)
		}

		if string(got) != want {
			t.Fatalf("got %q, wanted %q", got, want)
		}

		if len(hub.topics) != 1 {
			t.Fatalf("topic not registered, %d", len(hub.topics))
		}

		if len(hub.topics[1].clients) != 1 {
			t.Fatalf("client not subscribed to the topic")
		}
	})
}

func TestBroadcastEvent(t *testing.T) {
	t.Run("should verify request data", func(t *testing.T) {
		hub := NewHub()
		invalidData := `{"hello": "world"}`
		jsonData, _ := json.Marshal(invalidData)
		request, _ := http.NewRequest("POST", "/events/broadcast", bytes.NewBuffer(jsonData))
		response := httptest.NewRecorder()

		hub.BroadcastEvent(response, request)

		got := response.Body.String()
		want := "Invalid post data"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should not broadcast if topic is not yet created", func(t *testing.T) {
		hub := NewHub()
		data := struct {
			id  int    `json:"id"`
			msg string `json:"msg"`
		}{
			1,
			"hello",
		}
		jsonData, _ := json.Marshal(data)
		request, _ := http.NewRequest("POST", "/events/broadcast", bytes.NewBuffer(jsonData))
		response := httptest.NewRecorder()

		hub.BroadcastEvent(response, request)

		got := response.Body.String()
		want := "Incorrect topic id"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should broadcast event to correct topic", func(t *testing.T) {
		topic := &Topic{
			id:         1,
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}

		hub := NewHub()
		hub.topics[1] = topic

		data := struct {
			Id      int    `json:"id"`
			Message string `json:"msg"`
		}{
			1,
			"hello",
		}
		jsonData, _ := json.Marshal(data)
		request, _ := http.NewRequest("POST", "/events/broadcast", bytes.NewBuffer(jsonData))
		response := httptest.NewRecorder()

		go hub.BroadcastEvent(response, request)

		want := "1:hello"
		select {
		case got := <-topic.broadcast:
			if string(got) != want {
				t.Errorf("got %q, wanted %q", got, want)
			}
		case <-time.After(time.Second * 5):
			t.Errorf("event not broadcasted")
		}
	})
}

func TestGetOrCreateTopic(t *testing.T) {
	t.Run("should create a new topic if not present", func(t *testing.T) {
		hub := NewHub()
		topic := hub.GetOrCreateTopic(1)

		if len(hub.topics) != 1 {
			t.Fatalf("should only create a single topic")
		}

		if hub.topics[1] != topic {
			t.Fatalf("incorrect topic returned")
		}
	})
	t.Run("should return existing topic if present", func(t *testing.T) {
		topic := &Topic{
			id:         1,
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
		hub := NewHub()
		hub.topics[1] = topic

		got := hub.GetOrCreateTopic(1)

		if len(hub.topics) != 1 {
			t.Fatalf("should not create a new topic")
		}

		if got != topic {
			t.Fatalf("incorrect topic returned")
		}
	})
}
