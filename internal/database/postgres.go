package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
)

// PostgresConfig PostgreSQL 配置
type PostgresConfig struct {
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	User            string        `yaml:"user" json:"user"`
	Password        string        `yaml:"password" json:"password"`
	DBName          string        `yaml:"dbname" json:"dbname"`
	SSLMode         string        `yaml:"sslmode" json:"sslmode"`
	TimeZone        string        `yaml:"timezone" json:"timezone"`
	MaxIdleConns    int           `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns" json:"maxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime"`
	LogLevel        string        `yaml:"logLevel" json:"logLevel"`
}

// NewPostgresDB 创建一个新的 PostgreSQL 数据库连接
func NewPostgresDB(config PostgresConfig, appLogger logging.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
		config.TimeZone,
	)

	// 设置日志级别
	var logLevel logger.LogLevel
	switch config.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// 创建自定义日志记录器
	gormLogger := logger.New(
		&GormLogAdapter{logger: appLogger},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		// 禁用默认事务以提高性能，需要事务时使用 TransactionManager 显式开启
		SkipDefaultTransaction: true,
		// 统一使用 UTC 时间存储，便于云端同步和多时区支持
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("连接 PostgreSQL 数据库失败: %w", err)
	}

	// 获取底层 SQL DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 SQL DB 对象失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	appLogger.Info("msg", "PostgreSQL 数据库连接成功", "host", config.Host, "database", config.DBName)

	// 注册 JSON 序列化器
	dao.RegisterJSONSerializer(db)

	return db, nil
}

// GormLogAdapter 适配 GORM 日志到应用程序日志
type GormLogAdapter struct {
	logger logging.Logger
}

// Printf 实现 gorm logger.Writer 接口
func (l *GormLogAdapter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Debug("msg", msg, "source", "gorm")
}
