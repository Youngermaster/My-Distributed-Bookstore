package notification

import "sync"

type History struct {
	capacity int
	items    []Message
	mu       sync.RWMutex
}

func NewHistory(capacity int) *History {
	if capacity <= 0 {
		capacity = 1
	}
	return &History{
		capacity: capacity,
		items:    make([]Message, 0, capacity),
	}
}

func (h *History) Add(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.items) >= h.capacity {
		copy(h.items[0:], h.items[1:])
		h.items = h.items[:len(h.items)-1]
	}
	h.items = append(h.items, msg)
}

func (h *History) Recent(limit int) []Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if limit <= 0 || limit > len(h.items) {
		limit = len(h.items)
	}
	result := make([]Message, limit)
	for i := 0; i < limit; i++ {
		result[i] = h.items[len(h.items)-1-i]
	}
	return result
}
