package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	pg "github.com/jjaywhitaker/QuoteKeeper/internal/dao/postgres"
	"github.com/jjaywhitaker/QuoteKeeper/internal/dao/s3"
	r "github.com/jjaywhitaker/QuoteKeeper/internal/router"
	_ "github.com/joho/godotenv/autoload"
)

var (
	updateQuotesFromFile, _ = strconv.ParseBool(os.Getenv("UPDATE_QUOTES_FROM_FILE"))
)

func main() {
	fmt.Println("Starting Quote Keeper")
	pgdao, err := pg.NewPosgtresDao()
	if err != nil {
		log.Fatalf("Failed to create postgres dao. Error: %v", err)
	}

	s3Dao := s3.NewS3Dao()

	if updateQuotesFromFile {
		log.Printf("Attempting to Update quotes in DB from cloud.")
		quotes, err := s3Dao.ReadQuoteFile()
		if err != nil {
			log.Printf("WARNING - Failed to read in the quote file from S3. Error: %v", err)
		} else {
			log.Printf("INFO - Attempting to insert quote list into DB.")
			err = pgdao.UpsertQuoteList(quotes)
			if err != nil {
				log.Printf("WARNING - Failed to update the DB with quotes from S3. Error: %v", err)
			}
		}
	} else {
		log.Printf("Not Importing Quotes.")
	}

	//create routes and run router
	r.NewRouter(pgdao)

}
