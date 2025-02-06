package usecase

import (
	"context"
	"errors"

	"backend/internal/http/handler/model"
)

const (
	Limit = 3
)

type UseCase interface {
	CreatePing(ctx context.Context, req model.Address) (*model.Address, error)
	GetPing(ctx context.Context, req model.GetAddressListRequest) (*model.GetAddressListResponse, error)
	GetNumber() (int, error)
}

type Usecase struct {
	pgPepo UseCase
}

func New(pgPepo UseCase) UseCase {
	return &Usecase{pgPepo: pgPepo}
}

func (u *Usecase) GetPing(ctx context.Context, req model.GetAddressListRequest) (*model.GetAddressListResponse, error) {
	var pings *model.GetAddressListResponse
	n, err := u.GetNumber()
	if err != nil {
		return pings, err
	}

	req.Limit = Limit

	req.Ofset = req.Limit * (req.Page - 1)

	if (req.Page == 0) || (req.Page > getLast(n, req.Limit)) {
		return pings, errors.New("Invalid page")
	}

	pings, err = u.pgPepo.GetPing(ctx, req)
	if err != nil {
		return pings, err
	}

	// fmt.Println(req.Page, getLast(n, req.Limit))

	pings.Last = getLast(n, req.Limit)
	pings.Page = req.Page

	return pings, nil
}

func (u *Usecase) CreatePing(ctx context.Context, req model.Address) (*model.Address, error) {
	resp, err := u.pgPepo.CreatePing(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func getLast(n int, lim int) int {
	if n%lim != 0 {
		return (n / lim) + 1
	}
	return n / lim
}

func (u *Usecase) GetNumber() (int, error) {
	return u.pgPepo.GetNumber()
}

// func Authorization(login string, password string) bool {
// 	config, _ := config.LoadConfig()
// 	return login == config.AdminName && password == config.AdminPassword
// }
