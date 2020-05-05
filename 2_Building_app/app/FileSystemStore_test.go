package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, store.GetLeague(), want)

		// read again
		assertLeague(t, store.GetLeague(), want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		got := store.GetPlayerScore("Chris")
		assertScoreEquals(t, got, 33)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}
		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for new player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}
		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertScoreEquals(t, got, want)
	})
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file, %v", err)
	}

	if _, err := tmpFile.Write([]byte(initialData)); err != nil {
		assertNoError(t, err)
	}

	removeFile := func() {
		if err := tmpFile.Close(); err != nil {
			assertNoError(t, err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			assertNoError(t, err)
		}
	}

	return tmpFile, removeFile
}
