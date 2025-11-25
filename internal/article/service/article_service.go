package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pndwrzk/go-article/internal/article/dto"
	"github.com/pndwrzk/go-article/internal/article/model"
	"github.com/pndwrzk/go-article/internal/article/repository"
	dtoCommon "github.com/pndwrzk/go-article/internal/common/dto"
	"github.com/pndwrzk/go-article/internal/common/utils"
)

type ArticleService interface {
	GetAll(ctx context.Context, page, limit int) ([]dto.ArticleResponse, *dtoCommon.Meta, error)
	GetByID(ctx context.Context, id string) (*dto.ArticleResponse, error)
	Create(ctx context.Context, title, content string, files []*multipart.FileHeader) error
	Update(ctx context.Context, id, title, content string, keepPhotoIDs []string, newFiles []*multipart.FileHeader) error
	Delete(ctx context.Context, id string) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) GetAll(ctx context.Context, page, limit int) ([]dto.ArticleResponse, *dtoCommon.Meta, error) {
	articles, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]dto.ArticleResponse, 0, len(articles))
	for _, article := range articles {
		responses = append(responses, fromArticleModel(article, ctx))
	}

	totalPages := (int(total) + limit - 1) / limit
	meta := &dtoCommon.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return responses, meta, nil
}

func (s *articleService) GetByID(ctx context.Context, id string) (*dto.ArticleResponse, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := fromArticleModel(*article, ctx)
	return &resp, nil
}

func fromArticleModel(m model.Article, ctx context.Context) dto.ArticleResponse {
	hostInfo := ctx.Value("hostInfo").(string)
	photos := make([]dto.PhotoResponse, 0, len(m.Photos))
	for _, p := range m.Photos {
		photos = append(photos, dto.PhotoResponse{
			ID:  p.ID.String(),
			URL: fmt.Sprintf("%s/uploads/article/%s", hostInfo, filepath.Base(p.Path)),
		})
	}

	return dto.ArticleResponse{
		ID:      m.ID.String(),
		Title:   m.Title,
		Content: m.Content,
		Photos:  photos,
	}
}

func saveFiles(articleID uuid.UUID, files []*multipart.FileHeader) ([]model.Photo, []string, error) {
	if len(files) == 0 {
		return nil, nil, nil
	}

	if err := os.MkdirAll("uploads/article", os.ModePerm); err != nil {
		return nil, nil, err
	}

	photos := make([]model.Photo, 0, len(files))
	savedPaths := make([]string, 0, len(files))

	for _, file := range files {
		filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))
		savePath := filepath.Join("uploads/article", filename)

		if err := utils.SaveMultipartFile(file, savePath); err != nil {
			utils.CleanupFiles(savedPaths)
			return nil, nil, err
		}

		savedPaths = append(savedPaths, savePath)
		photos = append(photos, model.Photo{
			ArticleID: articleID,
			Path:      savePath,
		})
	}

	return photos, savedPaths, nil
}

func (s *articleService) Create(ctx context.Context, title, content string, files []*multipart.FileHeader) error {
	article := model.Article{
		Title:   title,
		Content: content,
	}

	return s.repo.Transaction(ctx, func(tx repository.ArticleRepository) error {
		if err := tx.CreateArticle(ctx, &article); err != nil {
			return err
		}

		photos, savedPaths, err := saveFiles(article.ID, files)
		if err != nil {
			return err
		}

		if len(photos) > 0 {
			if err := tx.CreatePhotos(ctx, photos); err != nil {
				utils.CleanupFiles(savedPaths)
				return err
			}
		}

		return nil
	})
}
func (s *articleService) Update(
	ctx context.Context,
	id, title, content string,
	keepPhotoIDs []string,
	newFiles []*multipart.FileHeader,
) error {
	return s.repo.Transaction(ctx, func(tx repository.ArticleRepository) error {

		article, err := tx.GetByID(ctx, id)
		if err != nil {
			return err
		}

		toDelete := filterPhotosToDelete(article.Photos, keepPhotoIDs)

		if len(toDelete) > 0 {
			delPaths := make([]string, 0, len(toDelete))
			toDeleteIDs := make([]string, 0, len(toDelete))

			for _, p := range toDelete {
				delPaths = append(delPaths, p.Path)
				toDeleteIDs = append(toDeleteIDs, p.ID.String())
			}

			if err := tx.DeletePhotosByIDs(ctx, toDeleteIDs); err != nil {
				return err
			}

			utils.CleanupFiles(delPaths)
		}

		newPhotos, savedPaths, err := saveFiles(article.ID, newFiles)
		if err != nil {
			return err
		}
		if len(newPhotos) > 0 {
			if err := tx.CreatePhotos(ctx, newPhotos); err != nil {
				utils.CleanupFiles(savedPaths)
				return err
			}
		}
		articleUpdate := model.Article{
			ID:      article.ID,
			Title:   title,
			Content: content,
		}
		if err := tx.UpdateArticle(ctx, &articleUpdate); err != nil {
			return err
		}

		return nil
	})
}

func filterPhotosToDelete(photos []model.Photo, keepIDs []string) []model.Photo {
	toDelete := []model.Photo{}
	for _, p := range photos {
		if !containsID(keepIDs, p.ID.String()) {
			toDelete = append(toDelete, p)
		}
	}
	return toDelete
}

func containsID(ids []string, id string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func extractPaths(photos []model.Photo) []string {
	paths := make([]string, 0, len(photos))
	for _, p := range photos {
		paths = append(paths, p.Path)
	}
	return paths
}

func (s *articleService) Delete(ctx context.Context, id string) error {
	return s.repo.Transaction(ctx, func(tx repository.ArticleRepository) error {
		article, err := tx.GetByID(ctx, id)
		if err != nil {
			return err
		}

		paths := make([]string, 0, len(article.Photos))
		for _, p := range article.Photos {
			paths = append(paths, p.Path)
		}

		if err := tx.DeletePhotosByArticleID(ctx, id); err != nil {
			return err
		}

		if err := tx.DeleteArticleByID(ctx, id); err != nil {
			return err
		}

		utils.CleanupFiles(paths)
		return nil
	})
}
