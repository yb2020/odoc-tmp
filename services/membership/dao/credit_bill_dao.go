package dao

import (
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/model"
	"gorm.io/gorm"
)

// CreditBillDAO GORM实现的会员积分流水DAO
type CreditBillDAO struct {
	*baseDao.GormBaseDAO[model.CreditBill]
	logger logging.Logger
}

// NewCreditBillDAO 创建一个新的会员积分流水DAO
func NewCreditBillDAO(db *gorm.DB, logger logging.Logger) *CreditBillDAO {
	return &CreditBillDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.CreditBill](db, logger),
		logger:      logger,
	}
}
