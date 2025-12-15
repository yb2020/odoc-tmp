package dao

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/model"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/paginate"
	"gorm.io/gorm"
)

// BaseDAO 定义通用的数据访问方法
type BaseDAO[T any] interface {
	// 根据ID查找数据
	FindById(ctx context.Context, id string) (*T, error)

	// 根据IDS查找数据
	FindByIds(ctx context.Context, ids []string) (*[]T, error)

	// 根据ID查找数据，只查找存在的数据（未软删除）
	FindExistById(ctx context.Context, id string) (*T, error)

	// 查找整张表，包含已删除，用ID反序
	FindListAll(ctx context.Context) ([]T, error)

	// 查找整张表，不包含已删除，用ID反序
	FindExistAll(ctx context.Context) ([]T, error)

	// 分页查询，返回结果、总数和可能的错误
	Paginate(ctx context.Context, page, size int32, options *paginate.PaginateOptions) ([]T, int64, error)

	// 批量保存数据
	SaveAll(ctx context.Context, entities *[]T) error

	// 插入表数据
	Save(ctx context.Context, entity *T) error

	// 插入表数据，如果有空值就不处理
	SaveExcludeNull(ctx context.Context, entity *T) error

	// 修改数据，如果有空值就不处理
	ModifyExcludeNull(ctx context.Context, entity *T) error

	// 批量保存数据，如果有空值就不处理
	SaveAllExcludeNull(ctx context.Context, entities *[]T) error

	// 修改数据
	Modify(ctx context.Context, entity *T) error

	// 逻辑删除数据
	DeleteById(ctx context.Context, id string) error

	// 逻辑删除数据
	DeleteByIds(ctx context.Context, ids []string) error

	// 物理删除数据
	RemoveById(ctx context.Context, id string) error
}

// DBType 数据库类型
type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeSQLite   DBType = "sqlite"
	DBTypeMySQL    DBType = "mysql"
)

// GormBaseDAO 是 BaseDAO 的 GORM 实现
type GormBaseDAO[T any] struct {
	db     *gorm.DB
	logger logging.Logger
	dbType DBType
}

// NewGormBaseDAO 创建一个新的 GormBaseDAO 实例
func NewGormBaseDAO[T any](db *gorm.DB, logger logging.Logger) *GormBaseDAO[T] {
	return &GormBaseDAO[T]{
		db:     db,
		logger: logger,
		dbType: detectDBType(db),
	}
}

// NewGormBaseDAOWithType 创建一个指定数据库类型的 GormBaseDAO 实例
func NewGormBaseDAOWithType[T any](db *gorm.DB, logger logging.Logger, dbType DBType) *GormBaseDAO[T] {
	return &GormBaseDAO[T]{
		db:     db,
		logger: logger,
		dbType: dbType,
	}
}

// detectDBType 自动检测数据库类型
func detectDBType(db *gorm.DB) DBType {
	dialectorName := db.Dialector.Name()
	switch dialectorName {
	case "postgres":
		return DBTypePostgres
	case "sqlite":
		return DBTypeSQLite
	case "mysql":
		return DBTypeMySQL
	default:
		return DBTypePostgres // 默认 PostgreSQL
	}
}

// notDeleted 返回未删除条件（兼容不同数据库）
func (d *GormBaseDAO[T]) notDeleted() string {
	if d.dbType == DBTypeSQLite {
		return "is_deleted = 0"
	}
	return "is_deleted = false"
}

// GetDBType 获取数据库类型
func (d *GormBaseDAO[T]) GetDBType() DBType {
	return d.dbType
}

// getDBFromContext 从上下文中获取数据库连接或事务（内部使用）
func (d *GormBaseDAO[T]) getDBFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TransactionContextKey).(*gorm.DB); ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

// GetDB 从上下文中获取数据库连接或事务（公开方法，供子类 DAO 使用）
// 重要：在事务中执行时，必须使用此方法获取 DB，否则 SQLite 会因为锁而阻塞
func (d *GormBaseDAO[T]) GetDB(ctx context.Context) *gorm.DB {
	return d.getDBFromContext(ctx)
}

// FindById 根据ID查找数据
func (d *GormBaseDAO[T]) FindById(ctx context.Context, id string) (*T, error) {
	var entity T
	result := d.db.WithContext(ctx).Where("id = ?", id).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据ID查找数据失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &entity, nil
}

// FindByIds 根据IDS查找数据
func (d *GormBaseDAO[T]) FindByIds(ctx context.Context, ids []string) ([]T, error) {
	var entities []T
	result := d.db.WithContext(ctx).Where("id IN(?)", ids).Order("id DESC").Find(&entities)
	if result.Error != nil {
		d.logger.Error("msg", "查找所有数据失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return entities, nil
}

// FindExistById 根据ID查找数据，只查找存在的数据（未软删除）
func (d *GormBaseDAO[T]) FindExistById(ctx context.Context, id string) (*T, error) {
	var entity T
	result := d.db.WithContext(ctx).Where("id = ? AND "+d.notDeleted(), id).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据ID查找存在数据失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &entity, nil
}

// FindListAll 查找整张表，包含已删除，用ID反序
func (d *GormBaseDAO[T]) FindListAll(ctx context.Context) ([]T, error) {
	var entities []T
	result := d.db.WithContext(ctx).Order("id DESC").Find(&entities)
	if result.Error != nil {
		d.logger.Error("msg", "查找所有数据失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return entities, nil
}

// FindExistAll 查找整张表，不包含已删除，用ID反序
func (d *GormBaseDAO[T]) FindExistAll(ctx context.Context) ([]T, error) {
	var entities []T
	result := d.db.WithContext(ctx).Where(d.notDeleted()).Order("id DESC").Find(&entities)
	if result.Error != nil {
		d.logger.Error("msg", "查找所有存在数据失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return entities, nil
}

// Paginate 使用更灵活的选项进行分页查询
func (d *GormBaseDAO[T]) Paginate(ctx context.Context, page, size int32, options *paginate.PaginateOptions) ([]T, int64, error) {
	var entities []T
	var total int64
	db := d.db.WithContext(ctx)

	// 应用查询条件
	if options != nil {
		// 精确匹配条件
		for field, value := range options.Equals {
			db = db.Where(field+" = ?", value)
		}

		// 模糊匹配条件
		for field, value := range options.Like {
			db = db.Where(field+" LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		}

		// 范围条件
		for field, rangeValue := range options.Ranges {
			if rangeValue.Min != nil {
				db = db.Where(field+" >= ?", rangeValue.Min)
			}
			if rangeValue.Max != nil {
				db = db.Where(field+" <= ?", rangeValue.Max)
			}
		}

		// 自定义条件
		for condition, args := range options.Custom {
			db = db.Where(condition, args...)
		}

		// 排序方式
		for field, direction := range options.OrderBy {
			db = db.Order(field + " " + direction)
		}
	}

	// 获取总记录数
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		d.logger.Error("msg", "获取总数失败", "error", err.Error())
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	result := db.Offset(int(offset)).Limit(int(size)).Find(&entities)
	if result.Error != nil {
		d.logger.Error("msg", "分页查询失败", "error", result.Error.Error())
		return nil, 0, result.Error
	}

	return entities, total, nil
}

// SaveAll 批量保存数据
func (d *GormBaseDAO[T]) SaveAll(ctx context.Context, entities *[]T) error {
	if len(*entities) == 0 {
		return nil
	}

	result := d.db.WithContext(ctx).Create(&entities)
	if result.Error != nil {
		d.logger.Error("msg", "批量保存数据失败", "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// SaveAllExcludeNull 批量保存数据，如果有空值就不处理
func (d *GormBaseDAO[T]) SaveAllExcludeNull(ctx context.Context, entities *[]T) error {
	if len(*entities) == 0 {
		return nil
	}

	// 使用事务处理批量插入
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 处理每个实体
		for i := range *entities {
			entity := &(*entities)[i]

			// 使用反射获取非空字段
			val := reflect.ValueOf(*entity)
			typ := val.Type()
			fields := make(map[string]interface{})
			var fieldNames []string

			// 收集非空字段
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				fieldType := typ.Field(i)
				fieldName := fieldType.Name

				// 处理非嵌入字段
				if !fieldType.Anonymous {
					if !shouldExcludeField(field) {
						// 使用数据库字段名作为键
						dbFieldName := getDBFieldName(fieldType)
						fields[dbFieldName] = field.Interface()
						fieldNames = append(fieldNames, fieldName)
					}
				} else if field.Kind() == reflect.Struct {
					// 处理嵌入字段
					embeddedType := fieldType.Type
					for j := 0; j < field.NumField(); j++ {
						embeddedField := field.Field(j)
						embeddedFieldType := embeddedType.Field(j)
						embeddedFieldName := embeddedFieldType.Name

						if !shouldExcludeField(embeddedField) {
							// 使用数据库字段名作为键
							dbFieldName := getDBFieldName(embeddedFieldType)
							fields[dbFieldName] = embeddedField.Interface()
							fieldNames = append(fieldNames, embeddedFieldName)
						}
					}
				}
			}

			if len(fields) == 0 {
				continue // 跳过没有非空字段的实体
			}

			// 动态获取 BaseModel 中的字段
			baseFieldNames := getBaseModelFieldNames()
			for _, field := range baseFieldNames {
				if !containsString(fieldNames, field) {
					fieldNames = append(fieldNames, field)
				}
			}

			// 使用 Select 指定要插入的字段
			result := tx.Select(fieldNames).Create(entity)
			if result.Error != nil {
				d.logger.Error("msg", "批量保存非空数据失败", "error", result.Error.Error())
				return result.Error
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Save 插入表数据
func (d *GormBaseDAO[T]) Save(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	result := d.getDBFromContext(ctx).Create(entity)
	if result.Error != nil {
		d.logger.Error("msg", "保存数据失败", "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// SaveExcludeNull 插入表数据，如果有空值就不处理
func (d *GormBaseDAO[T]) SaveExcludeNull(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	// 使用反射获取非空字段
	val := reflect.ValueOf(*entity)
	typ := val.Type()
	fields := make(map[string]interface{})
	var fieldNames []string

	// 收集非空字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// 处理非嵌入字段
		if !fieldType.Anonymous {
			if !shouldExcludeField(field) {
				// 使用数据库字段名作为键
				dbFieldName := getDBFieldName(fieldType)
				fields[dbFieldName] = field.Interface()
				fieldNames = append(fieldNames, fieldName)
			}
		} else if field.Kind() == reflect.Struct {
			// 处理嵌入字段
			embeddedType := fieldType.Type
			for j := 0; j < field.NumField(); j++ {
				embeddedField := field.Field(j)
				embeddedFieldType := embeddedType.Field(j)
				embeddedFieldName := embeddedFieldType.Name

				if !shouldExcludeField(embeddedField) {
					// 使用数据库字段名作为键
					dbFieldName := getDBFieldName(embeddedFieldType)
					fields[dbFieldName] = embeddedField.Interface()
					fieldNames = append(fieldNames, embeddedFieldName)
				}
			}
		}
	}

	if len(fields) == 0 {
		return errors.New("no non-null fields to save")
	}

	// 动态获取 BaseModel 中的字段
	baseFieldNames := getBaseModelFieldNames()
	for _, field := range baseFieldNames {
		if !containsString(fieldNames, field) {
			fieldNames = append(fieldNames, field)
		}
	}

	// 使用 Select 指定要插入的字段
	result := d.getDBFromContext(ctx).Select(fieldNames).Create(entity)
	if result.Error != nil {
		d.logger.Error("msg", "保存非空数据失败", "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// getBaseModelFieldNames 获取 BaseModel 中的所有字段名
func getBaseModelFieldNames() []string {
	// 使用反射获取 BaseModel 的字段
	baseModelType := reflect.TypeOf(model.BaseModel{})
	fieldCount := baseModelType.NumField()
	baseFields := make([]string, 0, fieldCount)

	for i := 0; i < fieldCount; i++ {
		field := baseModelType.Field(i)
		baseFields = append(baseFields, field.Name)
	}

	return baseFields
}

// containsString 检查字符串切片中是否包含指定字符串
func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// modifyInternal 内部修改数据方法，excludeNull 参数控制是否排除空值
func (d *GormBaseDAO[T]) modifyInternal(ctx context.Context, entity *T, excludeNull bool) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	// 获取实体ID
	id := d.getEntityId(entity)
	if id == "" {
		return errors.New("entity ID cannot be empty")
	}

	// 使用反射获取字段
	val := reflect.ValueOf(*entity)
	typ := val.Type()
	fields := make(map[string]interface{})

	// 收集字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// 跳过 ID 字段
		if fieldName == "Id" || fieldName == "ID" {
			continue
		}

		// 处理非嵌入字段
		if !fieldType.Anonymous {
			if !excludeNull || !shouldExcludeField(field) {
				// 使用数据库字段名作为键
				dbFieldName := getDBFieldName(fieldType)
				fields[dbFieldName] = field.Interface()
			}
		} else if field.Kind() == reflect.Struct {
			// 处理嵌入字段
			embeddedType := fieldType.Type
			for j := 0; j < field.NumField(); j++ {
				embeddedField := field.Field(j)
				embeddedFieldType := embeddedType.Field(j)
				embeddedFieldName := embeddedFieldType.Name

				// 跳过 ID 字段
				if embeddedFieldName == "Id" || embeddedFieldName == "ID" {
					continue
				}

				if !excludeNull || !shouldExcludeField(embeddedField) {
					// 使用数据库字段名作为键
					dbFieldName := getDBFieldName(embeddedFieldType)
					fields[dbFieldName] = embeddedField.Interface()
				}
			}
		}
	}

	if excludeNull && len(fields) == 0 {
		return errors.New("no non-null fields to update")
	}

	// 获取元数据更新
	updates := handlerBeforeUpdate(ctx, false)

	// 合并字段和元数据, 应该反过来，因为修改时间、修改人要更新
	for key, value := range updates {
		fields[key] = value
	}

	// 直接更新数据库
	result := d.getDBFromContext(ctx).Model(new(T)).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		msg := "修改数据失败"
		if excludeNull {
			msg = "更新非空数据失败"
		}
		d.logger.Error("msg", msg, "id", id, "error", result.Error.Error())
		return result.Error
	}

	// 如果没有行被影响（可能是因为数据没有变化），返回原始实体
	return nil
}

// ModifyExcludeNull 修改数据，如果有空值就不处理
func (d *GormBaseDAO[T]) ModifyExcludeNull(ctx context.Context, entity *T) error {
	return d.modifyInternal(ctx, entity, true)
}

// Modify 修改数据
func (d *GormBaseDAO[T]) Modify(ctx context.Context, entity *T) error {
	return d.modifyInternal(ctx, entity, false)
}

// 从实体中获取ID字段的值
func (d *GormBaseDAO[T]) getEntityId(entity *T) string {
	if entity == nil {
		return ""
	}

	val := reflect.ValueOf(*entity)
	typ := val.Type()
	var id string

	// 首先尝试从顶层字段获取 ID
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		if fieldName == "Id" || fieldName == "ID" {
			id = field.String()
			break
		}
	}

	// 如果顶层没有找到 ID，尝试在嵌入字段中查找
	if id == "" {
		// 检查是否有 BaseModel 嵌入字段
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			// 检查是否是嵌入字段
			if fieldType.Anonymous {
				// 如果是结构体类型，遍历其字段
				if field.Kind() == reflect.Struct {
					embeddedType := fieldType.Type
					for j := 0; j < field.NumField(); j++ {
						embeddedField := field.Field(j)
						embeddedFieldName := embeddedType.Field(j).Name

						if embeddedFieldName == "Id" || embeddedFieldName == "ID" {
							id = embeddedField.String()
							break
						}
					}
				}

				// 如果已找到 ID，跳出循环
				if id != "" {
					break
				}
			}
		}
	}

	return id
}

// DeleteById 逻辑删除数据
func (d *GormBaseDAO[T]) DeleteById(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	updates := handlerBeforeUpdate(ctx, true)
	// 直接更新数据库
	result := d.getDBFromContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		d.logger.Error("msg", "逻辑删除数据失败", "id", id, "error", result.Error.Error())
		return result.Error
	}

	return nil
}

// DeleteByIds 逻辑删除数据
func (d *GormBaseDAO[T]) DeleteByIds(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return errors.New("ids cannot be empty")
	}

	updates := handlerBeforeUpdate(ctx, true)
	// 直接更新数据库
	result := d.getDBFromContext(ctx).Model(new(T)).Where("id in ?", ids).Updates(updates)
	if result.Error != nil {
		d.logger.Error("msg", "逻辑删除数据失败", "ids", ids, "error", result.Error.Error())
		return result.Error
	}

	return nil
}

// 处理BeforeDelete, context中需要包含用户信息
func handlerBeforeUpdate(ctx context.Context, isDeleted bool) map[string]interface{} {
	// 创建更新映射
	updates := make(map[string]interface{})
	if isDeleted {
		updates["is_deleted"] = true
	}
	updates["updated_at"] = time.Now()

	// 从上下文中获取用户信息
	if uc := userContext.GetUserContext(ctx); uc != nil {
		if uc.UserId != "" {
			updates["modifier_id"] = uc.UserId
		}
		if uc.Username != "" {
			updates["modifier"] = uc.Username
		}
	}

	return updates
}

// RemoveById 物理删除数据
func (d *GormBaseDAO[T]) RemoveById(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	var entity T
	result := d.getDBFromContext(ctx).Where("id = ?", id).Delete(&entity)
	if result.Error != nil {
		d.logger.Error("msg", "物理删除数据失败", "id", id, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// 获取数据库字段名，考虑 GORM 标签
func getDBFieldName(field reflect.StructField) string {
	// 尝试从 GORM 标签获取字段名
	tag := field.Tag.Get("gorm")
	if tag != "" {
		// 解析 gorm 标签
		for _, option := range strings.Split(tag, ";") {
			if strings.HasPrefix(option, "column:") {
				return strings.TrimPrefix(option, "column:")
			}
		}
	}

	// 如果没有指定 column，则使用字段名转为蛇形命名
	return toSnakeCase(field.Name)
}

// shouldExcludeField 判断字段是否应该在 ExcludeNull 模式下被排除
// 对于 bool 类型，false 是有效值，不应该被排除
func shouldExcludeField(field reflect.Value) bool {
	// 对于 bool 类型，不排除任何值（包括 false）
	if field.Kind() == reflect.Bool {
		return false
	}

	// 对于指针类型，只有 nil 才排除
	if field.Kind() == reflect.Ptr {
		return field.IsNil()
	}

	// 对于其他类型，使用 IsZero 判断
	return field.IsZero()
}

// 将驼峰命名转换为蛇形命名
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
