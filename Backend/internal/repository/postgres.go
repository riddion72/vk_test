package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"backend/internal/config"
	"backend/internal/domain"
)

const (
	col1 string = "article_name"
	col2 string = "article_content"
)

type pgRepo struct {
	db *sqlx.DB
}

func New(conn *sqlx.DB) *pgRepo {
	return &pgRepo{db: conn}
}

func NewConnection(ctx context.Context, cfg config.DB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *pgRepo) CreatePing(ctx context.Context, address, time, latestData string) (*domain.Address, error) {

	res := &domain.Address{}

	err := r.db.QueryRowxContext(ctx, querySetPing, address, time, latestData).
		Scan(&res.Id, &res.IP, &res.ResponseTime, &res.LastSuccessfulPing)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *pgRepo) GetArticle(req models.GetArticleRequest) (models.GetArticleResponse, error) {
	var article models.GetArticleResponse
	id := ""

	row, err := r.db.Query(queryGetArticle, req.Limit, req.Ofset)
	if err != nil {
		log.Println("Error getting")
		return article, err
	}

	defer row.Close()

	for i := 0; i < req.Limit && row.Next(); i++ {

		var articleItem models.Article = models.Article{Article_name: "", Article_content: ""}

		err := row.Scan(&id, &articleItem.Article_name, &articleItem.Article_content)
		if err != nil {
			log.Println("Error scanning")
			return article, err
		}
		article.Response = append(article.Response, articleItem)
	}

	return article, nil
}
func (r *pgRepo) GetNumber() (int, error) {
	var number int

	row, err := r.db.Query(queryGetNumber)

	defer row.Close()
	if err != nil {
		log.Println("Error getting number")
		return number, err
	}
	for row.Next() {
		row.Scan(&number)
	}

	return number, nil
}
