package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

type FileSystemPlayerStore struct {
	mutex    sync.Mutex
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	if err := initializePlayerDBFile(file); err != nil {
		return nil, fmt.Errorf("problem init player db file, %v", err)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func (s *FileSystemPlayerStore) getLeague() League {
	sort.Slice(s.league, func(i, j int) bool {
		return s.league[i].Wins > s.league[j].Wins
	})
	return s.league
}

func (s *FileSystemPlayerStore) recordWin(name string) {
	if player := s.league.Find(name); player != nil {
		player.Wins++
	} else {
		s.league = append(s.league, Player{name, 1})
	}

	s.database.Encode(s.league)
}

func (s *FileSystemPlayerStore) getPlayerScore(name string) int {
	if player := s.getLeague().Find(name); player != nil {
		return player.Wins
	}
	return 0
}

func (s *FileSystemPlayerStore) GetLeague() League {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getLeague()
}

func (s *FileSystemPlayerStore) GetPlayerScore(name string) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getPlayerScore(name)
}

func (s *FileSystemPlayerStore) RecordWin(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.recordWin(name)
}

func initializePlayerDBFile(file *os.File) error {
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("problem seek file %s, %v", file.Name(), err)
	}

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		if _, err := file.Write([]byte("[]")); err != nil {
			return fmt.Errorf("problem write to file %s, %v", file.Name(), err)
		}
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("problem seek file %s, %v", file.Name(), err)
		}
	}

	return nil
}
