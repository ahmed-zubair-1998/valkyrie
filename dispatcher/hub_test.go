package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGetOrCreateServer(t *testing.T) {
	t.Run("should return new server for new connectio", func(t *testing.T) {
		hub := NewHub()
		var ws WebsocketConnectionInterface
		server := hub.GetOrCreateServer(ws)

		close(server.broadcast)

		got := len(hub.servers)
		want := 1

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should return existing server for existing connection", func(t *testing.T) {
		hub := NewHub()
		var ws WebsocketConnectionInterface
		hub.GetOrCreateServer(ws)
		server := hub.GetOrCreateServer(ws)

		close(server.broadcast)

		got := len(hub.servers)
		want := 1

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}

func TestBroadcastEvent(t *testing.T) {
	t.Run("should return error if invalid request data", func(t *testing.T) {
		hub := NewHub()

		var jsonData bytes.Buffer
		invalidData := `{"invalid_json": "error"}`
		json.NewEncoder(&jsonData).Encode(invalidData)

		request, _ := http.NewRequest(http.MethodPost, "/events/broadcast", &jsonData)
		response := httptest.NewRecorder()

		hub.BroadcastEvent(response, request)

		got := response.Body.String()
		want := "Invalid request data"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should return error if topic is not subscribed by any FE server", func(t *testing.T) {
		hub := NewHub()

		var jsonData bytes.Buffer
		data := &BroadcastEventData{
			TopicId: 1,
			Message: "Hello",
		}
		json.NewEncoder(&jsonData).Encode(data)

		request, _ := http.NewRequest(http.MethodPost, "/events/broadcast", &jsonData)
		response := httptest.NewRecorder()

		hub.BroadcastEvent(response, request)

		got := response.Body.String()
		want := "Incorrect topic id"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should return success message if broadcast is successful", func(t *testing.T) {
		channel := make(chan string)
		ws := &MockWebsocketClient{channel}

		go func() {
			<-channel
		}()

		hub := NewHub()
		server := hub.GetOrCreateServer(ws)
		server.topics[1] = true

		var jsonData bytes.Buffer
		data := &BroadcastEventData{
			TopicId: 1,
			Message: "Hello",
		}
		json.NewEncoder(&jsonData).Encode(data)

		request, _ := http.NewRequest(http.MethodPost, "/events/broadcast", &jsonData)
		response := httptest.NewRecorder()

		hub.BroadcastEvent(response, request)

		got := response.Body.String()
		want := "Event broadcasted successfully"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}

func TestFrontendServerConnection(t *testing.T) {
	t.Run("should create a websocket connection with FE server", func(t *testing.T) {
		hub := NewHub()

		mux := httptest.NewServer(http.HandlerFunc(hub.FrontendServerConnection))
		defer mux.Close()

		url := "ws" + strings.TrimPrefix(mux.URL, "http") + "/frontend/connect"
		client, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("could not connect to WebSocket server: %v", err)
		}
		defer client.Close()

		got := len(hub.servers)
		want := 1

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}
