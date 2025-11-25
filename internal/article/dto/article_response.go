package dto

type ArticleResponse struct {
	ID      string          `json:"id"`
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Photos  []PhotoResponse `json:"photos"`
}
