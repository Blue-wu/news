package api

import (
	"blog-backend/internal/app/config"
	"blog-backend/internal/app/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRSSFeed(c *gin.Context) {
	var articles []model.Article
	if err := config.DB.Where("status = ?", 1).
		Order("created_at DESC").
		Limit(20).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	// 生成RSS XML
	var rssContent string
	rssContent += "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	rssContent += "<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n"
	rssContent += "<channel>\n"
	rssContent += fmt.Sprintf("<title>博客系统</title>\n")
	rssContent += fmt.Sprintf("<link>%s</link>\n", config.LoadConfig().BaseURL)
	rssContent += fmt.Sprintf("<description>博客系统RSS订阅</description>\n")
	rssContent += fmt.Sprintf("<atom:link href=\"%s/rss\" rel=\"self\" type=\"application/rss+xml\" />\n", config.LoadConfig().BaseURL)
	rssContent += fmt.Sprintf("<lastBuildDate>%s</lastBuildDate>\n", time.Now().Format(time.RFC1123))
	rssContent += "<language>zh-CN</language>\n"

	for _, article := range articles {
		rssContent += "<item>\n"
		rssContent += fmt.Sprintf("<title>%s</title>\n", article.Title)
		rssContent += fmt.Sprintf("<link>%s/articles/%d</link>\n", config.LoadConfig().BaseURL, article.ID)
		rssContent += fmt.Sprintf("<pubDate>%s</pubDate>\n", article.CreatedAt.Format(time.RFC1123))
		rssContent += fmt.Sprintf("<description><![CDATA[%s]]></description>\n", article.Content)
		rssContent += fmt.Sprintf("<guid>%s/articles/%d</guid>\n", config.LoadConfig().BaseURL, article.ID)
		rssContent += "</item>\n"
	}

	rssContent += "</channel>\n"
	rssContent += "</rss>\n"

	c.Header("Content-Type", "application/rss+xml")
	c.String(http.StatusOK, rssContent)
}
