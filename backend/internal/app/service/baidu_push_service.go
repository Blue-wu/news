package service

import (
	"blog-backend/internal/app/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// BaiduPushResponse 百度链接推送响应结构体
type BaiduPushResponse struct {
	Success     int      `json:"success"`       // 成功推送的URL条数
	Remain      int      `json:"remain"`        // 当天剩余的可推送URL条数
	NotSameSite []string `json:"not_same_site"` // 由于不是本站URL而未处理的URL列表
	NotValid    []string `json:"not_valid"`     // 不合法的URL列表
}

// BaiduPushService 百度链接推送服务
type BaiduPushService struct {
	baseURL    string
	pushToken  string
	pushSite   string
	httpClient *http.Client
	retryCount int
}

// NewBaiduPushService 创建百度链接推送服务实例
func NewBaiduPushService(cfg *config.Config) *BaiduPushService {
	return &BaiduPushService{
		baseURL:   cfg.BaseURL,
		pushToken: cfg.BaiduPushToken,
		pushSite:  cfg.BaiduPushSite,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // 设置10秒超时
		},
		retryCount: 3, // 默认重试3次
	}
}

// PushURL 推送单个URL到百度
func (s *BaiduPushService) PushURL(url string) (*BaiduPushResponse, error) {
	return s.PushURLs([]string{url})
}

// PushURLs 推送多个URL到百度
func (s *BaiduPushService) PushURLs(urls []string) (*BaiduPushResponse, error) {
	// 检查配置是否完整
	if s.pushToken == "" || s.pushSite == "" {
		log.Printf("百度推送配置不完整，请检查环境变量 BAIDU_PUSH_TOKEN 和 BAIDU_PUSH_SITE")
		return nil, fmt.Errorf("百度推送配置不完整")
	}

	// 检查URL列表是否为空
	if len(urls) == 0 {
		return nil, fmt.Errorf("推送URL列表不能为空")
	}

	// 构建推送接口URL
	pushURL := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", s.pushSite, s.pushToken)

	// 构建请求体（每个URL一行）
	var requestBody bytes.Buffer
	for _, url := range urls {
		requestBody.WriteString(url + "\n")
	}

	// 发送请求并支持重试
	var resp *http.Response
	var err error
	for i := 0; i < s.retryCount; i++ {
		// 创建请求
		req, err := http.NewRequest("POST", pushURL, &requestBody)
		if err != nil {
			log.Fatalf("创建百度推送请求失败 (尝试 %d/%d): %v", i+1, s.retryCount, err)
			if i == s.retryCount-1 {
				return nil, fmt.Errorf("创建百度推送请求失败: %v", err)
			}
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}

		// 设置请求头
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("User-Agent", "blog-backend/1.0")

		// 发送请求
		resp, err = s.httpClient.Do(req)
		if err != nil {
			log.Printf("发送百度推送请求失败 (尝试 %d/%d): %v", i+1, s.retryCount, err)
			if i == s.retryCount-1 {
				return nil, fmt.Errorf("发送百度推送请求失败: %v", err)
			}
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}

		// 检查响应状态码
		if resp.StatusCode == http.StatusOK {
			break
		}

		resp.Body.Close()
		log.Printf("百度推送请求返回错误状态码 (尝试 %d/%d): %d", i+1, s.retryCount, resp.StatusCode)
		if i == s.retryCount-1 {
			return nil, fmt.Errorf("百度推送请求返回错误状态码: %d", resp.StatusCode)
		}
		time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
	}

	if resp == nil {
		return nil, fmt.Errorf("百度推送请求失败")
	}

	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取百度推送响应失败: %v", err)
		return nil, fmt.Errorf("读取百度推送响应失败: %v", err)
	}

	// 解析响应
	var response BaiduPushResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("解析百度推送响应失败，响应内容: %s, 错误: %v", string(body), err)
		return nil, fmt.Errorf("解析百度推送响应失败: %v", err)
	}

	log.Printf("百度推送结果: 成功推送 %d 条, 剩余 %d 条, 总推送 %d 条", response.Success, response.Remain, len(urls))
	if len(response.NotSameSite) > 0 {
		log.Printf("百度推送: 非本站URL %v", response.NotSameSite)
	}
	if len(response.NotValid) > 0 {
		log.Printf("百度推送: 不合法URL %v", response.NotValid)
	}
	if response.Success < len(urls) {
		log.Printf("部分URL推送失败，成功 %d/%d 条", response.Success, len(urls))
	}

	return &response, nil
}

// PushArticleToBaidu 推送文章URL到百度
func (s *BaiduPushService) PushArticleToBaidu(articleID uint) error {
	// 构建文章URL
	articleURL := fmt.Sprintf("%s/articles/%d", s.baseURL, articleID)

	log.Printf("准备推送文章URL到百度: %s", articleURL)

	// 推送URL
	resp, err := s.PushURL(articleURL)
	if err != nil {
		log.Printf("推送文章 %d 到百度失败: %v", articleID, err)
		return fmt.Errorf("推送文章到百度失败: %v", err)
	}

	log.Printf("文章 %d 推送成功，剩余推送次数: %d", articleID, resp.Remain)

	return nil
}

// SetRetryCount 设置重试次数
func (s *BaiduPushService) SetRetryCount(count int) {
	if count > 0 {
		s.retryCount = count
	}
}

// SetTimeout 设置超时时间
func (s *BaiduPushService) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		s.httpClient.Timeout = timeout
	}
}
