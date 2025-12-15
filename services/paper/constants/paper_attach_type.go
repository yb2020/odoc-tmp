package constants

// PaperAttachType 论文附件类型枚举
type PaperAttachType int

const (
	// Unknown 未知类型
	Unknown PaperAttachType = iota
	// PDF PDF文件
	PDF
	// Image 图片文件
	Image
	// Document 文档文件
	Document
	// Other 其他类型
	Other
)

// String 返回枚举值的字符串表示
func (t PaperAttachType) String() string {
	switch t {
	case PDF:
		return "PDF"
	case Image:
		return "IMAGE"
	case Document:
		return "DOCUMENT"
	case Other:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

// Value 返回枚举值的整数值
func (t PaperAttachType) Value() int {
	return int(t)
}

// PaperAttachTypeFromValue 根据整数值返回对应的PaperAttachType枚举值
func PaperAttachTypeFromValue(value int) PaperAttachType {
	switch value {
	case 1:
		return PDF
	case 2:
		return Image
	case 3:
		return Document
	case 4:
		return Other
	default:
		return Unknown
	}
}

// PaperAttachTypeFromString 根据字符串值返回对应的PaperAttachType枚举值
func PaperAttachTypeFromString(value string) PaperAttachType {
	switch value {
	case "PDF":
		return PDF
	case "IMAGE":
		return Image
	case "DOCUMENT":
		return Document
	case "OTHER":
		return Other
	default:
		return Unknown
	}
}
