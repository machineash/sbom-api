package api

import "sync"

type Component struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
	Source   string `json:"source"`
	License  string `json:"license"`
}

// in-memory store
type store struct {
	mu         sync.RWMutex
	components map[int]Component
	nextID     int
}

func newStore() *store {
	return &store{
		components: make(map[int]Component),
		nextID:     1,
	}
}
