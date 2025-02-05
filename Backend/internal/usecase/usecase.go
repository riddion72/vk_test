package usecase

import (
	"errors"

	"main/config"
	"main/internal/models"
	"main/internal/repository"
)

const (
	Limit = 3
)

type Usecase struct {
	pgPepo repository.Repository
}

type UseCase interface {
	GetArticle(models.GetArticleRequest) (models.GetArticleResponse, error)
	SetArticle(models.SetArticleRequest) error
}

func NewUsecase(pgPepo repository.Repository) UseCase {
	return &Usecase{pgPepo: pgPepo}
}

func (u *Usecase) GetArticle(req models.GetArticleRequest) (models.GetArticleResponse, error) {
	var article models.GetArticleResponse
	n, err := u.pgPepo.GetNumber()
	if err != nil {
		return article, err
	}

	req.Limit = Limit
	req.Ofset = req.Limit * (req.Page - 1)

	if (req.Page == 0) || (req.Page > getLast(n, req.Limit)) {
		return article, errors.New("Invalid page")
	}

	article, err = u.pgPepo.GetArticle(req)
	if err != nil {
		return article, err
	}

	article.Last = getLast(n, req.Limit)
	article.Page = req.Page

	return article, nil
}

func (u *Usecase) SetArticle(req models.SetArticleRequest) error {
	err := u.pgPepo.SetArticle(req)
	if err != nil {
		return err
	}
	return nil
}

func getLast(n int, lim int) int {
	if n%lim != 0 {
		return (n / lim) + 1
	}
	return n / lim
}

func Authorization(login string, password string) bool {
	config, _ := config.LoadConfig()
	return login == config.AdminName && password == config.AdminPassword
}
