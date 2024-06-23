package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	t.Run("should return heartbeat", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		response := httptest.NewRecorder()

		Heartbeat(response, request)

		got := response.Body.String()
		want := "heartbeat"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}

type MockHub struct {
	topics map[int]*Topic
}

func NewMockHub() *MockHub {
	return &MockHub{
		topics: make(map[int]*Topic),
	}
}

func (hub *MockHub) MockSubscribeToTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Received request on ", r.URL)
}

func (hub *MockHub) MockBroadcastEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Received request on ", r.URL)
}

func TestHubRoutes(t *testing.T) {
	hub := NewMockHub()
	tests := []struct {
		name        string
		url         string
		handlerFunc http.HandlerFunc
	}{
		{"should handle topic subscription", "/topics/subscribe?topicId=1", hub.MockSubscribeToTopic},
		{"should handle event broadcast", "/events/broadcast", hub.MockBroadcastEvent},
	}

	for _, test := range tests {
		t.Run("should handle event broadcast", func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, test.url, nil)
			response := httptest.NewRecorder()

			test.handlerFunc(response, request)

			got := response.Body.String()
			want := "Received request on " + test.url

			if got != want {
				t.Errorf("got %q, wanted %q", got, want)
			}
		})
	}
}
