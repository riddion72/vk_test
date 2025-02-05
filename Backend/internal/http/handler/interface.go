package handler

import (
	"context"

	"backend/internal/domain"
)

type PingService interface {
	CreatePing(ctx context.Context, address, time, latestData string) (*domain.Address, error)
}

type Handler struct {
	pingService PingService
}

func New(ping PingService) *Handler {
	return &Handler{
		pingService: ping,
	}
}
