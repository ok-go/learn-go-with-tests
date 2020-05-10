package main

import (
	poker "learn-go-with-tests/2_Building_app/app"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, fClose, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fClose()

	game := poker.NewTexasHoldem(store, poker.BlindAlerterFunc(poker.Alerter))

	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatalf("could not create player server: %v", err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
