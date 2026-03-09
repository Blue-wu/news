package service

import (
	"blog-backend/internal/app/model"
	"blog-backend/internal/app/repository"
	"errors"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

// CategoryService 分类服务接口
type CategoryService interface {
	CreateCategory(category *model.Category) error
	UpdateCategory(category *model.Category) error
	DeleteCategory(id uint) error
	GetCategoryByID(id uint) (*model.Category, error)
	GetCategoryBySlug(slug string) (*model.Category, error)
	GetAllCategories() ([]model.Category, error)
	GetRootCategories() ([]model.Category, error)
	GetCategoriesWithChildren() ([]model.Category, error)
}

// categoryService 分类服务实现
type categoryService struct {
	repo repository.CategoryRepository
	db   *gorm.DB
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(repo repository.CategoryRepository, db *gorm.DB) CategoryService {
	return &categoryService{
		repo: repo,
		db:   db,
	}
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(category *model.Category) error {
	// 验证分类名称
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("分类名称不能为空")
	}

	// 验证分类名称长度
	if len(category.Name) > 50 {
		return errors.New("分类名称长度不能超过50个字符")
	}

	// 生成Slug
	if category.Slug == "" {
		category.Slug = generateSlug(category.Name)
	}

	// 验证Slug
	if strings.TrimSpace(category.Slug) == "" {
		return errors.New("分类标识不能为空")
	}

	// 检查分类名称是否已存在
	exists, err := s.repo.ExistsByName(category.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("分类名称已存在")
	}

	// 检查Slug是否已存在
	exists, err = s.repo.ExistsBySlug(category.Slug)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("分类标识已存在")
	}

	return s.repo.Create(category)
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(category *model.Category) error {
	// 验证分类是否存在
	existing, err := s.repo.GetByID(uint(category.ID))
	if err != nil {
		return errors.New("分类不存在")
	}

	// 验证分类名称
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("分类名称不能为空")
	}

	// 验证分类名称长度
	if len(category.Name) > 50 {
		return errors.New("分类名称长度不能超过50个字符")
	}

	// 检查分类名称是否已存在（排除当前分类）
	if category.Name != existing.Name {
		exists, err := s.repo.ExistsByName(category.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("分类名称已存在")
		}
	}

	// 生成Slug
	if category.Slug == "" {
		category.Slug = generateSlug(category.Name)
	}

	// 验证Slug
	if strings.TrimSpace(category.Slug) == "" {
		return errors.New("分类标识不能为空")
	}

	// 检查Slug是否已存在（排除当前分类）
	if category.Slug != existing.Slug {
		exists, err := s.repo.ExistsBySlug(category.Slug)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("分类标识已存在")
		}
	}

	return s.repo.Update(category)
}

// DeleteCategory 删除分类
func (s *categoryService) DeleteCategory(id uint) error {
	// 验证分类是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	// 检查是否有子分类
	children, err := s.repo.GetChildrenByParentID(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("该分类下有子分类，请先删除子分类")
	}

	// 检查是否有关联的文章
	var count int64
	err = s.db.Model(&model.Article{}).Joins("JOIN article_categories ON article_categories.article_id = articles.id").Where("article_categories.category_id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该分类下有文章，请先移除文章分类关联")
	}

	return s.repo.Delete(id)
}

// GetCategoryByID 根据ID获取分类
func (s *categoryService) GetCategoryByID(id uint) (*model.Category, error) {
	return s.repo.GetByID(id)
}

// GetCategoryBySlug 根据Slug获取分类
func (s *categoryService) GetCategoryBySlug(slug string) (*model.Category, error) {
	return s.repo.GetBySlug(slug)
}

// GetAllCategories 获取所有分类
func (s *categoryService) GetAllCategories() ([]model.Category, error) {
	return s.repo.GetAll()
}

// GetRootCategories 获取根分类
func (s *categoryService) GetRootCategories() ([]model.Category, error) {
	return s.repo.GetRootCategories()
}

// GetCategoriesWithChildren 获取分类树
func (s *categoryService) GetCategoriesWithChildren() ([]model.Category, error) {
	// 获取所有分类
	allCategories, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// 构建分类映射
	categoryMap := make(map[uint]*model.Category)
	for i := range allCategories {
		categoryMap[uint(allCategories[i].ID)] = &allCategories[i]
	}

	// 构建分类树
	var rootCategories []model.Category
	for i := range allCategories {
		category := &allCategories[i]
		if category.ParentID == 0 {
			// 根分类
			rootCategories = append(rootCategories, *category)
		} else {
			// 子分类，添加到父分类
			if _, exists := categoryMap[uint(category.ParentID)]; exists {
				// 这里可以扩展Category结构体添加Children字段
				// 目前模型中没有Children字段，所以暂时只返回平级列表
			}
		}
	}

	return rootCategories, nil
}

// generateSlug 生成URL友好的Slug
func generateSlug(title string) string {
	// 移除前后空格
	slug := strings.TrimSpace(title)
	// 转换为小写
	slug = strings.ToLower(slug)
	// 替换非字母数字字符为连字符
	result := strings.Builder{}
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else if result.Len() > 0 && result.String()[result.Len()-1] != '-' {
			result.WriteRune('-')
		}
	}
	// 移除末尾的连字符
	slug = strings.TrimRight(result.String(), "-")
	return slug
}
