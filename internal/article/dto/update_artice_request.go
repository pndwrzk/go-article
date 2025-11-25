package dto

import "mime/multipart"

type UpdateArticleRequest struct {
	Title        string                  `form:"title" binding:"required"`
	Content      string                  `form:"content" binding:"required"`
	KeepPhotoIDs []string                `form:"keepPhotoIDs"`
	Photos       []*multipart.FileHeader `form:"photos"`
}
