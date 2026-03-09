package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
type Config struct {
	// 服务器配置
	Port    string
	BaseURL string

	// 数据库配置
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBCharset   string
	DBParseTime string
	DBLoc       string

	// 数据库连接池配置
	DatabaseMaxIdle     int
	DatabaseMaxOpen     int
	DatabaseMaxLifetime int

	// 日志配置
	LogLevel string

	// 百度链接推送配置
	BaiduPushToken string
	BaiduPushSite  string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用默认配置")
	}

	config := &Config{
		// 服务器配置
		Port:    getEnv("PORT", "8000"),
		BaseURL: getEnv("BASE_URL", "http://localhost:8000"),

		// 数据库配置
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "root123456"),
		DBName:      getEnv("DB_NAME", "blog_db"),
		DBCharset:   getEnv("DB_CHARSET", "utf8mb4"),
		DBParseTime: getEnv("DB_PARSE_TIME", "True"),
		DBLoc:       getEnv("DB_LOC", "Local"),

		// 数据库连接池配置
		DatabaseMaxIdle:     getEnvAsInt("DATABASE_MAX_IDLE", 10),
		DatabaseMaxOpen:     getEnvAsInt("DATABASE_MAX_OPEN", 100),
		DatabaseMaxLifetime: getEnvAsInt("DATABASE_MAX_LIFETIME", 3600), // 秒

		// 日志配置
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// 百度链接推送配置
		BaiduPushToken: getEnv("BAIDU_PUSH_TOKEN", ""),
		BaiduPushSite:  getEnv("BAIDU_PUSH_SITE", ""),
	}

	return config
}

// GetDSN 获取MySQL连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBCharset, c.DBParseTime, c.DBLoc)
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("环境变量 %s 转换为整数失败: %v, 使用默认值 %d", key, err, defaultValue)
		return defaultValue
	}

	return value
}
