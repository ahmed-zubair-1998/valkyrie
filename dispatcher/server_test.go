package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

type MockHub struct {
	servers map[*websocket.Conn]*FrontendServer
}

func NewMockHub() *MockHub {
	return &MockHub{
		servers: make(map[*websocket.Conn]*FrontendServer),
	}
}

func (hub *MockHub) FrontendServerConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Received request on ", r.URL)
}

func (hub *MockHub) BroadcastEvent(w http.ResponseWriter, r *http.Request) {
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
		{"should handle topic subscription", "/frontend/connect", "Received request on /frontend/connect"},
		{"should handle event broadcast", "/events/broadcast", "Received request on /events/broadcast"},
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
