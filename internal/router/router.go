package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pg "github.com/jjaywhitaker/QuoteKeeper/internal/dao/postgres"
	"github.com/jjaywhitaker/QuoteKeeper/internal/model"
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
	r.pgDao.GetQuote()

	quotes := model.QuoteResponse{
		Quote:      "Hello world",
		Author:     "Me",
		Categories: []string{"Random"},
	}
	c.IndentedJSON(http.StatusOK, quotes)
}
