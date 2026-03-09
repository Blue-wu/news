// backend/internal/app/repository/sitemap_repository.go
package repository

import (
	"blog-backend/internal/app/model"
	"encoding/json"
	"gorm.io/gorm"
)

// SitemapRepository 站点地图仓库接口
type SitemapRepository interface {
	GetSettings() (map[string]interface{}, error)
	SaveSettings(settings map[string]interface{}) error
}

// sitemapRepository 站点地图仓库实现
type sitemapRepository struct {
	db *gorm.DB
}

// NewSitemapRepository 创建站点地图仓库实例
func NewSitemapRepository(db *gorm.DB) SitemapRepository {
	return &sitemapRepository{db: db}
}

// GetSettings 获取站点地图设置（从数据库）
func (r *sitemapRepository) GetSettings() (map[string]interface{}, error) {
	var setting model.SiteSetting
	result := r.db.Where("key_name = ?", "sitemap_settings").First(&setting)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 返回默认配置
			return map[string]interface{}{
				"cacheEnabled":  true,
				"cacheDuration": 86400,
			}, nil
		}
		return nil, result.Error
	}

	// 解析JSON字符串为map
	var settingsMap map[string]interface{}
	if err := json.Unmarshal([]byte(setting.Value), &settingsMap); err != nil {
		// 如果解析失败，返回默认配置
		return map[string]interface{}{
			"cacheEnabled":  true,
			"cacheDuration": 86400,
		}, nil
	}

	return settingsMap, nil
}

// SaveSettings 保存站点地图设置（到数据库）
func (r *sitemapRepository) SaveSettings(settings map[string]interface{}) error {
	// 将map转换为JSON字符串
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// 使用upsert操作：如果存在则更新，不存在则插入
	var setting model.SiteSetting
	result := r.db.Where("key_name = ?", "sitemap_settings").First(&setting)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 插入新记录
			setting = model.SiteSetting{
				KeyName:     "sitemap_settings",
				Value:       string(settingsJSON),
				Description: "站点地图配置",
			}
			return r.db.Create(&setting).Error
		}
		return result.Error
	}

	// 更新现有记录
	setting.Value = string(settingsJSON)
	return r.db.Save(&setting).Error
}
