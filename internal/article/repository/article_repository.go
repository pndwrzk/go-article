package repository

import (
	"context"

	"github.com/pndwrzk/go-article/internal/article/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]model.Article, int64, error)
	GetByID(ctx context.Context, id string) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) error
	CreatePhotos(ctx context.Context, photos []model.Photo) error
	DeleteArticleByID(ctx context.Context, id string) error
	DeletePhotosByArticleID(ctx context.Context, articleID string) error
	UpdateArticle(ctx context.Context, article *model.Article) error
	DeletePhotosByIDs(ctx context.Context, ids []string) error
	Transaction(ctx context.Context, f func(txRepo ArticleRepository) error) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) GetAll(ctx context.Context, page, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Preload("Photos").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *articleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
	var article model.Article

	err := r.db.WithContext(ctx).
		Preload("Photos").
		Where("id = ?", id).
		First(&article).Error

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *articleRepository) DeletePhotosByArticleID(ctx context.Context, articleID string) error {
	return r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.Photo{}).Error
}

func (r *articleRepository) DeleteArticleByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Article{}).Error
}

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) CreatePhotos(ctx context.Context, photos []model.Photo) error {
	return r.db.WithContext(ctx).Create(&photos).Error
}

func (r *articleRepository) Transaction(ctx context.Context, f func(txRepo ArticleRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &articleRepository{db: tx}
		return f(txRepo)
	})
}

func (r *articleRepository) UpdateArticle(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepository) DeletePhotosByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Photo{}).Error
}
