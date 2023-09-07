package storage

import (
	"fmt"
)

type Storage struct {
	values map[string]string
}

func (s *Storage) GetValue(key string) (string, error) {
	if value, ok := s.values[key]; ok {
		return value, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", key)
}

func (s *Storage) AddValue(key string, value string) {
	s.values[key] = value
}

func NewStorage() Storage {
	return Storage{values: make(map[string]string)}
}
