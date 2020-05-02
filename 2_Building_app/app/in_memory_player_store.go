package main

import "sync"

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		store: map[string]int{},
	}
}

type InMemoryPlayerStore struct {
	store map[string]int
	m     sync.Mutex
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.m.Lock()
	defer i.m.Unlock()
	i.store[name]++
}
