package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	wins := 100

	wg := &sync.WaitGroup{}
	wg.Add(wins)
	for i := 0; i < wins; i++ {
		go serveHttp(wg, server, player)
	}
	wg.Wait()

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), fmt.Sprintf("%d", wins))
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{player, wins},
		}
		assertLeague(t, got, want)
	})
}

func serveHttp(wg *sync.WaitGroup, server *PlayerServer, player string) {
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	wg.Done()
}
