package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	Username  string         `json:"username" gorm:"type:varchar(50);not null;unique"`
	Email     string         `json:"email" gorm:"type:varchar(100);not null;unique"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"` // 不在JSON中返回密码
	Nickname  string         `json:"nickname" gorm:"type:varchar(50)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Bio       string         `json:"bio" gorm:"type:text"`
	Role      string         `json:"role" gorm:"type:varchar(20);default:'user'"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// SetPassword 设置密码（加密存储）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	// 首先尝试明文密码验证（仅用于开发环境）
	if u.Password == password {
		return true
	}

	// 然后尝试bcrypt验证（用于生产环境）
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
