package repository

import (
	"context"
	"sync"
)

type InMemoryRepo struct {
	db map[string]string
	mu sync.Mutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		db: make(map[string]string),
	}
}

func (r *InMemoryRepo) CreateShortLink(ctx context.Context, shortURL string, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[shortURL] = originalURL
	return nil
}

func (r *InMemoryRepo) GetOriginalByShort(ctx context.Context, shortURL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if originalURL, ok := r.db[shortURL]; ok {
		return originalURL, nil
	}
	return "", ErrLinkNotFound
}

func (r *InMemoryRepo) CheckDuplicate(ctx context.Context, originalURL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for shortURL, stored := range r.db {
		if stored == originalURL {
			return shortURL, nil
		}
	}
	return "", ErrLinkNotFound
}
