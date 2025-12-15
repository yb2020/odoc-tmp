package transport

import (
	"encoding/json"
	"fmt"
)

// UnmarshalEvent 解析事件数据
// 将JSON格式的事件数据解析为指定的事件结构体
func UnmarshalEvent(data []byte, event interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("empty event data")
	}

	err := json.Unmarshal(data, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return nil
}
