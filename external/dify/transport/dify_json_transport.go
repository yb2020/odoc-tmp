package transport

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ToJson 将对象转换为JSON字符串
func ToJson(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJson 将JSON字符串转换为对象
func FromJson(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

// IsValidJson 检查字符串是否为有效的JSON
func IsValidJson(str string) bool {
	var js interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

// ParseJsonMap 解析JSON字符串为map
func ParseJsonMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ExtractErrorInfo 从错误响应中提取错误信息
func ExtractErrorInfo(jsonStr string) (string, string, string) {
	errorCode := "unknown_error"
	errorMessage := jsonStr
	params := ""

	if jsonStr != "" && IsValidJson(jsonStr) {
		data, err := ParseJsonMap(jsonStr)
		if err == nil {
			// 提取错误代码
			if code, ok := data["error_code"]; ok && code != nil {
				errorCode = code.(string)
			} else if code, ok := data["code"]; ok && code != nil {
				errorCode = toString(code)
			}

			// 提取错误消息
			if message, ok := data["error_message"]; ok && message != nil {
				errorMessage = message.(string)
			} else if message, ok := data["message"]; ok && message != nil {
				errorMessage = message.(string)
			}

			// 提取参数
			if p, ok := data["params"]; ok && p != nil {
				params = toString(p)
			}
		}
	}

	if params != "" {
		errorMessage += " 【" + params + "】"
	}

	return errorCode, errorMessage, params
}

// toString 将任意类型转换为字符串
func toString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return strings.TrimSpace(strings.Trim(strings.Trim(strings.TrimSpace(fmt.Sprintf("%v", v)), "\""), "'"))
	default:
		jsonStr, err := ToJson(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return jsonStr
	}
}
