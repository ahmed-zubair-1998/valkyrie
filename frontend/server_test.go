package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

type MockHub struct {
	topics map[int]*Topic
}

func NewMockHub() *MockHub {
	return &MockHub{
		topics: make(map[int]*Topic),
	}
}

func (hub *MockHub) SubscribeToTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Received request on ", r.URL)
}

func TestSetupRoutes(t *testing.T) {
	hub := NewMockHub()
	mux := SetupRoutes(hub)
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"should return heartbeat", "/heartbeat", "heartbeat"},
		{"should handle topic subscription", "/topics/subscribe?topicId=1", "Received request on /topics/subscribe?topicId=1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, test.url, nil)
			response := httptest.NewRecorder()

			mux.ServeHTTP(response, request)

			got := response.Body.String()

			if got != test.want {
				t.Errorf("got %q, wanted %q", got, test.want)
			}
		})
	}
}

func TestSetupConnectionToDispatcher(t *testing.T) {
	t.Run("should return error if connection not successful", func(t *testing.T) {
		_, err := SetupConnectionToDispatcher("http://127.0.0.1:8080")
		got := err.Error()
		want := "dial tcp 127.0.0.1:8080: connect: connection refused"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
	t.Run("should not return error if connection not successful", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/frontend/connect", func(w http.ResponseWriter, r *http.Request) {
			websocketUpgrader := websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}
			conn, _ := websocketUpgrader.Upgrade(w, r, nil)
			defer conn.Close()
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		_, err := SetupConnectionToDispatcher(ts.URL)
		if err != nil {
			t.Errorf("should not return error")
		}
	})
}
