package proto

// NoteAnnotationStyleTransformer 标注样式转换器
type NoteAnnotationStyleTransformer struct {
}

// 默认样式常量
const (
	DefaultStyleID = 1
	DefaultColor   = "#338AFF"
)

// 样式ID到颜色的映射
var colorMap = map[uint32]string{
	1: "#338AFF",
	2: "#80E639",
	3: "#FFFF00",
	4: "#FF8C19",
	5: "#F24030",
}

// GetColorByStyleId 根据样式ID获取颜色
func (s *NoteAnnotationStyleTransformer) GetColorByStyleId(styleId uint32) string {
	color, exists := colorMap[styleId]
	if !exists {
		return DefaultColor
	}
	return color
}

// GetStyleIdByColor 根据颜色获取样式ID
func (s *NoteAnnotationStyleTransformer) GetStyleIdByColor(color string) uint32 {
	for styleId, mapColor := range colorMap {
		if mapColor == color {
			return styleId
		}
	}
	return DefaultStyleID
}
