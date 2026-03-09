// backend/internal/app/service/sitemap_service.go
package service

import (
	"blog-backend/internal/app/config"
	"blog-backend/internal/app/model"
	"fmt"
	"log"
	"sync"
	"time"
)

// SitemapSettings 站点地图配置
type SitemapSettings struct {
	CacheEnabled  bool
	CacheDuration int // 缓存时长（秒）
}

// SitemapStats 站点地图统计数据
type SitemapStats struct {
	VisitCount    int64
	LastVisitTime time.Time
}

// SitemapCache 站点地图缓存
type SitemapCache struct {
	Content string
	Time    time.Time
}

// SitemapService 站点地图服务
type SitemapService struct {
	baseURL         string
	cache           *SitemapCache
	stats           *SitemapStats
	cacheMutex      sync.RWMutex
	statsMutex      sync.Mutex
	cacheDuration   time.Duration
	baiduPushService *BaiduPushService
}

// NewSitemapService 创建站点地图服务实例
func NewSitemapService(baseURL string, baiduPushService *BaiduPushService) *SitemapService {
	return &SitemapService{
		baseURL:          baseURL,
		cache:            &SitemapCache{},
		stats:            &SitemapStats{},
		cacheDuration:    24 * time.Hour, // 默认缓存24小时
		baiduPushService: baiduPushService,
	}
}

// SetBaiduPushService 设置百度推送服务
func (s *SitemapService) SetBaiduPushService(baiduPushService *BaiduPushService) {
	s.baiduPushService = baiduPushService
}

// GetSitemapSettings 获取站点地图设置
func (s *SitemapService) GetSitemapSettings() SitemapSettings {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	return SitemapSettings{
		CacheEnabled:  true,
		CacheDuration: int(s.cacheDuration.Seconds()),
	}
}

// UpdateSitemapSettings 更新站点地图设置
func (s *SitemapService) UpdateSitemapSettings(settings SitemapSettings) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cacheDuration = time.Duration(settings.CacheDuration) * time.Second
}

// GetSitemapStats 获取站点地图统计数据
func (s *SitemapService) GetSitemapStats() SitemapStats {
	s.statsMutex.Lock()
	defer s.statsMutex.Unlock()

	return *s.stats
}

// RecordVisit 记录站点地图访问
func (s *SitemapService) RecordVisit() {
	s.statsMutex.Lock()
	defer s.statsMutex.Unlock()

	s.stats.VisitCount++
	s.stats.LastVisitTime = time.Now()
}

// GetSitemap 获取站点地图内容
func (s *SitemapService) GetSitemap() string {
	// 检查缓存
	if cachedContent, found := s.getFromCache(); found {
		return cachedContent
	}

	// 生成新内容
	content := s.generateSitemapXML()

	// 更新缓存
	s.updateCache(content)

	return content
}

// RegenerateSitemap 重新生成站点地图
func (s *SitemapService) RegenerateSitemap() string {
	content := s.generateSitemapXML()
	s.updateCache(content)
	
	// 如果设置了百度推送服务，异步推送所有已发布的文章链接到百度
	if s.baiduPushService != nil {
		go func() {
			// 获取所有已发布的文章
			var articles []model.Article
			if err := config.DB.Where("status = ?", 1).Find(&articles).Error; err == nil {
				// 构建文章URL列表
				var urls []string
				for _, article := range articles {
					url := fmt.Sprintf("%s/articles/%d", s.baseURL, article.ID)
					urls = append(urls, url)
				}
				
				// 推送所有URL到百度
				if len(urls) > 0 {
					if _, err := s.baiduPushService.PushURLs(urls); err != nil {
						log.Printf("批量推送文章到百度失败: %v\n", err)
					}
				}
			} else {
				log.Printf("获取已发布文章失败: %v\n", err)
			}
		}()
	}
	
	return content
}

// ClearCache 清除缓存
func (s *SitemapService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache.Content = ""
	s.cache.Time = time.Time{}
}

// GetCacheStatus 获取缓存状态
func (s *SitemapService) GetCacheStatus() map[string]interface{} {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	isCached := s.cache.Content != ""
	var cacheAge time.Duration
	var cachedAt string

	if isCached {
		cacheAge = time.Since(s.cache.Time)
		cachedAt = s.cache.Time.Format(time.RFC3339)
	}

	return map[string]interface{}{
		"isCached": isCached,
		"cacheAge": cacheAge.String(),
		"cachedAt": cachedAt,
	}
}

// getFromCache 从缓存获取站点地图
func (s *SitemapService) getFromCache() (string, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	if s.cache.Content == "" {
		return "", false
	}

	// 检查缓存是否过期
	if time.Since(s.cache.Time) > s.cacheDuration {
		return "", false
	}

	return s.cache.Content, true
}

// updateCache 更新缓存
func (s *SitemapService) updateCache(content string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache.Content = content
	s.cache.Time = time.Now()
}

// generateSitemapXML 生成站点地图XML
// backend/internal/app/service/sitemap_service.go
func (s *SitemapService) generateSitemapXML() string {
	var xmlContent string
	xmlContent += "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	xmlContent += "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n"

	// 添加首页URL
	xmlContent += s.generateUrlElement(s.baseURL, "", "daily", 1.0)

	// 添加所有分类页面
	var categories []model.Category
	if err := config.DB.Find(&categories).Error; err == nil {
		for _, category := range categories {
			categoryURL := fmt.Sprintf("%s/category/%s", s.baseURL, category.Slug)
			xmlContent += s.generateUrlElement(categoryURL, "", "weekly", 0.8)
		}
	}

	// 添加所有已发布的文章页面
	var articles []model.Article
	if err := config.DB.Where("status = ?", 1).Find(&articles).Error; err == nil {
		for _, article := range articles {
			articleURL := fmt.Sprintf("%s/articles/%d", s.baseURL, article.ID)
			lastMod := article.UpdatedAt.Format(time.RFC3339)
			xmlContent += s.generateUrlElement(articleURL, lastMod, "monthly", 0.9)
		}
	}

	// 关闭XML根标签
	xmlContent += "</urlset>\n"
	return xmlContent
}

// generateUrlElement 生成URL元素
func (s *SitemapService) generateUrlElement(loc, lastMod, changefreq string, priority float64) string {
	var element string
	element += "  <url>\n"
	element += fmt.Sprintf("    <loc>%s</loc>\n", loc)

	if lastMod != "" {
		element += fmt.Sprintf("    <lastmod>%s</lastmod>\n", lastMod)
	}

	element += fmt.Sprintf("    <changefreq>%s</changefreq>\n", changefreq)
	element += fmt.Sprintf("    <priority>%.1f</priority>\n", priority)
	element += "  </url>\n"

	return element
}
