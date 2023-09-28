package main

import (
	"fmt"
	"log"

	pg "github.com/jjaywhitaker/QuoteKeeper/internal/dao/postgres"
	r "github.com/jjaywhitaker/QuoteKeeper/internal/router"
)

func main() {
	fmt.Println("Starting Quote Keeper")
	pgdao, err := pg.NewPosgtresDao()
	if err != nil {
		log.Fatalf("Failed to create postgres dao. Error: %v", err)
	}

	//create routes and run router
	r.NewRouter(pgdao)

}
