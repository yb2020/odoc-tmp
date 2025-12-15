package constants

// ExtParmKeys 扩展参数键枚举
type ExtParmKeys string

const (
	// SendMessage 是否通知过群组了
	SendMessage ExtParmKeys = "sendMessage"
)

// String 返回枚举值的字符串表示
func (k ExtParmKeys) String() string {
	return string(k)
}
