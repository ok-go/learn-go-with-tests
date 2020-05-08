package poker_test

import (
	"bytes"
	"fmt"
	"io"
	poker "learn-go-with-tests/2_Building_app/app"
	"strconv"
	"strings"
	"testing"
)

var (
	dummyBlindAlerter = &poker.SpyBlindAlerter{}
	dummyPlayerStore  = &poker.StubPlayerStore{}
	dummyStdIn        = &bytes.Buffer{}
	dummyStdOut       = &bytes.Buffer{}
)

type userInput struct {
	amount string
	winner string
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
	FinishCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishCalled = true
}

func TestCLI(t *testing.T) {
	cases := []userInput{
		{"3", "Chris"},
		{"8", "Cleo"},
	}
	for _, c := range cases {
		name := fmt.Sprintf("start game with %q players finish game with %q as winner", c.amount, c.winner)
		t.Run(name, func(t *testing.T) {
			game := &GameSpy{}
			stdout := &bytes.Buffer{}

			in := userSends(c.amount, fmt.Sprintf("%s wins", c.winner))
			cli := poker.NewCLI(in, stdout, game)

			cli.PlayPoker()

			assertMessageSentToUser(t, stdout, poker.PlayerPrompt)
			assertGameStartedWith(t, game, c.amount)
			assertFinishCalledWith(t, game, c.winner)
		})
	}

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("Pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputMessage)
	})

	t.Run("it prints an error when incorrect winner is entered", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("3", "Pies is a killer")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMessage)
	})
}

func assertMessageSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertGameStartedWith(t *testing.T, game *GameSpy, want string) {
	t.Helper()
	amount, _ := strconv.Atoi(want)
	if game.StartedWith != amount {
		t.Errorf("got %d, want Start called with %d", game.StartedWith, amount)
	}
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, want string) {
	t.Helper()
	if game.FinishedWith != want {
		t.Errorf("got %q, want Finish called with %q", game.FinishedWith, want)
	}
}

func assertGameNotStarted(t *testing.T, game *GameSpy) {
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func assertGameNotFinished(t *testing.T, game *GameSpy) {
	if !game.FinishCalled {
		t.Errorf("game should not have finished")
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}
