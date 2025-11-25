package dto

import "mime/multipart"

type CreateArticleRequest struct {
	Title   string                  `form:"title" binding:"required"`
	Content string                  `form:"content" binding:"required"`
	Photos  []*multipart.FileHeader `form:"photos"`
}
