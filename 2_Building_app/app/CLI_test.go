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
	dummyGame         = &poker.GameSpy{}
	dummyBlindAlerter = &poker.SpyBlindAlerter{}
	dummyPlayerStore  = &poker.StubPlayerStore{}
	dummyStdIn        = &bytes.Buffer{}
	dummyStdOut       = &bytes.Buffer{}
)

type userInput struct {
	amount int
	winner string
}

func TestCLI(t *testing.T) {
	cases := []userInput{
		{3, "Chris"},
		{8, "Cleo"},
	}
	for _, c := range cases {
		name := fmt.Sprintf("start game with %q players finish game with %q as winner", c.amount, c.winner)
		t.Run(name, func(t *testing.T) {
			game := &poker.GameSpy{}

			out := &bytes.Buffer{}
			in := userSends(strconv.Itoa(c.amount), fmt.Sprintf("%s wins", c.winner))

			poker.NewCLI(in, out, game).PlayPoker()

			assertMessageSentToUser(t, out, poker.PlayerPrompt)
			poker.AssertGameStartedWith(t, game, c.amount)
			poker.AssertFinishCalledWith(t, game, c.winner)
		})
	}

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &poker.GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("Pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		poker.AssertGameNotStarted(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputMessage)
	})

	t.Run("it prints an error when incorrect winner is entered", func(t *testing.T) {
		game := &poker.GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("3", "Pies is a killer")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		poker.AssertGameNotFinished(t, game)
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

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}
