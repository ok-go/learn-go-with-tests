package main

import (
	"fmt"
	"log"
	"os"

	poker "learn-go-with-tests/2_Building_app/app"
)

const dbFileName = "game.db.json"

func main() {
	store, fClose, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fClose()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
