package mocks

import (
	"fmt"
	"io"
	"time"
)

const (
	finalWord      = "Go!"
	countdownStart = 3
)

type Sleeper interface {
	Sleep()
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (s ConfigurableSleeper) Sleep() {
	s.sleep(s.duration)
}

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
		if _, err := fmt.Fprintln(out, i); err != nil {
		}
	}
	sleeper.Sleep()
	if _, err := fmt.Fprint(out, finalWord); err != nil {
	}
}
