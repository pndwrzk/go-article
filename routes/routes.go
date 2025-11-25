package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pndwrzk/go-article/internal/article/handler"
	"github.com/pndwrzk/go-article/internal/article/repository"
	"github.com/pndwrzk/go-article/internal/article/service"
	"github.com/pndwrzk/go-article/pkg/database"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	RegisterArticleRoutes(api)
}

func RegisterArticleRoutes(rg *gin.RouterGroup) {
	repo := repository.NewArticleRepository(database.DB)
	svc := service.NewArticleService(repo)
	h := handler.NewArticleHandler(svc)

	article := rg.Group("/articles")
	{
		article.GET("", h.GetAll)
		article.GET("/:id", h.GetByID)
		article.PUT("/:id", h.Update)
		article.POST("", h.Create)
		article.DELETE("/:id", h.Delete)

	}
}
