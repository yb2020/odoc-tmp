package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// // ShapeType 形状类型枚举
// type ShapeType string

// // AnnotationColor 注释颜色枚举
// type AnnotationColor string

// // 形状类型常量
// const (
// 	ShapeTypeRectangle ShapeType = "RECTANGLE"
// 	ShapeTypeEllipse   ShapeType = "ELLIPSE"
// 	ShapeTypeLine      ShapeType = "LINE"
// 	// 可以根据实际需要添加更多类型
// )

// // 注释颜色常量
// const (
// 	AnnotationColorRed    AnnotationColor = "RED"
// 	AnnotationColorGreen  AnnotationColor = "GREEN"
// 	AnnotationColorBlue   AnnotationColor = "BLUE"
// 	AnnotationColorYellow AnnotationColor = "YELLOW"
// 	// 可以根据实际需要添加更多颜色
// )

// NoteShape 笔记形状实体
type NoteShape struct {
	model.BaseModel         // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UUID            string  `json:"uuid" gorm:"column:uuid; type:varchar(120)"`                // 唯一标识符
	Type            string  `json:"type" gorm:"column:type; type:varchar(120)"`                // 形状类型
	X               float64 `json:"x" gorm:"column:x"`                                         // X坐标
	Y               float64 `json:"y" gorm:"column:y"`                                         // Y坐标
	StrokeColor     string  `json:"strokeColor" gorm:"column:stroke_color; type:varchar(120)"` // 描边颜色
	Width           float64 `json:"width" gorm:"column:width"`                                 // 宽度
	Height          float64 `json:"height" gorm:"column:height"`                               // 高度
	RadiusX         float64 `json:"radiusX" gorm:"column:radius_x"`                            // X半径
	RadiusY         float64 `json:"radiusY" gorm:"column:radius_y"`                            // Y半径
	EndX            float64 `json:"endX" gorm:"column:end_x"`                                  // 结束X坐标
	EndY            float64 `json:"endY" gorm:"column:end_y"`                                  // 结束Y坐标
	PageNumber      int     `json:"pageNumber" gorm:"column:page_number"`                      // 页码
	NoteId          string  `json:"noteId" gorm:"column:note_id;index"`                        // 笔记ID
}

// TableName 返回表名
func (NoteShape) TableName() string {
	return "t_note_shape"
}
