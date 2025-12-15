package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/spf13/viper"
	"github.com/yb2020/odoc/config"
	configAdapter "github.com/yb2020/odoc/config/adapter"
	"github.com/yb2020/odoc/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	// Import model packages
	docmodel "github.com/yb2020/odoc/services/doc/model"
	membershipmodel "github.com/yb2020/odoc/services/membership/model"
	"github.com/yb2020/odoc/services/nav/model"
	notemodel "github.com/yb2020/odoc/services/note/model"
	oauth2model "github.com/yb2020/odoc/services/oauth2/model"
	ossmodel "github.com/yb2020/odoc/services/oss/model"
	papermodel "github.com/yb2020/odoc/services/paper/model"
	paymodel "github.com/yb2020/odoc/services/pay/model"
	pdfmodel "github.com/yb2020/odoc/services/pdf/model"
	translatemodel "github.com/yb2020/odoc/services/translate/model"
	usermodel "github.com/yb2020/odoc/services/user/model"
)

// TableInfo 存储表结构信息
type TableInfo struct {
	Name    string
	Columns []ColumnInfo
}

// ColumnInfo 存储列信息
type ColumnInfo struct {
	Name         string
	Type         string
	Nullable     bool
	IsPrimaryKey bool
	DefaultValue string
	Comment      string
	IndexType    string // "none", "index", "uniqueIndex"
}

// ModelInfo 存储模型信息
type ModelInfo struct {
	Type      reflect.Type
	TableName string
	Package   string
}

// 当前数据库类型（在 main 中设置）
var currentDBType string

// 获取数据库中表的结构
func getDBTableInfo(db *gorm.DB, tableName string) (*TableInfo, error) {
	if currentDBType == "sqlite" {
		return getDBTableInfoSQLite(db, tableName)
	}
	return getDBTableInfoPostgres(db, tableName)
}

// 获取 SQLite 表结构
func getDBTableInfoSQLite(db *gorm.DB, tableName string) (*TableInfo, error) {
	var columns []struct {
		CID        int    `gorm:"column:cid"`
		Name       string `gorm:"column:name"`
		Type       string `gorm:"column:type"`
		NotNull    int    `gorm:"column:notnull"`
		DefaultVal string `gorm:"column:dflt_value"`
		PK         int    `gorm:"column:pk"`
	}

	result := db.Raw(fmt.Sprintf("PRAGMA table_info(%s)", tableName)).Scan(&columns)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("表 %s 不存在", tableName)
	}

	// 获取索引信息
	var indexes []struct {
		Name   string `gorm:"column:name"`
		Unique int    `gorm:"column:unique"`
	}
	db.Raw(fmt.Sprintf("PRAGMA index_list(%s)", tableName)).Scan(&indexes)

	// 获取每个索引的列
	indexedColumnsMap := make(map[string]string)
	for _, idx := range indexes {
		if idx.Name == "" {
			continue
		}
		var indexCols []struct {
			Name string `gorm:"column:name"`
		}
		db.Raw(fmt.Sprintf("PRAGMA index_info(%s)", idx.Name)).Scan(&indexCols)
		for _, col := range indexCols {
			if idx.Unique == 1 {
				indexedColumnsMap[col.Name] = "uniqueIndex"
			} else {
				indexedColumnsMap[col.Name] = "index"
			}
		}
	}

	tableInfo := &TableInfo{
		Name:    tableName,
		Columns: make([]ColumnInfo, 0, len(columns)),
	}

	for _, col := range columns {
		indexType := "none"
		if col.PK == 0 {
			if idxType, exists := indexedColumnsMap[col.Name]; exists {
				indexType = idxType
			}
		}

		tableInfo.Columns = append(tableInfo.Columns, ColumnInfo{
			Name:         col.Name,
			Type:         col.Type,
			Nullable:     col.NotNull == 0,
			IsPrimaryKey: col.PK > 0,
			DefaultValue: col.DefaultVal,
			Comment:      "",
			IndexType:    indexType,
		})
	}

	return tableInfo, nil
}

// 获取 PostgreSQL 表结构
func getDBTableInfoPostgres(db *gorm.DB, tableName string) (*TableInfo, error) {
	var columns []struct {
		ColumnName    string `gorm:"column:column_name"`
		DataType      string `gorm:"column:data_type"`
		IsNullable    string `gorm:"column:is_nullable"`
		ColumnDefault string `gorm:"column:column_default"`
		Comment       string `gorm:"column:column_comment"`
	}

	// 获取列信息
	result := db.Raw(`
		SELECT 
			column_name, 
			data_type, 
			is_nullable,
			column_default,
			'' as column_comment
		FROM 
			information_schema.columns 
		WHERE 
			table_name = ? 
		ORDER BY 
			ordinal_position
	`, tableName).Find(&columns)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("表 %s 不存在", tableName)
	}

	// 获取主键信息
	var primaryKeys []string
	db.Raw(`
		SELECT 
			c.column_name 
		FROM 
			information_schema.table_constraints tc 
		JOIN 
			information_schema.constraint_column_usage AS ccu USING (constraint_schema, constraint_name) 
		JOIN 
			information_schema.columns AS c ON c.table_schema = tc.constraint_schema AND tc.table_name = c.table_name AND ccu.column_name = c.column_name
		WHERE 
			tc.constraint_type = 'PRIMARY KEY' AND tc.table_name = ?
	`, tableName).Pluck("column_name", &primaryKeys)

	// 获取索引信息
	var indexedColumns []struct {
		ColumnName string `gorm:"column:column_name"`
		IndexName  string `gorm:"column:indexname"`
		IsUnique   bool   `gorm:"column:is_unique"`
	}
	db.Raw(`
		SELECT 
			a.attname as column_name,
			i.relname as indexname,
			ix.indisunique as is_unique
		FROM pg_class t, pg_class i, pg_index ix, pg_attribute a
		WHERE t.oid = ix.indrelid AND i.oid = ix.indexrelid
		AND a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		AND t.relkind = 'r' AND t.relname = ?
	`, tableName).Scan(&indexedColumns)

	// 创建索引列映射
	indexedColumnsMap := make(map[string]string) // 列名 -> 索引类型
	for _, idx := range indexedColumns {
		if idx.IsUnique {
			indexedColumnsMap[idx.ColumnName] = "uniqueIndex"
		} else {
			indexedColumnsMap[idx.ColumnName] = "index"
		}
	}

	// 构建表信息
	tableInfo := &TableInfo{
		Name:    tableName,
		Columns: make([]ColumnInfo, 0, len(columns)),
	}

	for _, col := range columns {
		isPK := false
		for _, pk := range primaryKeys {
			if col.ColumnName == pk {
				isPK = true
				break
			}
		}

		// 确定索引类型
		indexType := "none"
		if isPK {
			// 主键默认有索引，但我们不把它当作普通索引
			indexType = "none"
		} else if idxType, exists := indexedColumnsMap[col.ColumnName]; exists {
			indexType = idxType
		}

		tableInfo.Columns = append(tableInfo.Columns, ColumnInfo{
			Name:         col.ColumnName,
			Type:         col.DataType,
			Nullable:     col.IsNullable == "YES",
			IsPrimaryKey: isPK,
			DefaultValue: col.ColumnDefault,
			Comment:      col.Comment,
			IndexType:    indexType,
		})
	}

	return tableInfo, nil
}

// 获取模型结构体定义的表结构
func getModelTableInfo(modelType reflect.Type) *TableInfo {
	// 创建一个模型实例
	modelValue := reflect.New(modelType).Interface()

	// 获取表名
	var tableName string
	if tabler, ok := modelValue.(schema.Tabler); ok {
		tableName = tabler.TableName()
	} else {
		tableName = strings.ToLower(modelType.Name()) + "s"
	}

	// 构建表信息
	tableInfo := &TableInfo{
		Name:    tableName,
		Columns: []ColumnInfo{},
	}

	// 递归处理结构体字段，包括嵌入字段
	processStructFields(modelType, "", &tableInfo.Columns)

	return tableInfo
}

// 递归处理结构体字段，包括嵌入字段
func processStructFields(structType reflect.Type, prefix string, columns *[]ColumnInfo) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// 跳过非导出字段
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		// 处理嵌入字段
		if field.Anonymous {
			// 如果是嵌入字段，递归处理其字段
			if field.Type.Kind() == reflect.Struct {
				processStructFields(field.Type, prefix, columns)
			}
			continue
		}

		// 获取 gorm 标签
		tag := field.Tag.Get("gorm")
		if tag == "-" {
			// 跳过标记为 "-" 的字段
			continue
		}

		// 如果没有 gorm 标签，跳过
		if tag == "" {
			continue
		}

		columnName := toSnakeCase(field.Name) // 默认转为蛇形命名
		isPK := false
		nullable := true
		indexType := "none"  // 默认无索引
		hasValidTag := false // 是否有有效的 gorm 标签

		// 解析标签
		tagParts := strings.Split(tag, ";")

		for _, part := range tagParts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// 检查是否有有效的 gorm 标签
			if strings.HasPrefix(part, "column:") ||
				part == "primaryKey" ||
				part == "not null" ||
				part == "index" ||
				strings.HasPrefix(part, "index:") ||
				part == "uniqueIndex" ||
				strings.HasPrefix(part, "uniqueIndex:") ||
				strings.HasPrefix(part, "type:") ||
				strings.HasPrefix(part, "size:") ||
				strings.HasPrefix(part, "default:") {
				hasValidTag = true
			}

			if strings.HasPrefix(part, "column:") {
				// 如果有明确的column标签，使用标签中指定的列名
				columnName = strings.TrimPrefix(part, "column:")
			} else if part == "primaryKey" {
				isPK = true
				// 主键默认有索引，但我们不把它当作普通索引
			} else if part == "not null" {
				nullable = false
			} else if part == "index" || strings.HasPrefix(part, "index:") {
				indexType = "index"
			} else if part == "uniqueIndex" || strings.HasPrefix(part, "uniqueIndex:") {
				indexType = "uniqueIndex"
			}
		}

		// 如果没有有效的 gorm 标签，跳过这个字段
		if !hasValidTag {
			continue
		}

		// 获取字段类型
		fieldType := field.Type.String()
		if field.Type.Kind() == reflect.Ptr {
			fieldType = field.Type.Elem().String()
			nullable = true
		}

		// 添加列信息
		*columns = append(*columns, ColumnInfo{
			Name:         columnName,
			Type:         fieldType,
			Nullable:     nullable,
			IsPrimaryKey: isPK,
			DefaultValue: "",
			Comment:      "",
			IndexType:    indexType,
		})
	}
}

// 更新表结构以匹配模型定义
func updateTable(db *gorm.DB, modelType reflect.Type) error {
	// 创建一个模型实例
	modelInstance := reflect.New(modelType).Interface()

	// 自动迁移
	if err := db.AutoMigrate(modelInstance); err != nil {
		return err
	}

	// 获取模型表信息
	modelTable := getModelTableInfo(modelType)

	// 获取数据库表信息
	dbTable, err := getDBTableInfo(db, modelTable.Name)
	if err != nil {
		// 如果表不存在，AutoMigrate 已经创建了表，不需要额外处理
		return nil
	}

	// 创建列映射
	dbColumns := make(map[string]ColumnInfo)
	for _, col := range dbTable.Columns {
		dbColumns[col.Name] = col
	}

	modelColumns := make(map[string]ColumnInfo)
	for _, col := range modelTable.Columns {
		modelColumns[col.Name] = col
	}

	// 检查并添加缺失的索引
	for modelColName, modelCol := range modelColumns {
		if modelCol.IndexType != "none" {
			// 查找对应的数据库列
			for dbColName, dbCol := range dbColumns {
				if normalizeColumnName(dbColName, modelColName) {
					// 如果模型中有索引但数据库中没有，或者索引类型不同，创建索引
					if dbCol.IndexType != modelCol.IndexType {
						var indexSQL string
						indexName := fmt.Sprintf("idx_%s_%s", modelTable.Name, dbColName)

						// 根据索引类型创建不同的索引
						if modelCol.IndexType == "uniqueIndex" {
							indexSQL = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s)",
								indexName, modelTable.Name, dbColName)
						} else {
							indexSQL = fmt.Sprintf("CREATE INDEX %s ON %s (%s)",
								indexName, modelTable.Name, dbColName)
						}

						err := db.Exec(indexSQL).Error
						if err != nil {
							return fmt.Errorf("创建索引失败: %v", err)
						}
						// 使用if-else替代三元运算符
						indexTypeStr := "索引"
						if modelCol.IndexType == "uniqueIndex" {
							indexTypeStr = "唯一索引"
						}

						fmt.Printf("已创建%s: %s 在表 %s 的列 %s 上\n",
							indexTypeStr, indexName, modelTable.Name, dbColName)
					}
					break
				}
			}
		}
	}

	// 在 updateTable 函数中，在处理完索引后添加
	// 检查并删除未使用的列
	for dbColName := range dbColumns {
		found := false
		for modelColName := range modelColumns {
			if normalizeColumnName(dbColName, modelColName) {
				found = true
				break
			}
		}

		// 如果数据库中的列不在模型中，询问是否删除
		if !found && !isSystemColumn(dbColName) {
			if askForConfirmation(fmt.Sprintf("是否删除未使用的列 %s.%s?", modelTable.Name, dbColName)) {
				dropColumnSQL := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", modelTable.Name, dbColName)
				err := db.Exec(dropColumnSQL).Error
				if err != nil {
					return fmt.Errorf("删除列失败: %v", err)
				}
				fmt.Printf("已删除未使用的列: %s.%s\n", modelTable.Name, dbColName)
			}
		}
	}

	return nil
}

// 比较数据库表和模型定义的差异
func compareTableInfo(dbTable, modelTable *TableInfo) []string {
	differences := []string{}

	// 检查表名
	if dbTable.Name != modelTable.Name {
		differences = append(differences, fmt.Sprintf("表名不同: 数据库=%s, 模型=%s", dbTable.Name, modelTable.Name))
	}

	// 创建列映射
	dbColumns := make(map[string]ColumnInfo)
	for _, col := range dbTable.Columns {
		dbColumns[col.Name] = col
	}

	modelColumns := make(map[string]ColumnInfo)
	for _, col := range modelTable.Columns {
		modelColumns[col.Name] = col
	}

	// 检查模型中定义但数据库中不存在的列
	for modelColName := range modelColumns {
		found := false
		for dbColName := range dbColumns {
			if normalizeColumnName(dbColName, modelColName) {
				found = true
				break
			}
		}
		if !found {
			differences = append(differences, fmt.Sprintf("列 %s 在模型中定义但在数据库中不存在", modelColName))
		}
	}

	// 检查数据库中存在但模型中未定义的列
	for dbColName := range dbColumns {
		found := false
		for modelColName := range modelColumns {
			if normalizeColumnName(dbColName, modelColName) {
				found = true
				break
			}
		}
		if !found {
			differences = append(differences, fmt.Sprintf("列 %s 在数据库中存在但在模型中未定义", dbColName))
		}
	}

	// 检查列属性差异
	for modelColName, modelCol := range modelColumns {
		for dbColName, dbCol := range dbColumns {
			if normalizeColumnName(dbColName, modelColName) {
				// 检查类型
				if !isTypeCompatible(dbCol.Type, modelCol.Type) {
					differences = append(differences, fmt.Sprintf("列 %s 类型不兼容: 数据库=%s, 模型=%s", dbColName, dbCol.Type, modelCol.Type))
				}

				// 检查可空性
				if dbCol.Nullable != modelCol.Nullable && !dbCol.IsPrimaryKey && !modelCol.IsPrimaryKey {
					differences = append(differences, fmt.Sprintf("列 %s 可空性不同: 数据库=%v, 模型=%v", dbColName, dbCol.Nullable, modelCol.Nullable))
				}

				// 检查主键
				if dbCol.IsPrimaryKey != modelCol.IsPrimaryKey {
					differences = append(differences, fmt.Sprintf("列 %s 主键属性不同: 数据库=%v, 模型=%v", dbColName, dbCol.IsPrimaryKey, modelCol.IsPrimaryKey))
				}

				// 检查索引
				if dbCol.IndexType != modelCol.IndexType {
					differences = append(differences, fmt.Sprintf("列 %s 索引属性不同: 数据库=%s, 模型=%s", dbColName, dbCol.IndexType, modelCol.IndexType))
				}

				break
			}
		}
	}

	return differences
}

// 检查是否为系统列（如id, created_at, updated_at等）
func isSystemColumn(columnName string) bool {
	systemColumns := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, name := range systemColumns {
		if name == columnName {
			return true
		}
	}
	return false
}

// 将驼峰式命名转换为蛇形命名
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(c))
		} else {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// 将蛇形命名转换为驼峰式命名
func toCamelCase(s string) string {
	var result strings.Builder
	upperNext := false

	for _, c := range s {
		if c == '_' {
			upperNext = true
		} else if upperNext {
			result.WriteRune(unicode.ToUpper(c))
			upperNext = false
		} else {
			result.WriteRune(c)
		}
	}

	return result.String()
}

// 规范化列名，处理命名风格差异
func normalizeColumnName(dbName, modelName string) bool {
	// 直接比较
	if dbName == modelName {
		return true
	}

	// 转换为小写后比较
	if strings.ToLower(dbName) == strings.ToLower(modelName) {
		return true
	}

	// 将数据库名转为驼峰后比较
	if toCamelCase(dbName) == modelName {
		return true
	}

	// 将模型名转为蛇形后比较
	if dbName == toSnakeCase(modelName) {
		return true
	}

	// 将两者都转为小写后比较
	if strings.ToLower(dbName) == strings.ToLower(toSnakeCase(modelName)) {
		return true
	}

	// 将两者都转为小写后比较（模型名转为蛇形）
	if strings.ToLower(dbName) == strings.ToLower(modelName) {
		return true
	}

	return false
}

// 简单检查类型是否兼容
func isTypeCompatible(dbType, modelType string) bool {
	// 转换为小写进行比较
	dbType = strings.ToLower(dbType)
	modelType = strings.ToLower(modelType)

	// 直接匹配
	if dbType == modelType {
		return true
	}

	// 处理自定义枚举类型
	if (strings.Contains(modelType, "userstatus") || strings.Contains(modelType, "user.userstatus") ||
		strings.Contains(modelType, "pb.userstatus")) && dbType == "integer" {
		return true
	}

	if (strings.Contains(modelType, "userroles") || strings.Contains(modelType, "model.userroles")) &&
		(dbType == "json" || dbType == "jsonb") {
		return true
	}

	// 特殊类型映射
	dbToGoTypeMap := map[string][]string{
		"int64":     {"bigint", "int8"},
		"int32":     {"integer", "int", "int4"},
		"int16":     {"smallint", "int2"},
		"int8":      {"smallint", "int2"},
		"int":       {"integer", "int4"},
		"uint64":    {"bigint", "int8"},
		"uint32":    {"integer", "int4"},
		"uint16":    {"smallint", "int2"},
		"uint8":     {"smallint", "int2"},
		"uint":      {"integer", "int4"},
		"float64":   {"double precision", "float8"},
		"float32":   {"real", "float4"},
		"bool":      {"boolean", "bool"},
		"string":    {"character varying", "varchar", "text", "character", "char"},
		"time.time": {"timestamp", "timestamptz", "timestamp with time zone", "timestamp without time zone", "date"},
		"[]byte":    {"bytea"},
	}

	// 检查类型兼容性
	if compatibleTypes, ok := dbToGoTypeMap[modelType]; ok {
		for _, t := range compatibleTypes {
			if strings.Contains(dbType, t) {
				return true
			}
		}
	}

	return false
}

// 全局变量，用于控制是否自动确认所有操作
var globalAutoApprove bool
var globalCreateAllTables bool

// 用户确认函数
func askForConfirmation(prompt string) bool {
	// 只有在设置了自动确认时，才自动确认所有操作
	if globalAutoApprove {
		fmt.Printf("%s [自动确认: y]\n", prompt)
		return true
	}

	// 检查是否在非交互环境中运行
	if isNonInteractive() {
		fmt.Printf("%s [非交互环境，默认: n]\n", prompt)
		return false
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入失败，默认为否")
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		} else {
			fmt.Println("请输入 y 或 n")
		}
	}
}

// 检查是否在非交互环境中运行（如容器、CI/CD）
func isNonInteractive() bool {
	// 检查常见的环境变量
	if os.Getenv("CI") != "" || os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	// 检查标准输入是否是TTY
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		// 如果无法确定，假设是非交互式
		return true
	}

	// 检查是否是字符设备（终端）
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

// 扫描所有模型并返回模型信息列表
func scanAllModels() []ModelInfo {
	var models []ModelInfo

	// 添加 User 模型
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(usermodel.User{}),
		TableName: usermodel.User{}.TableName(),
		Package:   "user",
	})

	// 添加 OAuth2 模型
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(oauth2model.OAuth2Clients{}),
		TableName: oauth2model.OAuth2Clients{}.TableName(),
		Package:   "oauth2",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(oauth2model.OAuth2AuthCode{}),
		TableName: oauth2model.OAuth2AuthCode{}.TableName(),
		Package:   "oauth2",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(oauth2model.OAuth2Token{}),
		TableName: oauth2model.OAuth2Token{}.TableName(),
		Package:   "oauth2",
	})

	// 添加 Translate 模型
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(translatemodel.Glossary{}),
		TableName: translatemodel.Glossary{}.TableName(),
		Package:   "translate",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(translatemodel.OCRTranslateLog{}),
		TableName: translatemodel.OCRTranslateLog{}.TableName(),
		Package:   "translate",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(translatemodel.TextTranslateLog{}),
		TableName: translatemodel.TextTranslateLog{}.TableName(),
		Package:   "translate",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(translatemodel.WordPronunciation{}),
		TableName: translatemodel.WordPronunciation{}.TableName(),
		Package:   "translate",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(translatemodel.FullTextTranslate{}),
		TableName: translatemodel.FullTextTranslate{}.TableName(),
		Package:   "translate",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(translatemodel.FullTextTranslateFix{}),
		TableName: translatemodel.FullTextTranslateFix{}.TableName(),
		Package:   "translate",
	})

	// ----- Note 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.PaperNote{}),
		TableName: notemodel.PaperNote{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.PaperNoteAccess{}),
		TableName: notemodel.PaperNoteAccess{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteDrawEntity{}),
		TableName: notemodel.NoteDrawEntity{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteShape{}),
		TableName: notemodel.NoteShape{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteSummary{}),
		TableName: notemodel.NoteSummary{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteWord{}),
		TableName: notemodel.NoteWord{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteWordConfig{}),
		TableName: notemodel.NoteWordConfig{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteLatestRead{}),
		TableName: notemodel.NoteLatestRead{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteReadLocation{}),
		TableName: notemodel.NoteReadLocation{}.TableName(),
		Package:   "note",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(notemodel.NoteExportHistory{}),
		TableName: notemodel.NoteExportHistory{}.TableName(),
		Package:   "note",
	})
	// ----- Note 模块---//

	// ----- PDF 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PaperPdf{}),
		TableName: pdfmodel.PaperPdf{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PaperPdfSelectRecord{}),
		TableName: pdfmodel.PaperPdfSelectRecord{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfAnnotation{}),
		TableName: pdfmodel.PdfAnnotation{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfComment{}),
		TableName: pdfmodel.PdfComment{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfMark{}),
		TableName: pdfmodel.PdfMark{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfMarkBackup{}),
		TableName: pdfmodel.PdfMarkBackup{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfMarkTag{}),
		TableName: pdfmodel.PdfMarkTag{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfMarkTagRelation{}),
		TableName: pdfmodel.PdfMarkTagRelation{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfReaderSetting{}),
		TableName: pdfmodel.PdfReaderSetting{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfThumb{}),
		TableName: pdfmodel.PdfThumb{}.TableName(),
		Package:   "pdf",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(pdfmodel.PdfSummary{}),
		TableName: pdfmodel.PdfSummary{}.TableName(),
		Package:   "pdf",
	})
	// ----- PDF 模块---//

	// ----- Paper 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.Paper{}),
		TableName: papermodel.Paper{}.TableName(),
		Package:   "paper",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperAccess{}),
		TableName: papermodel.PaperAccess{}.TableName(),
		Package:   "paper",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperAttachment{}),
		TableName: papermodel.PaperAttachment{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperAnswer{}),
		TableName: papermodel.PaperAnswer{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperComment{}),
		TableName: papermodel.PaperComment{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperCommentApproval{}),
		TableName: papermodel.PaperCommentApproval{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperQuestion{}),
		TableName: papermodel.PaperQuestion{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperResources{}),
		TableName: papermodel.PaperResources{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperPdfParsed{}),
		TableName: papermodel.PaperPdfParsed{}.TableName(),
		Package:   "paper",
	})

	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(papermodel.PaperJcrEntity{}),
		TableName: papermodel.PaperJcrEntity{}.TableName(),
		Package:   "paper",
	})

	// ----- Paper 模块---//

	// ----- Doc 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.UserDoc{}),
		TableName: docmodel.UserDoc{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.UserDocClassify{}),
		TableName: docmodel.UserDocClassify{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.UserDocFolder{}),
		TableName: docmodel.UserDocFolder{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.UserDocFolderRelation{}),
		TableName: docmodel.UserDocFolderRelation{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.UserDocAttachment{}),
		TableName: docmodel.UserDocAttachment{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.DocClassifyRelation{}),
		TableName: docmodel.DocClassifyRelation{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.Csl{}),
		TableName: docmodel.Csl{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.UserCslRelation{}),
		TableName: docmodel.UserCslRelation{}.TableName(),
		Package:   "doc",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(docmodel.DoiMetaInfo{}),
		TableName: docmodel.DoiMetaInfo{}.TableName(),
		Package:   "doc",
	})
	// ----- Doc 模块---//

	// ----- Membership 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(membershipmodel.UserMembership{}),
		TableName: membershipmodel.UserMembership{}.TableName(),
		Package:   "membership",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(membershipmodel.Credit{}),
		TableName: membershipmodel.Credit{}.TableName(),
		Package:   "membership",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(membershipmodel.CreditBill{}),
		TableName: membershipmodel.CreditBill{}.TableName(),
		Package:   "membership",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(membershipmodel.Order{}),
		TableName: membershipmodel.Order{}.TableName(),
		Package:   "membership",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(membershipmodel.CreditPaymentRecord{}),
		TableName: membershipmodel.CreditPaymentRecord{}.TableName(),
		Package:   "membership",
	})
	// ----- Membership 模块---//

	// ----- Pay 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(paymodel.PaymentRecord{}),
		TableName: paymodel.PaymentRecord{}.TableName(),
		Package:   "pay",
	})
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(paymodel.PaymentSubscription{}),
		TableName: paymodel.PaymentSubscription{}.TableName(),
		Package:   "pay",
	})
	// ----- Pay 模块---//
	// ----- Nav 模块---//
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(model.Website{}),
		TableName: model.Website{}.TableName(),
		Package:   "nav",
	})
	// ----- Nav 模块---//

	// oss模块 //
	models = append(models, ModelInfo{
		Type:      reflect.TypeOf(ossmodel.OssRecord{}),
		TableName: ossmodel.OssRecord{}.TableName(),
		Package:   "oss",
	})

	return models
}

// 检查表是否存在（支持 PostgreSQL 和 SQLite）
func tableExists(db *gorm.DB, tableName string) (bool, error) {
	// 使用 GORM 的 Migrator 来检查表是否存在，这是跨数据库兼容的方式
	return db.Migrator().HasTable(tableName), nil
}

func main() {
	// 解析命令行参数
	var configPath string
	var specificModel string
	var listOnly bool
	var generateCode bool
	var autoApprove bool
	var createAllTables bool
	var env string

	// 从环境变量获取环境名称
	envFromEnv := os.Getenv("APP_ENV")
	if envFromEnv == "" {
		// 默认使用develop环境
		envFromEnv = "develop"
	}

	flag.StringVar(&configPath, "config", "", "配置文件路径，不指定则根据环境变量或-env参数生成")
	flag.StringVar(&env, "env", envFromEnv, "环境名称 (develop/staging/release)，用于确定配置文件路径")
	flag.StringVar(&specificModel, "model", "", "指定要处理的模型名称，不指定则处理所有模型")
	flag.BoolVar(&listOnly, "list", false, "仅列出所有可用模型，不执行数据库操作")
	flag.BoolVar(&generateCode, "gen", false, "是否生成模型代码")
	flag.BoolVar(&autoApprove, "y", false, "自动确认所有操作，适用于CI/CD环境")
	flag.BoolVar(&createAllTables, "all", false, "一键创建或更新所有数据库表，自动确认所有操作")
	flag.Parse()

	// 设置全局自动确认变量
	globalAutoApprove = autoApprove
	globalCreateAllTables = createAllTables

	// 如果没有指定配置文件路径，则根据环境生成
	if configPath == "" {
		configPath = utils.GetConfigPath(env)
		fmt.Printf("使用配置文件: %s\n", configPath)
	}

	// 加载配置，并支持环境变量覆盖
	var cfg *config.Config
	var err error

	// 使用 LoadConfigWithViper 而不是 LoadConfig，以支持环境变量覆盖
	cfg, err = configAdapter.LoadConfigWithViper(configPath)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v，尝试使用默认配置\n", err)
		cfg = config.GetConfig()

		// 即使使用默认配置，也尝试应用环境变量
		v := viper.New()
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// 处理数据库相关的环境变量（PostgreSQL）
		if host := v.GetString("DATABASE_HOST"); host != "" {
			cfg.Database.Postgres.Host = host
		}
		if port := v.GetInt("DATABASE_PORT"); port != 0 {
			cfg.Database.Postgres.Port = port
		}
		if user := v.GetString("DATABASE_USER"); user != "" {
			cfg.Database.Postgres.User = user
		}
		if password := v.GetString("DATABASE_PASSWORD"); password != "" {
			cfg.Database.Postgres.Password = password
		}
		if dbName := v.GetString("DATABASE_DBNAME"); dbName != "" {
			cfg.Database.Postgres.DBName = dbName
		}
		if sslMode := v.GetString("DATABASE_SSLMODE"); sslMode != "" {
			cfg.Database.Postgres.SSLMode = sslMode
		}
	}

	// 更新全局配置，确保其他地方通过 config.GetConfig() 获取的是最新配置
	config.SetConfig(cfg)

	// 扫描所有模型
	allModels := scanAllModels()

	// 如果只是列出模型，则打印后退出
	if listOnly {
		fmt.Println("可用模型列表:")
		for i, model := range allModels {
			modelName := model.Type.Name()
			fmt.Printf("%d. %s (表名: %s, 包: %s)\n", i+1, modelName, model.TableName, model.Package)
		}
		return
	}

	// 根据配置的数据库类型连接数据库
	var db *gorm.DB
	dbType := cfg.Database.Type
	currentDBType = dbType // 设置全局变量供其他函数使用
	fmt.Printf("数据库类型: %s\n", dbType)

	switch dbType {
	case "sqlite":
		sq := cfg.Database.SQLite
		dbPath := sq.DBPath
		// 确保目录存在
		if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				panic(fmt.Errorf("创建数据库目录失败: %w", err))
			}
		}
		fmt.Printf("SQLite 数据库路径: %s\n", dbPath)
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	case "postgres":
		pg := cfg.Database.Postgres
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			pg.Host,
			pg.Port,
			pg.User,
			pg.Password,
			pg.DBName,
			pg.SSLMode,
			pg.TimeZone,
		)
		fmt.Printf("PostgreSQL 连接: %s:%d/%s\n", pg.Host, pg.Port, pg.DBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		panic(fmt.Errorf("不支持的数据库类型: %s", dbType))
	}
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	// 过滤模型（如果指定了特定模型或-all参数）
	var modelsToProcess []ModelInfo
	if specificModel != "" {
		// 处理指定的单个模型
		for _, model := range allModels {
			if strings.EqualFold(model.Type.Name(), specificModel) {
				modelsToProcess = append(modelsToProcess, model)
				break
			}
		}
		if len(modelsToProcess) == 0 {
			panic(fmt.Errorf("未找到指定的模型: %s", specificModel))
		}
	} else if createAllTables {
		// 只有明确指定-all参数时才处理所有模型
		modelsToProcess = allModels
		fmt.Println("使用-all参数，将处理所有模型")
	} else {
		// 既没有指定模型名称，也没有使用-all参数，显示模型列表并退出
		fmt.Println("请指定要处理的模型名称(-model 参数)或使用-all参数处理所有模型")
		fmt.Println("\n可用模型列表:")
		for i, model := range allModels {
			modelName := model.Type.Name()
			fmt.Printf("%d. %s (表名: %s, 包: %s)\n", i+1, modelName, model.TableName, model.Package)
		}
		return
	}

	// 处理每个模型
	for _, model := range modelsToProcess {
		fmt.Printf("\n处理模型: %s (表名: %s)\n", model.Type.Name(), model.TableName)

		// 检查表是否存在
		exists, err := tableExists(db, model.TableName)
		if err != nil {
			fmt.Printf("检查表 %s 是否存在时出错: %v\n", model.TableName, err)
			continue
		}

		if exists {
			fmt.Printf("表 %s 已存在，正在检查结构差异...\n", model.TableName)

			// 获取数据库中的表结构
			dbTableInfo, err := getDBTableInfo(db, model.TableName)
			if err != nil {
				fmt.Printf("获取数据库表结构失败: %v\n", err)
				continue
			}

			// 获取模型定义的表结构
			modelTableInfo := getModelTableInfo(model.Type)

			// 比较差异
			differences := compareTableInfo(dbTableInfo, modelTableInfo)

			if len(differences) > 0 {
				fmt.Println("发现以下差异:")
				for _, diff := range differences {
					fmt.Printf("  - %s\n", diff)
				}

				// 询问用户是否要更新表结构
				if askForConfirmation("是否要更新表结构以匹配模型定义?") {
					fmt.Println("正在更新表结构...")
					if err := updateTable(db, model.Type); err != nil {
						fmt.Printf("更新表结构失败: %v\n", err)
					} else {
						fmt.Println("表结构更新成功")
					}
				} else {
					fmt.Println("保持表结构不变")
				}
			} else {
				fmt.Println("表结构与模型定义一致，无需更改")
			}
		} else {
			fmt.Printf("表 %s 不存在\n", model.TableName)
			if askForConfirmation("是否创建表?") {
				fmt.Println("正在创建表...")
				if err := updateTable(db, model.Type); err != nil {
					fmt.Printf("创建表失败: %v\n", err)
				} else {
					fmt.Println("表创建成功")
				}
			} else {
				fmt.Println("跳过创建表")
			}
		}
	}

	// 设置全局自动确认标志
	globalAutoApprove = autoApprove
	// 设置全局一键创建所有表标志
	globalCreateAllTables = createAllTables

	// 询问用户是否要生成代码
	if generateCode && askForConfirmation("是否要生成模型代码?") {
		// 按包分组模型
		modelsByPackage := make(map[string][]ModelInfo)
		for _, model := range modelsToProcess {
			modelsByPackage[model.Package] = append(modelsByPackage[model.Package], model)
		}

		// 为每个包生成模型代码
		for pkg, models := range modelsByPackage {
			// 初始化生成器
			outPath := fmt.Sprintf("./services/%s/model", pkg)
			fmt.Printf("为 %s 包生成模型代码，输出路径: %s\n", pkg, outPath)

			g := gen.NewGenerator(gen.Config{
				// 输出路径 - 根据包名动态设置
				OutPath: outPath,
				// 只生成模型定义，不生成查询代码
				Mode: gen.WithoutContext,
				// 只生成模型结构体
				ModelPkgPath: "model",
			})

			// 使用数据库连接
			g.UseDB(db)

			// 收集该包中所有模型的表名
			var tableNames []string
			for _, model := range models {
				tableNames = append(tableNames, model.TableName)
			}

			// 为每个表生成模型
			for _, tableName := range tableNames {
				g.ApplyBasic(g.GenerateModel(tableName))
			}

			// 执行生成
			g.Execute()
			fmt.Printf("%s 包的模型代码生成完成\n", pkg)
		}

		fmt.Println("所有模型代码生成完成")
		fmt.Println("注意: 只生成了模型定义，请使用手动编写的DAO层代码进行数据库操作")
	} else {
		fmt.Println("取消代码生成")
	}
}
