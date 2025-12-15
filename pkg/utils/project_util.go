package utils

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 缓存项目根目录，避免重复查找
var (
	projectRootCache string
	projectRootOnce  sync.Once
)

// IsProjectRoot 检查目录是否为项目根目录
func IsProjectRoot(dir string) bool {
	// 检查是否存在典型的项目文件和目录
	markers := []string{
		"go.mod", // Go 模块文件
		"config", // 配置目录
		"cmd",    // 命令目录
		"pkg",    // 包目录
	}

	for _, marker := range markers {
		path := filepath.Join(dir, marker)
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}

	return true
}

// HasGoMod 检查目录是否包含 go.mod 文件
func HasGoMod(dir string) bool {
	goModPath := filepath.Join(dir, "go.mod")
	_, err := os.Stat(goModPath)
	return err == nil
}

// IsTemporaryDir 检查目录是否为临时目录（go run 或 dlv 调试模式）
func IsTemporaryDir(dir string) bool {
	// macOS 临时目录
	if strings.Contains(dir, "/var/folders/") {
		return true
	}
	// Linux 临时目录
	if strings.HasPrefix(dir, "/tmp/") {
		return true
	}
	// Windows 临时目录
	if strings.Contains(strings.ToLower(dir), "\\temp\\") || strings.Contains(strings.ToLower(dir), "\\tmp\\") {
		return true
	}
	// go-build 目录（go run 模式）
	if strings.Contains(dir, "go-build") {
		return true
	}
	return false
}

// FindProjectRootFromDir 从指定目录向上查找项目根目录（包含 go.mod 的目录）
func FindProjectRootFromDir(startDir string) string {
	dir := startDir
	for {
		if HasGoMod(dir) {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到达文件系统根目录
			return ""
		}
		dir = parent
	}
}

// FindProjectRoot 查找项目根目录（带缓存）
func FindProjectRoot() string {
	projectRootOnce.Do(func() {
		projectRootCache = findProjectRootInternal()
	})
	return projectRootCache
}

// findProjectRootInternal 内部实现，查找项目根目录
func findProjectRootInternal() string {
	// 1. 优先使用环境变量
	if projectRoot := os.Getenv("GO_SEA_ROOT"); projectRoot != "" {
		return projectRoot
	}

	// 2. 尝试通过可执行文件路径确定
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)

		// 如果不是临时目录（非 go run / dlv 模式）
		if !IsTemporaryDir(execDir) {
			// 容器环境
			if filepath.Base(execDir) == "app" {
				return execDir
			}
			// 从可执行文件目录向上查找
			if root := FindProjectRootFromDir(execDir); root != "" {
				return root
			}
		}
	}

	// 3. 从当前工作目录向上查找
	wd, err := os.Getwd()
	if err == nil {
		if root := FindProjectRootFromDir(wd); root != "" {
			return root
		}
	}

	// 4. 容器环境备选
	if _, err := os.Stat("/app/go.mod"); err == nil {
		return "/app"
	}

	// 5. 返回当前目录作为最后备选
	if wd != "" {
		return wd
	}

	return "."
}

// ResolveRelativePath 将相对路径解析为基于项目根目录的绝对路径
// 如果已经是绝对路径，则直接返回
func ResolveRelativePath(relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(FindProjectRoot(), relativePath)
}

// GetConfigPath 获取配置文件路径
// env: 环境名称，如 "local.develop", "production" 等
func GetConfigPath(env string) string {
	return filepath.Join(FindProjectRoot(), "config", "config."+env+".yaml")
}
