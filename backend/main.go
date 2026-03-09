package main

import (
	"blog-backend/internal/app/config"
	"blog-backend/internal/app/model"
	"blog-backend/internal/app/router"
	"blog-backend/internal/app/service"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 自动迁移数据库表，按正确顺序迁移以避免外键约束错误
	// 1. 先迁移基础表（无外键依赖的表）
	err := config.DB.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{})
	if err != nil {
		log.Fatal("基础表迁移失败:", err)
	}

	// 2. 再迁移关联表
	err = config.DB.AutoMigrate(&model.ArticleCategory{}, &model.ArticleTag{})
	if err != nil {
		log.Fatal("关联表迁移失败:", err)
	}

	// 3. 最后迁移文章表（有外键依赖的表）
	err = config.DB.AutoMigrate(&model.Article{})
	if err != nil {
		log.Fatal("文章表迁移失败:", err)
	}

	// 设置路由
	r := router.SetupRouter(cfg)
	// 创建百度推送服务
	baiduPushService := service.NewBaiduPushService(cfg)

	// 创建站点地图服务
	sitemapService := service.NewSitemapService(cfg.BaseURL, baiduPushService)
	sitemapService.SetBaiduPushService(baiduPushService)

	// 设置定时任务
	setupCronJobs(sitemapService)
	// 启动服务器
	log.Printf("服务器启动在 http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
func setupCronJobs(sitemapService *service.SitemapService) {
	c := cron.New()

	// 每天凌晨2点重新生成站点地图
	_, err := c.AddFunc("*/1 * * * *", func() {
		log.Println("定时任务: 开始重新生成站点地图")
		sitemapService.RegenerateSitemap()
		log.Println("定时任务: 站点地图重新生成完成")
	})

	if err != nil {
		log.Printf("添加站点地图定时任务失败: %v", err)
		return
	}

	c.Start()
	log.Println("定时任务已启动")
}
