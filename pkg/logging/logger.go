package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is our custom logger interface
type Logger interface {
	Debug(keyvals ...interface{}) error
	Info(keyvals ...interface{}) error
	Warn(keyvals ...interface{}) error
	Error(keyvals ...interface{}) error
	With(keyvals ...interface{}) Logger
}

type logger struct {
	logger log.Logger
	format string
}

// LogFormat 定义支持的日志格式类型
const (
	LogFormatJSON       = "json"
	LogFormatLogfmt     = "logfmt"
	LogFormatSpringBoot = "springboot"
)

// LogOption 定义日志选项
type LogOption func(*logOptions)

type logOptions struct {
	filePath   string // 日志文件路径
	maxSize    int    // 单个日志文件最大大小（MB）
	maxAge     int    // 日志文件保留天数
	maxBackups int    // 保留的旧日志文件数量
	useFile    bool   // 是否使用文件日志
}

// WithLogFile 设置日志输出到文件
func WithLogFile(filePath string, maxSize, maxAge, maxBackups int) LogOption {
	return func(o *logOptions) {
		o.filePath = filePath
		o.maxSize = maxSize
		o.maxAge = maxAge
		o.maxBackups = maxBackups
		o.useFile = true
	}
}

// 自定义调用者函数，获取实际的调用位置
func customCaller() interface{} {
	// 从第3帧开始查找，跳过日志库内部调用
	for i := 3; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 排除日志库内部文件
		if !strings.Contains(file, "go-kit/log") &&
			!strings.Contains(file, "pkg/logging") {
			// 只保留相对路径，使日志更简洁
			file = filepath.Base(file)
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return "unknown:0"
}

// NewLogger creates a new logger instance
func NewLogger(logLevel string, logFormat string, options ...LogOption) Logger {
	var logLogger log.Logger
	
	// 处理选项
	opts := &logOptions{}
	for _, option := range options {
		option(opts)
	}
	
	// 确定日志输出目标
	var writer io.Writer = os.Stdout // 默认输出到标准输出
	
	// 如果配置了文件日志，同时输出到文件和控制台
	if opts.useFile && opts.filePath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(opts.filePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
			// 如果创建目录失败，只输出到标准输出
		} else {
			// 使用 lumberjack 进行日志轮转
			fileWriter := &lumberjack.Logger{
				Filename:   opts.filePath,
				MaxSize:    opts.maxSize,    // 单个文件最大大小，单位是 MB
				MaxBackups: opts.maxBackups, // 保留的旧日志文件数量
				MaxAge:     opts.maxAge,     // 保留天数
				Compress:   true,           // 是否压缩
			}
			// 同时输出到文件和控制台
			writer = io.MultiWriter(os.Stdout, fileWriter)
		}
	}

	// 设置日志格式
	switch logFormat {
	case LogFormatJSON:
		logLogger = log.NewJSONLogger(log.NewSyncWriter(writer))
	case LogFormatSpringBoot:
		// 使用自定义的Spring Boot风格日志格式
		logLogger = NewSpringBootLogger(log.NewSyncWriter(writer))
	default:
		logLogger = log.NewLogfmtLogger(log.NewSyncWriter(writer))
	}

	// 使用自定义的调用者函数
	logLogger = log.With(logLogger, "ts", log.TimestampFormat(time.Now, "2006-01-02T15:04:05.000000Z07:00"), "caller", customCaller())

	// Set the log level filter
	switch logLevel {
	case "debug":
		logLogger = level.NewFilter(logLogger, level.AllowDebug())
	case "info":
		logLogger = level.NewFilter(logLogger, level.AllowInfo())
	case "warn":
		logLogger = level.NewFilter(logLogger, level.AllowWarn())
	case "error":
		logLogger = level.NewFilter(logLogger, level.AllowError())
	default:
		logLogger = level.NewFilter(logLogger, level.AllowInfo())
	}

	return &logger{
		logger: logLogger,
		format: logFormat,
	}
}

// Debug logs a debug message
func (l *logger) Debug(keyvals ...interface{}) error {
	// 每次调用时更新调用者信息
	return level.Debug(log.With(l.logger, "caller", customCaller())).Log(keyvals...)
}

// Info logs an info message
func (l *logger) Info(keyvals ...interface{}) error {
	// 每次调用时更新调用者信息
	return level.Info(log.With(l.logger, "caller", customCaller())).Log(keyvals...)
}

// Warn logs a warning message
func (l *logger) Warn(keyvals ...interface{}) error {
	// 每次调用时更新调用者信息
	return level.Warn(log.With(l.logger, "caller", customCaller())).Log(keyvals...)
}

// Error logs an error message
func (l *logger) Error(keyvals ...interface{}) error {
	// 每次调用时更新调用者信息
	return level.Error(log.With(l.logger, "caller", customCaller())).Log(keyvals...)
}

// With returns a new logger with the given keyvals added to each log message
func (l *logger) With(keyvals ...interface{}) Logger {
	return &logger{
		logger: log.With(l.logger, keyvals...),
		format: l.format,
	}
}

// SpringBootLogger 实现Spring Boot风格的日志格式
type SpringBootLogger struct {
	writer io.Writer
}

// NewSpringBootLogger 创建一个新的Spring Boot风格日志记录器
func NewSpringBootLogger(w io.Writer) log.Logger {
	return &SpringBootLogger{
		writer: w,
	}
}

// Log 实现Spring Boot风格的日志格式
func (l *SpringBootLogger) Log(keyvals ...interface{}) error {
	// 提取时间戳、级别、消息和其他字段
	var timestamp, levelStr, msg, caller string
	var module string
	otherKVs := make([]string, 0)

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 >= len(keyvals) {
			break
		}

		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}

		value := keyvals[i+1]

		switch key {
		case "ts":
			if timeVal, ok := value.(time.Time); ok {
				timestamp = timeVal.Format("2006-01-02 15:04:05.000")
			} else {
				timestamp = fmt.Sprintf("%v", value)
			}
		case "level":
			levelStr = strings.ToUpper(fmt.Sprintf("%v", value))
		case "msg":
			msg = fmt.Sprintf("%v", value)
		case "caller":
			caller = fmt.Sprintf("%v", value)
		case "module":
			module = fmt.Sprintf("%v", value)
		default:
			otherKVs = append(otherKVs, fmt.Sprintf("%s=%v", key, value))
		}
	}

	// 格式化为Spring Boot风格
	var logLine string
	if module != "" {
		logLine = fmt.Sprintf("%s %5s [%s] (%s) - %s", timestamp, levelStr, module, caller, msg)
	} else {
		logLine = fmt.Sprintf("%s %5s (%s) - %s", timestamp, levelStr, caller, msg)
	}

	// 添加其他键值对
	if len(otherKVs) > 0 {
		logLine += " " + strings.Join(otherKVs, " ")
	}

	// 输出日志
	_, err := fmt.Fprintln(l.writer, logLine)
	return err
}
