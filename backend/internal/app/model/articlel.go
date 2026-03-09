package model

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	UserID       uint64         `json:"user_id" gorm:"not null;type:bigint(20);index"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Title        string         `json:"title" gorm:"not null"`
	Content      string         `json:"content" gorm:"type:text;not null"`
	Summary      string         `json:"summary" gorm:"size:500"`
	CoverImage   string         `json:"cover_image" gorm:"size:255"`
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	LikeCount    int            `json:"like_count" gorm:"default:0"`
	Status       int            `json:"status" gorm:"default:1"`
	IsTop        int            `json:"is_top" gorm:"default:0"`
	Categories   []Category     `json:"categories,omitempty" gorm:"many2many:article_categories;"`
	Tags         []Tag          `json:"tags,omitempty" gorm:"many2many:article_tags;"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Category 分类模型
type Category struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:50;not null;unique"`
	Slug        string    `json:"slug" gorm:"size:50;not null;unique"`
	Description string    `json:"description" gorm:"type:text"`
	ParentID    uint64    `json:"parent_id" gorm:"default:0"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag 标签模型
type Tag struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"size:50;not null;unique"`
	Slug         string    `json:"slug" gorm:"size:50;not null;unique"`
	ArticleCount int       `json:"article_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ArticleCategory 文章分类关联模型
type ArticleCategory struct {
	ArticleID  uint64 `gorm:"primaryKey"`
	CategoryID uint64 `gorm:"primaryKey"`
}

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	ArticleID uint64 `gorm:"primaryKey"`
	TagID     uint64 `gorm:"primaryKey"`
}

// 表名
func (Article) TableName() string {
	return "articles"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}

func (ArticleTag) TableName() string {
	return "article_tags"
}
