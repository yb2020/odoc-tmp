# 消息队列消费者管理器 (ConsumerManager)

## 问题背景

在分布式环境下运行消息队列消费者时，我们面临以下问题：

1. 使用相同消费组ID (groupId) 时，只有一个消费者能消费数据，其他消费者空闲
2. 使用不同消费组ID时，会导致重复消费同一消息

## 解决方案

`ConsumerManager` 通过 Redis 分布式锁实现领导者选举机制，确保在多实例部署环境中，只有一个实例处理消息。

### 核心功能

- **领导者选举**：使用 Redis 分布式锁确保只有一个实例成为领导者
- **自动故障转移**：当领导者实例崩溃时，锁会自动过期，其他实例会接管消费
- **锁自动续约**：领导者定期续约锁，确保持续消费
- **优雅退出**：实例停止时会释放锁，允许其他实例接管

## 使用方法

### 1. 创建 ConsumerManager

```go
// 创建配置
config := &mq.ConsumerManagerConfig{
    LockKey:         "mq:consumer:leader:your-consumer-group",
    LockExpiry:      30 * time.Second,
    RefreshInterval: 10 * time.Second,
}

// 获取 Redis 客户端
// 注意：ConsumerManager 需要 redis.UniversalClient 类型
// 如果使用项目中的 database.RedisClient 接口，需要获取原始的 redis.Client 实例
// 例如：如果您的 Redis 客户端是 redis.Client 类型，可以直接使用
// 如果是其他类型，可能需要进行转换

// 创建 ConsumerManager
consumerManager := mq.NewConsumerManager(
    redisClient,        // Redis 客户端 (redis.UniversalClient 类型)
    consumer,           // 消息队列消费者实例
    "your-consumer-group", // 消费者组 ID
    logger,             // 日志记录器
    tracer,             // 分布式追踪器（可选）
    config,             // 配置（可选，传 nil 使用默认值）
)
```

### 2. 订阅主题

```go
// 订阅主题
consumerManager.Subscribe("your-topic", "*", func(ctx context.Context, msg mq_interface.Message) error {
    // 处理消息
    return nil
})
```

### 3. 启动和停止

```go
// 启动 ConsumerManager
consumerManager.Start(context.Background())

// 应用关闭时停止 ConsumerManager
defer consumerManager.Stop(context.Background())
```

### 4. 检查领导者状态

```go
// 检查当前实例是否是领导者
if consumerManager.IsLeader() {
    // 只有领导者才执行的操作
}
```

## 配置选项

`ConsumerManagerConfig` 提供以下配置选项：

- **LockKey**: 分布式锁的键名，默认为 `"mq:consumer:leader:" + consumerGroupID`
- **LockExpiry**: 锁的过期时间，默认为 30 秒
- **RefreshInterval**: 锁的刷新间隔，默认为 10 秒
- **RetryDelay**: 锁获取失败后的重试间隔，默认为 100 毫秒
- **LockTimeout**: 锁获取的超时时间，默认为 5 秒

## 实现原理

1. **启动流程**：
   - 创建 ConsumerManager 实例
   - 调用 Start 方法启动领导者选举循环
   - 尝试获取分布式锁
   - 获取成功后成为领导者，启动消费者并应用所有订阅

2. **领导者维护**：
   - 定期刷新锁，维持领导权
   - 如果刷新失败，停止消费者并释放锁
   - 其他实例会尝试获取锁成为新的领导者

3. **停止流程**：
   - 调用 Stop 方法发送停止信号
   - 释放锁（如果持有）
   - 如果是领导者，停止消费者

## 依赖

- `github.com/redis/go-redis/v9`: Redis 客户端
- `github.com/yb2020/go-sea/pkg/distlock`: 分布式锁实现
- `github.com/yb2020/go-sea/pkg/logging`: 日志记录
- `github.com/yb2020/go-sea/pkg/mq/interface`: 消息队列接口

## Redis 客户端要求

ConsumerManager 使用 `redis.UniversalClient` 接口与 Redis 交互。在项目中，可能需要从现有的 Redis 客户端获取原始的 `redis.Client` 或 `redis.ClusterClient` 实例。

例如，如果您使用项目中的 `database.RedisClient` 接口，您可能需要获取原始的 Redis 客户端实例：

```go
// 如果您有获取原始 Redis 客户端的方法，使用它
// 例如：
// redisClient := redisClientInstance.GetOriginalClient()

// 或者，如果您直接使用 redis.Client，可以直接使用
// redisClient := redis.NewClient(&redis.Options{
//     Addr: "localhost:6379",
// })
```

## 注意事项

1. 确保所有服务实例能访问同一个 Redis 服务器
2. 适当调整锁的过期时间和续约间隔，以平衡可靠性和故障恢复速度
3. 在生产环境中，建议使用 Redis 集群以提高可用性
