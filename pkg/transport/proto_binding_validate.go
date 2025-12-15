package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/yb2020/odoc/config"
	pkgi18n "github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// 全局 logger 实例
	logger logging.Logger
)

// InitLogger 初始化 logger
func InitLogger(l logging.Logger) {
	logger = l
}

// ProtoBinding 是一个自定义的 Gin 绑定器，用于将 JSON 请求绑定到 proto 消息
type ProtoBinding struct{}

// Name 返回绑定器的名称
func (ProtoBinding) Name() string {
	return "protojson"
}

// Bind 将请求绑定到 proto 消息
func (ProtoBinding) Bind(req *http.Request, obj interface{}) error {
	// 检查对象是否为 proto.Message
	msg, ok := obj.(proto.Message)
	if !ok {
		return binding.JSON.Bind(req, obj) // 如果不是 proto 消息，使用标准 JSON 绑定
	}

	// 读取请求体
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// 重置请求体，以便其他中间件可以读取
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	// 预处理 JSON 数据，转换字符串格式的数字
	processedBody, err := preprocessJSON(body, obj.(proto.Message))
	if err != nil {
		// 如果预处理失败，使用原始 JSON
		processedBody = body
	}

	// 创建一个临时的 proto message 来存储请求体数据
	tmpMsg := proto.Clone(obj.(proto.Message))
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshaler.Unmarshal(processedBody, tmpMsg); err != nil {
		logger.Debug("msg", "绑定请求体失败", "error", err.Error(), "path", req.URL.Path)
		return err
	}

	// 只更新请求体中存在的字段
	proto.Merge(obj.(proto.Message), tmpMsg)

	// 验证消息
	if v, ok := msg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// BindProto 将请求绑定到 proto 消息
func BindProto(c *gin.Context, obj proto.Message) error {
	// 首先从URL查询参数绑定数据
	// 获取proto消息的反射值
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 只处理结构体类型
	if v.Kind() == reflect.Struct {
		// 创建一个map来存储map类型字段的值
		mapValues := make(map[string]map[string]string)

		// 遍历所有查询参数
		for key, values := range c.Request.URL.Query() {
			if len(values) == 0 {
				continue
			}

			// 检查是否是map类型的参数 (格式如: mapField[key]=value)
			mapMatch := regexp.MustCompile(`^([A-Za-z0-9_]+)\[([^]]+)\]$`).FindStringSubmatch(key)
			if len(mapMatch) == 3 {
				// 提取map字段名和key
				mapFieldName := toCamelCase(mapMatch[1])
				mapKey := mapMatch[2]
				mapValue := values[0]

				// 初始化map字段的值map
				if _, ok := mapValues[mapFieldName]; !ok {
					mapValues[mapFieldName] = make(map[string]string)
				}

				// 存储值
				mapValues[mapFieldName][mapKey] = mapValue
				continue
			}

			// 处理普通字段
			fieldName := toCamelCase(key)
			field := v.FieldByName(fieldName)
			if !field.IsValid() || !field.CanSet() {
				continue
			}

			// 处理数组类型
			if field.Kind() == reflect.Slice {
				elemType := field.Type().Elem()

				// 创建一个新的切片
				newSlice := reflect.MakeSlice(field.Type(), 0, len(values))

				// 遍历所有值并添加到切片中
				for _, value := range values {
					var elemValue reflect.Value

					switch elemType.Kind() {
					case reflect.String:
						elemValue = reflect.ValueOf(value)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
							elemValue = reflect.ValueOf(intVal).Convert(elemType)
						} else {
							continue
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
							elemValue = reflect.ValueOf(uintVal).Convert(elemType)
						} else {
							continue
						}
					case reflect.Float32, reflect.Float64:
						if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
							elemValue = reflect.ValueOf(floatVal).Convert(elemType)
						} else {
							continue
						}
					case reflect.Bool:
						if boolVal, err := strconv.ParseBool(value); err == nil {
							elemValue = reflect.ValueOf(boolVal)
						} else {
							continue
						}
					default:
						continue
					}

					newSlice = reflect.Append(newSlice, elemValue)
				}

				// 设置字段值为新的切片
				field.Set(newSlice)
				continue
			}

			// 处理指针类型
			if field.Kind() == reflect.Ptr {
				// 获取指针指向的类型
				elemType := field.Type().Elem()

				// 创建一个新的元素
				newElem := reflect.New(elemType)

				// 获取元素的值
				elemValue := newElem.Elem()

				// 根据类型设置值
				value := values[0]
				switch elemType.Kind() {
				case reflect.String:
					elemValue.SetString(value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
						elemValue.SetInt(intVal)
					} else {
						continue
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
						elemValue.SetUint(uintVal)
					} else {
						continue
					}
				case reflect.Float32, reflect.Float64:
					if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
						elemValue.SetFloat(floatVal)
					} else {
						continue
					}
				case reflect.Bool:
					if boolVal, err := strconv.ParseBool(value); err == nil {
						elemValue.SetBool(boolVal)
					} else {
						continue
					}
				default:
					continue
				}

				// 设置指针字段
				field.Set(newElem)
				continue
			}

			// 处理单个值的参数
			value := values[0]
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
					field.SetInt(intVal)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
					field.SetUint(uintVal)
				}
			case reflect.Float32, reflect.Float64:
				if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
					field.SetFloat(floatVal)
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(value); err == nil {
					field.SetBool(boolVal)
				}
			}
		}

		// 处理map类型字段
		for mapFieldName, keyValues := range mapValues {
			field := v.FieldByName(mapFieldName)
			if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.Map {
				continue
			}

			// 获取map的key和value类型
			mapType := field.Type()
			keyType := mapType.Key()
			valueType := mapType.Elem()

			// 只处理string类型的key和value
			if keyType.Kind() != reflect.String || valueType.Kind() != reflect.String {
				continue
			}

			// 创建一个新的map
			mapValue := reflect.MakeMap(mapType)

			// 添加所有键值对
			for k, v := range keyValues {
				mapValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}

			// 设置字段值为新的map
			field.Set(mapValue)
		}
	}

	// 如果开启了调试模式，记录请求信息
	cfg := config.GetConfig()
	if cfg != nil && cfg.Debug.EnableRequestLogging {
		clientIP := c.ClientIP()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 记录基本请求信息
		logFields := []interface{}{
			"msg", "请求详情",
			"client_ip", clientIP,
			"path", path,
			"method", method,
		}

		// 如果需要记录请求体
		if cfg.Debug.LogRequestBody {
			// 读取请求体并保存副本
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			// 恢复请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// 限制记录的请求体大小
			maxSize := cfg.Debug.MaxRequestBodySize
			if maxSize <= 0 {
				maxSize = 1024 // 默认1KB
			}

			bodyStr := string(bodyBytes)
			if len(bodyStr) > maxSize {
				bodyStr = bodyStr[:maxSize] + "... (截断)"
			}

			logFields = append(logFields, "request_body", bodyStr)
		}

		logger.Debug(logFields...)
	}

	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if bodyBytes != nil && len(bodyBytes) == 0 {
		return nil
	}
	if err != nil {
		logger.Debug("msg", "读取请求体失败", "error", err.Error(), "path", c.Request.URL.Path)
		return err
	}

	// 恢复请求体，以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 预处理 JSON 数据，转换字符串格式的数字
	processedBody, err := preprocessJSON(bodyBytes, obj.(proto.Message))
	if err != nil {
		// 如果预处理失败，使用原始 JSON
		processedBody = bodyBytes
	}

	// 创建一个临时的 proto message 来存储请求体数据
	tmpMsg := proto.Clone(obj.(proto.Message))
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshaler.Unmarshal(processedBody, tmpMsg); err != nil {
		logger.Debug("msg", "绑定请求体失败", "error", err.Error(), "path", c.Request.URL.Path)
		return err
	}

	// 只更新请求体中存在的字段
	proto.Merge(obj.(proto.Message), tmpMsg)

	// 验证参数
	return ValidateProto(obj, c)
}

// toCamelCase 将蛇形命名转换为驼峰命名
// 例如：user_id -> UserId
func toCamelCase(s string) string {
	// 处理空字符串
	if s == "" {
		return s
	}

	// 分割字符串
	parts := strings.Split(s, "_")

	// 转换每个部分的首字母为大写
	for i := 0; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}

	// 拼接结果
	return strings.Join(parts, "")
}

// ValidateProto 验证 proto 消息中的参数
func ValidateProto(obj proto.Message, c *gin.Context) error {
	// 使用 protoc-gen-validate 生成的验证方法
	if validator, ok := obj.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			// 使用 i18n 本地化错误消息
			return localizeValidationError(err, c)
		}
		return nil
	}

	// 当 protoc-gen-validate 生成的验证不可用时，回退到基于反射的验证
	// 注意：这个分支仅作为备用方案，建议使用 protoc-gen-validate 生成的验证代码
	return validateProtoMessage(obj.ProtoReflect(), c)
}

// validateProtoMessage 使用反射验证 proto 消息中的字段
func validateProtoMessage(msg protoreflect.Message, c *gin.Context) error {
	fields := msg.Descriptor().Fields()

	// 遍历所有字段
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := string(field.Name())
		fieldValue := msg.Get(field)

		// 从字段选项和注释中提取验证规则
		options := field.Options()
		validationRules := extractValidationRules(options)

		// 从注释中提取验证规则
		commentRules := extractValidationRulesFromComment(field)

		// 合并规则，注释中的规则优先级更高
		for k, v := range commentRules {
			validationRules[k] = v
		}

		// 如果有 required 规则，检查字段是否为空
		if validationRules["required"] == "true" {
			if isEmptyValue(fieldValue, field.Kind()) {
				return localizeError("validation.field.required", map[string]interface{}{
					"Field": fieldName,
				}, c)
			}
		}

		// 如果有 min_len 规则，检查字符串长度是否符合要求
		if minLen, ok := validationRules["min_len"]; ok && field.Kind() == protoreflect.StringKind {
			strValue := fieldValue.String()
			if len(strValue) < parseInt(minLen) {
				return localizeError("validation.field.min_length", map[string]interface{}{
					"Field":     fieldName,
					"MinLength": minLen,
				}, c)
			}
		}

		// 如果有 max_len 规则，检查字符串长度是否符合要求
		if maxLen, ok := validationRules["max_len"]; ok && field.Kind() == protoreflect.StringKind {
			strValue := fieldValue.String()
			if len(strValue) > parseInt(maxLen) {
				return localizeError("validation.field.max_length", map[string]interface{}{
					"Field":     fieldName,
					"MaxLength": maxLen,
				}, c)
			}
		}

		// 如果有 email 规则，检查是否是有效的邮箱格式
		if email, ok := validationRules["email"]; ok && email == "true" && field.Kind() == protoreflect.StringKind {
			strValue := fieldValue.String()
			if !isValidEmail(strValue) {
				return localizeError("validation.field.email", map[string]interface{}{
					"Field": fieldName,
				}, c)
			}
		}

		// 如果有 pattern 规则，检查字符串是否符合正则表达式
		if pattern, ok := validationRules["pattern"]; ok && field.Kind() == protoreflect.StringKind {
			strValue := fieldValue.String()
			matched, err := regexp.MatchString(pattern, strValue)
			if err != nil || !matched {
				return localizeError("validation.field.pattern", map[string]interface{}{
					"Field": fieldName,
				}, c)
			}
		}

		// 如果有 not_in 规则，检查字段值是否在禁止列表中
		if notIn, ok := validationRules["not_in"]; ok {
			strValue := fmt.Sprintf("%v", fieldValue.Interface())
			// 解析 not_in 值，格式如 ["0", "none", "null"]
			notInValues := parseStringList(notIn)
			for _, v := range notInValues {
				if strValue == v {
					return localizeError("validation.field.not_in", map[string]interface{}{
						"Field": fieldName,
					}, c)
				}
			}
		}

		// 递归验证嵌套消息
		// if field.Kind() == protoreflect.MessageKind && !field.IsMap() && !isWellKnownType(field) {
		// 	if !msg.Has(field) {
		// 		continue // 跳过未设置的嵌套消息
		// 	}
		// 	nestedMsg := fieldValue.Message()
		// 	if err := validateProtoMessage(nestedMsg, c); err != nil {
		// 		return fmt.Errorf("%s.%s", fieldName, err)
		// 	}
		// }
		// 递归验证嵌套消息   2025-06-20 日增加，修复json中对象缺少字段问题的验证
		if field.Kind() == protoreflect.MessageKind && !field.IsMap() && !isWellKnownType(field) {
			if !msg.Has(field) {
				continue // 跳过未设置的嵌套消息
			}

			// 添加类型检查，确保字段值是消息类型而不是列表
			if fieldValue.IsValid() {
				// 尝试安全地获取消息
				defer func() {
					if r := recover(); r != nil {
						// 捕获可能的 panic，如 "cannot convert list to message"
						// 这里可以记录日志，但不中断处理
					}
				}()

				nestedMsg := fieldValue.Message()
				if err := validateProtoMessage(nestedMsg, c); err != nil {
					return fmt.Errorf("%s.%s", fieldName, err)
				}
			}
		}
	}

	return nil
}

// parseStringList 解析字符串列表，格式如 ["0", "none", "null"]
func parseStringList(listStr string) []string {
	// 去除前后的 [ ]
	listStr = strings.TrimSpace(listStr)
	if strings.HasPrefix(listStr, "[") && strings.HasSuffix(listStr, "]") {
		listStr = listStr[1 : len(listStr)-1]
	}

	// 分割元素
	var result []string
	elements := strings.Split(listStr, ",")
	for _, e := range elements {
		e = strings.TrimSpace(e)
		// 去除引号
		if (strings.HasPrefix(e, "\"") && strings.HasSuffix(e, "\"")) ||
			(strings.HasPrefix(e, "'") && strings.HasSuffix(e, "'")) {
			e = e[1 : len(e)-1]
		}
		result = append(result, e)
	}

	return result
}

// extractValidationRules 从字段选项中提取验证规则
func extractValidationRules(options protoreflect.ProtoMessage) map[string]string {
	// 这里简化处理，实际上需要通过反射获取自定义选项
	// 在实际项目中，您可能需要使用 protoc-gen-validate 生成的代码
	// 或者在 proto 文件中使用自定义选项

	// 默认返回空的验证规则集合
	return map[string]string{}
}

// extractValidationRulesFromComment 从字段注释中提取验证规则
func extractValidationRulesFromComment(field protoreflect.FieldDescriptor) map[string]string {
	// 初始化空的验证规则集合
	rules := map[string]string{}

	// 获取字段注释
	comments := field.ParentFile().SourceLocations().ByDescriptor(field)
	if len(comments.LeadingComments) == 0 {
		return rules
	}

	// 解析注释中的验证规则
	comment := comments.LeadingComments

	// 寻找 @validate 标记
	validateIdx := strings.Index(comment, "@validate")
	if validateIdx == -1 {
		return rules
	}

	// 提取验证规则部分
	validateStr := comment[validateIdx+len("@validate"):]
	// 如果有多行注释，只取当前行
	if newlineIdx := strings.Index(validateStr, "\n"); newlineIdx != -1 {
		validateStr = validateStr[:newlineIdx]
	}

	// 去除前后空白
	validateStr = strings.TrimSpace(validateStr)

	// 解析规则对，格式如：min_len=3 max_len=50 pattern="^[a-zA-Z0-9_]+$"
	rulePairs := strings.Fields(validateStr)
	for _, pair := range rulePairs {
		// 分割键值对
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		// 如果值是引号包裹的，去除引号
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		rules[key] = value
	}

	return rules
}

// isEmptyValue 检查字段值是否为空
func isEmptyValue(value protoreflect.Value, kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.BoolKind:
		return !value.Bool()
	case protoreflect.StringKind:
		return value.String() == ""
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int() == 0
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint() == 0
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return value.Float() == 0
	case protoreflect.BytesKind:
		return len(value.Bytes()) == 0
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return !value.IsValid() || value.Message().IsValid()
	default:
		return false
	}
}

// parseInt 将字符串转换为整数
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// isValidEmail 检查字符串是否是有效的邮箱格式
func isValidEmail(email string) bool {
	// 简单的邮箱格式验证，实际项目中可能需要更复杂的验证
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// isWellKnownType 检查字段是否是 Protocol Buffers 的内置类型
func isWellKnownType(field protoreflect.FieldDescriptor) bool {
	// 检查是否是 google.protobuf.* 包中的类型
	msgName := string(field.Message().FullName())
	return strings.HasPrefix(msgName, "google.protobuf.")
}

// LocalizedError 自定义本地化错误类型
type LocalizedError struct {
	MessageID string                 // 原始消息ID
	Message   string                 // 本地化后的消息
	Data      map[string]interface{} // 模板数据
}

// Error 实现 error 接口
func (e *LocalizedError) Error() string {
	return e.Message
}

// GetMessageID 获取原始消息ID
func (e *LocalizedError) GetMessageID() string {
	return e.MessageID
}

// GetData 获取模板数据
func (e *LocalizedError) GetData() map[string]interface{} {
	return e.Data
}

// localizeError 本地化错误消息
func localizeError(messageID string, data map[string]interface{}, c *gin.Context) error {
	// 打印调试信息
	if logger != nil {
		logger.Debug("msg", "localizeError called", "messageID", messageID, "data", fmt.Sprintf("%v", data))
	}

	localizer := pkgi18n.GetLocalizer()
	if localizer == nil {
		if logger != nil {
			logger.Debug("msg", "localizer is nil, using default error message")
		}
		// 如果本地化器不可用，使用默认错误消息
		var msg string
		switch messageID {
		case "validation.field.required":
			msg = fmt.Sprintf("字段 %s 不能为空", data["Field"])
		case "validation.field.min_length":
			msg = fmt.Sprintf("字段 %s 长度不能小于 %s", data["Field"], data["MinLength"])
		case "validation.field.max_length":
			msg = fmt.Sprintf("字段 %s 长度不能大于 %s", data["Field"], data["MaxLength"])
		case "validation.field.email":
			msg = fmt.Sprintf("字段 %s 不是有效的邮箱格式", data["Field"])
		default:
			msg = fmt.Sprintf("验证错误: %s", messageID)
		}
		return &LocalizedError{
			MessageID: messageID,
			Message:   msg,
			Data:      data,
		}
	}

	// 打印请求头信息
	if c != nil && logger != nil {
		logger.Debug("msg", "Request headers", "Accept-Language", c.GetHeader("Accept-Language"))
	}

	// 使用本地化器本地化错误消息
	msg := localizer.LocalizeWithData(messageID, data, c)
	if logger != nil {
		logger.Debug("msg", "Localized message", "messageID", messageID, "localizedMessage", msg)
	}

	return &LocalizedError{
		MessageID: messageID,
		Message:   msg,
		Data:      data,
	}
}

// localizeValidationError 本地化 protoc-gen-validate 生成的验证错误
func localizeValidationError(err error, c *gin.Context) error {
	// 打印调试信息
	if logger != nil {
		logger.Debug("msg", "localizeValidationError called", "error", err)
	}

	// 解析错误消息，提取字段名和错误类型
	errMsg := err.Error()

	// 处理新的错误消息格式："invalid CreateUserRequest.User: embedded message failed validation | caused by: invalid User.Username: value length must be between 3 and 50 runes, inclusive"
	// 先检查是否包含 "caused by:" 这个模式
	if strings.Contains(errMsg, "caused by:") {
		// 提取 "caused by:" 后面的部分
		causedByParts := strings.Split(errMsg, "caused by:")
		if len(causedByParts) > 1 {
			detailMsg := strings.TrimSpace(causedByParts[1])

			// 处理邮箱验证错误
			if strings.Contains(detailMsg, "email") {
				// 尝试从错误消息中提取字段名
				fieldName := "Email" // 默认值

				// 从形如 "invalid User.Email: value must be a valid email address" 中提取字段名
				fieldParts := strings.Split(detailMsg, ":")
				if len(fieldParts) > 0 {
					// 提取字段名（如 "User.Email"）
					fieldNamePart := strings.TrimPrefix(fieldParts[0], "invalid ")
					fieldName = fieldNamePart
				}

				if logger != nil {
					logger.Debug("msg", "Found email validation error", "field", fieldName)
				}
				return localizeError("validation.field.email", map[string]interface{}{
					"Field": fieldName,
				}, c)
			}

			// 处理最小长度验证错误
			if strings.Contains(detailMsg, "length must be between") {
				// 尝试从错误消息中提取字段名和最小长度
				fieldName := "" // 默认值
				minLen := "3"   // 默认值

				// 从形如 "invalid User.Username: value length must be between 3 and 50 runes, inclusive" 中提取字段名
				fieldParts := strings.Split(detailMsg, ":")
				if len(fieldParts) > 0 {
					// 提取字段名（如 "User.Username"）
					fieldNamePart := strings.TrimPrefix(fieldParts[0], "invalid ")
					fieldName = fieldNamePart

					// 尝试提取最小长度
					if len(fieldParts) > 1 {
						// 使用正则表达式提取数字
						lenMatch := regexp.MustCompile(`between (\d+) and`).FindStringSubmatch(fieldParts[1])
						if len(lenMatch) > 1 {
							minLen = lenMatch[1]
						}
					}
				}

				if logger != nil {
					logger.Debug("msg", "Found min length validation error", "field", fieldName, "minLength", minLen)
				}
				return localizeError("validation.field.min_length", map[string]interface{}{
					"Field":     fieldName,
					"MinLength": minLen,
				}, c)
			}
		}
	}

	// 如果新的解析逻辑失败，回退到原来的解析逻辑
	// 对于邮箱验证错误
	if strings.Contains(errMsg, "valid email address") {
		// 尝试从错误消息中提取字段名
		fieldName := "Email" // 默认值

		// 尝试提取字段名
		parts := strings.Split(errMsg, "|")
		for _, part := range parts {
			if strings.Contains(part, "invalid") && strings.Contains(part, "Email") {
				// 从形如 "invalid User.Email: value must be a valid email address" 中提取字段名
				fieldParts := strings.Split(strings.TrimSpace(part), ":")
				if len(fieldParts) > 0 {
					// 提取字段名（如 "User.Email"）
					fieldNamePart := strings.TrimPrefix(fieldParts[0], "invalid ")
					fieldName = fieldNamePart
				}
				break
			}
		}

		if logger != nil {
			logger.Debug("msg", "Found email validation error", "field", fieldName)
		}
		return localizeError("validation.field.email", map[string]interface{}{
			"Field": fieldName,
		}, c)
	}

	// 对于最小长度验证错误
	if strings.Contains(errMsg, "min_len") {
		// 尝试从错误消息中提取字段名和最小长度
		fieldName := "" // 默认值
		minLen := "3"   // 默认值

		// 尝试提取字段名
		parts := strings.Split(errMsg, "|")
		for _, part := range parts {
			if strings.Contains(part, "invalid") {
				// 从形如 "invalid User.Username: value length must be at least 3 runes" 中提取字段名
				fieldParts := strings.Split(strings.TrimSpace(part), ":")
				if len(fieldParts) > 0 {
					// 提取字段名（如 "User.Username"）
					fieldNamePart := strings.TrimPrefix(fieldParts[0], "invalid ")
					fieldName = fieldNamePart

					// 尝试提取最小长度
					if len(fieldParts) > 1 && strings.Contains(fieldParts[1], "at least") {
						lenMatch := regexp.MustCompile(`at least (\d+)`).FindStringSubmatch(fieldParts[1])
						if len(lenMatch) > 1 {
							minLen = lenMatch[1]
						}
					}
				}
				break
			}
		}

		if logger != nil {
			logger.Debug("msg", "Found min length validation error", "field", fieldName, "minLength", minLen)
		}
		return localizeError("validation.field.min_length", map[string]interface{}{
			"Field":     fieldName,
			"MinLength": minLen,
		}, c)
	}

	// 如果无法解析错误，返回原始错误
	if logger != nil {
		logger.Debug("msg", "Could not parse validation error, returning original error", "error", err)
	}
	return err
}

// cleanJSONString 清理 JSON 字符串，移除可能导致解析错误的字符
func cleanJSONString(jsonStr string) string {
	// 移除所有换行符、回车符和制表符
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\r", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\t", "")

	// 移除多余的空格
	jsonStr = strings.TrimSpace(jsonStr)

	// 移除可能导致问题的 Unicode 字符
	jsonStr = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1 // 删除控制字符
		}
		return r
	}, jsonStr)

	// 修复常见的 JSON 格式错误

	// 1. 修复对象末尾多余的逗号
	// 将 ,} 替换为 }
	jsonStr = regexp.MustCompile(`,\s*}`).ReplaceAllString(jsonStr, "}")

	// 2. 修复数组末尾多余的逗号
	// 将 ,] 替换为 ]
	jsonStr = regexp.MustCompile(`,\s*\]`).ReplaceAllString(jsonStr, "]")

	return jsonStr
}

// preprocessJSON 预处理 JSON 数据，基于 protobuf 类型定义转换字符串格式的数字
func preprocessJSON(data []byte, msg proto.Message) ([]byte, error) {
	// 解析 JSON 到 map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return nil, err
	}

	// 获取消息描述符
	msgDesc := msg.ProtoReflect().Descriptor()

	// 处理 JSON 数据
	processJSONMap(jsonMap, msgDesc)

	// 将处理后的 map 转回 JSON
	processedData, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}

	return processedData, nil
}

// processJSONMap 递归处理 JSON map，基于 protobuf 类型定义转换字符串格式的数字
func processJSONMap(jsonMap map[string]interface{}, msgDesc protoreflect.MessageDescriptor) {
	fields := msgDesc.Fields()

	// 遍历所有字段
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		protoFieldName := string(field.Name())

		// 将 proto 字段名转换为 JSON 字段名（小驼峰式）
		jsonFieldName := protoToJSONFieldName(protoFieldName)

		// 检查 JSON 中是否有该字段
		if value, ok := jsonMap[jsonFieldName]; ok {
			// 根据字段类型处理
			switch field.Kind() {
			case protoreflect.Int32Kind, protoreflect.Int64Kind:
				// 如果是字符串格式的数字，转换为 int64
				if strValue, ok := value.(string); ok {
					if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
						jsonMap[jsonFieldName] = intValue
					}
				}
			case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
				// 如果是字符串格式的数字，转换为 uint64
				if strValue, ok := value.(string); ok {
					// 移除前导零，避免被误解为八进制
					cleanValue := strings.TrimLeft(strValue, "0")
					if cleanValue == "" {
						cleanValue = "0" // 如果全是零，保留一个零
					}
					if uintValue, err := strconv.ParseUint(cleanValue, 10, 64); err == nil {
						jsonMap[jsonFieldName] = uintValue
					}
				}
			case protoreflect.FloatKind, protoreflect.DoubleKind:
				// 如果是字符串格式的数字，转换为 float64
				if strValue, ok := value.(string); ok {
					if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
						jsonMap[jsonFieldName] = floatValue
					}
				}
			case protoreflect.BoolKind:
				// 如果是字符串格式的布尔值，转换为 bool
				if strValue, ok := value.(string); ok {
					if boolValue, err := strconv.ParseBool(strValue); err == nil {
						jsonMap[jsonFieldName] = boolValue
					}
				}
			case protoreflect.MessageKind:
				// 递归处理嵌套消息
				if nestedMap, ok := value.(map[string]interface{}); ok {
					processJSONMap(nestedMap, field.Message())
				} else if nestedSlice, ok := value.([]interface{}); ok {
					// 检查是否是重复的消息字段
					if field.IsList() && field.Message() != nil {
						// 处理数组中的每个元素
						for i, item := range nestedSlice {
							if itemMap, ok := item.(map[string]interface{}); ok {
								// 递归处理数组中的对象
								processJSONMap(itemMap, field.Message())

								// 特别处理数组元素中的整型字段
								processArrayItemFields(itemMap, field.Message())

								// 更新数组中的对象
								nestedSlice[i] = itemMap
							}
						}
						// 更新数组
						jsonMap[jsonFieldName] = nestedSlice
					}
				}
			}
		}
	}
}

// processArrayItemFields 处理数组元素中的整型字段，确保前导零被正确处理
func processArrayItemFields(itemMap map[string]interface{}, msgDesc protoreflect.MessageDescriptor) {
	fields := msgDesc.Fields()

	// 遍历所有字段
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		jsonFieldName := string(field.Name())
		// 将 proto 字段名转换为 JSON 字段名（小驼峰式）
		// jsonFieldName := protoToJSONFieldName(protoFieldName)
		// 检查元素中是否有该字段
		if value, ok := itemMap[jsonFieldName]; ok {
			// 根据字段类型处理
			switch field.Kind() {
			case protoreflect.Int32Kind, protoreflect.Int64Kind:
				// 如果是字符串格式的数字，转换为 int64
				if strValue, ok := value.(string); ok {
					if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
						itemMap[jsonFieldName] = intValue
					}
				}
			case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
				// 如果是字符串格式的数字，转换为 uint64
				if strValue, ok := value.(string); ok {
					// 移除前导零，避免被误解为八进制
					cleanValue := strings.TrimLeft(strValue, "0")
					if cleanValue == "" {
						cleanValue = "0" // 如果全是零，保留一个零
					}
					if uintValue, err := strconv.ParseUint(cleanValue, 10, 64); err == nil {
						itemMap[jsonFieldName] = uintValue
					}
				}
			case protoreflect.FloatKind, protoreflect.DoubleKind:
				// 如果是字符串格式的数字，转换为 float64
				if strValue, ok := value.(string); ok {
					if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
						itemMap[jsonFieldName] = floatValue
					}
				}
			case protoreflect.BoolKind:
				// 如果是字符串格式的布尔值，转换为 bool
				if strValue, ok := value.(string); ok {
					if boolValue, err := strconv.ParseBool(strValue); err == nil {
						itemMap[jsonFieldName] = boolValue
					}
				}
			}
		} else { // 这里兼容json中缺少字段的情况
			// 字段在JSON中不存在，根据字段类型设置默认值
			switch field.Kind() {
			case protoreflect.Int32Kind, protoreflect.Int64Kind:
				// 为整型字段设置默认值0
				itemMap[jsonFieldName] = int64(0)
			case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
				// 为无符号整型字段设置默认值0
				itemMap[jsonFieldName] = uint64(0)
			case protoreflect.FloatKind, protoreflect.DoubleKind:
				// 为浮点型字段设置默认值0.0
				itemMap[jsonFieldName] = float64(0)
			case protoreflect.BoolKind:
				// 为布尔型字段设置默认值false
				itemMap[jsonFieldName] = false
			case protoreflect.StringKind:
				// 为字符串字段设置默认值空字符串
				itemMap[jsonFieldName] = ""
			}
		}
	}
}

// protoToJSONFieldName 将 proto 字段名转换为 JSON 字段名
func protoToJSONFieldName(name string) string {
	// 如果字段名已经是小驼峰式，直接返回
	if len(name) == 0 {
		return name
	}
	// 将首字母转为小写
	return strings.ToLower(name[:1]) + name[1:]
}
