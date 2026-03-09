package service

import (
	"blog-backend/internal/app/model"
	"blog-backend/internal/app/repository"
)

type ArticleService struct {
	repo *repository.ArticleRepository
}

type PaginationResult struct {
	Data  []model.Article `json:"data"`
	Total int64           `json:"total"`
}

func NewArticleService(repo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

// 创建文章
func (s *ArticleService) CreateArticle(article *model.Article) error {
	return s.repo.Create(article)
}

// 根据ID获取文章
func (s *ArticleService) GetArticleByID(id uint) (*model.Article, error) {
	return s.repo.GetByID(id)
}

// 获取所有文章
func (s *ArticleService) GetAllArticles() ([]model.Article, error) {
	return s.repo.GetAll()
}

// 分页获取所有文章
func (s *ArticleService) GetAllArticlesPaged(page, pageSize int) (*PaginationResult, error) {
	articles, total, err := s.repo.GetAllPaged(page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PaginationResult{Data: articles, Total: total}, nil
}

// 更新文章
func (s *ArticleService) UpdateArticle(article *model.Article) error {
	// 先检查文章是否存在
	existingArticle, err := s.repo.GetByID(uint(article.ID))
	if err != nil {
		return err
	}

	// 更新字段
	existingArticle.Title = article.Title
	existingArticle.Content = article.Content
	existingArticle.Summary = article.Summary
	// 添加封面图片的更新
	existingArticle.CoverImage = article.CoverImage
	// 也应该更新其他可能的字段
	existingArticle.Status = article.Status
	existingArticle.IsTop = article.IsTop

	return s.repo.Update(existingArticle)
}

// 删除文章
func (s *ArticleService) DeleteArticle(id uint) error {
	// 先检查文章是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// GetDraftArticles 获取所有草稿文章（状态为0）
func (s *ArticleService) GetDraftArticles() ([]model.Article, error) {
	var articles []model.Article
	err := s.repo.Db.Where("status = ?", 0).Find(&articles).Error
	return articles, err
}

// UpdateArticleStatus 更新文章状态
func (s *ArticleService) UpdateArticleStatus(id uint, status int) error {
	return s.repo.Db.Model(&model.Article{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// 根据分类slug获取文章
func (s *ArticleService) GetArticlesByCategorySlug(slug string) ([]model.Article, error) {
	return s.repo.GetByCategorySlug(slug)
}

// 分页根据分类slug获取文章
func (s *ArticleService) GetArticlesByCategorySlugPaged(slug string, page, pageSize int) (*PaginationResult, error) {
	articles, total, err := s.repo.GetByCategorySlugPaged(slug, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PaginationResult{Data: articles, Total: total}, nil
}

// ClearArticleCategories 清除文章的所有分类关联
func (s *ArticleService) ClearArticleCategories(articleID uint) error {
	return s.repo.Db.Where("article_id = ?", articleID).Delete(&model.ArticleCategory{}).Error
}

// AddArticleCategory 添加文章分类关联
func (s *ArticleService) AddArticleCategory(articleID, categoryID uint) error {
	articleCategory := &model.ArticleCategory{
		ArticleID:  uint64(articleID),
		CategoryID: uint64(categoryID),
	}
	return s.repo.Db.Create(articleCategory).Error
}

// IncrementViewCount 增加文章阅读量
func (s *ArticleService) IncrementViewCount(id uint) error {
	return s.repo.IncrementViewCount(id)
}

// SearchArticles 搜索文章
func (s *ArticleService) SearchArticles(keyword string) ([]model.Article, error) {
	return s.repo.SearchArticles(keyword)
}

// SearchArticlesPaged 分页搜索文章
func (s *ArticleService) SearchArticlesPaged(keyword string, page, pageSize int) (*PaginationResult, error) {
	articles, total, err := s.repo.SearchArticlesPaged(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PaginationResult{Data: articles, Total: total}, nil
}

// GetRandomArticles 获取随机文章
func (s *ArticleService) GetRandomArticles(limit int) ([]model.Article, error) {
	return s.repo.GetRandomArticles(limit)
}

// GetPopularArticles 获取热门文章
func (s *ArticleService) GetPopularArticles(limit int) ([]model.Article, error) {
	return s.repo.GetPopularArticles(limit)
}
