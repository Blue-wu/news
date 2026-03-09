package repository

import (
	"blog-backend/internal/app/model"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id uint) error
	GetByID(id uint) (*model.Category, error)
	GetBySlug(slug string) (*model.Category, error)
	GetAll() ([]model.Category, error)
	GetRootCategories() ([]model.Category, error)
	GetChildrenByParentID(parentID uint) ([]model.Category, error)
	ExistsByName(name string) (bool, error)
	ExistsBySlug(slug string) (bool, error)
}

// categoryRepository 分类仓库实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库实例
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create 创建分类
func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// GetByID 根据ID获取分类
func (r *categoryRepository) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug 根据Slug获取分类
func (r *categoryRepository) GetBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类
func (r *categoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}

// GetRootCategories 获取根分类
func (r *categoryRepository) GetRootCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("parent_id = ?", 0).Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}

// GetChildrenByParentID 根据父ID获取子分类
func (r *categoryRepository) GetChildrenByParentID(parentID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("parent_id = ?", parentID).Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}

// ExistsByName 检查分类名称是否存在
func (r *categoryRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsBySlug 检查分类Slug是否存在
func (r *categoryRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}
