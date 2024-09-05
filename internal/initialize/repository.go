package initialize

import "urlshortner/internal/repository"

type Repositoryes struct {
	Sql      *repository.ShortnerRepo
	InMemory *repository.InMemoryRepo
}

func InitAllRepo(db *DB) *Repositoryes {
	return &Repositoryes{
		Sql:      repository.NewShortnerRepo(db.DB),
		InMemory: repository.NewInMemoryRepo(),
	}
}
