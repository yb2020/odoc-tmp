package event

import "github.com/yb2020/odoc/pkg/eventbus"

// 用户模块下的事件类型
const (
	UserRegisterEvent eventbus.EventType = "user.register"
	UserDeletedEvent  eventbus.EventType = "user.deleted"

	// 其他事件类型...
)
