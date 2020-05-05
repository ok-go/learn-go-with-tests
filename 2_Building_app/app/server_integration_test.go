package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := NewPlayerServer(store)
	player := "Pepper"

	wins := 100
	runConcurrently(t, wins, func() {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), fmt.Sprintf("%d", wins))
	})

	t.Run("get league concurrently", func(t *testing.T) {
		runConcurrently(t, wins, func() {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, newLeagueRequest())

			assertStatus(t, response.Code, http.StatusOK)

			got := getLeagueFromResponse(t, response.Body)
			want := []Player{
				{player, wins},
			}
			assertLeague(t, got, want)
		})
	})
}

func runConcurrently(t *testing.T, count int, f func()) {
	t.Helper()
	wg := &sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}
