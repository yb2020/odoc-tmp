package idgen

import (
	"context"
	"strconv"
	"sync"
)

// IDGenerator 定义ID生成器接口
type IDGenerator interface {
	// GenerateID 生成一个唯一ID
	GenerateID(ctx context.Context) (int64, error)
	// GenerateStringID 生成一个字符串格式的唯一ID
	// 这对于前端处理大数字时避免精度丢失很有用
	GenerateStringID(ctx context.Context) (string, error)
}

// SnowflakeGenerator 雪花算法ID生成器的实现
type SnowflakeGenerator struct {
	generator *SnowflakeIDGenerator
}

var (
	snowflakeInstance *SnowflakeGenerator
	once              sync.Once
)

// NewSnowflakeGenerator 创建一个新的雪花算法ID生成器
func NewSnowflakeGenerator() (*SnowflakeGenerator, error) {
	var err error
	once.Do(func() {
		var generator *SnowflakeIDGenerator
		generator, err = NewSnowflakeIDGenerator()
		if err != nil {
			return
		}
		snowflakeInstance = &SnowflakeGenerator{
			generator: generator,
		}
	})

	if err != nil {
		return nil, err
	}

	return snowflakeInstance, nil
}

// GetSnowflakeGenerator 获取雪花算法ID生成器的单例实例
func GetSnowflakeGenerator() (*SnowflakeGenerator, error) {
	if snowflakeInstance == nil {
		return NewSnowflakeGenerator()
	}
	return snowflakeInstance, nil
}

// GenerateId 实现IDGenerator接口，生成一个唯一ID
func (s *SnowflakeGenerator) GenerateId(ctx context.Context) (int64, error) {
	return s.generator.NextID()
}

// GenerateStringId 实现IDGenerator接口，生成一个字符串格式的唯一ID
func (s *SnowflakeGenerator) GenerateStringId(ctx context.Context) (string, error) {
	id, err := s.generator.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// GenerateSnowflakeID 获取一个唯一ID，如果出错则panic
// 这是一个便捷方法，用于简化ID生成代码
// 注意：此函数仅在确定不会出错的情况下使用，否则请使用标准的错误处理方式
func GenerateSnowflakeID(ctx context.Context) int64 {
	generator, err := GetSnowflakeGenerator()
	if err != nil {
		panic("failed to get snowflake generator: " + err.Error())
	}

	id, err := generator.GenerateId(ctx)
	if err != nil {
		panic("failed to generate ID: " + err.Error())
	}

	return id
}

// GenerateBase62ID 生成一个Base62编码的唯一ID
func (s *SnowflakeGenerator) GenerateBase62ID(ctx context.Context) (string, error) {
	uuID, err := s.generator.NextID()
	if err != nil {
		return "", err
	}
	return EncodeBase62(uuID)
}

func (s *SnowflakeGenerator) GenerateShardDir(ctx context.Context, id string) (string, error) {
	return GetShardDirectory(id), nil
}
