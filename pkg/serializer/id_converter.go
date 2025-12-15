package serializer

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// Int64ToUint64JSONMarshaler 是一个自定义的JSON编码器，
// 它会自动将结构体中名为"ID"的int64字段转换为uint64
type Int64ToUint64JSONMarshaler struct {
	Value interface{}
}

// MarshalJSON 实现json.Marshaler接口，自动将int64 ID转换为uint64
func (m Int64ToUint64JSONMarshaler) MarshalJSON() ([]byte, error) {
	// 如果值为nil，直接返回null
	if m.Value == nil {
		return []byte("null"), nil
	}

	// 获取值的反射
	v := reflect.ValueOf(m.Value)
	
	// 如果是指针，获取其指向的值
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return []byte("null"), nil
		}
		v = v.Elem()
	}
	
	// 只处理结构体
	if v.Kind() != reflect.Struct {
		return json.Marshal(m.Value)
	}
	
	// 创建一个映射来存储结构体的字段
	result := make(map[string]interface{})
	
	// 遍历结构体的字段
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		
		// 获取JSON标签
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		
		// 解析JSON标签
		name := field.Name
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}
		
		// 获取字段值
		fieldValue := v.Field(i)
		
		// 如果是嵌套结构体，递归处理
		if fieldValue.Kind() == reflect.Struct {
			result[name] = Int64ToUint64JSONMarshaler{Value: fieldValue.Interface()}
			continue
		}
		
		// 如果是结构体切片，递归处理每个元素
		if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.Struct {
			slice := make([]interface{}, fieldValue.Len())
			for j := 0; j < fieldValue.Len(); j++ {
				slice[j] = Int64ToUint64JSONMarshaler{Value: fieldValue.Index(j).Interface()}
			}
			result[name] = slice
			continue
		}
		
		// 处理ID字段，将int64转换为uint64
		if (name == "id" || name == "ID") && fieldValue.Kind() == reflect.Int64 {
			result[name] = uint64(fieldValue.Int())
		} else {
			result[name] = fieldValue.Interface()
		}
	}
	
	return json.Marshal(result)
}

// JSONResponse 用于包装JSON响应
type JSONResponse struct {
	Data interface{}
}

// MarshalJSON 实现json.Marshaler接口，自动转换ID
func (r JSONResponse) MarshalJSON() ([]byte, error) {
	return Int64ToUint64JSONMarshaler{Value: r.Data}.MarshalJSON()
}

// StringToInt64 将字符串转换为int64
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Int64ToString 将int64转换为字符串
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
