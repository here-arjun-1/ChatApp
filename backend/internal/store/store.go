package store

import (
	"sync"
	"time"

	"chatapp/internal/models"
)

type Store struct {
	mu       sync.RWMutex
	messages []models.Message
	nextID   int
}

func New() *Store {
	return &Store{
		messages: make([]models.Message, 0),
		nextID:   1,
	}
}

func (s *Store) Add(from, content string) models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := models.Message{
		ID:        s.nextID,
		From:      from,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.messages = append(s.messages, msg)
	return msg
}

func (s *Store) History() []models.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Message, len(s.messages))
	copy(out, s.messages)
	return out
}

// MarkRead flags the given message IDs as read and returns the IDs that
// actually existed and were newly marked (so callers can skip broadcasting
// no-op read receipts).
func (s *Store) MarkRead(ids []int) []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	want := make(map[int]bool, len(ids))
	for _, id := range ids {
		want[id] = true
	}

	updated := make([]int, 0, len(ids))
	for i := range s.messages {
		if want[s.messages[i].ID] && !s.messages[i].Read {
			s.messages[i].Read = true
			updated = append(updated, s.messages[i].ID)
		}
	}
	return updated
}
