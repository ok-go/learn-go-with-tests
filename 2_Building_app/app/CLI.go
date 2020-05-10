package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	PlayerPrompt          = "Please enter the number of players: "
	BadPlayerInputMessage = "bad value received for number of players, please try again with a number"
	BadWinnerInputMessage = "incorrect winner"
)

type CLI struct {
	out  io.Writer
	in   *bufio.Scanner
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := cli.readNumberOfPlayers()
	if err != nil {
		fmt.Fprint(cli.out, err)
		return
	}
	cli.game.Start(numberOfPlayers, cli.out)

	winner, err := cli.readWinner()
	if err != nil {
		fmt.Fprint(cli.out, err)
	}
	cli.game.Finish(winner)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) readWinner() (string, error) {
	input := cli.readLine()
	if !strings.Contains(input, " wins") {
		return "", errors.New(BadWinnerInputMessage)
	}
	return strings.Replace(input, " wins", "", 1), nil
}

func (cli *CLI) readNumberOfPlayers() (int, error) {
	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))
	if err != nil {
		return 0, errors.New(BadPlayerInputMessage)
	}

	return numberOfPlayers, nil
}
