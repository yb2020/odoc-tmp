package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	i18n "github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
)

var (
	// 全局本地化器
	globalLocalizer     i18n.Localizer
	globalLocalizerOnce sync.Once
)

// InitI18n 初始化国际化
func InitI18n(basePath string, defaultLang string, supportedLangs []string, fallbackLang string, logger logging.Logger) (i18n.Localizer, error) {
	var initErr error

	globalLocalizerOnce.Do(func() {
		// 确保使用绝对路径
		absBasePath, err := filepath.Abs(basePath)
		if err != nil {
			if logger != nil {
				logger.Error("msg", "Failed to get absolute path", "error", err.Error())
			}
			initErr = err
			return
		}

		// 尝试多个可能的路径
		possiblePaths := []string{
			filepath.Join(absBasePath, "internal/i18n"),       // 正常运行时路径
			filepath.Join(absBasePath, "../../internal/i18n"), // Debug模式路径
			filepath.Join(absBasePath, "../internal/i18n"),    // 其他可能路径
		}

		var i18nBaseDir string
		var pathFound bool

		// 检查每个可能的路径
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				i18nBaseDir = path
				pathFound = true
				if logger != nil {
					logger.Info("msg", "Found i18n base directory", "dir", i18nBaseDir)
				}
				break
			}
		}

		// 如果没有找到有效路径
		if !pathFound {
			errMsg := "i18n base directory not found in any of the expected locations"
			if logger != nil {
				logger.Error("msg", errMsg, "tried_paths", strings.Join(possiblePaths, ", "))
			}
			initErr = fmt.Errorf(errMsg)
			return
		}

		// 自动发现所有资源目录
		resourceDirs := []string{}

		// 首先添加框架目录
		frameworkDir := filepath.Join(i18nBaseDir, "framework")
		if _, err := os.Stat(frameworkDir); !os.IsNotExist(err) {
			if logger != nil {
				logger.Debug("msg", "Adding framework directory", "dir", frameworkDir)
			}
			resourceDirs = append(resourceDirs, frameworkDir)
		} else if logger != nil {
			logger.Warn("msg", "Framework directory does not exist", "dir", frameworkDir)
		}

		// 添加业务目录
		businessDir := filepath.Join(i18nBaseDir, "business")

		// 检查业务目录是否存在
		if _, err := os.Stat(businessDir); os.IsNotExist(err) {
			if logger != nil {
				logger.Warn("msg", "Business directory does not exist", "dir", businessDir)
			}
		} else {
			// 遍历业务目录下的所有子目录
			err := filepath.Walk(businessDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// 只处理目录，且不是 businessDir 自身
				if info.IsDir() && path != businessDir {
					// 检查目录中是否有 JSON 文件
					hasJsonFiles := false
					subDirEntries, err := os.ReadDir(path)
					if err == nil {
						for _, entry := range subDirEntries {
							if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
								hasJsonFiles = true
								break
							}
						}
					}

					if hasJsonFiles {
						if logger != nil {
							logger.Debug("msg", "Adding business directory with JSON files", "dir", path)
						}
						resourceDirs = append(resourceDirs, path)
					} else if logger != nil {
						logger.Debug("msg", "Skipping directory without JSON files", "dir", path)
					}
				}

				return nil
			})

			if err != nil && logger != nil {
				logger.Warn("msg", "Error scanning i18n directories", "error", err.Error())
			}
		}

		if len(resourceDirs) == 0 {
			if logger != nil {
				logger.Error("msg", "No resource directories found", "basePath", absBasePath)
			}
			initErr = fmt.Errorf("no i18n resource directories found in %s", i18nBaseDir)
			return
		}

		if logger != nil {
			logger.Info("msg", "Loading i18n resources from directories", "dirs", strings.Join(resourceDirs, ", "))
		}

		// 创建加载器
		loader := NewLoader(
			defaultLang,
			supportedLangs,
			fallbackLang,
			resourceDirs,
		)

		// 加载所有消息
		if logger != nil {
			logger.Info("msg", "Starting to load i18n messages")
		}
		if err := loader.LoadMessages(); err != nil {
			if logger != nil {
				// 检查错误是否包含重复键的信息
				if strings.Contains(err.Error(), "duplicate keys found") {
					// 只记录警告，不返回错误
					logger.Warn("msg", "Found duplicate keys in i18n messages, service will continue", "warning", err.Error())
				} else {
					// 其他错误仍然返回
					logger.Error("msg", "Failed to load i18n messages", "error", err.Error())
					initErr = err
					return
				}
			} else {
				// 如果没有 logger，检查是否是重复键错误
				if !strings.Contains(err.Error(), "duplicate keys found") {
					initErr = err
					return
				}
			}
		}
		if logger != nil {
			logger.Info("msg", "Successfully loaded i18n messages")
		}

		// 创建本地化器
		globalLocalizer = i18n.NewLocalizer(loader, logger)

		// 将本地化器设置到 pkg/i18n 的全局变量中
		i18n.SetLocalizer(globalLocalizer)

		if logger != nil {
			logger.Info("msg", "国际化资源初始化成功", "defaultLang", defaultLang)
		}
	})

	return globalLocalizer, initErr
}

// GetLocalizer 获取全局本地化器
func GetLocalizer() i18n.Localizer {
	return globalLocalizer
}
