package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// StringSlice 是一个字符串切片类型，实现了 sql.Scanner 和 driver.Valuer 接口
// 用于在数据库和 Go 类型之间转换数组数据
type StringSlice []string

// Scan 实现 sql.Scanner 接口，用于从数据库值转换为 Go 值
func (ss *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*ss = StringSlice{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// 处理 PostgreSQL 数组格式 {value1,value2}
		if len(v) > 0 && v[0] == '{' && v[len(v)-1] == '}' {
			// PostgreSQL 数组格式
			str := string(v[1 : len(v)-1])
			if str == "" {
				*ss = StringSlice{}
				return nil
			}

			// 分割并处理引号
			parts := strings.Split(str, ",")
			result := make([]string, len(parts))
			for i, part := range parts {
				// 去除引号
				part = strings.Trim(part, "\"")
				result[i] = part
			}
			*ss = result
			return nil
		}

		// 尝试 JSON 解析（适用于 MySQL JSON 类型）
		var slice []string
		if err := json.Unmarshal(v, &slice); err == nil {
			*ss = slice
			return nil
		}

		// 尝试作为普通字符串处理（逗号分隔）
		str := string(v)
		if str == "" {
			*ss = StringSlice{}
			return nil
		}
		*ss = strings.Split(str, ",")
		return nil

	case string:
		if v == "" {
			*ss = StringSlice{}
			return nil
		}

		// 尝试 JSON 解析
		var slice []string
		if err := json.Unmarshal([]byte(v), &slice); err == nil {
			*ss = slice
			return nil
		}

		// 作为普通字符串处理（逗号分隔）
		*ss = strings.Split(v, ",")
		return nil

	default:
		return fmt.Errorf("不支持的类型: %T", value)
	}
}

// Value 实现 driver.Valuer 接口，用于从 Go 值转换为数据库值
func (ss StringSlice) Value() (driver.Value, error) {
	if len(ss) == 0 {
		// 返回空 JSON 数组
		return "[]", nil
	}

	// 使用 json.Marshal 生成正确的 JSON 数组
	bytes, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}

	// 返回 JSON 字符串
	return string(bytes), nil
}
