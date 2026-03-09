package repository

import (
	"blog-backend/internal/app/model"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	Db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{Db: db}
}

// 创建文章
func (r *ArticleRepository) Create(article *model.Article) error {
	return r.Db.Create(article).Error
}

// 根据ID获取文章
func (r *ArticleRepository) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.Db.Preload("User").Preload("Categories").Preload("Tags").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// 获取所有文章
func (r *ArticleRepository) GetAll() ([]model.Article, error) {
	var articles []model.Article
	err := r.Db.Preload("User").
		Preload("Categories").
		Preload("Tags").
		Order("created_at desc").
		Find(&articles).Error
	return articles, err
}

// 分页获取所有文章
func (r *ArticleRepository) GetAllPaged(page, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	// 计算总数
	r.Db.Model(&model.Article{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.Db.Preload("User").
		Preload("Categories").
		Preload("Tags").
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&articles).Error

	return articles, total, err
}

// 更新文章
func (r *ArticleRepository) Update(article *model.Article) error {
	return r.Db.Save(article).Error
}

// 删除文章
func (r *ArticleRepository) Delete(id uint) error {
	return r.Db.Delete(&model.Article{}, id).Error
}

// 根据分类slug获取文章
func (r *ArticleRepository) GetByCategorySlug(slug string) ([]model.Article, error) {
	var articles []model.Article
	err := r.Db.Joins("JOIN article_categories ON article_categories.article_id = articles.id").
		Joins("JOIN categories ON categories.id = article_categories.category_id").
		Where("categories.slug = ?", slug).
		Preload("User").
		Preload("Categories").
		Order("articles.created_at desc").
		Find(&articles).Error
	return articles, err
}

// 分页根据分类slug获取文章
func (r *ArticleRepository) GetByCategorySlugPaged(slug string, page, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	// 构建基础查询
	query := r.Db.Model(&model.Article{}).
		Joins("JOIN article_categories ON article_categories.article_id = articles.id").
		Joins("JOIN categories ON categories.id = article_categories.category_id").
		Where("categories.slug = ?", slug)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Preload("Categories").
		Order("articles.created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&articles).Error

	return articles, total, err
}

// IncrementViewCount 增加文章阅读量
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.Db.Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// SearchArticles 搜索文章
func (r *ArticleRepository) SearchArticles(keyword string) ([]model.Article, error) {
	var articles []model.Article
	query := r.Db.Model(&model.Article{}).
		Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
					"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Where("status = ?", 1). // 只搜索已发布的文章
		Preload("User").
		Preload("Categories").
		Preload("Tags").
		Order("created_at desc")

	err := query.Find(&articles).Error
	return articles, err
}

// SearchArticlesPaged 分页搜索文章
func (r *ArticleRepository) SearchArticlesPaged(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	// 构建基础查询
	baseQuery := r.Db.Model(&model.Article{}).
		Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Where("status = ?", 1)

	// 计算总数
	baseQuery.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := baseQuery.Preload("User").
		Preload("Categories").
		Preload("Tags").
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&articles).Error

	return articles, total, err
}

// GetRandomArticles 获取随机文章
func (r *ArticleRepository) GetRandomArticles(limit int) ([]model.Article, error) {
	var articles []model.Article
	query := r.Db.Model(&model.Article{}).
		Where("status = ?", 1). // 只获取已发布的文章
		Preload("User").
		Preload("Categories").
		Preload("Tags").
		Order("RAND()"). // 使用数据库的随机函数
		Limit(limit)

	err := query.Find(&articles).Error
	return articles, err
}

// GetPopularArticles 获取热门文章（按阅读量排序）
func (r *ArticleRepository) GetPopularArticles(limit int) ([]model.Article, error) {
	var articles []model.Article
	query := r.Db.Model(&model.Article{}).
		Where("status = ?", 1).  // 只获取已发布的文章
		Where("view_count > 0"). // 筛选有阅读量的文章
		Preload("User").
		Preload("Categories").
		Preload("Tags").
		Order("view_count DESC, created_at DESC"). // 按阅读量降序，创建时间降序
		Limit(limit)

	err := query.Find(&articles).Error
	return articles, err
}
