package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// 学术网站
type Website struct {
	model.BaseModel        // Embedded BaseModel
	UserId          string `json:"user_id" gorm:"column:user_id;not null"`             // 用户ID
	Source          int32  `json:"source" gorm:"column:source;default:2;not null"`     // 来源 1:系统 2:用户
	IconUrl         string `json:"icon_url" gorm:"column:icon_url;type:text"`          // 网站图标URL
	Name            string `json:"name" gorm:"column:name;type:varchar(255);not null"` // 网站名称
	Url             string `json:"url" gorm:"column:url;type:text;not null"`           // 网站URL
	OpenType        int32  `json:"open_type" gorm:"column:open_type;default:1"`        // 打开方式 1: 新窗口打开, 2: 当前窗口打开
	SortOrder       int32  `json:"sort_order" gorm:"column:sort_order;default:0"`      // 排序
}

// TableName returns the database table name for the Website model.
func (Website) TableName() string {
	return "t_nav_website"
}
