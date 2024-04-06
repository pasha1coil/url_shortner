package repository

import (
	"context"
	"database/sql"
	"errors"
)

type ShortnerRepo struct {
	db *sql.DB
}

func NewShortnerRepo(db *sql.DB) *ShortnerRepo {
	return &ShortnerRepo{
		db: db,
	}
}

func (s *ShortnerRepo) CreateShortLink(ctx context.Context, shortURL string, originalURL string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO short_links (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	return err
}

func (s *ShortnerRepo) GetOriginalByShort(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRowContext(ctx, "SELECT original_url FROM short_links WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return originalURL, nil
}

func (s *ShortnerRepo) CheckDuplicate(ctx context.Context, originalURL string) (string, error) {
	var shortURL string
	err := s.db.QueryRowContext(ctx, "SELECT short_url FROM short_links WHERE original_url = $1", originalURL).Scan(&shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return shortURL, nil
}
