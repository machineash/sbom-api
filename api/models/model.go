package models

import (
	"errors"
	"strings"
	"sync"
)

// define what a component looks like
type Component struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
	Source   string `json:"source"`
	License  string `json:"license"`
}

// in-memory store (shared reference)
type Store struct {
	Mu         sync.RWMutex
	Components map[int]Component
	NextID     int
}

// Exported constructor
func NewStore() *Store {
	return &Store{
		Components: make(map[int]Component),
		NextID:     1,
	}
}

// simple validation for phase 1 (update later as needed)
func (c *Component) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(c.Version) == "" {
		return errors.New("version is required")
	}
	if strings.TrimSpace(c.Checksum) == "" { // CHECK
		return errors.New("checksum is required")
	}
	if strings.TrimSpace(c.Source) == "" {
		return errors.New("source is required")
	}
	if strings.TrimSpace(c.License) == "" {
		return errors.New("license is required")
	}
	return nil
}
