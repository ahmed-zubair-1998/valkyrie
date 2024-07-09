package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type MockWebsocketConnection struct{}

func (ws *MockWebsocketConnection) WriteMessage(messageType int, data []byte) error {
	return nil
}

func (ws *MockWebsocketConnection) ReadMessage() (messageType int, p []byte, err error) {
	return websocket.TextMessage, []byte("Hello"), nil
}

func TestSubscribeToTopic(t *testing.T) {
	t.Run("should fail if invalid topic id", func(t *testing.T) {
		var conn MockWebsocketClient
		hub := &Hub{
			conn:   &conn,
			topics: make(map[int]*Topic),
		}
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
		var conn MockWebsocketConnection
		hub := &Hub{
			conn:   &conn,
			topics: make(map[int]*Topic),
		}
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
	t.Run("should not broadcast if topic is not yet created", func(t *testing.T) {
		var conn MockWebsocketClient
		hub := &Hub{
			conn:   &conn,
			topics: make(map[int]*Topic),
		}
		data := BroadcastEventData{
			1,
			"Hello",
		}

		got := hub.BroadcastEvent(data)
		want := "incorrect topic id"

		if got.Error() != want {
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

		var conn MockWebsocketClient
		hub := &Hub{
			conn:   &conn,
			topics: make(map[int]*Topic),
		}
		hub.topics[1] = topic
		done := make(chan bool)

		go func() {
			want := "1:Hello"
			select {
			case got := <-topic.broadcast:
				if string(got) != want {
					t.Errorf("got %q, wanted %q", got, want)
				}
			case <-time.After(time.Second * 5):
				t.Errorf("event not broadcasted")
			}
			done <- true
		}()

		data := BroadcastEventData{
			1,
			"Hello",
		}
		err := hub.BroadcastEvent(data)
		if err != nil {
			t.Fatalf("Broadcast event should not return error")
		}
		<-done
	})
}

func TestGetOrCreateTopic(t *testing.T) {
	t.Run("should create a new topic if not present", func(t *testing.T) {
		var conn MockWebsocketClient
		hub := &Hub{
			conn:   &conn,
			topics: make(map[int]*Topic),
		}
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
		var conn MockWebsocketClient
		hub := &Hub{
			conn:   &conn,
			topics: make(map[int]*Topic),
		}
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

func TestNotifyDispatcher(t *testing.T) {
	t.Run("should successfully send topic id to dispatcher", func(t *testing.T) {
		channel := make(chan string)
		hub := &Hub{
			conn:   &MockWebsocketClient{channel},
			topics: make(map[int]*Topic),
		}
		go hub.NotifyDispatcher(1)
		want := `{"topic_id":1}`

		if got := <-channel; got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should log error if there sending message fails", func(t *testing.T) {
		channel := make(chan string)
		hub := &Hub{
			conn:   &MockWebsocketClient{channel},
			topics: make(map[int]*Topic),
		}
		go hub.NotifyDispatcher(-1)
		want := "error"

		if got := <-channel; got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}

func TestHandleMessagesFromDispatcher(t *testing.T) {
	data := &BroadcastEventData{
		Id:      1,
		Message: "Hello",
	}
	jsonData, _ := json.Marshal(data)

	tests := []struct {
		name    string
		message string
		want    int
	}{
		{"should return when connection is closed by dispatcher", "close", 1},
		{"should return when err is received from dispatcher", "error", 1},
		{"should return when invalid message is received from dispatcher", "invalid message", 1},
		{"should broadcast when event is received from dispatcher", string(jsonData), 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			channel := make(chan string)
			hub := &Hub{
				conn:   &MockWebsocketClient{channel},
				topics: make(map[int]*Topic),
			}

			msgCount := 0
			go func() {
				msgCount++
				channel <- test.message
				time.Sleep(time.Millisecond * 5)
				msgCount++
				channel <- "close"
			}()

			hub.HandleMessagesFromDispatcher()

			if msgCount != test.want {
				t.Errorf("got %q, wanted %q", msgCount, test.want)
			}
		})
	}
}

func TestNewHub(t *testing.T) {
	var conn MockWebsocketConnection
	hub := NewHub(&conn)

	got := len(hub.topics)
	want := 0

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
