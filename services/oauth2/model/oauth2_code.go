package model

import (
	"time"
)

// AuthCode 表示OAuth2授权码
type OAuth2AuthCode struct {
	Code        string    `gorm:"column:code;primaryKey"`
	ClientId    string    `gorm:"column:client_id;index"`
	UserId      string    `gorm:"column:user_id;index"`
	RedirectUri string    `gorm:"column:redirect_uri"`
	Scope       string    `gorm:"column:scope"`
	ExpiresAt   time.Time `gorm:"column:expires_at;index"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	Used        bool      `gorm:"column:used;default:false"`
}

// TableName 返回表名
func (OAuth2AuthCode) TableName() string {
	return "t_oauth2_auth_code"
}
