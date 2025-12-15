package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/utils"
)

// SQLiteConfig SQLite 配置
type SQLiteConfig struct {
	DBPath          string        `yaml:"dbPath" json:"dbPath"`                   // 数据库文件路径
	MaxIdleConns    int           `yaml:"maxIdleConns" json:"maxIdleConns"`       // 最大空闲连接数
	MaxOpenConns    int           `yaml:"maxOpenConns" json:"maxOpenConns"`       // 最大打开连接数
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime"` // 连接最大生命周期
	LogLevel        string        `yaml:"logLevel" json:"logLevel"`               // 日志级别
}

// NewSQLiteDB 创建一个新的 SQLite 数据库连接
func NewSQLiteDB(config SQLiteConfig, appLogger logging.Logger) (*gorm.DB, error) {
	// 使用公共方法解析相对路径
	dbPath := utils.ResolveRelativePath(config.DBPath)

	appLogger.Info("msg", "SQLite 数据库路径", "configPath", config.DBPath, "resolvedPath", dbPath)

	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if dbDir != "" && dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %w", err)
		}
	}

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

	// SQLite 连接参数
	// _journal_mode=WAL: 使用 WAL 模式提高并发性能
	// _busy_timeout=30000: 设置忙等待超时为 30 秒（本地上传等操作可能需要较长时间）
	// _foreign_keys=ON: 启用外键约束
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=30000&_foreign_keys=ON", dbPath)

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		// 禁用默认事务以提高性能（SQLite 单连接）
		// 重要：SQLite 是数据库级锁，开启默认事务会导致读写互相阻塞
		SkipDefaultTransaction: true,
		// 统一使用 UTC 时间存储，便于云端同步
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("连接 SQLite 数据库失败: %w", err)
	}

	// 获取底层 SQL DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 SQL DB 对象失败: %w", err)
	}

	// 设置连接池参数
	// SQLite 建议使用较小的连接池
	maxIdleConns := config.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 1
	}
	maxOpenConns := config.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 1 // SQLite 单文件，建议单连接
	}
	connMaxLifetime := config.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	appLogger.Info("msg", "SQLite 数据库连接成功", "path", dbPath)

	// 注册 JSON 序列化器
	dao.RegisterJSONSerializer(db)

	return db, nil
}
