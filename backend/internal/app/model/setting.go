// backend/internal/app/model/setting.go
package model

import "time"

// SiteSetting 网站设置模型
type SiteSetting struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	KeyName     string    `json:"key_name" gorm:"column:key_name;type:varchar(50);not null;unique;index"`
	Value       string    `json:"value" gorm:"type:text"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SiteSetting) TableName() string {
	return "site_settings"
}
