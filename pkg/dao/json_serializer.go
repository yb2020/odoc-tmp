package dao

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// JSONField 标记需要自动序列化/反序列化的字段
type JSONField struct {
	Field     string // 结构体字段名
	JSONField string // JSON 存储字段名
}

// RegisterJSONSerializer 注册 JSON 序列化器
func RegisterJSONSerializer(db *gorm.DB) {
	// 注册创建前回调
	db.Callback().Create().Before("gorm:create").Register("json_serializer:before_create", func(db *gorm.DB) {
		if db.Statement.Schema == nil || db.Statement.Dest == nil {
			return
		}

		// 获取模型类型和值
		destValue := reflect.ValueOf(db.Statement.Dest)
		if destValue.Kind() == reflect.Ptr {
			destValue = destValue.Elem()
		}

		// 如果不是结构体，直接返回
		if destValue.Kind() != reflect.Struct {
			return
		}

		// 处理结构体的每个字段
		destType := destValue.Type()
		for i := 0; i < destValue.NumField(); i++ {
			field := destType.Field(i)
			jsonTag := field.Tag.Get("json_serialize")
			if jsonTag != "" {
				// 找到对应的 JSON 存储字段
				jsonField := field.Tag.Get("json_field")
				if jsonField != "" {
					// 获取字段值
					fieldValue := destValue.Field(i).Interface()

					// 序列化字段值
					jsonData, err := json.Marshal(fieldValue)
					if err != nil {
						db.AddError(fmt.Errorf("序列化字段 %s 失败: %w", field.Name, err))
						continue
					}

					// 设置 JSON 存储字段的值
					jsonFieldValue := destValue.FieldByName(jsonField)
					if jsonFieldValue.IsValid() && jsonFieldValue.CanSet() {
						jsonFieldValue.SetString(string(jsonData))
						fmt.Printf("已序列化字段 %s 到 %s: %s\n", field.Name, jsonField, string(jsonData))
					} else {
						db.AddError(fmt.Errorf("无法设置字段 %s 的值", jsonField))
					}
				}
			}
		}
	})

	// 注册查询后回调
	db.Callback().Query().After("gorm:after_query").Register("json_serializer:after_query", func(db *gorm.DB) {
		if db.Statement.Schema == nil || db.Statement.Dest == nil {
			return
		}

		// 处理单个对象
		destValue := reflect.ValueOf(db.Statement.Dest)
		if destValue.Kind() == reflect.Ptr {
			destValue = destValue.Elem()

			// 如果是结构体，直接处理
			if destValue.Kind() == reflect.Struct {
				deserializeJSONFields(destValue)
				return
			}

			// 如果是切片，处理每个元素
			if destValue.Kind() == reflect.Slice {
				for i := 0; i < destValue.Len(); i++ {
					item := destValue.Index(i)
					if item.Kind() == reflect.Struct {
						deserializeJSONFields(item)
					} else if item.Kind() == reflect.Ptr && item.Elem().Kind() == reflect.Struct {
						deserializeJSONFields(item.Elem())
					}
				}
			}
		}
	})
}

// deserializeJSONFields 反序列化 JSON 字段
func deserializeJSONFields(value reflect.Value) {
	// 如果不是结构体，直接返回
	if value.Kind() != reflect.Struct {
		return
	}

	// 处理结构体的每个字段
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		jsonTag := field.Tag.Get("json_serialize")
		if jsonTag != "" {
			// 找到对应的 JSON 存储字段
			jsonField := field.Tag.Get("json_field")
			if jsonField != "" {
				// 获取 JSON 存储字段的值
				jsonFieldValue := value.FieldByName(jsonField)
				if jsonFieldValue.IsValid() && jsonFieldValue.String() != "" {
					// 获取目标字段
					targetField := value.Field(i)
					if targetField.CanSet() {
						// 创建目标类型的新实例
						targetType := targetField.Type()
						newValue := reflect.New(targetType).Interface()

						// 反序列化 JSON 数据
						jsonData := jsonFieldValue.String()
						if err := json.Unmarshal([]byte(jsonData), newValue); err != nil {
							fmt.Printf("反序列化字段 %s 失败: %v\n", field.Name, err)
							continue
						}

						// 设置目标字段的值
						targetField.Set(reflect.ValueOf(newValue).Elem())
						fmt.Printf("已反序列化字段 %s 从 %s\n", field.Name, jsonField)
					}
				}
			}
		}
	}
}
