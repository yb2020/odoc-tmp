package idgen

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// UUIDV7Generator UUID v7 ID生成器
// UUID v7 是基于时间戳的UUID，具有以下特点：
// 1. 时间有序：生成的ID按时间递增，适合数据库索引
// 2. 全局唯一：无需中心化协调，适合分布式系统
// 3. 无冲突：不同客户端生成的ID不会冲突，适合离线场景
type UUIDV7Generator struct{}

var (
	uuidv7Instance *UUIDV7Generator
	uuidv7Once     sync.Once
)

// NewUUIDV7Generator 创建一个新的 UUID v7 生成器
func NewUUIDV7Generator() *UUIDV7Generator {
	uuidv7Once.Do(func() {
		uuidv7Instance = &UUIDV7Generator{}
	})
	return uuidv7Instance
}

// GetUUIDV7Generator 获取 UUID v7 生成器的单例实例
func GetUUIDV7Generator() *UUIDV7Generator {
	if uuidv7Instance == nil {
		return NewUUIDV7Generator()
	}
	return uuidv7Instance
}

// GenerateID 生成一个 UUID v7 格式的唯一ID
func (g *UUIDV7Generator) GenerateID(ctx context.Context) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// GenerateIDWithoutDash 生成一个不带横线的 UUID v7
func (g *UUIDV7Generator) GenerateIDWithoutDash(ctx context.Context) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	// 返回不带横线的格式：32个字符
	return id.String()[:8] + id.String()[9:13] + id.String()[14:18] + id.String()[19:23] + id.String()[24:], nil
}

// MustGenerateID 生成一个 UUID v7，如果出错则 panic
// 注意：此函数仅在确定不会出错的情况下使用
func (g *UUIDV7Generator) MustGenerateID(ctx context.Context) string {
	id, err := g.GenerateID(ctx)
	if err != nil {
		panic("failed to generate UUID v7: " + err.Error())
	}
	return id
}

// GenerateShardDir 根据 UUID 生成分片目录
func (g *UUIDV7Generator) GenerateShardDir(ctx context.Context, id string) string {
	return GetShardDirectory(id)
}

// GenerateUUIDV7 便捷函数：生成一个 UUID v7
func generateUUIDV7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// UUID v7 生成失败的概率极低，降级为 UUID v4
		return uuid.New().String()
	}
	return id.String()
}

// GenerateUUID 生成一个 UUID（兼容旧代码，实际生成 UUID v7）
func GenerateUUID() string {
	return strings.ReplaceAll(GeneratorUUIDUnderline(), "-", "")
}

func GeneratorUUIDUnderline() string {
	return generateUUIDV7()
}

// GenerateBase62ID 生成一个Base62编码的UUID v7
// UUID v7 是 128 位，编码后约为 22 个字符
func (g *UUIDV7Generator) GenerateBase62ID(ctx context.Context) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return EncodeUUIDToBase62(id), nil
}

// GenerateBase62UUIDV7 便捷函数：生成一个Base62编码的UUID v7
func GenerateBase62UUIDV7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// UUID v7 生成失败的概率极低，降级为 UUID v4
		return EncodeUUIDToBase62(uuid.New())
	}
	return EncodeUUIDToBase62(id)
}
