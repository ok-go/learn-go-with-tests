package poker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	AssertNoError(t, err)

	server := NewPlayerServer(store)
	player := "Pepper"

	wins := 100
	runConcurrently(t, wins, func() {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)
	})

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		AssertIntEquals(t, response.Code, http.StatusOK)

		AssertStringEquals(t, response.Body.String(), fmt.Sprintf("%d", wins))
	})

	t.Run("get league concurrently", func(t *testing.T) {
		runConcurrently(t, wins, func() {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, newLeagueRequest())

			AssertIntEquals(t, response.Code, http.StatusOK)

			got := GetLeagueFromResponse(t, response.Body)
			want := []Player{
				{player, wins},
			}
			AssertLeague(t, got, want)
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
