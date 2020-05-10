package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips At %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{At: duration, Amount: amount})
}

type GameSpy struct {
	StartedWith int
	StartCalled bool
	BlindAlert  []byte

	FinishedWith string
	FinishCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
	out.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishCalled = true
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.WinCalls) != 1 {
		t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
	}
	if len(store.WinCalls) > 0 && store.WinCalls[0] != winner {
		t.Errorf("did not playerStore correct winner got %q want %q", store.WinCalls[0], winner)
	}
}

func GetLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func AssertStringEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func AssertStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Errorf("got status %d, want %d", response.Code, want)
	}
}

func AssertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %q, got %v", want, response.Result().Header)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func AssertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func AssertGameStartedWith(t *testing.T, game *GameSpy, want int) {
	t.Helper()
	if game.StartedWith != want {
		t.Errorf("got %d, want Start called with %d", game.StartedWith, want)
	}
}

func AssertFinishCalledWith(t *testing.T, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == winner
	})

	if !passed {
		t.Errorf("got %q, winner Finish called with %q", game.FinishedWith, winner)
	}
}

func AssertGameNotStarted(t *testing.T, game *GameSpy) {
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func AssertGameNotFinished(t *testing.T, game *GameSpy) {
	if !game.FinishCalled {
		t.Errorf("game should not have finished")
	}
}

func CreateTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file, %v", err)
	}

	if _, err := tmpFile.Write([]byte(initialData)); err != nil {
		AssertNoError(t, err)
	}

	removeFile := func() {
		if err := tmpFile.Close(); err != nil {
			AssertNoError(t, err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			AssertNoError(t, err)
		}
	}

	return tmpFile, removeFile
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
