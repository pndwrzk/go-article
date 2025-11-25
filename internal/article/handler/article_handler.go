package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pndwrzk/go-article/internal/article/dto"
	"github.com/pndwrzk/go-article/internal/article/service"
	"github.com/pndwrzk/go-article/internal/common/constants"
	"github.com/pndwrzk/go-article/internal/common/response"
	"github.com/pndwrzk/go-article/internal/common/utils"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service}
}

func (h *ArticleHandler) GetAll(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	articles, meta, err := h.service.GetAll(ctx, page, limit)
	if err != nil {
		response.Error(ctx, constants.CodeInternalServerError, "Failed to get articles")
		return
	}

	response.SuccessWithMeta(
		ctx,
		"Successfully get article list",
		articles,
		meta,
	)
}

func (h *ArticleHandler) GetByID(ctx *gin.Context) {

	id := ctx.Param("id")

	article, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.Error(ctx, constants.CodeNotFound, "Article not found")
		return
	}

	response.Success(ctx, "Successfully get article detail", article)
}

func (h *ArticleHandler) Create(ctx *gin.Context) {

	var req dto.CreateArticleRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	err := h.service.Create(ctx, req.Title, req.Content, req.Photos)
	if err != nil {
		response.Error(ctx, constants.CodeInternalServerError, "Failed to create article")
		return
	}

	response.Success(ctx, "Article created successfully", nil)
}

func (h *ArticleHandler) Delete(ctx *gin.Context) {

	id := ctx.Param("id")

	if err := h.service.Delete(ctx, id); err != nil {
		response.Error(ctx, constants.CodeInternalServerError, "Failed to delete article")
		return
	}

	response.Success(ctx, "Article deleted successfully", nil)
}

func (h *ArticleHandler) Update(ctx *gin.Context) {

	id := ctx.Param("id")

	var req dto.UpdateArticleRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	if err := h.service.Update(ctx, id, req.Title, req.Content, req.KeepPhotoIDs, req.Photos); err != nil {
		response.Error(ctx, constants.CodeInternalServerError, "Failed to update article")
		return
	}

	response.Success(ctx, "Article updated successfully", nil)
}
