package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"urlshortner/internal/initialize"
	"urlshortner/internal/model"
	"urlshortner/internal/repository"
	"urlshortner/internal/utils"
)

type ShortnerService struct {
	repo        *repository.ShortnerRepo
	logger      *zap.Logger
	config      *initialize.Config
	redisClient *redis.Client
}

type Deps struct {
	Repo        *repository.ShortnerRepo
	Logger      *zap.Logger
	Config      *initialize.Config
	RedisClient *redis.Client
}

func NewShortnerService(deps Deps) *ShortnerService {
	return &ShortnerService{
		repo:        deps.Repo,
		logger:      deps.Logger,
		config:      deps.Config,
		redisClient: deps.RedisClient,
	}
}

func (s *ShortnerService) CreateShortLink(ctx context.Context, url string) (*model.Response, error) {
	existURL, err := s.repo.CheckDuplicate(ctx, url)
	if err != nil && err != repository.ErrLinkNotFound {
		return nil, err
	}

	if existURL != "" {
		return &model.Response{
			URL: "http://" + s.config.HTTPHost + ":" + s.config.HTTPPort + "/" + existURL,
		}, nil
	}

	count, err := s.redisClient.Get(ctx, "count").Result()
	if err == redis.Nil {
		err := s.redisClient.Set(ctx, "count", 1, 0).Err()
		if err != nil {
			return nil, err
		}
		count = "1"
	} else if err != nil {
		return nil, err
	}

	_, err = s.redisClient.Incr(ctx, "count").Result()
	if err != nil {
		s.logger.Error("error increment count to redis", zap.Error(err))
		return nil, err
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		s.logger.Error("error convert count from redis", zap.Error(err))
		return nil, err
	}

	shortenURL := utils.GenShort(countInt)

	err = s.repo.CreateShortLink(ctx, shortenURL, url)
	if err != nil {
		s.logger.Error("error save shortenURL to db", zap.Error(err))
		return nil, err
	}

	return &model.Response{
		URL: "http://" + s.config.HTTPHost + ":" + s.config.HTTPPort + "/" + shortenURL,
	}, nil
}

func (s *ShortnerService) GetOriginalByShort(ctx context.Context, url string) (*model.Response, error) {
	originalURL, err := s.repo.GetOriginalByShort(ctx, url)
	if err != nil {
		s.logger.Error("error getting originalURL from db", zap.Error(err))
		return nil, err
	}

	return &model.Response{
		URL: originalURL,
	}, err
}
