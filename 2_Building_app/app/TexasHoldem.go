package poker

import (
	"io"
	"time"
)

type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}

func NewTexasHoldem(store PlayerStore, alerter BlindAlerter) Game {
	return &TexasHoldem{
		alerter: alerter,
		store:   store,
	}
}

func (g *TexasHoldem) Start(numberOfPlayers int, alertDestination io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, amount := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, amount, alertDestination)
		blindTime = blindTime + blindIncrement
	}
}

func (g *TexasHoldem) Finish(winner string) {
	g.store.RecordWin(winner)
}
