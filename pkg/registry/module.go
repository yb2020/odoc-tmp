package registry

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yb2020/odoc/pkg/scheduler"
	"google.golang.org/grpc"
)

// Module 定义统一的模块接口
type Module interface {
	// 基本信息
	Name() string // 模块名称,不要重复

	// 路由注册
	RegisterRoutes(*gin.Engine) // 注册HTTP路由

	RegisterGRPC(server *grpc.Server) //注册grpc

	RegisterJobSchedulers(scheduler *scheduler.Scheduler) //注册Job定时任务

	// 初始化和销毁
	RegisterProviders() //注册服务提供者
	Initialize() error  // 模块初始化
	Shutdown() error    // 模块销毁（可选）
}

// ProviderRegistrar 可选接口：支持服务提供者的模块
// 模块实现此接口后，框架会在 Initialize 之前调用 RegisterProviders
type ProviderRegistrar interface {
	RegisterProviders() // 注册服务提供者
}

// ModuleGroup 定义模块组接口
type ModuleGroup interface {
	Module                // 继承基本模块接口
	SubModules() []Module // 返回所有子模块
}

// GRPCModule 定义支持gRPC的模块接口
type GRPCModule interface {
	// RegisterGRPC 注册gRPC服务
	RegisterGRPC(*grpc.Server)
}

// 全局模块注册表
var (
	moduleInstances = make(map[string]interface{})
	modulesMutex    sync.RWMutex
)

// RegisterModule 注册模块实例
func RegisterModule(name string, instance interface{}) {
	modulesMutex.Lock()
	defer modulesMutex.Unlock()
	moduleInstances[name] = instance
}

// GetModuleInstance 获取指定名称的模块实例
func GetModuleInstance(name string) interface{} {
	modulesMutex.RLock()
	defer modulesMutex.RUnlock()
	return moduleInstances[name]
}

// GetAllModuleInstances 获取所有已注册的模块实例
func GetAllModuleInstances() []interface{} {
	modulesMutex.RLock()
	defer modulesMutex.RUnlock()

	instances := make([]interface{}, 0, len(moduleInstances))
	for _, instance := range moduleInstances {
		instances = append(instances, instance)
	}
	return instances
}

// GetModules 返回所有注册的模块
func GetModules() []Module {
	var result []Module
	instances := GetAllModuleInstances()
	for _, instance := range instances {
		if module, ok := instance.(Module); ok {
			result = append(result, module)
		}
	}
	return result
}

// GetGRPCModules 返回所有支持gRPC的模块
func GetGRPCModules() []GRPCModule {
	var result []GRPCModule
	instances := GetAllModuleInstances()
	for _, instance := range instances {
		if module, ok := instance.(GRPCModule); ok {
			result = append(result, module)
		}
	}
	return result
}

// InitializeAllModules 初始化所有模块
func InitializeAllModules() error {
	for _, module := range GetModules() {
		if err := module.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownAllModules 关闭所有模块
func ShutdownAllModules() error {
	for _, module := range GetModules() {
		if err := module.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}
