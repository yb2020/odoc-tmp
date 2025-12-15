package event

import "github.com/yb2020/odoc/pkg/eventbus"

// 文献事件
const (
	// 文献创建事件
	DocCreatedEvent eventbus.EventType = "doc.created"
	// 文献删除事件
	DocDeletedEvent eventbus.EventType = "doc.deleted"
)
