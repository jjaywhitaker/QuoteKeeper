package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaywhitaker/QuoteKeeper/internal/model"
)

type PostgresDAO interface {
	GetRandomQuoteByCategory(category string) (quote model.QuoteResponse, err error)
	InsertQuote(body string, author string, categories []string) error
	UpsertQuoteList(quotes []model.QuoteResponse) error
}

type pgDao struct {
	pool *pgxpool.Pool
}

var (
	dbURL  = os.Getenv("DATABASE_URL")
	dbUser = os.Getenv("DATABASE_USER")
	dbPass = os.Getenv("DATABASE_PASSWORD")
	port   = os.Getenv("DATABASE_PORT")
	dbName = os.Getenv("DATABASE_NAME")
)

func NewPosgtresDao() (PostgresDAO, error) {
	//postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=require&pool_max_conns=10", dbUser, dbPass, dbURL, port, dbName)
	//connString := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=verify-ca pool_max_conns=10", dbUser, dbPass, dbURL, port, dbName)
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return pgDao{}, err
	}
	log.Printf("Postgres connection created successfully")
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
		err = pg.pool.QueryRow(context.Background(), queryString).Scan(&quote.Body, &quote.Author)
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

		err = pg.pool.QueryRow(context.Background(), queryString, category).Scan(&quote.Body, &quote.Author)
		if err != nil {
			log.Printf("Query with category failed: %v\n", err)
		}

		return quote, err
	}

	return quote, err
}

func (pg pgDao) UpsertQuoteList(quotes []model.QuoteResponse) (err error) {
	log.Printf("Upserting quote list into DB.")

	for _, quoteObj := range quotes {
		err = pg.InsertQuote(quoteObj.Body, quoteObj.Author, quoteObj.Categories)
	}
	//TODO:: this error will only return the last error
	return err
}

func (pg pgDao) InsertQuote(body string, author string, categories []string) (err error) {
	for _, category := range categories {
		//insert category
		categoryQuery := `
		insert into categories(category, created_date) values ($1, now()) 
			on conflict (category) do nothing;
		`
		_, err = pg.pool.Exec(context.Background(), categoryQuery, category)
		if err != nil {
			log.Printf("Error - Counldn't add category. Error: %v", err)
		} else {
			log.Printf("Added category successfully.")
		}
		//insert quote
		quoteQueryString := `
			insert into quotes(body, author, created_date) values ($1, $2, now()) 
				on conflict (body, author) do nothing;
			`
		_, err = pg.pool.Exec(context.Background(), quoteQueryString, body, author)
		if err != nil {
			log.Printf("Error - Counldn't add quote. Error: %v", err)
		} else {
			log.Printf("Added quote successfully.")
		}
		//Connect the two
		quoteCatQueryString := `
		insert into quotes_category (quote_id, category_id) 
			values ((select id from quotes q where q.body = $2), (select id from categories c where c.category = $1))
			on conflict (quote_id, category_id) do nothing; 
		`
		_, err = pg.pool.Exec(context.Background(), quoteCatQueryString, category, body)
		if err != nil {
			log.Printf("Error - Counldn't add quote/category. Error: %v", err)
		} else {
			log.Printf("Added quote/category successfully.")
		}
	}
	//TODO:: this error will only return the last error
	return err
}
