package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperJcrEntity 论文JCR实体
type PaperJcrEntity struct {
	model.BaseModel          // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	ImpactOfFactor float32   `json:"impactOfFactor" gorm:"column:impact_of_factor"` // 影响因子
	JcrPartion     string    `json:"jcrPartion" gorm:"column:jcr_partion"`          // JCR分区
	Source         string    `json:"source" gorm:"column:source"`                   // 来源
	Venue          string    `json:"venue" gorm:"column:venue;index"`               // 发表场所
}

// TableName 返回表名
func (PaperJcrEntity) TableName() string {
	return "t_paper_jcr_entity"
}

// NewPaperJcrEntityWithVenue 根据发表场所创建PaperJcrEntity实例
func NewPaperJcrEntityWithVenue(venue string) *PaperJcrEntity {
	return &PaperJcrEntity{
		Venue: venue,
	}
}
