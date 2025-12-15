package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ConvertObject2Bytes(object interface{}) ([]byte, error) {
	// 检查对象是否为空
	if object == nil {
		return nil, fmt.Errorf("object cannot be nil")
	}

	// 检查对象是否为空切片或空数组
	objectValue := reflect.ValueOf(object)
	if (objectValue.Kind() == reflect.Slice || objectValue.Kind() == reflect.Array) && objectValue.Len() == 0 {
		return nil, fmt.Errorf("slice or array object cannot be empty")
	}

	// 将对象序列化为JSON
	body, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 解析userMetadata JSON字符串为map
func ParseUserMetadata(userMetadata string) (map[string]string, error) {
	var userMetadataMap map[string]string
	if err := json.Unmarshal([]byte(userMetadata), &userMetadataMap); err == nil {
		return userMetadataMap, nil
	}
	return nil, fmt.Errorf("failed to parse user metadata")
}

// 通过userMetadata获取key对应的value
func GetByUserMetadata(userMetadata string, key string) (string, error) {
	userMetadataMap, err := ParseUserMetadata(userMetadata)
	if err != nil {
		return "", err
	}
	return userMetadataMap[key], nil
}

// 通过map获取key对应的value
func GetByKey(userMetadataMap map[string]string, key string) (string, error) {
	return userMetadataMap[key], nil
}
