// backend/internal/app/service/crawler_service.go
package service

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type CrawlerStats struct {
	UserAgent string
	IP        string
	Path      string
	Timestamp time.Time
}

func RecordCrawlerVisit(r *http.Request) {
	// 获取用户代理
	userAgent := r.UserAgent()

	// 检查是否为常见爬虫
	crawlers := []string{"Googlebot", "Bingbot", "BaiduSpider", "Sogou Spider", "Yahoo Slurp"}
	isCrawler := false
	for _, crawler := range crawlers {
		if strings.Contains(userAgent, crawler) {
			isCrawler = true
			break
		}
	}

	if !isCrawler {
		return
	}

	// 记录爬虫访问
	stats := CrawlerStats{
		UserAgent: userAgent,
		IP:        r.RemoteAddr,
		Path:      r.URL.Path,
		Timestamp: time.Now(),
	}

	// 这里可以将统计数据存储到数据库或日志文件
	//config.DB.Create(&model.CrawlerLog{
	//	UserAgent: stats.UserAgent,
	//	IP:        stats.IP,
	//	Path:      stats.Path,
	//	CreatedAt: stats.Timestamp,
	//})
	log.Default().Print(stats)
}
