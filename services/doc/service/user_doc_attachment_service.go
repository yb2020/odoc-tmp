package service

import (
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/dao"
)

// UserDocAttachmentService 用户文档附件服务实现
type UserDocAttachmentService struct {
	userDocAttachmentDAO *dao.UserDocAttachmentDAO
	logger               logging.Logger
	tracer               opentracing.Tracer
}

// NewUserDocAttachmentService 创建新的用户文档附件服务
func NewUserDocAttachmentService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	userDocAttachmentDAO *dao.UserDocAttachmentDAO,
) *UserDocAttachmentService {
	return &UserDocAttachmentService{
		logger:               logger,
		tracer:               tracer,
		userDocAttachmentDAO: userDocAttachmentDAO,
	}
}

