package poker_test

import (
	"fmt"
	poker "learn-go-with-tests/2_Building_app/app"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := poker.CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)

	server := mustMakePlayerServer(t, store, dummyGame)

	player := "Pepper"

	wins := 100
	runConcurrently(t, wins, func() {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)
	})

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		poker.AssertStatus(t, response, http.StatusOK)

		poker.AssertStringEquals(t, response.Body.String(), fmt.Sprintf("%d", wins))
	})

	t.Run("get League concurrently", func(t *testing.T) {
		runConcurrently(t, wins, func() {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, newLeagueRequest())

			poker.AssertStatus(t, response, http.StatusOK)

			got := poker.GetLeagueFromResponse(t, response.Body)
			want := []poker.Player{
				{player, wins},
			}
			poker.AssertLeague(t, got, want)
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
