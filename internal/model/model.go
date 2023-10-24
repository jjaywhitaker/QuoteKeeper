package model

type QuoteResponse struct {
	Body       string   `json:"body"`
	Author     string   `json:"author"`
	Categories []string `json:"categories"`
}
