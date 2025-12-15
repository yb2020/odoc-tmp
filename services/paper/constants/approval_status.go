package constants

// ApprovalStatus 赞同状态枚举
type ApprovalStatus string

const (
	// Approve 赞同
	Approve ApprovalStatus = "APPROVE"
	// Reject 反对
	Reject ApprovalStatus = "REJECT"
	// Neutral 中立
	Neutral ApprovalStatus = "NEUTRAL"
)

// String 返回枚举值的字符串表示
func (s ApprovalStatus) String() string {
	return string(s)
}

// IsValid 检查状态是否有效
func (s ApprovalStatus) IsValid() bool {
	switch s {
	case Approve, Reject, Neutral:
		return true
	default:
		return false
	}
}
