package poker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}
	if store.winCalls[0] != winner {
		t.Errorf("did not playerStore correct winner got %q want %q", store.winCalls[0], winner)
	}
}

func GetLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
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
