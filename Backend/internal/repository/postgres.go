package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"backend/internal/config"
	"backend/internal/http/handler/model"
)

const (
	col1 string = "article_name"
	col2 string = "article_content"
)

// type Repository interface {
// 	CreatePing(ctx context.Context, req model.Address) (*model.Address, error)
// 	GetPing(ctx context.Context, req model.GetAddressListRequest) (*model.GetAddressListResponse, error)
// 	GetNumber() (int, error)
// }

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

func (r *pgRepo) CreatePing(ctx context.Context, req model.Address) (*model.Address, error) {

	res := &model.Address{}
	var id int
	var err error

	if req.ResponseTime == "no answer" {
		err = r.db.QueryRowxContext(ctx, querySetNoAnswer, req.IP, req.ResponseTime).
			Scan(&id, &res.IP, &res.ResponseTime)
	} else {
		err = r.db.QueryRowxContext(ctx, querySetPing, req.IP, req.ResponseTime, req.LastSuccessfulPing).
			Scan(&id, &res.IP, &res.ResponseTime, &res.LastSuccessfulPing)

	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *pgRepo) GetPing(ctx context.Context, req model.GetAddressListRequest) (*model.GetAddressListResponse, error) {

	rows, err := r.db.QueryContext(ctx, queryGetPing, req.Limit, req.Ofset)
	if err != nil {
		return nil, err
	}

	pings := make([]model.Address, 0, req.Limit)

	pingList := model.GetAddressListResponse{Addresses: pings}

	for rows.Next() {
		var ping model.Address
		var id int
		err = rows.Scan(&id, &ping.IP, &ping.ResponseTime, &ping.LastSuccessfulPing)
		if err != nil {
			return nil, err
		}
		pingList.Addresses = append(pingList.Addresses, ping)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pingList, nil
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
