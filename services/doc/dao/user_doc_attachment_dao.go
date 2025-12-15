package dao

import (
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocAttachmentDAO GORM实现的用户文档附件DAO
type UserDocAttachmentDAO struct {
	*baseDao.GormBaseDAO[model.UserDocAttachment]
	logger logging.Logger
}

// NewUserDocAttachmentDAO 创建一个新的用户文档附件DAO
func NewUserDocAttachmentDAO(db *gorm.DB, logger logging.Logger) *UserDocAttachmentDAO {
	return &UserDocAttachmentDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserDocAttachment](db, logger),
		logger:      logger,
	}
}
