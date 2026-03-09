package api

import (
	"blog-backend/internal/app/model"
	"blog-backend/internal/app/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service          *service.ArticleService
	userService      *service.UserService
	baiduPushService *service.BaiduPushService
}

func NewArticleHandler(service *service.ArticleService, userService *service.UserService, baiduPushService *service.BaiduPushService) *ArticleHandler {
	return &ArticleHandler{
		service:          service,
		userService:      userService,
		baiduPushService: baiduPushService,
	}
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateArticle(&article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	// 如果文章状态是已发布，则推送链接到百度
	if article.Status == 1 {
		go func() {
			if h.baiduPushService != nil {
				if err := h.baiduPushService.PushArticleToBaidu(uint(article.ID)); err != nil {
					// 非阻塞推送，只记录错误日志
					fmt.Printf("推送文章 %d 到百度失败: %v\n", article.ID, err)
				}
			}
		}()
	}

	c.JSON(http.StatusCreated, article)
}

// GetArticle 获取单个文章
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	article, err := h.service.GetArticleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// GetAllArticles 获取所有文章（支持分页）
func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	// 获取分页参数，默认为第1页，每页10条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.service.GetAllArticlesPaged(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// 更新文章请求结构体
type UpdateArticleRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary"`
	CoverImage  string   `json:"cover_image"`
	Status      int      `json:"status"`
	IsTop       int      `json:"is_top"`
	CategoryIds []uint64 `json:"category_ids"`
	TagIds      []uint64 `json:"tag_ids"`
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 使用自定义结构体接收请求
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 先检查文章是否存在
	existingArticle, err := h.service.GetArticleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 更新基本字段
	existingArticle.Title = req.Title
	existingArticle.Content = req.Content
	existingArticle.Summary = req.Summary
	existingArticle.CoverImage = req.CoverImage
	existingArticle.Status = req.Status
	existingArticle.IsTop = req.IsTop

	// 处理分类关联
	if len(req.CategoryIds) > 0 {
		// 清除现有分类关联
		h.service.ClearArticleCategories(uint(id))
		// 添加新的分类关联
		for _, categoryID := range req.CategoryIds {
			h.service.AddArticleCategory(uint(id), uint(categoryID))
		}
	}

	// 保存文章基本信息
	if err := h.service.UpdateArticle(existingArticle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, existingArticle)
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	err = h.service.DeleteArticle(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

// GetDraftArticles 获取草稿文章（管理员权限）
func (h *ArticleHandler) GetDraftArticles(c *gin.Context) {
	articles, err := h.service.GetDraftArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取草稿文章失败"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// UpdateArticleStatus 更新文章状态（管理员权限）
func (h *ArticleHandler) UpdateArticleStatus(c *gin.Context) {
	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := h.service.UpdateArticleStatus(uint(articleID), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章状态失败"})
		return
	}

	// 如果文章状态更新为已发布，推送链接到百度
	if req.Status == 1 {
		go func() {
			// 获取完整的文章信息
			article, err := h.service.GetArticleByID(uint(articleID))
			if err != nil {
				log.Printf("获取文章信息失败: %v\n", err)
				return
			}

			err = h.baiduPushService.PushArticleToBaidu(uint(article.ID))
			if err != nil {
				log.Printf("百度链接推送失败: %v\n", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章状态更新成功"})
}

// UploadImage 处理文章内容中的图片上传
func (h *ArticleHandler) UploadImage(c *gin.Context) {
	// 确保目录存在
	uploadDir := "./uploads/images"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供图片文件"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		return
	}

	// 构建文件URL
	url := fmt.Sprintf("/uploads/images/%s", filename)

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// 修改 GetArticlesByCategory 方法支持分页
func (h *ArticleHandler) GetArticlesByCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类参数不能为空"})
		return
	}

	// 获取分页参数，默认为第1页，每页10条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.service.GetArticlesByCategorySlugPaged(slug, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类文章失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// IncrementViewCount 增加文章阅读量
func (h *ArticleHandler) IncrementViewCount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	err = h.service.IncrementViewCount(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新阅读量失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "阅读量更新成功"})
}

// 添加搜索文章的处理函数
func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.service.SearchArticlesPaged(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索文章失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRandomArticles 获取随机文章
func (h *ArticleHandler) GetRandomArticles(c *gin.Context) {
	// 获取limit参数，默认为6
	limitStr := c.DefaultQuery("limit", "6")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 20 {
		limit = 6 // 如果参数无效，使用默认值
	}

	articles, err := h.service.GetRandomArticles(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取随机文章失败"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetPopularArticles 获取热门文章
func (h *ArticleHandler) GetPopularArticles(c *gin.Context) {
	// 获取limit参数，默认为8
	limitStr := c.DefaultQuery("limit", "8")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 20 {
		limit = 8 // 如果参数无效，使用默认值
	}

	articles, err := h.service.GetPopularArticles(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门文章失败"})
		return
	}

	c.JSON(http.StatusOK, articles)
}
