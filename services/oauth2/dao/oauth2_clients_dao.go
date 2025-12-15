package dao

import (
	"context"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oauth2/model"
	"gorm.io/gorm"
)

// OAuth2ClientsDAO PostgreSQL实现的客户端DAO
type OAuth2ClientsDAO struct {
	db     *gorm.DB
	logger logging.Logger
	config *config.Config
}

// NewOAuth2ClientsDAO 创建客户端DAO
func NewOAuth2ClientsDAO(db *gorm.DB, logger logging.Logger, config *config.Config) OAuth2ClientsDAO {
	return OAuth2ClientsDAO{
		db:     db,
		logger: logger,
		config: config,
	}
}

// SaveClient 保存客户端信息
func (dao *OAuth2ClientsDAO) SaveClient(ctx context.Context, client *model.OAuth2Clients) error {
	return dao.db.WithContext(ctx).Create(client).Error
}

// GetClientByID 根据客户端ID获取客户端信息
func (dao *OAuth2ClientsDAO) GetClientByID(ctx context.Context, clientID string) (*model.OAuth2Clients, error) {
	var client model.OAuth2Clients
	err := dao.db.WithContext(ctx).Where("id = ?", clientID).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// GetActiveClients 获取所有活跃的客户端
func (dao *OAuth2ClientsDAO) GetActiveClients(ctx context.Context) ([]*model.OAuth2Clients, error) {
	var clients []*model.OAuth2Clients
	err := dao.db.WithContext(ctx).Where("active = ?", true).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// UpdateClient 更新客户端信息
func (dao *OAuth2ClientsDAO) UpdateClient(ctx context.Context, client *model.OAuth2Clients) error {
	return dao.db.WithContext(ctx).Save(client).Error
}

// DeleteClient 删除客户端
func (dao *OAuth2ClientsDAO) DeleteClient(ctx context.Context, clientID string) error {
	return dao.db.WithContext(ctx).Where("id = ?", clientID).Delete(&model.OAuth2Clients{}).Error
}
