package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// ToolType 工具类型枚举
type ToolType string

// DrawShapeType 绘制形状类型枚举
type DrawShapeType string

// 工具类型常量
const (
	ToolTypePen    ToolType = "PEN"
	ToolTypeEraser ToolType = "ERASER"
	// 可以根据实际需要添加更多类型
)

// 绘制形状类型常量
const (
	DrawShapeTypeFreehand DrawShapeType = "FREEHAND"
	DrawShapeTypeLine     DrawShapeType = "LINE"
	// 可以根据实际需要添加更多类型
)

// NoteDrawEntity 笔记绘制实体
type NoteDrawEntity struct {
	model.BaseModel               // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	LineHexColor    string        `json:"lineHexColor" gorm:"column:line_hex_color"`   // 线条十六进制颜色
	LineAlpha       float32       `json:"lineAlpha" gorm:"column:line_alpha"`          // 线条透明度
	Points          string        `json:"points" gorm:"column:points"`                 // 点集合
	PathTag         int           `json:"pathTag" gorm:"column:path_tag"`              // iPad专用路径标签
	LineWidth       float32       `json:"lineWidth" gorm:"column:line_width"`          // 线条宽度
	ToolType        ToolType      `json:"toolType" gorm:"column:tool_type"`            // 工具类型
	DrawShapeType   DrawShapeType `json:"drawShapeType" gorm:"column:draw_shape_type"` // 绘制形状类型
	NoteId          string        `json:"noteId" gorm:"column:note_id;index"`          // 笔记ID
	PageNumber      int           `json:"pageNumber" gorm:"column:page_number"`        // 页码
}

// TableName 返回表名
func (NoteDrawEntity) TableName() string {
	return "t_note_draw_entity"
}
