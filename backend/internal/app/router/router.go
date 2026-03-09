package router

import (
	"blog-backend/internal/app/api"
	"blog-backend/internal/app/config"
	"blog-backend/internal/app/middleware"
	"blog-backend/internal/app/repository"
	"blog-backend/internal/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		// 设置允许的源，可以根据需要设置为特定域名
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 设置是否允许发送Cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 初始化依赖
	articleRepo := repository.NewArticleRepository(config.DB)
	articleService := service.NewArticleService(articleRepo)
	// 添加用户相关依赖
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo, config.DB)
	userHandler := api.NewUserHandler(userService)
	// 初始化百度推送服务并将其传递给ArticleHandler
	baiduPushService := service.NewBaiduPushService(cfg)
	articleHandler := api.NewArticleHandler(articleService, userService, baiduPushService)
	// 添加分类相关依赖
	categoryRepo := repository.NewCategoryRepository(config.DB)
	categoryService := service.NewCategoryService(categoryRepo, config.DB)
	categoryHandler := api.NewCategoryHandler(categoryService)

	// 初始化站点地图服务
	sitemapService := service.NewSitemapService(cfg.BaseURL, baiduPushService)
	sitemapHandler := api.NewSitemapHandler(sitemapService, baiduPushService)
	// API路由组
	apiGroup := r.Group("/api")
	{
		// 健康检查
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "博客后端服务运行正常",
			})
		})

		// 认证相关路由（无需认证）
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.Register)
		}

		// ... existing code ...
		// 公开路由（无需认证）
		// 分类相关公开路由
		apiGroup.GET("/categories", categoryHandler.GetAllCategories)
		apiGroup.GET("/categories/root", categoryHandler.GetRootCategories) // 这个接口对前端很有用，可以直接获取顶级分类
		apiGroup.GET("/categories/tree", categoryHandler.GetCategoriesWithChildren)
		apiGroup.GET("/categories/slug/:slug", categoryHandler.GetCategoryBySlug)
		// 在router.go文件中添加
		apiGroup.POST("/articles/:id/view", articleHandler.IncrementViewCount)
		// 在公开API路由组中添加搜索路由
		apiGroup.GET("/articles/search", articleHandler.SearchArticles)
		apiGroup.GET("/articles/random", articleHandler.GetRandomArticles)
		apiGroup.GET("/articles/popular", articleHandler.GetPopularArticles)
		// 添加按分类获取文章的路由
		apiGroup.GET("/categories/:slug/articles", articleHandler.GetArticlesByCategory)
		// 公开的文章API - 移出需要认证的路由组
		apiGroup.GET("/articles", articleHandler.GetAllArticles) // 获取所有文章（公开）
		apiGroup.GET("/articles/:id", articleHandler.GetArticle) // 获取单个文章（公开）
		// 需要认证的路由
		protected := apiGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			protected.GET("/user/profile", userHandler.GetProfile)
			// 图片上传路由（需要认证）
			protected.POST("/upload/image", articleHandler.UploadImage)
			// 文章相关路由（需要认证）
			articles := protected.Group("/articles")
			{
				articles.POST("", articleHandler.CreateArticle)       // 创建文章
				articles.PUT("/:id", articleHandler.UpdateArticle)    // 更新文章
				articles.DELETE("/:id", articleHandler.DeleteArticle) // 删除文章
			}

			// 管理员路由（需要管理员权限）
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// 用户管理
				admin.GET("/users", userHandler.GetAllUsers)                 // 获取所有用户列表
				admin.PUT("/users/:id/status", userHandler.UpdateUserStatus) // 更新用户状态
				admin.DELETE("/users/:id", userHandler.DeleteUser)           // 删除用户
				// 分类管理
				admin.POST("/categories", categoryHandler.CreateCategory)       // 创建分类
				admin.GET("/categories", categoryHandler.GetAllCategories)      // 获取所有分类
				admin.GET("/categories/:id", categoryHandler.GetCategory)       // 获取单个分类
				admin.PUT("/categories/:id", categoryHandler.UpdateCategory)    // 更新分类
				admin.DELETE("/categories/:id", categoryHandler.DeleteCategory) // 删除分类
				// 文章管理（管理员功能）
				// 文章管理（管理员功能）
				admin.GET("/articles/draft", articleHandler.GetDraftArticles)         // 获取草稿文章
				admin.GET("/articles/:id", articleHandler.GetArticle)                 // 添加这一行，使用现有的GetArticle处理函数
				admin.PUT("/articles/:id/status", articleHandler.UpdateArticleStatus) // 更新文章状态
				// 站点地图管理路由
				admin.GET("/settings/sitemap", sitemapHandler.GetSettings)
				admin.PUT("/settings/sitemap", sitemapHandler.UpdateSettings)
				admin.POST("/sitemap/generate", sitemapHandler.RegenerateSitemap)
				admin.POST("/sitemap/clear-cache", sitemapHandler.ClearCache)
				admin.GET("/sitemap/stats", sitemapHandler.GetStats)
				// 系统设置路由（临时实现在router中）
				admin.GET("/settings", func(c *gin.Context) {
					// 默认系统设置
					defaultSettings := map[string]interface{}{
						"siteName":         "我的博客",
						"siteDescription":  "这是一个博客网站",
						"siteKeywords":     "博客,技术,分享",
						"footerText":       "© 2024 我的博客. 保留所有权利.",
						"articlesPerPage":  10,
						"allowComments":    true,
						"autoPublish":      false,
						"contactEmail":     "",
						"socialLinks":      []interface{}{},
						"sitemap_settings": "{\"cacheEnabled\":true,\"cacheDuration\":86400}",
					}
					c.JSON(http.StatusOK, defaultSettings)
				})
				admin.PUT("/settings", func(c *gin.Context) {
					var req map[string]interface{}
					if err := c.ShouldBindJSON(&req); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"success": false,
							"error":   "请求参数无效",
						})
						return
					}
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"message": "系统设置已成功更新",
					})
				})
			}
		}

	}
	// 在返回router之前添加静态文件服务配置
	// 提供静态文件访问
	r.Static("/uploads", "./uploads")

	// 添加站点地图路由 - 公开访问，无需认证
	// 添加站点地图路由 - 公开访问，无需认证
	r.GET("/sitemap.xml", sitemapHandler.GetSitemap)

	return r
}
