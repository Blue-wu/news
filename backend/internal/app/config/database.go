package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *Config) error {
	var err error

	// 配置GORM日志
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 获取MySQL连接字符串
	dsn := cfg.GetDSN()
	log.Printf("正在连接MySQL数据库: %s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取底层SQL数据库连接以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DatabaseMaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DatabaseMaxLifetime) * time.Second)

	log.Println("MySQL数据库连接成功")
	return nil
}
