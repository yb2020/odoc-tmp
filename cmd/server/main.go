package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	// 使用新的自动模块注册机制
	"github.com/yb2020/odoc/pkg/app"
	"github.com/yb2020/odoc/pkg/utils"
)

func main() {
	// 解析命令行参数
	var env string
	var configPath string
	var httpPort int
	var grpcPort int

	// 从环境变量中获取默认值
	defaultHTTPPort := 8081
	defaultGRPCPort := 50052

	// 如果设置了环境变量，使用环境变量的值
	if portStr := os.Getenv("HTTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			defaultHTTPPort = port
		}
	}
	if portStr := os.Getenv("GRPC_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			defaultGRPCPort = port
		}
	}
	// 兼容 PORT 环境变量
	if portStr := os.Getenv("PORT"); portStr != "" && os.Getenv("HTTP_PORT") == "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			defaultHTTPPort = port
		}
	}

	flag.StringVar(&configPath, "config", configPath, "配置文件路径")
	flag.StringVar(&env, "env", env, "运行环境 (dev, staging, prod)")
	flag.IntVar(&httpPort, "http-port", defaultHTTPPort, "HTTP 服务端口")
	flag.IntVar(&grpcPort, "grpc-port", defaultGRPCPort, "gRPC 服务端口")
	flag.Parse()

	// 如果没有指定配置文件路径，则根据环境生成
	if configPath == "" {
		if env == "" {
			env = "develop" // 默认环境
		}
		configPath = utils.GetConfigPath(env)
	}

	fmt.Printf("使用配置文件: %s\n", configPath)

	// 创建应用程序实例
	options := []app.Option{
		app.WithConfigPath(configPath),
	}

	// 添加端口配置选项
	if httpPort != defaultHTTPPort {
		options = append(options, app.WithHTTPPort(httpPort))
	}

	if grpcPort != defaultGRPCPort {
		options = append(options, app.WithGRPCPort(grpcPort))
	}

	app, err := app.NewApp(options...)
	if err != nil {
		fmt.Printf("创建应用程序失败: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	// 设置应用程序
	if err := app.Setup(); err != nil {
		fmt.Printf("设置应用程序失败: %v\n", err)
		os.Exit(1)
	}

	// 运行应用程序
	if err := app.Run(); err != nil {
		fmt.Printf("应用程序运行失败: %v\n", err)
		os.Exit(1)
	}
}
