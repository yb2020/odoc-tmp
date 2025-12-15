package model

import (
	"strings"
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// OAuth2Clients 表示OAuth2客户端
type OAuth2Clients struct {
	ID            string            `gorm:"column:id;primaryKey"`
	Secret        string            `gorm:"column:secret"`
	Name          string            `gorm:"column:name"`
	RedirectUris  model.StringSlice `gorm:"column:redirect_uris;type:json"`
	AllowedScopes model.StringSlice `gorm:"column:allowed_scopes;type:json"`
	Active        bool              `gorm:"column:active;default:true"`
	CreatedAt     time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time         `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 返回表名
func (OAuth2Clients) TableName() string {
	return "t_oauth2_clients"
}

// IsValidScope 检查作用域是否有效
func (c *OAuth2Clients) IsValidScope(scope string) bool {
	// 如果客户端没有限制作用域，则所有作用域都有效
	if len(c.AllowedScopes) == 0 {
		return true
	}

	// 检查作用域是否在允许列表中
	for _, allowedScope := range c.AllowedScopes {
		if allowedScope == scope {
			return true
		}
	}

	return false
}

// IsValidRedirectURI 检查重定向URI是否有效
func (c *OAuth2Clients) IsValidRedirectURI(uri string) bool {
	for _, redirectURI := range c.RedirectUris {
		if redirectURI == uri {
			return true
		}
	}
	return false
}

// GetRedirectURIsAsString 将重定向URI数组转换为字符串
func (c *OAuth2Clients) GetRedirectURIsAsString() string {
	return strings.Join(c.RedirectUris, ",")
}

// SetRedirectURIsFromString 从字符串设置重定向URI数组
func (c *OAuth2Clients) SetRedirectURIsFromString(uris string) {
	if uris == "" {
		c.RedirectUris = []string{}
		return
	}
	c.RedirectUris = strings.Split(uris, ",")
}

// GetAllowedScopesAsString 将允许的作用域数组转换为字符串
func (c *OAuth2Clients) GetAllowedScopesAsString() string {
	return strings.Join(c.AllowedScopes, " ")
}

// SetAllowedScopesFromString 从字符串设置允许的作用域数组
func (c *OAuth2Clients) SetAllowedScopesFromString(scopes string) {
	if scopes == "" {
		c.AllowedScopes = []string{}
		return
	}
	c.AllowedScopes = strings.Split(scopes, " ")
}
