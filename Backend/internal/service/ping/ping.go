package ping

import (
	"context"

	"backend/internal/domain"
)

type PingRepo interface {
	CreatePing(ctx context.Context, address, time, latestData string) (*domain.Address, error)
}

type Service struct {
	repo PingRepo
}

func New(repo PingRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreatePing(ctx context.Context, address, time, latestData string) (*domain.Address, error) {
	ping, err := s.repo.CreatePing(ctx, address, time, latestData)
	if err != nil {
		return nil, err
	}

	return ping, nil
}
