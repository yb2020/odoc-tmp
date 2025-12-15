package registry

import (
	"fmt"
	"sync"
)

// 服务提供者类型常量
const (
	ProviderLocal = "local" // 本地实现
	ProviderCloud = "cloud" // 云端实现
)

type ServiceFactory[T any] struct {
	mu           sync.RWMutex
	constructors map[string]func(any) T
	defaultImpl  string
}

// 全局服务工厂注册表
var (
	factoriesMu sync.RWMutex
	factories   = make(map[string]any) // map[serviceName]*ServiceFactory[T]
)

// getOrCreateFactory 获取或创建服务工厂
func getOrCreateFactory[T any](serviceName, defaultImpl string) *ServiceFactory[T] {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()

	if f, ok := factories[serviceName]; ok {
		return f.(*ServiceFactory[T])
	}

	factory := &ServiceFactory[T]{
		constructors: make(map[string]func(any) T),
		defaultImpl:  defaultImpl,
	}
	factories[serviceName] = factory
	return factory
}

// Register 注册服务构造函数
// 各实现在 init() 中调用，开源版本只注册 local，私有版本可额外注册 cloud
func Register[T any](serviceName, implName string, constructor func(any) T) {
	factory := getOrCreateFactory[T](serviceName, "local")
	factory.mu.Lock()
	defer factory.mu.Unlock()

	if _, exists := factory.constructors[implName]; exists {
		panic(fmt.Sprintf("service '%s': implementation '%s' already registered", serviceName, implName))
	}
	factory.constructors[implName] = constructor
}

// Create 根据配置创建服务实例
// 如果找不到指定实现，自动回退到默认实现（local）
func Create[T any](serviceName, implName string, deps any) T {
	factoriesMu.RLock()
	f, ok := factories[serviceName]
	factoriesMu.RUnlock()

	if !ok {
		panic(fmt.Sprintf("service '%s' not registered", serviceName))
	}

	factory := f.(*ServiceFactory[T])
	factory.mu.RLock()
	defer factory.mu.RUnlock()

	// 尝试获取指定实现
	if constructor, ok := factory.constructors[implName]; ok {
		return constructor(deps)
	}

	// 回退到默认实现
	if implName != factory.defaultImpl {
		if constructor, ok := factory.constructors[factory.defaultImpl]; ok {
			return constructor(deps)
		}
	}

	panic(fmt.Sprintf("service '%s': no implementation found for '%s', available: %v", serviceName, implName, listImpl[T](serviceName)))
}

// SetDefault 设置默认实现
func SetDefault[T any](serviceName, implName string) {
	factory := getOrCreateFactory[T](serviceName, implName)
	factory.mu.Lock()
	defer factory.mu.Unlock()
	factory.defaultImpl = implName
}

// listImpl 列出所有已注册的实现
func listImpl[T any](serviceName string) []string {
	factoriesMu.RLock()
	f, ok := factories[serviceName]
	factoriesMu.RUnlock()

	if !ok {
		return nil
	}

	factory := f.(*ServiceFactory[T])
	factory.mu.RLock()
	defer factory.mu.RUnlock()

	names := make([]string, 0, len(factory.constructors))
	for name := range factory.constructors {
		names = append(names, name)
	}
	return names
}
