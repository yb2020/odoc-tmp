package model

import (
	"time"

	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/idgen"
	"gorm.io/gorm"
)

// BaseModel 包含所有模型共享的基础字段和方法
// 注意：时间字段使用 TEXT 类型存储 RFC3339 格式字符串，便于 SQLite 和 PostgreSQL 兼容
type BaseModel struct {
	// 主键 - 在数据库中使用text类型存储雪花Id
	Id string `json:"id" gorm:"primaryKey;size:36"`

	// 软删除标记
	IsDeleted bool `json:"is_deleted" gorm:"default:false"`

	// 创建信息
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	Creator   string    `json:"creator" gorm:"type:varchar(200)"`
	CreatorId string    `json:"creator_id" gorm:"index;size:36"`

	// 修改信息
	UpdatedAt  time.Time `json:"updated_at" gorm:"index"`
	Modifier   string    `json:"modifier" gorm:"type:varchar(200)"`
	ModifierId string    `json:"modifier_id" gorm:"index;size:36"`
}

// BeforeCreate GORM钩子，在创建记录前自动生成ID和设置时间戳
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果ID已经设置，则不生成新ID
	if b.Id == "" {
		// 生成 UUID v7 格式的ID
		b.Id = idgen.GenerateUUID()
	}

	// 设置ID和初始值
	b.IsDeleted = false

	// 从上下文中获取用户信息
	if uc := userContext.GetUserContext(tx.Statement.Context); uc != nil {
		if uc.UserId != "" {
			b.CreatorId = uc.UserId
			b.ModifierId = uc.UserId
		}
		if uc.Username != "" {
			b.Creator = uc.Username
			b.Modifier = uc.Username
		} else if uc.ServiceName != "" {
			b.Creator = uc.ServiceName
			b.Modifier = uc.ServiceName
		}
	}
	// 统一使用 UTC 时间，便于云端同步和 SQLite 字符串排序
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now

	return nil
}

// BeforeUpdate GORM钩子，在更新记录前自动更新时间戳
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	// 更新修改时间（统一使用 UTC）
	b.UpdatedAt = time.Now().UTC()

	// 从上下文中获取用户信息
	if uc := userContext.GetUserContext(tx.Statement.Context); uc != nil {
		if uc.UserId != "" {
			b.ModifierId = uc.UserId
		}
		if uc.Username != "" {
			b.Modifier = uc.Username
		} else if uc.ServiceName != "" {
			b.Modifier = uc.ServiceName
		}
	}
	return nil
}
