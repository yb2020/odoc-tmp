package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
)

// DBType 数据库类型
type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeSQLite   DBType = "sqlite"
	DBTypeMySQL    DBType = "mysql"
)

// NewDB 根据配置创建数据库连接
// 通过 config.Database.Type 自动选择数据库类型
func NewDB(cfg *config.Config, logger logging.Logger) (*gorm.DB, error) {
	if !cfg.Database.Enabled {
		return nil, fmt.Errorf("数据库未启用")
	}

	dbType := DBType(cfg.Database.Type)

	switch dbType {
	case DBTypePostgres:
		pg := cfg.Database.Postgres
		return NewPostgresDB(PostgresConfig{
			Host:            pg.Host,
			Port:            pg.Port,
			User:            pg.User,
			Password:        pg.Password,
			DBName:          pg.DBName,
			SSLMode:         pg.SSLMode,
			TimeZone:        pg.TimeZone,
			MaxIdleConns:    pg.MaxIdleConns,
			MaxOpenConns:    pg.MaxOpenConns,
			ConnMaxLifetime: time.Duration(pg.ConnMaxLifetime) * time.Second,
			LogLevel:        pg.LogLevel,
		}, logger)

	case DBTypeSQLite:
		sq := cfg.Database.SQLite
		return NewSQLiteDB(SQLiteConfig{
			DBPath:          sq.DBPath,
			MaxIdleConns:    sq.MaxIdleConns,
			MaxOpenConns:    sq.MaxOpenConns,
			ConnMaxLifetime: time.Duration(sq.ConnMaxLifetime) * time.Second,
			LogLevel:        sq.LogLevel,
		}, logger)

	case DBTypeMySQL:
		return nil, fmt.Errorf("MySQL 暂未实现")

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", dbType)
	}
}

// GetDBType 从配置获取数据库类型
func GetDBType(cfg *config.Config) DBType {
	return DBType(cfg.Database.Type)
}
