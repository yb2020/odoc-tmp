# 雪花算法ID生成器

这个包提供了基于Twitter雪花算法(Snowflake)的分布式ID生成器实现，适用于go-sea微服务框架。

## 特性

- **分布式友好**：自动从机器IP地址获取数据中心ID，支持最多65535个数据中心
- **高性能**：单机每秒可生成数十万个ID
- **时间有序**：生成的ID按照时间自增排序
- **全局唯一**：在分布式系统中保证ID不会冲突
- **前端友好**：提供字符串格式的ID，避免JavaScript中的数字精度丢失问题

## 位分配结构

```
0 - 0000000000 0000000000 0000000000 0000000000 0 - 0000000000 0000 - 00000000
|                          40位时间戳                       | 16位数据中心ID | 8位序列号 |
```

- 40位时间戳：精确到毫秒，以2025-03-14为起始时间，可使用约34年
- 16位数据中心ID：支持65535个数据中心，默认使用机器IP地址的后两段计算
- 8位序列号：每毫秒支持生成256个ID

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/yb2020/go-sea/pkg/idgen"
)

func main() {
    // 获取雪花算法ID生成器
    generator, err := idgen.GetSnowflakeGenerator()
    if err != nil {
        log.Fatalf("获取ID生成器失败: %v", err)
    }
    
    // 生成uint64格式的ID
    id, err := generator.GenerateID(context.Background())
    if err != nil {
        log.Fatalf("生成ID失败: %v", err)
    }
    fmt.Printf("生成的ID (uint64): %d\n", id)
    
    // 生成字符串格式的ID（适用于前端）
    strID, err := generator.GenerateStringID(context.Background())
    if err != nil {
        log.Fatalf("生成字符串ID失败: %v", err)
    }
    fmt.Printf("生成的ID (string): %s\n", strID)
}
```

## 在用户服务中使用

在用户服务中，可以将ID生成器注入到服务层，自动为新创建的用户生成ID：

```go
// 在用户服务中使用雪花算法生成ID
func (s *DefaultUserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
    // 获取雪花算法ID生成器
    generator, err := idgen.GetSnowflakeGenerator()
    if err != nil {
        return nil, fmt.Errorf("获取ID生成器失败: %v", err)
    }
    
    // 生成字符串格式的ID
    id, err := generator.GenerateStringID(ctx)
    if err != nil {
        return nil, fmt.Errorf("生成用户ID失败: %v", err)
    }
    
    // 设置用户ID
    user.ID = id
    
    // 创建用户
    if err := s.userDAO.CreateUser(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

## 注意事项

1. 雪花算法依赖系统时钟，如果发生时钟回拨，将拒绝生成ID
2. 数据中心ID默认从机器IP地址获取，确保网络配置正确
3. 在高并发场景下，建议使用单例模式获取ID生成器，避免重复创建
