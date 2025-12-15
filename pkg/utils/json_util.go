package utils

import "encoding/json"

// DumpJson 将对象转换为 JSON 字符串
func DumpJson(sourceJsonStr interface{}) string {
	jsonBytes, err := json.Marshal(sourceJsonStr)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// ParseJsonByString 将 JSON 字符串解析为对象
func ParseJsonByString(sourceJsonStr string, targetObj interface{}) error {
	err := json.Unmarshal([]byte(sourceJsonStr), targetObj)
	if err != nil {
		return err
	}
	return nil
}

// ParseJsonByObj 将 对象为对象
func ParseJsonByObj(source interface{}, target interface{}) error {
	// 如果源是字符串，直接解析
	if str, ok := source.(string); ok {
		return ParseJsonByString(str, target)
	}
	jsonStrString := DumpJson(source)
	err := ParseJsonByString(jsonStrString, target)
	if err != nil {
		return err
	}
	return nil
}
