package handler

import (
	"context"

	"backend/internal/http/handler/model"
)

type PingService interface {
	CreatePing(ctx context.Context, req model.Address) (*model.Address, error)
	GetPing(ctx context.Context, req model.GetAddressListRequest) (*model.GetAddressListResponse, error)
	GetNumber() (int, error)
}

type Handler struct {
	pingService PingService
}

func New(ping PingService) *Handler {
	return &Handler{
		pingService: ping,
	}
}
