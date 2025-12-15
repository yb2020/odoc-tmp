package dao

import (
	"context"

	"github.com/opentracing/opentracing-go"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserCslRelationDAO 用户引用样式关联数据访问对象
type UserCslRelationDAO struct {
	*baseDao.GormBaseDAO[model.UserCslRelation]
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewUserCslRelationDAO 创建用户引用样式关联数据访问对象
func NewUserCslRelationDAO(db *gorm.DB, logger logging.Logger, tracer opentracing.Tracer) *UserCslRelationDAO {
	return &UserCslRelationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserCslRelation](db, logger),
		logger:      logger,
		tracer:      tracer,
	}
}

// IsAddedCsl 检查用户是否已添加引用样式
func (dao *UserCslRelationDAO) IsAddedCsl(ctx context.Context, userId string) (*model.UserCslRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, dao.tracer, "UserCslRelationDAO.IsAddedCsl")
	defer span.Finish()
	var relation model.UserCslRelation
	query := dao.GetDB(ctx).Where("user_id = ?", userId)
	err := query.Limit(1).Find(&relation).Error
	if err != nil {
		dao.logger.Error("check user csl relation failed", "error", err.Error())
		return nil, err
	}
	// 如果没有找到记录，返回nil
	if relation.Id == "" {
		return nil, nil
	}
	return &relation, nil
}
