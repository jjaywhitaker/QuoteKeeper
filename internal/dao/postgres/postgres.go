package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaywhitaker/QuoteKeeper/internal/model"
)

type PostgresDAO interface {
	GetRandomQuoteByCategory(category string) (quote model.QuoteResponse, err error)
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

func (pg pgDao) GetRandomQuoteByCategory(category string) (quote model.QuoteResponse, err error) {

	if category == "Random" {
		queryString := `select q.body, q.author from 
						quotes q 
						order by random()
						limit 1`
		err = pg.pool.QueryRow(context.Background(), queryString).Scan(&quote.Quote, quote.Author)
		if err != nil {
			log.Printf("Query with no category failed: %v\n", err)
		}

	} else {
		queryString := `select q.body, q.author from 
						quotes q 
						join quotes_category qc on q.id = qc.quote_id
						join categories c on qc.category_id = c.id
						where c.category = $1
						order by random()
						limit 1`

		err = pg.pool.QueryRow(context.Background(), queryString, category).Scan(&quote.Quote, &quote.Author)
		if err != nil {
			log.Printf("Query with category failed: %v\n", err)
		}

		return quote, err
	}

	return quote, err
}

func (pg pgDao) InsertQuote(body string, author string, categories []string) error {
	//TODO: unmarshal categories
	_, err := pg.pool.Exec(context.Background(), "insert into quotes (body, author, categories, created_date) values ($1, $2, $3, now())", body, author, categories)
	return err
}
