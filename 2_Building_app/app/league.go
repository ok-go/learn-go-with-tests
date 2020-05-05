package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func NewLeague(r io.Reader) (league League, err error) {
	if err = json.NewDecoder(r).Decode(&league); err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
