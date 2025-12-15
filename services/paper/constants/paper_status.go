package constants

// PaperStatus 论文状态枚举
type PaperStatus string

const (
	// Private 私有
	Private PaperStatus = "private"
	// Public 公开
	Public PaperStatus = "public"
)

// String 返回枚举值的字符串表示
func (s PaperStatus) String() string {
	return string(s)
}

// IsValid 检查状态是否有效
func (s PaperStatus) IsValid() bool {
	switch s {
	case Private, Public:
		return true
	default:
		return false
	}
}
