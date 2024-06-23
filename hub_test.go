package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
}
