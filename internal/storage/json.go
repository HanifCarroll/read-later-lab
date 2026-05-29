package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go-embedded-js/internal/app"
)

type JSONStore struct {
	path string
	mu   sync.Mutex
}

func NewJSONStore(path string) (*JSONStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(path, []byte("[]\n"), 0o644); err != nil {
			return nil, err
		}
	}
	return &JSONStore{path: path}, nil
}

func (s *JSONStore) List(context.Context) ([]app.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readLocked()
}

func (s *JSONStore) Create(_ context.Context, item app.Item) (app.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := s.readLocked()
	if err != nil {
		return app.Item{}, err
	}
	items = append(items, item)
	return item, s.writeLocked(items)
}

func (s *JSONStore) UpdateStatus(_ context.Context, id string, status app.ItemStatus) (app.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := s.readLocked()
	if err != nil {
		return app.Item{}, err
	}
	for i := range items {
		if items[i].ID == id {
			items[i].Status = status
			items[i].UpdatedAt = time.Now().UTC()
			return items[i], s.writeLocked(items)
		}
	}
	return app.Item{}, errors.New("item not found")
}

func (s *JSONStore) readLocked() ([]app.Item, error) {
	content, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var items []app.Item
	if err := json.Unmarshal(content, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *JSONStore) writeLocked(items []app.Item) error {
	content, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, content, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
