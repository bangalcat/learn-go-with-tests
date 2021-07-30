package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
	Wraping up
	* select
	  * Helps you wait on multiple channels.
 	  * Sometimes you'll want to include time.
		  After in one of your cases to prevent your system blocking forever.

	* httptest
		* A convenient way of creating test servers so you can have reliable and controllable tests.
		* Using the same interfaces as the "real" net/http servers which is consistent and less for you to learn.
*/

func TestRacer(t *testing.T) {

	t.Run("compares speeds of servers, returning the url of the faster", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, _ := Racer(slowURL, fastURL)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		server := makeDelayedServer(20 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected and error but didn't get one")
		}
	})
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		rw.WriteHeader(http.StatusOK)
	}))
}
