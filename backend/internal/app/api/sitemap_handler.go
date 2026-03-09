// backend/internal/app/api/sitemap_handler.go
package api

import (
	"blog-backend/internal/app/config"
	"blog-backend/internal/app/model"
	"blog-backend/internal/app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SitemapHandler 站点地图处理器
type SitemapHandler struct {
	SitemapService   *service.SitemapService
	BaiduPushService *service.BaiduPushService
}

// NewSitemapHandler 创建站点地图处理器实例
func NewSitemapHandler(sitemapService *service.SitemapService, baiduPushService *service.BaiduPushService) *SitemapHandler {
	return &SitemapHandler{
		SitemapService:   sitemapService,
		BaiduPushService: baiduPushService,
	}
}

// GetSitemap 获取站点地图XML
func (h *SitemapHandler) GetSitemap(c *gin.Context) {
	// 记录访问
	h.SitemapService.RecordVisit()

	// 获取站点地图内容
	sitemapContent := h.SitemapService.GetSitemap()

	// 设置响应头
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemapContent)
}

// GetSettings 获取站点地图设置
func (h *SitemapHandler) GetSettings(c *gin.Context) {
	settings := h.SitemapService.GetSitemapSettings()
	stats := h.SitemapService.GetSitemapStats()
	cacheStatus := h.SitemapService.GetCacheStatus()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"settings": gin.H{
			"cacheEnabled":  settings.CacheEnabled,
			"cacheDuration": settings.CacheDuration,
			"cacheStatus":   cacheStatus,
		},
		"statistics": gin.H{
			"visitCount":    stats.VisitCount,
			"lastVisitTime": stats.LastVisitTime.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

// UpdateSettings 更新站点地图设置
func (h *SitemapHandler) UpdateSettings(c *gin.Context) {
	var req struct {
		CacheEnabled  bool `json:"cacheEnabled"`
		CacheDuration int  `json:"cacheDuration" binding:"min=60"` // 最小60秒
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	h.SitemapService.UpdateSitemapSettings(service.SitemapSettings{
		CacheEnabled:  req.CacheEnabled,
		CacheDuration: req.CacheDuration,
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "站点地图设置已更新",
	})
}

// RegenerateSitemap 重新生成站点地图
func (h *SitemapHandler) RegenerateSitemap(c *gin.Context) {
	// 重新生成站点地图
	h.SitemapService.RegenerateSitemap()

	// 异步推送所有已发布的文章到百度
	go func() {
		if h.BaiduPushService != nil {
			// 获取所有已发布的文章
			var articles []model.Article
			if err := config.DB.Where("status = ?", 1).Find(&articles).Error; err == nil {
				// 构建文章URL列表
				var urls []string
				for _, article := range articles {
					url := fmt.Sprintf("%s/articles/%d", config.LoadConfig().BaseURL, article.ID)
					urls = append(urls, url)
				}

				// 推送所有URL到百度
				if len(urls) > 0 {
					if _, err := h.BaiduPushService.PushURLs(urls); err != nil {
						// 非阻塞推送，只记录错误日志
						fmt.Printf("批量推送文章到百度失败: %v\n", err)
					}
				}
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "站点地图已重新生成",
	})
}

// ClearCache 清除站点地图缓存
func (h *SitemapHandler) ClearCache(c *gin.Context) {
	h.SitemapService.ClearCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "站点地图缓存已清除",
	})
}

// GetStats 获取站点地图统计
func (h *SitemapHandler) GetStats(c *gin.Context) {
	stats := h.SitemapService.GetSitemapStats()
	cacheStatus := h.SitemapService.GetCacheStatus()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"statistics": gin.H{
			"visitCount":    stats.VisitCount,
			"lastVisitTime": stats.LastVisitTime.Format("2006-01-02T15:04:05Z07:00"),
		},
		"cache": cacheStatus,
	})
}
