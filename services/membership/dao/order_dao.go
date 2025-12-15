package dao

import (
	"context"
	"database/sql"
	"time"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/model"
	"gorm.io/gorm"
)

// OrderDAO GORM实现的会员订单DAO
type OrderDAO struct {
	*baseDao.GormBaseDAO[model.Order]
	logger logging.Logger
}

// NewOrderDAO 创建一个新的会员订单DAO
func NewOrderDAO(db *gorm.DB, logger logging.Logger) *OrderDAO {
	return &OrderDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.Order](db, logger),
		logger:      logger,
	}
}

// GetTotalNum 获取订单numberCount字段的统计数量
func (d *OrderDAO) GetTotalNum(ctx context.Context, userId string, orderStatus int32, orderType int32, startAt time.Time, endAt time.Time) (int32, error) {
	var totalNum *int32
	err := d.GetDB(ctx).Model(&model.Order{}).Where("user_id = ? and is_deleted = false and order_type = ? and order_status = ? and created_at >= ? and created_at <= ?", userId, orderType, orderStatus, startAt, endAt).Select("sum(number_count)").Row().Scan(&totalNum)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	if totalNum == nil {
		return 0, nil
	}
	return *totalNum, nil
}
