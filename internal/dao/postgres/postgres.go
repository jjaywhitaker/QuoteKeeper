package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaywhitaker/QuoteKeeper/internal/model"
)

type PostgresDAO interface {
	GetQuote(id int) (quote model.QuoteResponse, err error)
	InsertQuote(body string, author string, categories []string) error
}

type pgDao struct {
	pool *pgxpool.Pool
}

func NewPosgtresDao() (PostgresDAO, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return pgDao{}, err
	}
	return pgDao{
		pool: dbpool,
	}, nil
}

func (pg pgDao) GetQuote(id int) (quote model.QuoteResponse, err error) { // (body string, author string, categories []string, err error) {
	err = pg.pool.QueryRow(context.Background(), "select body, author, categories from quotes where id=$1", id).Scan(&quote.Quote, &quote.Author, &quote.Categories)
	//TODO do we need to marshal categories
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
	}
	return quote, err
}

func (pg pgDao) InsertQuote(body string, author string, categories []string) error {
	//TODO: unmarshal categories
	_, err := pg.pool.Exec(context.Background(), "insert into quotes (body, author, categories, created_date) values ($1, $2, $3, now())", body, author, categories)
	return err
}
