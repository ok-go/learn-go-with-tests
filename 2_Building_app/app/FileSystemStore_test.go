package poker

import (
	"log"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("/League from a in", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player playerStore, %v", err)
		}

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, store.GetLeague(), want)

		// read again
		AssertLeague(t, store.GetLeague(), want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player playerStore, %v", err)
		}

		got := store.GetPlayerScore("Chris")
		AssertScoreEquals(t, got, 33)
	})

	t.Run("playerStore wins for existing player", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player playerStore, %v", err)
		}
		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		AssertScoreEquals(t, got, want)
	})

	t.Run("playerStore wins for new player", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player playerStore, %v", err)
		}
		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		AssertScoreEquals(t, got, want)
	})
}
