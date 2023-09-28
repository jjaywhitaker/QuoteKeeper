package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pg "github.com/jjaywhitaker/QuoteKeeper/internal/dao/postgres"
)

type router interface {
}

type apiRouter struct {
	pgDao pg.PostgresDAO
}

func NewRouter(pg pg.PostgresDAO) apiRouter {

	r := apiRouter{pgDao: pg}

	router := gin.Default()
	router.GET("/quotes", r.getQuotes)
	router.GET("/health", r.healthCheck)

	router.Run("localhost:8080")
	return r
}

func (r apiRouter) healthCheck(c *gin.Context) {
	c.String(200, "OK")
}

func (r apiRouter) getQuotes(c *gin.Context) {
	cat := c.Query("category")
	log.Printf("Getting quote for category: %v", cat)
	quote, err := r.pgDao.GetRandomQuoteByCategory(cat)
	if err != nil {

	}

	c.IndentedJSON(http.StatusOK, quote)
}
